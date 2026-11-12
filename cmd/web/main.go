package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	_ "time/tzdata"

	"asetra/internal/config"
	"asetra/internal/database"
	"asetra/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := runtimeRoot()
	if err != nil {
		return err
	}
	if err := os.Chdir(root); err != nil {
		return fmt.Errorf("set runtime directory: %w", err)
	}
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return fmt.Errorf("load Jakarta timezone: %w", err)
	}
	time.Local = jakarta
	if err := godotenv.Load(filepath.Join(root, ".env")); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	var db *sql.DB
	if cfg.Database.Enabled {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		db, err = database.Open(ctx, cfg.Database)
		cancel()
		if err != nil {
			return err
		}
		defer db.Close()
	}
	r, err := router.New(cfg, db, filepath.Join(root, "web"))
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr: cfg.HTTPAddr, Handler: r,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	slog.Info("server starting", "address", cfg.HTTPAddr, "database_enabled", cfg.Database.Enabled)
	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()
			return err
		}
		return nil
	}
}

func runtimeRoot() (string, error) {
	workingDirectory, err := os.Getwd()
	if err == nil {
		if info, statErr := os.Stat(filepath.Join(workingDirectory, "web", "templates")); statErr == nil && info.IsDir() {
			return workingDirectory, nil
		}
	}
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate executable: %w", err)
	}
	executableDirectory := filepath.Dir(executable)
	if info, statErr := os.Stat(filepath.Join(executableDirectory, "web", "templates")); statErr == nil && info.IsDir() {
		return executableDirectory, nil
	}
	return "", fmt.Errorf("runtime folder not found: place the executable beside the web folder")
}

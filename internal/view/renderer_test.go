package view

import (
	"path/filepath"
	"testing"
)

func TestTemplatesParse(t *testing.T) {
	root := filepath.Join("..", "..", "web", "templates")
	if _, err := New(root); err != nil {
		t.Fatalf("parse templates: %v", err)
	}
}

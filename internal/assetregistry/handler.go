package assetregistry

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"asetra/internal/auth"
	"asetra/internal/navigation"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

type Handler struct {
	appName    string
	service    *Service
	navigation *navigation.Service
}

func NewHandler(appName string, service *Service, navigationService *navigation.Service) *Handler {
	return &Handler{appName, service, navigationService}
}

func (h *Handler) Index(c *gin.Context) {
	assets, err := h.service.List(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, http.StatusOK, "assetregistry/index.html", gin.H{"Title": "Data Aset", "Section": "Manajemen Aset", "Assets": assets, "Summary": Summarize(assets)})
}

func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.Redirect(http.StatusSeeOther, "/errors/404")
		return
	}
	asset, err := h.service.Find(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.Redirect(http.StatusSeeOther, "/errors/404")
		return
	}
	if err != nil {
		h.fail(c, err)
		return
	}
	depreciation, err := h.service.Depreciation(c.Request.Context(), asset, time.Now())
	if err != nil {
		h.fail(c, err)
		return
	}
	assetURL := absoluteAssetURL(c, asset.ID)
	h.render(c, http.StatusOK, "assetregistry/detail.html", gin.H{"Title": asset.Name, "Section": "Data Aset", "Asset": asset, "Depreciation": depreciation, "AssetURL": assetURL})
}

func (h *Handler) QRCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err = h.service.Find(c.Request.Context(), id); errors.Is(err, sql.ErrNoRows) {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		h.fail(c, err)
		return
	}
	png, err := qrcode.Encode(absoluteAssetURL(c, id), qrcode.Medium, 256)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.Header("Cache-Control", "private, max-age=3600")
	c.Data(http.StatusOK, "image/png", png)
}

func absoluteAssetURL(c *gin.Context, id uint64) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if forwarded := strings.TrimSpace(strings.Split(c.GetHeader("X-Forwarded-Proto"), ",")[0]); forwarded == "http" || forwarded == "https" {
		scheme = forwarded
	}
	return fmt.Sprintf("%s://%s/assets/%d", scheme, c.Request.Host, id)
}

func (h *Handler) render(c *gin.Context, status int, template string, data gin.H) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, "assets")
	if err != nil {
		h.fail(c, err)
		return
	}
	token, err := auth.CSRFToken(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	data["AppName"], data["Authenticated"], data["IsDashboard"] = h.appName, true, true
	data["UserName"], data["UserInitials"] = user.Name, initials(user.Name)
	data["StoreName"], data["RoleName"], data["Menus"], data["CSRFToken"] = user.StoreName, user.RoleName, menus, token
	c.HTML(status, template, data)
}

func (h *Handler) fail(c *gin.Context, err error) {
	slog.Error("asset registry failed", "error", err)
	c.Redirect(http.StatusSeeOther, "/errors/500")
}
func initials(name string) string {
	p := strings.Fields(name)
	if len(p) == 0 {
		return "US"
	}
	x := string([]rune(strings.ToUpper(p[0]))[0])
	if len(p) > 1 {
		x += string([]rune(strings.ToUpper(p[len(p)-1]))[0])
	}
	return x
}

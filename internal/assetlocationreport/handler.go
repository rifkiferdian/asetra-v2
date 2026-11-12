package assetlocationreport

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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

func NewHandler(appName string, service *Service, navigation *navigation.Service) *Handler {
	return &Handler{appName: appName, service: service, navigation: navigation}
}

func (h *Handler) Index(c *gin.Context) {
	report, err := h.service.Report(c.Request.Context())
	if err != nil {
		slog.Error("asset location report failed", "error", err)
		c.Redirect(http.StatusSeeOther, "/errors/500")
		return
	}
	h.render(c, http.StatusOK, "assetlocationreport/index.html", gin.H{"Title": "Rekap Aset per Lokasi", "Section": "Report", "Report": report})
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := locationID(c)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/errors/404")
		return
	}
	location, area, err := h.service.Area(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.Redirect(http.StatusSeeOther, "/errors/404")
		return
	}
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, http.StatusOK, "assetlocationreport/detail.html", gin.H{
		"Title": area.Name, "Section": "Report", "Location": location, "Area": area, "LocationURL": absoluteLocationURL(c, id),
	})
}

func (h *Handler) QRCode(c *gin.Context) {
	id, ok := locationID(c)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	if _, _, err := h.service.Area(c.Request.Context(), id); errors.Is(err, sql.ErrNoRows) {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		h.fail(c, err)
		return
	}
	png, err := qrcode.Encode(absoluteLocationURL(c, id), qrcode.Medium, 256)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.Header("Cache-Control", "private, max-age=3600")
	c.Data(http.StatusOK, "image/png", png)
}

func locationID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return id, err == nil && id != 0
}

func absoluteLocationURL(c *gin.Context, id uint64) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if forwarded := strings.TrimSpace(strings.Split(c.GetHeader("X-Forwarded-Proto"), ",")[0]); forwarded == "http" || forwarded == "https" {
		scheme = forwarded
	}
	return fmt.Sprintf("%s://%s/reports/assets-by-location/%d", scheme, c.Request.Host, id)
}

func (h *Handler) render(c *gin.Context, status int, template string, data gin.H) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, "asset-location-report")
	if err != nil {
		h.fail(c, err)
		return
	}
	data["AppName"], data["Authenticated"], data["IsDashboard"] = h.appName, true, true
	data["UserName"], data["UserInitials"] = user.Name, initials(user.Name)
	data["StoreName"], data["RoleName"], data["Menus"] = user.StoreName, user.RoleName, menus
	c.HTML(status, template, data)
}

func (h *Handler) fail(c *gin.Context, err error) {
	slog.Error("asset location report failed", "error", err)
	c.Redirect(http.StatusSeeOther, "/errors/500")
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "US"
	}
	value := string([]rune(strings.ToUpper(parts[0]))[0])
	if len(parts) > 1 {
		value += string([]rune(strings.ToUpper(parts[len(parts)-1]))[0])
	}
	return value
}

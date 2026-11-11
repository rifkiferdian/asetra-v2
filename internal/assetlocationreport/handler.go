package assetlocationreport

import (
	"log/slog"
	"net/http"
	"strings"

	"asetra/internal/auth"
	"asetra/internal/navigation"
	"github.com/gin-gonic/gin"
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
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, "asset-location-report")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/errors/500")
		return
	}
	c.HTML(http.StatusOK, "assetlocationreport/index.html", gin.H{
		"AppName": h.appName, "Authenticated": true, "IsDashboard": true,
		"Title": "Rekap Aset per Lokasi", "Section": "Report", "Report": report,
		"UserName": user.Name, "UserInitials": initials(user.Name), "StoreName": user.StoreName,
		"RoleName": user.RoleName, "Menus": menus,
	})
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

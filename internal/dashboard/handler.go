package dashboard

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
	navigation *navigation.Service
}

func NewHandler(appName string, navigationService *navigation.Service) *Handler {
	return &Handler{appName: appName, navigation: navigationService}
}

func (h *Handler) Index(c *gin.Context) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, "dashboard")
	if err != nil {
		slog.Error("prepare dashboard navigation", "error", err)
		c.Redirect(http.StatusSeeOther, "/errors/500")
		return
	}
	csrfToken, err := auth.CSRFToken(c)
	if err != nil {
		slog.Error("prepare dashboard csrf", "error", err)
		c.Redirect(http.StatusSeeOther, "/errors/500")
		return
	}
	c.HTML(http.StatusOK, "dashboard/index.html", gin.H{
		"Title":         "Dashboard",
		"Section":       "Ruang kerja",
		"AppName":       h.appName,
		"Authenticated": true,
		"IsDashboard":   true,
		"UserName":      user.Name,
		"UserInitials":  initials(user.Name),
		"StoreName":     user.StoreName,
		"RoleName":      user.RoleName,
		"Menus":         menus,
		"CSRFToken":     csrfToken,
	})
}

func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "US"
	}
	result := []rune(strings.ToUpper(parts[0]))[:1]
	if len(parts) > 1 {
		result = append(result, []rune(strings.ToUpper(parts[len(parts)-1]))[0])
	}
	return string(result)
}

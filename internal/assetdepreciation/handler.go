package assetdepreciation

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
	return &Handler{appName, service, navigation}
}

func (h *Handler) Index(c *gin.Context) {
	period := month(c.Query("period"))
	companyID, _ := strconv.ParseUint(c.Query("company_id"), 10, 64)
	page, err := h.service.Page(c.Request.Context(), companyID, period)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, http.StatusOK, gin.H{"Title": "Penyusutan Aset", "Section": "Manajemen Aset", "Page": page, "SuccessMessage": message(c.Query("success")), "ErrorMessage": message(c.Query("error"))})
}

func (h *Handler) Post(c *gin.Context) {
	assetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || assetID == 0 {
		h.redirect(c, "invalid")
		return
	}
	period := month(c.PostForm("period"))
	user, _ := auth.CurrentUser(c)
	userID, _ := strconv.ParseUint(user.ID, 10, 64)
	err = h.service.Post(c.Request.Context(), assetID, userID, period)
	if err != nil {
		code := "post-failed"
		if errors.Is(err, ErrDuplicate) {
			code = "already-posted"
		}
		if errors.Is(err, ErrInvalidPeriod) {
			code = "invalid-period"
		}
		if errors.Is(err, ErrFullyDepreciated) {
			code = "fully-depreciated"
		}
		h.redirect(c, code)
		return
	}
	q := url.Values{"period": {period.Format("2006-01")}, "company_id": {c.PostForm("company_id")}, "success": {"posted"}}
	c.Redirect(http.StatusSeeOther, "/asset-depreciations?"+q.Encode())
}

func (h *Handler) redirect(c *gin.Context, code string) {
	q := url.Values{"period": {c.PostForm("period")}, "company_id": {c.PostForm("company_id")}, "error": {code}}
	c.Redirect(http.StatusSeeOther, "/asset-depreciations?"+q.Encode())
}
func month(value string) time.Time {
	t, err := time.Parse("2006-01", value)
	if err != nil {
		return time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	}
	return t
}
func message(code string) string {
	m := map[string]string{"posted": "Penyusutan aset berhasil diposting.", "already-posted": "Aset tersebut sudah diposting pada periode ini.", "invalid-period": "Periode penyusutan tidak valid.", "fully-depreciated": "Aset sudah habis disusutkan.", "post-failed": "Penyusutan gagal diposting.", "invalid": "Data aset tidak valid."}
	return m[code]
}

func (h *Handler) render(c *gin.Context, status int, data gin.H) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, "asset-depreciations")
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
	c.HTML(status, "assetdepreciation/index.html", data)
}
func (h *Handler) fail(c *gin.Context, err error) {
	slog.Error("asset depreciation failed", "error", err)
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

package assetmaster

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
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

func NewHandler(appName string, service *Service, navigationService *navigation.Service) *Handler {
	return &Handler{appName: appName, service: service, navigation: navigationService}
}

var pages = map[string]struct{ Title, Single, ParentLabel, ActivePage string }{
	"categories": {"Asset Category", "category", "", "asset-categories"},
	"types":      {"Asset Type", "type", "Asset Category", "asset-types"},
	"brands":     {"Asset Brand", "brand", "", "asset-brands"},
	"models":     {"Asset Model", "model", "Asset Type", "asset-models"},
	"variants":   {"Asset Variant", "variant", "Asset Model", "asset-variants"},
}

func (h *Handler) Index(c *gin.Context) {
	kind := pageKind(c)
	page, ok := pages[kind]
	if !ok {
		c.Redirect(http.StatusSeeOther, "/errors/404")
		return
	}
	items, err := h.service.Items(c.Request.Context(), kind)
	if err != nil {
		h.fail(c, err)
		return
	}
	groups, err := h.service.Groups(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	parents, err := h.service.Parents(c.Request.Context(), kind)
	if err != nil {
		h.fail(c, err)
		return
	}
	brands, err := h.service.Brands(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	typeBrands, err := h.service.TypeBrands(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, page.ActivePage, gin.H{"Title": page.Title, "Single": page.Single, "Kind": kind, "ParentLabel": page.ParentLabel, "Items": items, "Groups": groups, "Parents": parents, "Brands": brands, "TypeBrands": typeBrands})
}
func (h *Handler) Create(c *gin.Context) { h.mutate(c, 0) }
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.redirect(c, "error", "invalid")
		return
	}
	h.mutate(c, id)
}
func (h *Handler) mutate(c *gin.Context, id uint64) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	kind := pageKind(c)
	if _, ok := pages[kind]; !ok {
		h.redirect(c, "error", "invalid")
		return
	}
	x := Input{GroupID: uintForm(c, "group_id"), ParentID: uintForm(c, "parent_id"), BrandID: uintForm(c, "brand_id"), Code: c.PostForm("code"), Name: c.PostForm("name"), Description: c.PostForm("description"), IsActive: c.PostForm("is_active") == "1"}
	var err error
	if id == 0 {
		err = h.service.Create(c.Request.Context(), kind, x)
	} else {
		err = h.service.Update(c.Request.Context(), kind, id, x)
	}
	h.result(c, err, map[bool]string{true: "created", false: "updated"}[id == 0])
}
func (h *Handler) Delete(c *gin.Context) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.redirect(c, "error", "invalid")
		return
	}
	h.result(c, h.service.Delete(c.Request.Context(), pageKind(c), id), "deleted")
}
func (h *Handler) UpdateTypeBrands(c *gin.Context) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.redirect(c, "error", "invalid")
		return
	}
	values := c.PostFormArray("brand_ids")
	brandIDs := make([]uint64, 0, len(values))
	for _, value := range values {
		parsed, parseErr := strconv.ParseUint(value, 10, 64)
		if parseErr != nil {
			h.redirect(c, "error", "invalid")
			return
		}
		brandIDs = append(brandIDs, parsed)
	}
	h.result(c, h.service.ReplaceTypeBrands(c.Request.Context(), id, brandIDs), "brands-updated")
}
func (h *Handler) UpdateBrandTypes(c *gin.Context) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.redirect(c, "error", "invalid")
		return
	}
	values := c.PostFormArray("type_ids")
	typeIDs := make([]uint64, 0, len(values))
	for _, value := range values {
		parsed, parseErr := strconv.ParseUint(value, 10, 64)
		if parseErr != nil {
			h.redirect(c, "error", "invalid")
			return
		}
		typeIDs = append(typeIDs, parsed)
	}
	h.result(c, h.service.ReplaceBrandTypes(c.Request.Context(), id, typeIDs), "types-updated")
}
func uintForm(c *gin.Context, key string) uint64 {
	v, _ := strconv.ParseUint(c.PostForm(key), 10, 64)
	return v
}
func (h *Handler) render(c *gin.Context, active string, data gin.H) {
	current, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), current, active)
	if err != nil {
		h.fail(c, err)
		return
	}
	token, err := auth.CSRFToken(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	data["AppName"], data["IsDashboard"], data["Authenticated"] = h.appName, true, true
	data["Section"] = "Master Data Asset"
	data["UserName"], data["UserInitials"] = current.Name, initials(current.Name)
	data["StoreName"], data["RoleName"], data["CSRFToken"] = current.StoreName, current.RoleName, token
	data["Menus"] = menus
	data["SuccessMessage"] = message(c.Query("success"))
	data["ErrorMessage"] = message(c.Query("error"))
	c.HTML(http.StatusOK, "assetmaster/index.html", data)
}
func (h *Handler) result(c *gin.Context, err error, success string) {
	if err == nil {
		h.redirect(c, "success", success)
		return
	}
	code := "failed"
	switch {
	case errors.Is(err, ErrInvalid):
		code = "invalid"
	case errors.Is(err, ErrDuplicate):
		code = "duplicate"
	case errors.Is(err, ErrInUse):
		code = "in-use"
	case errors.Is(err, sql.ErrNoRows):
		code = "not-found"
	default:
		slog.Error("asset master operation failed", "error", err)
	}
	h.redirect(c, "error", code)
}
func (h *Handler) redirect(c *gin.Context, key, value string) {
	kind := pageKind(c)
	if _, ok := pages[kind]; !ok {
		kind = "categories"
	}
	c.Redirect(http.StatusSeeOther, "/master/assets/"+kind+"?"+key+"="+url.QueryEscape(value))
}
func pageKind(c *gin.Context) string {
	kind, _ := c.Get("assetMasterKind")
	value, _ := kind.(string)
	return value
}
func (h *Handler) fail(c *gin.Context, err error) {
	slog.Error("asset master page failed", "error", err)
	c.Redirect(http.StatusSeeOther, "/errors/500")
}
func initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "US"
	}
	v := string([]rune(strings.ToUpper(parts[0]))[0])
	if len(parts) > 1 {
		v += string([]rune(strings.ToUpper(parts[len(parts)-1]))[0])
	}
	return v
}
func message(code string) string {
	return map[string]string{"created": "Data berhasil dibuat.", "updated": "Data berhasil diperbarui.", "deleted": "Data berhasil dihapus.", "brands-updated": "Relasi brand untuk tipe aset berhasil diperbarui.", "types-updated": "Relasi tipe untuk brand berhasil diperbarui.", "csrf": "Sesi form berakhir. Silakan muat ulang halaman.", "invalid": "Data belum lengkap atau tidak valid.", "duplicate": "Kode sudah digunakan dalam company group tersebut.", "in-use": "Relasi data tidak valid atau data masih digunakan.", "not-found": "Data tidak ditemukan.", "failed": "Operasi gagal diproses."}[code]
}

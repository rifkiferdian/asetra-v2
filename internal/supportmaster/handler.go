package supportmaster

import (
	"asetra/internal/auth"
	"asetra/internal/navigation"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Handler struct {
	appName    string
	service    *Service
	navigation *navigation.Service
}

func NewHandler(a string, s *Service, n *navigation.Service) *Handler { return &Handler{a, s, n} }

var Pages = map[string]struct{ Title, Section, Active string }{
	"conditions": {"Kondisi Aset", "Master Operasional", "asset-conditions"}, "statuses": {"Status Aset", "Master Operasional", "asset-statuses"}, "vendors": {"Vendor", "Master Operasional", "vendors"},
	"depreciation-methods": {"Metode Penyusutan", "Master Accounting", "depreciation-methods"}, "depreciation-policies": {"Kebijakan Penyusutan", "Master Accounting", "depreciation-policies"}, "gl-accounts": {"Chart of Accounts (GL)", "Master Accounting", "gl-accounts"}, "accounting-groups": {"Kelompok Accounting Aset", "Master Accounting", "asset-accounting-groups"}, "accounting-mappings": {"Mapping Akun Accounting", "Master Accounting", "asset-accounting-mappings"}}

func (h *Handler) Index(c *gin.Context) {
	kind := kind(c)
	p, ok := Pages[kind]
	if !ok {
		c.Redirect(303, "/errors/404")
		return
	}
	items, err := h.service.List(c, kind)
	if err != nil {
		h.fail(c, err)
		return
	}
	opts, err := h.service.Options(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	data := gin.H{"Title": p.Title, "Section": p.Section, "Kind": kind, "Items": items}
	for k, v := range opts {
		data[k] = v
	}
	h.render(c, p.Active, data)
}
func (h *Handler) Create(c *gin.Context) { h.save(c, 0) }
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.save(c, id)
}
func (h *Handler) save(c *gin.Context, id uint64) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	x := Input{CompanyGroupID: u(c, "company_group_id"), CompanyID: u(c, "company_id"), ParentID: u(c, "parent_id"), AssetTypeID: u(c, "asset_type_id"), MethodID: u(c, "method_id"), AccountingGroupID: u(c, "accounting_group_id"), AssetGLID: u(c, "asset_gl_id"), AccumulatedGLID: u(c, "accumulated_gl_id"), ExpenseGLID: u(c, "expense_gl_id"), Code: c.PostForm("code"), Name: c.PostForm("name"), Description: c.PostForm("description"), VendorType: c.PostForm("vendor_type"), TaxNumber: c.PostForm("tax_number"), Address: c.PostForm("address"), Phone: c.PostForm("phone"), Email: c.PostForm("email"), ContactPerson: c.PostForm("contact_person"), AccountType: c.PostForm("account_type"), ExternalAccountID: c.PostForm("external_account_id"), BookType: c.PostForm("book_type"), UsefulLifeMonths: u(c, "useful_life_months"), ResidualValuePercent: c.PostForm("residual_value_percent"), EffectiveFrom: c.PostForm("effective_from"), EffectiveTo: c.PostForm("effective_to"), IsActive: c.PostForm("is_active") == "1"}
	h.result(c, h.service.Save(c, kind(c), id, x), map[bool]string{true: "created", false: "updated"}[id == 0])
}
func (h *Handler) Delete(c *gin.Context) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, "error", "csrf")
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.result(c, h.service.Delete(c, kind(c), id), "deleted")
}
func u(c *gin.Context, k string) uint64 { v, _ := strconv.ParseUint(c.PostForm(k), 10, 64); return v }
func kind(c *gin.Context) string        { v, _ := c.Get("supportMasterKind"); s, _ := v.(string); return s }
func (h *Handler) render(c *gin.Context, active string, d gin.H) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c, user, active)
	if err != nil {
		h.fail(c, err)
		return
	}
	token, err := auth.CSRFToken(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	d["AppName"], d["IsDashboard"], d["Authenticated"] = h.appName, true, true
	d["UserName"], d["UserInitials"] = user.Name, initials(user.Name)
	d["StoreName"], d["RoleName"], d["CSRFToken"] = user.StoreName, user.RoleName, token
	d["Menus"] = menus
	d["SuccessMessage"] = message(c.Query("success"))
	d["ErrorMessage"] = message(c.Query("error"))
	c.HTML(200, "supportmaster/index.html", d)
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
	case errors.Is(err, ErrRelation):
		code = "relation"
	case errors.Is(err, ErrOverlap):
		code = "overlap"
	case errors.Is(err, sql.ErrNoRows):
		code = "not-found"
	default:
		slog.Error("support master operation failed", "error", err)
	}
	h.redirect(c, "error", code)
}
func (h *Handler) redirect(c *gin.Context, k, v string) {
	c.Redirect(http.StatusSeeOther, "/master/support/"+kind(c)+"?"+k+"="+url.QueryEscape(v))
}
func (h *Handler) fail(c *gin.Context, e error) {
	slog.Error("support master page failed", "error", e)
	c.Redirect(303, "/errors/500")
}
func initials(n string) string {
	p := strings.Fields(n)
	if len(p) == 0 {
		return "US"
	}
	v := string([]rune(strings.ToUpper(p[0]))[0])
	if len(p) > 1 {
		v += string([]rune(strings.ToUpper(p[len(p)-1]))[0])
	}
	return v
}
func message(c string) string {
	return map[string]string{"created": "Data berhasil dibuat.", "updated": "Data berhasil diperbarui.", "deleted": "Data berhasil dihapus.", "csrf": "Sesi form berakhir.", "invalid": "Data belum lengkap atau tidak valid.", "duplicate": "Kode atau kombinasi data sudah digunakan.", "in-use": "Data masih digunakan atau relasinya tidak tersedia.", "relation": "Relasi company, group, atau tipe akun tidak sesuai.", "overlap": "Periode kebijakan bertumpang tindih dengan kebijakan lain.", "not-found": "Data tidak ditemukan.", "failed": "Operasi gagal diproses."}[c]
}

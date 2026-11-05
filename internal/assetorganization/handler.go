package assetorganization

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

func (h *Handler) Companies(c *gin.Context) {
	groups, err := h.service.Groups(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	companies, err := h.service.Companies(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, "assetorganization/companies.html", "companies", gin.H{"Title": "Company", "Groups": groups, "Companies": companies, "GroupCount": len(groups), "CompanyCount": len(companies)})
}
func (h *Handler) Branches(c *gin.Context) {
	companies, err := h.service.Companies(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	branches, err := h.service.Branches(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, "assetorganization/branches.html", "branches", gin.H{"Title": "Branch", "Companies": companies, "Branches": branches, "BranchCount": len(branches)})
}
func (h *Handler) Departments(c *gin.Context) {
	groups, err := h.service.Groups(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	departments, err := h.service.Departments(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, "assetorganization/departments.html", "asset-departments", gin.H{"Title": "Department", "Groups": groups, "Departments": departments, "DepartmentCount": len(departments)})
}
func (h *Handler) Locations(c *gin.Context) {
	branches, err := h.service.Branches(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	locations, err := h.service.Locations(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	h.render(c, "assetorganization/locations.html", "asset-locations", gin.H{"Title": "Asset Location", "Branches": branches, "Locations": locations, "LocationCount": len(locations)})
}

func (h *Handler) CreateGroup(c *gin.Context) {
	h.mutate(c, "companies", func() error { return h.service.CreateGroup(c.Request.Context(), groupInput(c)) }, "group-created")
}
func (h *Handler) UpdateGroup(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "companies", "error", "invalid")
		return
	}
	h.mutate(c, "companies", func() error { return h.service.UpdateGroup(c.Request.Context(), id, groupInput(c)) }, "group-updated")
}
func (h *Handler) DeleteGroup(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "companies", "error", "invalid")
		return
	}
	h.mutate(c, "companies", func() error { return h.service.DeleteGroup(c.Request.Context(), id) }, "group-deleted")
}
func (h *Handler) CreateCompany(c *gin.Context) {
	h.mutate(c, "companies", func() error { return h.service.CreateCompany(c.Request.Context(), companyInput(c)) }, "company-created")
}
func (h *Handler) UpdateCompany(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "companies", "error", "invalid")
		return
	}
	h.mutate(c, "companies", func() error { return h.service.UpdateCompany(c.Request.Context(), id, companyInput(c)) }, "company-updated")
}
func (h *Handler) DeleteCompany(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "companies", "error", "invalid")
		return
	}
	h.mutate(c, "companies", func() error { return h.service.DeleteCompany(c.Request.Context(), id) }, "company-deleted")
}
func (h *Handler) CreateBranch(c *gin.Context) {
	h.mutate(c, "branches", func() error { return h.service.CreateBranch(c.Request.Context(), branchInput(c)) }, "branch-created")
}
func (h *Handler) UpdateBranch(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "branches", "error", "invalid")
		return
	}
	h.mutate(c, "branches", func() error { return h.service.UpdateBranch(c.Request.Context(), id, branchInput(c)) }, "branch-updated")
}
func (h *Handler) DeleteBranch(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "branches", "error", "invalid")
		return
	}
	h.mutate(c, "branches", func() error { return h.service.DeleteBranch(c.Request.Context(), id) }, "branch-deleted")
}
func (h *Handler) CreateDepartment(c *gin.Context) {
	h.mutate(c, "departments", func() error { return h.service.CreateDepartment(c.Request.Context(), departmentInput(c)) }, "department-created")
}
func (h *Handler) UpdateDepartment(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "departments", "error", "invalid")
		return
	}
	h.mutate(c, "departments", func() error { return h.service.UpdateDepartment(c.Request.Context(), id, departmentInput(c)) }, "department-updated")
}
func (h *Handler) DeleteDepartment(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "departments", "error", "invalid")
		return
	}
	h.mutate(c, "departments", func() error { return h.service.DeleteDepartment(c.Request.Context(), id) }, "department-deleted")
}
func (h *Handler) CreateLocation(c *gin.Context) {
	h.mutate(c, "asset-locations", func() error { return h.service.CreateLocation(c.Request.Context(), locationInput(c)) }, "location-created")
}
func (h *Handler) UpdateLocation(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "asset-locations", "error", "invalid")
		return
	}
	h.mutate(c, "asset-locations", func() error { return h.service.UpdateLocation(c.Request.Context(), id, locationInput(c)) }, "location-updated")
}
func (h *Handler) DeleteLocation(c *gin.Context) {
	id, ok := id(c)
	if !ok {
		h.redirect(c, "asset-locations", "error", "invalid")
		return
	}
	h.mutate(c, "asset-locations", func() error { return h.service.DeleteLocation(c.Request.Context(), id) }, "location-deleted")
}

func groupInput(c *gin.Context) GroupInput {
	return GroupInput{Code: c.PostForm("code"), Name: c.PostForm("name"), Description: c.PostForm("description"), Address: c.PostForm("address"), Phone: c.PostForm("phone"), Email: c.PostForm("email"), IsActive: c.PostForm("is_active") == "1"}
}
func companyInput(c *gin.Context) CompanyInput {
	return CompanyInput{GroupID: u64(c.PostForm("group_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), LegalName: c.PostForm("legal_name"), TaxNumber: c.PostForm("tax_number"), Address: c.PostForm("address"), Phone: c.PostForm("phone"), Email: c.PostForm("email"), IsActive: c.PostForm("is_active") == "1"}
}
func branchInput(c *gin.Context) BranchInput {
	return BranchInput{CompanyID: u64(c.PostForm("company_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), BranchType: c.PostForm("branch_type"), Address: c.PostForm("address"), City: c.PostForm("city"), Province: c.PostForm("province"), PostalCode: c.PostForm("postal_code"), Phone: c.PostForm("phone"), IsActive: c.PostForm("is_active") == "1"}
}
func departmentInput(c *gin.Context) DepartmentInput {
	return DepartmentInput{GroupID: u64(c.PostForm("group_id")), ParentID: u64(c.PostForm("parent_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), Description: c.PostForm("description"), IsActive: c.PostForm("is_active") == "1"}
}
func locationInput(c *gin.Context) LocationInput {
	return LocationInput{BranchID: u64(c.PostForm("branch_id")), ParentID: u64(c.PostForm("parent_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), LocationType: c.PostForm("location_type"), Description: c.PostForm("description"), IsActive: c.PostForm("is_active") == "1"}
}
func u64(v string) uint64 { n, _ := strconv.ParseUint(v, 10, 64); return n }
func id(c *gin.Context) (uint64, bool) {
	n, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return n, err == nil && n > 0
}

func (h *Handler) mutate(c *gin.Context, page string, operation func() error, success string) {
	if !auth.ValidCSRF(c, c.PostForm("csrf_token")) {
		h.redirect(c, page, "error", "csrf")
		return
	}
	err := operation()
	if err == nil {
		h.redirect(c, page, "success", success)
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
		slog.Error("organization mutation failed", "error", err)
	}
	h.redirect(c, page, "error", code)
}
func (h *Handler) render(c *gin.Context, template, active string, data gin.H) {
	user, _ := auth.CurrentUser(c)
	menus, err := h.navigation.Menus(c.Request.Context(), user, active)
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
	data["Section"] = "Struktur Organisasi"
	data["UserName"], data["UserInitials"] = user.Name, initials(user.Name)
	data["StoreName"], data["RoleName"], data["CSRFToken"] = user.StoreName, user.RoleName, token
	data["Menus"] = menus
	data["SuccessMessage"], data["ErrorMessage"] = message(c.Query("success")), message(c.Query("error"))
	c.HTML(http.StatusOK, template, data)
}
func (h *Handler) redirect(c *gin.Context, page, key, value string) {
	c.Redirect(http.StatusSeeOther, "/organization/"+page+"?"+key+"="+url.QueryEscape(value))
}
func (h *Handler) fail(c *gin.Context, err error) {
	slog.Error("asset organization failed", "error", err)
	c.Redirect(http.StatusSeeOther, "/errors/500")
}
func initials(name string) string {
	p := strings.Fields(name)
	if len(p) == 0 {
		return "US"
	}
	v := string([]rune(strings.ToUpper(p[0]))[0])
	if len(p) > 1 {
		v += string([]rune(strings.ToUpper(p[len(p)-1]))[0])
	}
	return v
}
func message(code string) string {
	return map[string]string{"group-created": "Company Group berhasil dibuat.", "group-updated": "Company Group berhasil diperbarui.", "group-deleted": "Company Group berhasil dihapus.", "company-created": "Company berhasil dibuat.", "company-updated": "Company berhasil diperbarui.", "company-deleted": "Company berhasil dihapus.", "branch-created": "Branch berhasil dibuat.", "branch-updated": "Branch berhasil diperbarui.", "branch-deleted": "Branch berhasil dihapus.", "department-created": "Department berhasil dibuat.", "department-updated": "Department berhasil diperbarui.", "department-deleted": "Department berhasil dihapus.", "location-created": "Asset Location berhasil dibuat.", "location-updated": "Asset Location berhasil diperbarui.", "location-deleted": "Asset Location berhasil dihapus.", "csrf": "Sesi form berakhir. Muat ulang halaman.", "invalid": "Data belum lengkap atau tidak valid.", "duplicate": "Kode sudah digunakan pada induk yang sama.", "in-use": "Data masih digunakan atau relasi induknya tidak sesuai.", "not-found": "Data tidak ditemukan.", "failed": "Operasi gagal diproses."}[code]
}

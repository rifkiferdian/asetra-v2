package assetorganization

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrInvalid   = errors.New("invalid organization data")
	ErrDuplicate = errors.New("duplicate organization data")
	ErrInUse     = errors.New("organization data in use")
)
var codePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{1,29}$`)
var branchTypes = map[string]bool{"HEAD_OFFICE": true, "STORE": true, "WAREHOUSE": true, "OFFICE": true, "DISTRIBUTION_CENTER": true, "OTHER": true}
var locationTypes = map[string]bool{"BUILDING": true, "FLOOR": true, "ROOM": true, "AREA": true, "WAREHOUSE": true, "OTHER": true}

type Service struct{ repository *Repository }

func NewService(r *Repository) *Service { return &Service{repository: r} }
func (s *Service) Groups(ctx context.Context) ([]CompanyGroup, error) {
	return s.repository.Groups(ctx)
}
func (s *Service) Companies(ctx context.Context) ([]Company, error) {
	return s.repository.Companies(ctx)
}
func (s *Service) Branches(ctx context.Context) ([]Branch, error) { return s.repository.Branches(ctx) }
func (s *Service) Departments(ctx context.Context) ([]Department, error) {
	return s.repository.Departments(ctx)
}
func (s *Service) Locations(ctx context.Context) ([]AssetLocation, error) {
	return s.repository.Locations(ctx)
}

func (s *Service) CreateGroup(ctx context.Context, x GroupInput) error {
	x = normalizeGroup(x)
	if !validBase(x.Code, x.Name) || !validEmail(x.Email) {
		return ErrInvalid
	}
	return mapDB(s.repository.CreateGroup(ctx, x))
}
func (s *Service) UpdateGroup(ctx context.Context, id uint64, x GroupInput) error {
	x = normalizeGroup(x)
	if id == 0 || !validBase(x.Code, x.Name) || !validEmail(x.Email) {
		return ErrInvalid
	}
	return mapDB(s.repository.UpdateGroup(ctx, id, x))
}
func (s *Service) DeleteGroup(ctx context.Context, id uint64) error {
	return mapDB(s.repository.DeleteGroup(ctx, id))
}
func (s *Service) CreateCompany(ctx context.Context, x CompanyInput) error {
	x = normalizeCompany(x)
	if x.GroupID == 0 || !validBase(x.Code, x.Name) || !validEmail(x.Email) {
		return ErrInvalid
	}
	return mapDB(s.repository.CreateCompany(ctx, x))
}
func (s *Service) UpdateCompany(ctx context.Context, id uint64, x CompanyInput) error {
	x = normalizeCompany(x)
	if id == 0 || x.GroupID == 0 || !validBase(x.Code, x.Name) || !validEmail(x.Email) {
		return ErrInvalid
	}
	return mapDB(s.repository.UpdateCompany(ctx, id, x))
}
func (s *Service) DeleteCompany(ctx context.Context, id uint64) error {
	return mapDB(s.repository.DeleteCompany(ctx, id))
}
func (s *Service) CreateBranch(ctx context.Context, x BranchInput) error {
	x = normalizeBranch(x)
	if x.CompanyID == 0 || !validBase(x.Code, x.Name) || !branchTypes[x.BranchType] {
		return ErrInvalid
	}
	return mapDB(s.repository.CreateBranch(ctx, x))
}
func (s *Service) UpdateBranch(ctx context.Context, id uint64, x BranchInput) error {
	x = normalizeBranch(x)
	if id == 0 || x.CompanyID == 0 || !validBase(x.Code, x.Name) || !branchTypes[x.BranchType] {
		return ErrInvalid
	}
	return mapDB(s.repository.UpdateBranch(ctx, id, x))
}
func (s *Service) DeleteBranch(ctx context.Context, id uint64) error {
	return mapDB(s.repository.DeleteBranch(ctx, id))
}
func (s *Service) CreateDepartment(ctx context.Context, x DepartmentInput) error {
	x = normalizeDepartment(x)
	if x.GroupID == 0 || !validBase(x.Code, x.Name) {
		return ErrInvalid
	}
	return mapDB(s.repository.CreateDepartment(ctx, x))
}
func (s *Service) UpdateDepartment(ctx context.Context, id uint64, x DepartmentInput) error {
	x = normalizeDepartment(x)
	if id == 0 || x.GroupID == 0 || x.ParentID == id || !validBase(x.Code, x.Name) {
		return ErrInvalid
	}
	return mapDB(s.repository.UpdateDepartment(ctx, id, x))
}
func (s *Service) DeleteDepartment(ctx context.Context, id uint64) error {
	return mapDB(s.repository.DeleteDepartment(ctx, id))
}
func (s *Service) CreateLocation(ctx context.Context, x LocationInput) error {
	x = normalizeLocation(x)
	if x.BranchID == 0 || !validBase(x.Code, x.Name) || !locationTypes[x.LocationType] {
		return ErrInvalid
	}
	return mapDB(s.repository.CreateLocation(ctx, x))
}
func (s *Service) UpdateLocation(ctx context.Context, id uint64, x LocationInput) error {
	x = normalizeLocation(x)
	if id == 0 || x.BranchID == 0 || x.ParentID == id || !validBase(x.Code, x.Name) || !locationTypes[x.LocationType] {
		return ErrInvalid
	}
	return mapDB(s.repository.UpdateLocation(ctx, id, x))
}
func (s *Service) DeleteLocation(ctx context.Context, id uint64) error {
	return mapDB(s.repository.DeleteLocation(ctx, id))
}

func validBase(code, name string) bool {
	return codePattern.MatchString(code) && name != "" && len(name) <= 250
}
func validEmail(v string) bool {
	if v == "" {
		return true
	}
	_, err := mail.ParseAddress(v)
	return err == nil && len(v) <= 150
}
func clean(v string) string { return strings.TrimSpace(v) }
func upper(v string) string { return strings.ToUpper(clean(v)) }
func normalizeGroup(x GroupInput) GroupInput {
	x.Code = upper(x.Code)
	x.Name = clean(x.Name)
	x.Description = clean(x.Description)
	x.Address = clean(x.Address)
	x.Phone = clean(x.Phone)
	x.Email = strings.ToLower(clean(x.Email))
	return x
}
func normalizeCompany(x CompanyInput) CompanyInput {
	x.Code = upper(x.Code)
	x.Name = clean(x.Name)
	x.LegalName = clean(x.LegalName)
	x.TaxNumber = clean(x.TaxNumber)
	x.Address = clean(x.Address)
	x.Phone = clean(x.Phone)
	x.Email = strings.ToLower(clean(x.Email))
	return x
}
func normalizeBranch(x BranchInput) BranchInput {
	x.Code = upper(x.Code)
	x.Name = clean(x.Name)
	x.BranchType = upper(x.BranchType)
	x.Address = clean(x.Address)
	x.City = clean(x.City)
	x.Province = clean(x.Province)
	x.PostalCode = clean(x.PostalCode)
	x.Phone = clean(x.Phone)
	return x
}
func normalizeDepartment(x DepartmentInput) DepartmentInput {
	x.Code = upper(x.Code)
	x.Name = clean(x.Name)
	x.Description = clean(x.Description)
	return x
}
func normalizeLocation(x LocationInput) LocationInput {
	x.Code = upper(x.Code)
	x.Name = clean(x.Name)
	x.LocationType = upper(x.LocationType)
	x.Description = clean(x.Description)
	return x
}
func mapDB(err error) error {
	var e *mysql.MySQLError
	if errors.As(err, &e) {
		if e.Number == 1062 {
			return ErrDuplicate
		}
		if e.Number == 1451 || e.Number == 1452 {
			return ErrInUse
		}
	}
	return err
}
func branchTypeName(v string) string {
	return map[string]string{"HEAD_OFFICE": "Head Office", "STORE": "Store", "WAREHOUSE": "Warehouse", "OFFICE": "Office", "DISTRIBUTION_CENTER": "Distribution Center", "OTHER": "Lainnya"}[v]
}
func locationTypeName(v string) string {
	return map[string]string{"BUILDING": "Gedung", "FLOOR": "Lantai", "ROOM": "Ruangan", "AREA": "Area", "WAREHOUSE": "Gudang", "OTHER": "Lainnya"}[v]
}

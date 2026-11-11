package supportmaster

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalid   = errors.New("invalid support master")
	ErrDuplicate = errors.New("duplicate support master")
	ErrInUse     = errors.New("support master in use")
	ErrRelation  = errors.New("invalid relation")
	ErrOverlap   = errors.New("overlapping policy")
)
var codePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,39}$`)

type Service struct{ repository *Repository }

func NewService(r *Repository) *Service { return &Service{repository: r} }
func (s *Service) List(ctx context.Context, kind string) ([]map[string]any, error) {
	return s.repository.List(ctx, kind)
}
func (s *Service) Options(ctx context.Context) (map[string][]Option, error) {
	return s.repository.Options(ctx)
}
func (s *Service) Save(ctx context.Context, kind string, id uint64, x Input) error {
	x = normalize(x)
	if !valid(kind, x) {
		return ErrInvalid
	}
	if err := s.repository.Validate(ctx, kind, id, x); err != nil {
		return err
	}
	return mapDB(s.repository.Save(ctx, kind, id, x))
}
func (s *Service) Delete(ctx context.Context, kind string, id uint64) error {
	if id == 0 {
		return ErrInvalid
	}
	return mapDB(s.repository.Delete(ctx, kind, id))
}
func normalize(x Input) Input {
	x.Code = strings.ToUpper(strings.TrimSpace(x.Code))
	x.Name = strings.TrimSpace(x.Name)
	x.Description = strings.TrimSpace(x.Description)
	x.VendorType = strings.ToUpper(strings.TrimSpace(x.VendorType))
	x.AccountType = strings.ToUpper(strings.TrimSpace(x.AccountType))
	x.BookType = strings.ToUpper(strings.TrimSpace(x.BookType))
	x.Email = strings.ToLower(strings.TrimSpace(x.Email))
	x.ResidualValuePercent = strings.TrimSpace(x.ResidualValuePercent)
	x.AnnualRatePercent = strings.TrimSpace(x.AnnualRatePercent)
	return x
}
func valid(kind string, x Input) bool {
	base := func() bool { return codePattern.MatchString(x.Code) && x.Name != "" }
	switch kind {
	case "conditions", "statuses", "depreciation-methods":
		return base()
	case "depreciation-rates":
		rate, err := strconv.ParseFloat(x.AnnualRatePercent, 64)
		if !base() || x.UsefulLifeMonths == 0 || err != nil || rate <= 0 || rate > 100 {
			return false
		}
		expected := 1200 / float64(x.UsefulLifeMonths)
		return rate >= expected-0.01 && rate <= expected+0.01
	case "vendors":
		return x.CompanyGroupID > 0 && base() && (x.VendorType == "SUPPLIER" || x.VendorType == "SERVICE" || x.VendorType == "BOTH")
	case "gl-accounts":
		return x.CompanyID > 0 && base() && map[string]bool{"ASSET": true, "CONTRA_ASSET": true, "LIABILITY": true, "EQUITY": true, "REVENUE": true, "EXPENSE": true}[x.AccountType]
	case "depreciation-policies":
		if x.CompanyGroupID == 0 || x.AssetTypeID == 0 || x.MethodID == 0 || x.DepreciationRateID == 0 || !map[string]bool{"ACCOUNTING": true, "TAX": true}[x.BookType] {
			return false
		}
		rv, err := strconv.ParseFloat(x.ResidualValuePercent, 64)
		if err != nil || rv < 0 || rv > 100 {
			return false
		}
		from, e := time.Parse("2006-01-02", x.EffectiveFrom)
		if e != nil {
			return false
		}
		if x.EffectiveTo != "" {
			to, e := time.Parse("2006-01-02", x.EffectiveTo)
			if e != nil || to.Before(from) {
				return false
			}
		}
		return true
	case "accounting-mappings":
		return x.CompanyID > 0 && x.DepreciationPolicyID > 0 && x.GLAccountID > 0
	}
	return false
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

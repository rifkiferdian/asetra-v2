package assetregistry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service             { return &Service{repository: repository} }
func (s *Service) List(ctx context.Context) ([]Asset, error) { return s.repository.List(ctx) }
func (s *Service) Find(ctx context.Context, id uint64) (Asset, error) {
	return s.repository.Find(ctx, id)
}

func (s *Service) FormOptions(ctx context.Context) (AssetFormOptions, error) {
	return s.repository.FormOptions(ctx)
}

func (s *Service) Create(ctx context.Context, x AssetInput) (uint64, error) {
	x.AssetNumber = strings.TrimSpace(x.AssetNumber)
	x.Name = strings.TrimSpace(x.Name)
	x.SerialNumber = strings.TrimSpace(x.SerialNumber)
	x.Barcode = strings.TrimSpace(x.Barcode)
	x.Description = strings.TrimSpace(x.Description)
	x.WarrantyUntil = strings.TrimSpace(x.WarrantyUntil)
	x.Notes = strings.TrimSpace(x.Notes)
	if x.CompanyID == 0 || x.BranchID == 0 || x.LocationID == 0 || x.TypeID == 0 || x.ConditionID == 0 || x.StatusID == 0 || x.AssetNumber == "" || x.Name == "" {
		return 0, fmt.Errorf("field wajib belum lengkap")
	}
	if len(x.AssetNumber) > 50 || len(x.Name) > 200 || len(x.SerialNumber) > 120 || len(x.Barcode) > 120 {
		return 0, fmt.Errorf("nomor aset, nama, serial number, atau barcode terlalu panjang")
	}
	if (x.ModelID > 0 && x.BrandID == 0) || (x.VariantID > 0 && x.ModelID == 0) {
		return 0, fmt.Errorf("brand harus dipilih sebelum model, dan model sebelum varian")
	}
	if _, err := time.Parse("2006-01-02", x.AcquisitionDate); err != nil {
		return 0, fmt.Errorf("tanggal perolehan tidak valid")
	}
	if x.WarrantyUntil != "" {
		warranty, err := time.Parse("2006-01-02", x.WarrantyUntil)
		if err != nil {
			return 0, fmt.Errorf("tanggal garansi tidak valid")
		}
		acquired, _ := time.Parse("2006-01-02", x.AcquisitionDate)
		if warranty.Before(acquired) {
			return 0, fmt.Errorf("tanggal garansi tidak boleh sebelum tanggal perolehan")
		}
	}
	cost, err := strconv.ParseFloat(x.AcquisitionCost, 64)
	if err != nil || cost < 0 {
		return 0, fmt.Errorf("nilai perolehan tidak valid")
	}
	id, err := s.repository.Create(ctx, x)
	if err == nil {
		return id, nil
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return 0, fmt.Errorf("nomor aset atau barcode sudah digunakan")
	}
	if err.Error() == "invalid asset relations" {
		return 0, fmt.Errorf("pilihan organisasi atau master aset tidak sesuai")
	}
	return 0, err
}

func (s *Service) Depreciation(ctx context.Context, asset Asset, asOf time.Time) (DepreciationDetail, error) {
	policy, err := s.repository.FindDepreciationPolicy(ctx, asset.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return DepreciationDetail{}, nil
	}
	if err != nil {
		return DepreciationDetail{}, err
	}
	return calculateDepreciation(asset, policy, asOf)
}

func calculateDepreciation(asset Asset, policy DepreciationPolicy, asOf time.Time) (DepreciationDetail, error) {
	cost, err := strconv.ParseFloat(asset.AcquisitionCost, 64)
	if err != nil {
		return DepreciationDetail{}, err
	}
	residualPercent, err := strconv.ParseFloat(policy.ResidualPercent, 64)
	if err != nil || policy.UsefulLifeMonths == 0 {
		return DepreciationDetail{}, fmt.Errorf("invalid depreciation policy")
	}
	acquired, err := time.Parse("2006-01-02", asset.AcquisitionDate)
	if err != nil {
		return DepreciationDetail{}, err
	}
	elapsed := completedMonths(acquired, asOf)
	if elapsed > policy.UsefulLifeMonths {
		elapsed = policy.UsefulLifeMonths
	}
	residualValue := cost * residualPercent / 100
	depreciable := math.Max(0, cost-residualValue)
	monthly := depreciable / float64(policy.UsefulLifeMonths)
	accumulated := math.Min(depreciable, monthly*float64(elapsed))
	period := policy.EffectiveFrom + " — seterusnya"
	if policy.EffectiveTo != "" {
		period = policy.EffectiveFrom + " — " + policy.EffectiveTo
	}
	return DepreciationDetail{
		Available: true, PolicyName: policy.AssetTypeName, MethodName: policy.MethodName,
		BookType: policy.BookType, GLAccount: policy.GLAccount,
		UsefulLife: fmt.Sprintf("%d bulan (%s)", policy.UsefulLifeMonths, monthYears(policy.UsefulLifeMonths)),
		AnnualRate: formatPercent(policy.AnnualRatePercent), ResidualPercent: formatPercent(policy.ResidualPercent),
		EffectivePeriod: period, MonthlyDepreciation: rupiah(monthly),
		AccumulatedDepreciation: rupiah(accumulated), BookValue: rupiah(cost - accumulated),
		AsOfDate: asOf.Format("2006-01-02"), ElapsedMonths: elapsed,
	}, nil
}

func completedMonths(from, to time.Time) uint64 {
	if to.Before(from) {
		return 0
	}
	months := (to.Year()-from.Year())*12 + int(to.Month()-from.Month())
	if to.Day() < from.Day() {
		months--
	}
	if months < 0 {
		return 0
	}
	return uint64(months)
}

func rupiah(value float64) string {
	n := strconv.FormatInt(int64(math.Round(value)), 10)
	for i := len(n) - 3; i > 0; i -= 3 {
		n = n[:i] + "." + n[i:]
	}
	return "Rp" + n
}

func formatPercent(value string) string {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value + "%"
	}
	formatted := strconv.FormatFloat(number, 'f', 2, 64)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return strings.ReplaceAll(formatted, ".", ",") + "%"
}

func monthYears(months uint64) string {
	years, remainder := months/12, months%12
	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d tahun", years))
	}
	if remainder > 0 {
		parts = append(parts, fmt.Sprintf("%d bulan", remainder))
	}
	return strings.Join(parts, " ")
}

func Summarize(assets []Asset) Summary {
	x := Summary{Total: len(assets)}
	for _, asset := range assets {
		switch asset.StatusCode {
		case "MAINTENANCE":
			x.Maintenance++
		case "RETIRED", "DISPOSED", "INACTIVE":
			x.Inactive++
		default:
			x.Active++
		}
	}
	return x
}

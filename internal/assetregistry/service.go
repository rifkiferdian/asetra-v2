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
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service             { return &Service{repository: repository} }
func (s *Service) List(ctx context.Context) ([]Asset, error) { return s.repository.List(ctx) }
func (s *Service) Find(ctx context.Context, id uint64) (Asset, error) {
	return s.repository.Find(ctx, id)
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

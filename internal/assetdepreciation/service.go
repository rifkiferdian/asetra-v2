package assetdepreciation

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrDuplicate        = errors.New("depreciation already posted")
	ErrInvalidPeriod    = errors.New("invalid depreciation period")
	ErrFullyDepreciated = errors.New("asset fully depreciated")
)

type Service struct{ repository *Repository }

func NewService(repository *Repository) *Service { return &Service{repository: repository} }

func (s *Service) Page(ctx context.Context, companyID uint64, period time.Time) (Page, error) {
	companies, err := s.repository.Companies(ctx)
	if err != nil {
		return Page{}, err
	}
	rows, err := s.repository.List(ctx, companyID, period)
	if err != nil {
		return Page{}, err
	}
	var posted, totalBook float64
	postedCount := 0
	for i := range rows {
		monthly, _ := strconv.ParseFloat(rows[i].MonthlyAmount, 64)
		book, _ := strconv.ParseFloat(rows[i].BookValue, 64)
		totalBook += book
		if rows[i].PostedAmount != "" {
			v, _ := strconv.ParseFloat(rows[i].PostedAmount, 64)
			posted += v
			postedCount++
		}
		rows[i].CanPost = rows[i].PostedAmount == "" && monthly > 0
		rows[i].AcquisitionCost = formatMoney(rows[i].AcquisitionCost)
		rows[i].Accumulated = formatMoney(rows[i].Accumulated)
		rows[i].BookValue = formatMoney(rows[i].BookValue)
		rows[i].MonthlyAmount = formatMoney(rows[i].MonthlyAmount)
		if rows[i].PostedAmount != "" {
			rows[i].PostedAmount = formatMoney(rows[i].PostedAmount)
		}
	}
	return Page{Companies: companies, Rows: rows, Period: period.Format("2006-01"), PeriodLabel: period.Format("01/2006"), MaxPeriod: time.Now().Format("2006-01"), CompanyID: companyID, TotalAssets: len(rows), PostedAssets: postedCount, TotalPosted: formatMoney(strconv.FormatFloat(posted, 'f', 2, 64)), TotalBookValue: formatMoney(strconv.FormatFloat(totalBook, 'f', 2, 64))}, nil
}

func (s *Service) Post(ctx context.Context, assetID, userID uint64, period time.Time) error {
	_, err := s.repository.Post(ctx, assetID, userID, period)
	if err == nil {
		return nil
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		return ErrDuplicate
	}
	switch err.Error() {
	case "future period", "period before acquisition":
		return ErrInvalidPeriod
	case "fully depreciated":
		return ErrFullyDepreciated
	}
	return err
}

func monthlyAmount(b postingBasis, accumulated float64) float64 {
	base := math.Max(0, b.AcquisitionCost-(b.AcquisitionCost*b.ResidualPercent/100))
	remaining := math.Max(0, base-accumulated)
	if b.UsefulLifeMonths == 0 {
		return 0
	}
	monthly := math.Round((base/float64(b.UsefulLifeMonths))*100) / 100
	return math.Min(monthly, remaining)
}

func formatMoney(value string) string {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}
	n := strconv.FormatInt(int64(math.Round(v)), 10)
	for i := len(n) - 3; i > 0; i -= 3 {
		n = n[:i] + "." + n[i:]
	}
	return "Rp " + strings.TrimSpace(n)
}

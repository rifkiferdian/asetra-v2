package assetdepreciation

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Companies(ctx context.Context) ([]Company, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name FROM companies WHERE is_active=1 ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Company
	for rows.Next() {
		var x Company
		if err := rows.Scan(&x.ID, &x.Name); err != nil {
			return nil, err
		}
		result = append(result, x)
	}
	return result, rows.Err()
}

func (r *Repository) List(ctx context.Context, companyID uint64, period time.Time) ([]Row, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT a.id,a.asset_number,COALESCE(NULLIF(CONCAT_WS(' ',NULLIF(br.name,''),NULLIF(m.name,''),NULLIF(vr.name,'')),''),t.name),c.name,
		DATE_FORMAT(a.acquisition_date,'%Y-%m-%d'),CAST(a.acquisition_cost AS CHAR),
		p.id,r.useful_life_months,CAST(p.residual_value_percent AS CHAR),
		CAST(COALESCE(hist.accumulated,0) AS CHAR),
		CAST(GREATEST(a.acquisition_cost-COALESCE(hist.accumulated,0),0) AS CHAR),
		CAST(LEAST(ROUND((a.acquisition_cost-(a.acquisition_cost*p.residual_value_percent/100))/r.useful_life_months,2),GREATEST((a.acquisition_cost-(a.acquisition_cost*p.residual_value_percent/100))-COALESCE(hist.accumulated,0),0)) AS CHAR),
		COALESCE(CAST(current_period.depreciation_amount AS CHAR),''),COALESCE(DATE_FORMAT(current_period.posted_at,'%Y-%m-%d %H:%i'),'')
	FROM assets a
	JOIN companies c ON c.id=a.company_id
	JOIN asset_types t ON t.id=a.asset_type_id
	LEFT JOIN asset_brands br ON br.id=a.asset_brand_id
	LEFT JOIN asset_models m ON m.id=a.asset_model_id
	LEFT JOIN asset_variants vr ON vr.id=a.asset_variant_id
	JOIN asset_depreciation_policies p ON p.id=(
		SELECT p2.id FROM asset_depreciation_policies p2
		WHERE p2.asset_type_id=a.asset_type_id AND p2.company_group_id=c.group_id
		  AND (p2.company_id=a.company_id OR p2.company_id IS NULL)
		  AND p2.book_type='ACCOUNTING' AND p2.is_active=1
		  AND p2.effective_from<=a.acquisition_date
		  AND (p2.effective_to IS NULL OR p2.effective_to>=a.acquisition_date)
		ORDER BY (p2.company_id=a.company_id) DESC,p2.effective_from DESC,p2.id DESC LIMIT 1)
	JOIN depreciation_rates r ON r.id=p.depreciation_rate_id
	LEFT JOIN (SELECT asset_id,SUM(depreciation_amount) accumulated FROM asset_depreciations WHERE status='POSTED' GROUP BY asset_id) hist ON hist.asset_id=a.id
	LEFT JOIN asset_depreciations current_period ON current_period.asset_id=a.id AND current_period.period_date=? AND current_period.status='POSTED'
	WHERE (?=0 OR a.company_id=?) AND a.acquisition_date<=LAST_DAY(?)
	ORDER BY c.name,a.asset_number`, period, companyID, companyID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Row
	for rows.Next() {
		var x Row
		if err := rows.Scan(&x.AssetID, &x.AssetNumber, &x.AssetName, &x.CompanyName, &x.AcquisitionDate, &x.AcquisitionCost, &x.PolicyID, &x.UsefulLifeMonths, &x.ResidualPercent, &x.Accumulated, &x.BookValue, &x.MonthlyAmount, &x.PostedAmount, &x.PostedAt); err != nil {
			return nil, err
		}
		result = append(result, x)
	}
	return result, rows.Err()
}

func (r *Repository) Post(ctx context.Context, assetID, userID uint64, period time.Time) (float64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var b postingBasis
	err = tx.QueryRowContext(ctx, `SELECT p.id,r.useful_life_months,a.acquisition_date,a.acquisition_cost,p.residual_value_percent
		FROM assets a JOIN companies c ON c.id=a.company_id
		JOIN asset_depreciation_policies p ON p.id=(SELECT p2.id FROM asset_depreciation_policies p2
		 WHERE p2.asset_type_id=a.asset_type_id AND p2.company_group_id=c.group_id
		 AND (p2.company_id=a.company_id OR p2.company_id IS NULL) AND p2.book_type='ACCOUNTING' AND p2.is_active=1
		 AND p2.effective_from<=a.acquisition_date AND (p2.effective_to IS NULL OR p2.effective_to>=a.acquisition_date)
		 ORDER BY (p2.company_id=a.company_id) DESC,p2.effective_from DESC,p2.id DESC LIMIT 1)
		JOIN depreciation_rates r ON r.id=p.depreciation_rate_id WHERE a.id=? FOR UPDATE`, assetID).Scan(&b.PolicyID, &b.UsefulLifeMonths, &b.AcquisitionDate, &b.AcquisitionCost, &b.ResidualPercent)
	if err != nil {
		return 0, err
	}
	if period.After(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)) {
		return 0, fmt.Errorf("future period")
	}
	acquisitionMonth := time.Date(b.AcquisitionDate.Year(), b.AcquisitionDate.Month(), 1, 0, 0, 0, 0, time.Local)
	if period.Before(acquisitionMonth) {
		return 0, fmt.Errorf("period before acquisition")
	}
	var accumulated float64
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(depreciation_amount),0) FROM asset_depreciations WHERE asset_id=? AND status='POSTED'`, assetID).Scan(&accumulated); err != nil {
		return 0, err
	}
	amount := monthlyAmount(b, accumulated)
	if amount <= 0 {
		return 0, fmt.Errorf("fully depreciated")
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO asset_depreciations(asset_id,depreciation_policy_id,period_date,depreciation_amount,status,posted_by) VALUES(?,?,?,?, 'POSTED',?)`, assetID, b.PolicyID, period, amount, userID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return amount, nil
}

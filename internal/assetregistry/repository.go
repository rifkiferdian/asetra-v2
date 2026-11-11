package assetregistry

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const assetSelect = `SELECT a.id,a.asset_number,a.name,COALESCE(a.serial_number,''),COALESCE(a.barcode,''),COALESCE(a.description,''),
c.name,b.name,COALESCE(d.name,'Belum ditentukan'),l.name,COALESCE(CONCAT(parent.name,' / ',l.name),l.name),ac.name,t.name,COALESCE(br.name,''),COALESCE(m.name,''),COALESCE(vr.name,''),
cond.code,cond.name,st.code,st.name,COALESCE(v.name,''),DATE_FORMAT(a.acquisition_date,'%Y-%m-%d'),CAST(a.acquisition_cost AS CHAR),COALESCE(DATE_FORMAT(a.warranty_until,'%Y-%m-%d'),''),COALESCE(a.notes,''),
DATE_FORMAT(a.created_at,'%Y-%m-%d %H:%i'),DATE_FORMAT(a.updated_at,'%Y-%m-%d %H:%i')
FROM assets a
JOIN companies c ON c.id=a.company_id
JOIN branches b ON b.id=a.branch_id
LEFT JOIN departments d ON d.id=a.department_id
JOIN asset_locations l ON l.id=a.asset_location_id
LEFT JOIN asset_locations parent ON parent.id=l.parent_id
JOIN asset_types t ON t.id=a.asset_type_id
JOIN asset_categories ac ON ac.id=t.asset_category_id
LEFT JOIN asset_brands br ON br.id=a.asset_brand_id
LEFT JOIN asset_models m ON m.id=a.asset_model_id
LEFT JOIN asset_variants vr ON vr.id=a.asset_variant_id
JOIN asset_conditions cond ON cond.id=a.condition_id
JOIN asset_statuses st ON st.id=a.status_id
LEFT JOIN vendors v ON v.id=a.vendor_id`

func (r *Repository) List(ctx context.Context) ([]Asset, error) {
	rows, err := r.db.QueryContext(ctx, assetSelect+` ORDER BY a.asset_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []Asset
	for rows.Next() {
		var x Asset
		if err := scanAsset(rows, &x); err != nil {
			return nil, err
		}
		assets = append(assets, x)
	}
	return assets, rows.Err()
}

func (r *Repository) Find(ctx context.Context, id uint64) (Asset, error) {
	var x Asset
	err := scanAsset(r.db.QueryRowContext(ctx, assetSelect+` WHERE a.id=?`, id), &x)
	return x, err
}

func (r *Repository) FindDepreciationPolicy(ctx context.Context, assetID uint64) (DepreciationPolicy, error) {
	var x DepreciationPolicy
	err := r.db.QueryRowContext(ctx, `SELECT p.id,t.name,m.code,m.name,p.book_type,r.useful_life_months,
		CAST(p.residual_value_percent AS CHAR),CAST(r.annual_rate_percent AS CHAR),
		DATE_FORMAT(p.effective_from,'%Y-%m-%d'),COALESCE(DATE_FORMAT(p.effective_to,'%Y-%m-%d'),''),
		COALESCE(CONCAT(gl.code,' — ',gl.name),'Belum dimapping')
		FROM assets a
		JOIN companies c ON c.id=a.company_id
		JOIN asset_depreciation_policies p ON p.asset_type_id=a.asset_type_id
		 AND p.company_group_id=c.group_id
		 AND (p.company_id=a.company_id OR p.company_id IS NULL)
		 AND p.book_type='ACCOUNTING' AND p.is_active=1
		 AND p.effective_from<=a.acquisition_date
		 AND (p.effective_to IS NULL OR p.effective_to>=a.acquisition_date)
		JOIN asset_types t ON t.id=p.asset_type_id
		JOIN depreciation_methods m ON m.id=p.depreciation_method_id
		JOIN depreciation_rates r ON r.id=p.depreciation_rate_id
		LEFT JOIN asset_accounting_mappings am ON am.company_id=a.company_id AND am.depreciation_policy_id=p.id AND am.is_active=1
		LEFT JOIN gl_accounts gl ON gl.id=am.gl_account_id
		WHERE a.id=?
		ORDER BY (p.company_id=a.company_id) DESC,p.effective_from DESC,p.id DESC LIMIT 1`, assetID).Scan(
		&x.ID, &x.AssetTypeName, &x.MethodCode, &x.MethodName, &x.BookType, &x.UsefulLifeMonths,
		&x.ResidualPercent, &x.AnnualRatePercent, &x.EffectiveFrom, &x.EffectiveTo, &x.GLAccount)
	return x, err
}

type scanner interface{ Scan(...any) error }

func scanAsset(row scanner, x *Asset) error {
	return row.Scan(&x.ID, &x.AssetNumber, &x.Name, &x.SerialNumber, &x.Barcode, &x.Description,
		&x.CompanyName, &x.BranchName, &x.DepartmentName, &x.LocationName, &x.LocationPath, &x.CategoryName, &x.TypeName, &x.BrandName, &x.ModelName, &x.VariantName,
		&x.ConditionCode, &x.ConditionName, &x.StatusCode, &x.StatusName, &x.VendorName, &x.AcquisitionDate, &x.AcquisitionCost, &x.WarrantyUntil, &x.Notes, &x.CreatedAt, &x.UpdatedAt)
}

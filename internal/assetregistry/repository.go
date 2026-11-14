package assetregistry

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const assetSelect = `SELECT a.id,a.asset_number,COALESCE(NULLIF(CONCAT_WS(' ',NULLIF(br.name,''),NULLIF(m.name,''),NULLIF(vr.name,'')),''),t.name),COALESCE(a.serial_number,''),COALESCE(a.barcode,''),COALESCE(a.description,''),
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

func (r *Repository) FormOptions(ctx context.Context) (AssetFormOptions, error) {
	var result AssetFormOptions
	queries := []struct {
		destination *[]FormOption
		query       string
	}{
		{&result.Companies, `SELECT id,group_id,0,0,code,name,name FROM companies WHERE is_active=1 ORDER BY name`},
		{&result.Branches, `SELECT id,company_id,0,0,code,name,CONCAT(code,' — ',name) FROM branches WHERE is_active=1 ORDER BY name`},
		{&result.Departments, `SELECT id,0,group_id,0,code,name,CONCAT(code,' — ',name) FROM departments WHERE is_active=1 ORDER BY name`},
		{&result.Locations, `SELECT l.id,l.branch_id,0,0,l.code,l.name,CONCAT(b.code,' / ',l.code,' — ',l.name) FROM asset_locations l JOIN branches b ON b.id=l.branch_id WHERE l.is_active=1 ORDER BY b.name,l.name`},
		{&result.Types, `SELECT t.id,0,c.group_id,0,t.code,t.name,CONCAT(c.name,' — ',t.name) FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE t.is_active=1 AND c.is_active=1 ORDER BY c.name,t.name`},
		{&result.Brands, `SELECT b.id,tb.asset_type_id,b.group_id,0,b.code,b.name,b.name FROM asset_type_brands tb JOIN asset_brands b ON b.id=tb.asset_brand_id WHERE b.is_active=1 ORDER BY b.name`},
		{&result.Models, `SELECT m.id,m.asset_type_id,b.group_id,m.asset_brand_id,m.code,m.name,CONCAT(t.name,' — ',b.name,' ',m.name) FROM asset_models m JOIN asset_types t ON t.id=m.asset_type_id JOIN asset_brands b ON b.id=m.asset_brand_id JOIN asset_type_brands tb ON tb.asset_type_id=m.asset_type_id AND tb.asset_brand_id=m.asset_brand_id WHERE m.is_active=1 ORDER BY t.name,m.name`},
		{&result.Variants, `SELECT v.id,v.asset_model_id,b.group_id,0,v.code,v.name,CONCAT(m.name,' — ',v.name) FROM asset_variants v JOIN asset_models m ON m.id=v.asset_model_id JOIN asset_brands b ON b.id=m.asset_brand_id WHERE v.is_active=1 ORDER BY m.name,v.name`},
		{&result.Conditions, `SELECT id,0,0,0,code,name,name FROM asset_conditions WHERE is_active=1 ORDER BY name`},
		{&result.Statuses, `SELECT id,0,0,0,code,name,name FROM asset_statuses WHERE is_active=1 ORDER BY name`},
		{&result.Vendors, `SELECT id,0,0,0,code,name,name FROM vendors WHERE is_active=1 ORDER BY name`},
	}
	for _, item := range queries {
		rows, err := r.db.QueryContext(ctx, item.query)
		if err != nil {
			return AssetFormOptions{}, err
		}
		for rows.Next() {
			var x FormOption
			if err := rows.Scan(&x.ID, &x.ParentID, &x.GroupID, &x.RelatedID, &x.Code, &x.Name, &x.Label); err != nil {
				rows.Close()
				return AssetFormOptions{}, err
			}
			*item.destination = append(*item.destination, x)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return AssetFormOptions{}, err
		}
	}
	policyRows, err := r.db.QueryContext(ctx, `SELECT COALESCE(p.company_id,0),p.company_group_id,p.asset_type_id,r.useful_life_months,m.name,
		CAST(r.annual_rate_percent AS CHAR),CAST(p.residual_value_percent AS CHAR),DATE_FORMAT(p.effective_from,'%Y-%m-%d'),COALESCE(DATE_FORMAT(p.effective_to,'%Y-%m-%d'),'')
		FROM asset_depreciation_policies p JOIN depreciation_rates r ON r.id=p.depreciation_rate_id JOIN depreciation_methods m ON m.id=p.depreciation_method_id
		WHERE p.book_type='ACCOUNTING' AND p.is_active=1 ORDER BY p.effective_from DESC,p.id DESC`)
	if err != nil {
		return AssetFormOptions{}, err
	}
	defer policyRows.Close()
	for policyRows.Next() {
		var x DepreciationPreview
		if err := policyRows.Scan(&x.CompanyID, &x.GroupID, &x.TypeID, &x.UsefulLifeMonths, &x.MethodName, &x.Rate, &x.Residual, &x.EffectiveFrom, &x.EffectiveTo); err != nil {
			return AssetFormOptions{}, err
		}
		result.DepreciationPolicies = append(result.DepreciationPolicies, x)
	}
	if err := policyRows.Err(); err != nil {
		return AssetFormOptions{}, err
	}
	return result, nil
}

func (r *Repository) Create(ctx context.Context, x AssetInput) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	var valid bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(
		SELECT 1 FROM companies c
		JOIN branches b ON b.id=? AND b.company_id=c.id
		JOIN asset_locations l ON l.id=? AND l.branch_id=b.id
		JOIN asset_types t ON t.id=? JOIN asset_categories ac ON ac.id=t.asset_category_id AND ac.group_id=c.group_id
		JOIN asset_conditions cond ON cond.id=? JOIN asset_statuses st ON st.id=?
		WHERE c.id=?
		AND (?=0 OR EXISTS(SELECT 1 FROM departments d WHERE d.id=? AND d.group_id=c.group_id))
		AND (?=0 OR EXISTS(SELECT 1 FROM asset_type_brands tb JOIN asset_brands br ON br.id=tb.asset_brand_id WHERE tb.asset_type_id=t.id AND tb.asset_brand_id=? AND br.group_id=c.group_id))
		AND (?=0 OR EXISTS(SELECT 1 FROM asset_models m WHERE m.id=? AND m.asset_type_id=t.id AND (?=0 OR m.asset_brand_id=?)))
		AND (?=0 OR EXISTS(SELECT 1 FROM asset_variants v WHERE v.id=? AND v.asset_model_id=?))
		AND (?=0 OR EXISTS(SELECT 1 FROM vendors vnd WHERE vnd.id=?))
	)`, x.BranchID, x.LocationID, x.TypeID, x.ConditionID, x.StatusID, x.CompanyID,
		x.DepartmentID, x.DepartmentID, x.BrandID, x.BrandID, x.ModelID, x.ModelID, x.BrandID, x.BrandID,
		x.VariantID, x.VariantID, x.ModelID, x.VendorID, x.VendorID).Scan(&valid)
	if err != nil {
		return 0, err
	}
	if !valid {
		return 0, fmt.Errorf("invalid asset relations")
	}
	var companyCode, typeCode string
	if err := tx.QueryRowContext(ctx, `SELECT c.code,t.code,COALESCE(NULLIF(CONCAT_WS(' ',NULLIF(br.name,''),NULLIF(m.name,''),NULLIF(v.name,'')),''),t.name)
		FROM companies c JOIN asset_types t ON t.id=?
		LEFT JOIN asset_brands br ON br.id=NULLIF(?,0)
		LEFT JOIN asset_models m ON m.id=NULLIF(?,0)
		LEFT JOIN asset_variants v ON v.id=NULLIF(?,0)
		WHERE c.id=?`, x.TypeID, x.BrandID, x.ModelID, x.VariantID, x.CompanyID).Scan(&companyCode, &typeCode, &x.Name); err != nil {
		return 0, err
	}
	year := x.acquisitionYear()
	if _, err := tx.ExecContext(ctx, `INSERT INTO asset_number_sequences(company_id,asset_type_id,sequence_year,last_number)
		VALUES(?,?,?,LAST_INSERT_ID(1)) ON DUPLICATE KEY UPDATE last_number=LAST_INSERT_ID(last_number+1)`, x.CompanyID, x.TypeID, year); err != nil {
		return 0, err
	}
	var sequence uint64
	if err := tx.QueryRowContext(ctx, `SELECT LAST_INSERT_ID()`).Scan(&sequence); err != nil {
		return 0, err
	}
	x.AssetNumber = generatedAssetNumber(companyCode, typeCode, year, sequence)
	x.Barcode = generatedBarcode(x.AssetNumber)
	result, err := tx.ExecContext(ctx, `INSERT INTO assets(
		company_id,branch_id,asset_location_id,department_id,asset_type_id,asset_brand_id,asset_model_id,asset_variant_id,
		condition_id,status_id,vendor_id,asset_number,name,serial_number,barcode,description,acquisition_date,acquisition_cost,warranty_until,notes
	) VALUES(?,?,?,NULLIF(?,0),?,NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),?,?,NULLIF(?,0),?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),?,?,NULLIF(?,''),NULLIF(?,''))`,
		x.CompanyID, x.BranchID, x.LocationID, x.DepartmentID, x.TypeID, x.BrandID, x.ModelID, x.VariantID,
		x.ConditionID, x.StatusID, x.VendorID, x.AssetNumber, x.Name, x.SerialNumber, x.Barcode, x.Description,
		x.AcquisitionDate, x.AcquisitionCost, x.WarrantyUntil, x.Notes)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return uint64(id), nil
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

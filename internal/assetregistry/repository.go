package assetregistry

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

const assetSelect = `SELECT a.id,a.asset_number,a.name,COALESCE(a.serial_number,''),COALESCE(a.barcode,''),COALESCE(a.description,''),
c.name,b.name,l.name,COALESCE(CONCAT(parent.name,' / ',l.name),l.name),ac.name,t.name,COALESCE(br.name,''),COALESCE(m.name,''),COALESCE(vr.name,''),
cond.code,cond.name,st.code,st.name,COALESCE(v.name,''),DATE_FORMAT(a.acquisition_date,'%Y-%m-%d'),CAST(a.acquisition_cost AS CHAR),COALESCE(DATE_FORMAT(a.warranty_until,'%Y-%m-%d'),''),COALESCE(a.notes,''),
DATE_FORMAT(a.created_at,'%Y-%m-%d %H:%i'),DATE_FORMAT(a.updated_at,'%Y-%m-%d %H:%i')
FROM assets a
JOIN companies c ON c.id=a.company_id
JOIN branches b ON b.id=a.branch_id
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

type scanner interface{ Scan(...any) error }

func scanAsset(row scanner, x *Asset) error {
	return row.Scan(&x.ID, &x.AssetNumber, &x.Name, &x.SerialNumber, &x.Barcode, &x.Description,
		&x.CompanyName, &x.BranchName, &x.LocationName, &x.LocationPath, &x.CategoryName, &x.TypeName, &x.BrandName, &x.ModelName, &x.VariantName,
		&x.ConditionCode, &x.ConditionName, &x.StatusCode, &x.StatusName, &x.VendorName, &x.AcquisitionDate, &x.AcquisitionCost, &x.WarrantyUntil, &x.Notes, &x.CreatedAt, &x.UpdatedAt)
}

package assetlocationreport

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context) ([]Asset, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT a.id,a.asset_number,a.name,t.name,
		COALESCE(br.name,''),COALESCE(m.name,''),cond.name,st.name,b.name,
		COALESCE(parent.name,l.name),CASE WHEN parent.id IS NULL THEN 'Area Utama' ELSE l.name END
		FROM assets a
		JOIN branches b ON b.id=a.branch_id
		JOIN asset_locations l ON l.id=a.asset_location_id
		LEFT JOIN asset_locations parent ON parent.id=l.parent_id
		JOIN asset_types t ON t.id=a.asset_type_id
		LEFT JOIN asset_brands br ON br.id=a.asset_brand_id
		LEFT JOIN asset_models m ON m.id=a.asset_model_id
		JOIN asset_conditions cond ON cond.id=a.condition_id
		JOIN asset_statuses st ON st.id=a.status_id
		ORDER BY b.name,COALESCE(parent.name,l.name),CASE WHEN parent.id IS NULL THEN '' ELSE l.name END,a.asset_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []Asset
	for rows.Next() {
		var asset Asset
		if err := rows.Scan(&asset.ID, &asset.Number, &asset.Name, &asset.TypeName, &asset.BrandName,
			&asset.ModelName, &asset.Condition, &asset.Status, &asset.BranchName, &asset.Location, &asset.Area); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, rows.Err()
}

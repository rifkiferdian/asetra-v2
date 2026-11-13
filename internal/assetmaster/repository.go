package assetmaster

import (
	"context"
	"database/sql"
	"strings"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Groups(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT id,id,0,name,name FROM company_groups WHERE is_active=1 ORDER BY name`)
}
func (r *Repository) Categories(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT c.id,c.group_id,0,c.name,g.name FROM asset_categories c JOIN company_groups g ON g.id=c.group_id WHERE c.is_active=1 ORDER BY g.name,c.name`)
}
func (r *Repository) Types(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT t.id,c.group_id,0,t.name,g.name FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id JOIN company_groups g ON g.id=c.group_id WHERE t.is_active=1 ORDER BY g.name,t.name`)
}
func (r *Repository) Brands(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT b.id,b.group_id,0,b.name,g.name FROM asset_brands b JOIN company_groups g ON g.id=b.group_id WHERE b.is_active=1 ORDER BY g.name,b.name`)
}
func (r *Repository) TypeBrands(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT b.id,b.group_id,tb.asset_type_id,b.name,g.name FROM asset_type_brands tb JOIN asset_brands b ON b.id=tb.asset_brand_id JOIN company_groups g ON g.id=b.group_id WHERE b.is_active=1 ORDER BY g.name,b.name`)
}
func (r *Repository) Models(ctx context.Context) ([]Option, error) {
	return r.options(ctx, `SELECT m.id,c.group_id,0,m.name,g.name FROM asset_models m JOIN asset_types t ON t.id=m.asset_type_id JOIN asset_categories c ON c.id=t.asset_category_id JOIN company_groups g ON g.id=c.group_id WHERE m.is_active=1 ORDER BY g.name,m.name`)
}
func (r *Repository) options(ctx context.Context, query string) ([]Option, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Option
	for rows.Next() {
		var x Option
		if err := rows.Scan(&x.ID, &x.GroupID, &x.ParentID, &x.Name, &x.GroupName); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}

func (r *Repository) RelationsValid(ctx context.Context, kind string, x Input) (bool, error) {
	var count int
	var err error
	switch kind {
	case "categories", "brands":
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM company_groups WHERE id=?`, x.GroupID).Scan(&count)
	case "types":
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_categories WHERE id=?`, x.ParentID).Scan(&count)
	case "models":
		if x.BrandID == 0 {
			err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_types WHERE id=?`, x.ParentID).Scan(&count)
		} else {
			err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_type_brands WHERE asset_type_id=? AND asset_brand_id=?`, x.ParentID, x.BrandID).Scan(&count)
		}
	case "variants":
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_models WHERE id=?`, x.ParentID).Scan(&count)
	default:
		return false, ErrInvalid
	}
	return count == 1, err
}

func (r *Repository) Items(ctx context.Context, kind string) ([]Item, error) {
	queries := map[string]string{
		"categories": `SELECT x.id,x.group_id,0,0,g.name,'','',x.code,x.name,COALESCE(x.description,''),'',x.is_active FROM asset_categories x JOIN company_groups g ON g.id=x.group_id ORDER BY g.name,x.name`,
		"types":      `SELECT x.id,c.group_id,x.asset_category_id,0,g.name,c.name,COALESCE(GROUP_CONCAT(b.name ORDER BY b.name SEPARATOR ', '),''),x.code,x.name,COALESCE(x.description,''),COALESCE(GROUP_CONCAT(b.id ORDER BY b.name),''),x.is_active FROM asset_types x JOIN asset_categories c ON c.id=x.asset_category_id JOIN company_groups g ON g.id=c.group_id LEFT JOIN asset_type_brands tb ON tb.asset_type_id=x.id LEFT JOIN asset_brands b ON b.id=tb.asset_brand_id GROUP BY x.id,c.group_id,x.asset_category_id,g.name,c.name,x.code,x.name,x.description,x.is_active ORDER BY g.name,c.name,x.name`,
		"brands":     `SELECT x.id,x.group_id,0,0,g.name,'','',x.code,x.name,COALESCE(x.description,''),'',x.is_active FROM asset_brands x JOIN company_groups g ON g.id=x.group_id ORDER BY g.name,x.name`,
		"models":     `SELECT x.id,c.group_id,x.asset_type_id,COALESCE(x.asset_brand_id,0),g.name,t.name,COALESCE(b.name,''),x.code,x.name,COALESCE(x.description,''),'',x.is_active FROM asset_models x JOIN asset_types t ON t.id=x.asset_type_id JOIN asset_categories c ON c.id=t.asset_category_id JOIN company_groups g ON g.id=c.group_id LEFT JOIN asset_brands b ON b.id=x.asset_brand_id ORDER BY g.name,t.name,x.name`,
		"variants":   `SELECT x.id,c.group_id,x.asset_model_id,0,g.name,m.name,'',x.code,x.name,COALESCE(x.description,''),'',x.is_active FROM asset_variants x JOIN asset_models m ON m.id=x.asset_model_id JOIN asset_types t ON t.id=m.asset_type_id JOIN asset_categories c ON c.id=t.asset_category_id JOIN company_groups g ON g.id=c.group_id ORDER BY g.name,m.name,x.name`,
	}
	query, ok := queries[kind]
	if !ok {
		return nil, ErrInvalid
	}
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var x Item
		if err := rows.Scan(&x.ID, &x.GroupID, &x.ParentID, &x.BrandID, &x.GroupName, &x.ParentName, &x.BrandName, &x.Code, &x.Name, &x.Description, &x.BrandIDs, &x.IsActive); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}

func (r *Repository) ReplaceTypeBrands(ctx context.Context, typeID uint64, brandIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var groupID uint64
	if err := tx.QueryRowContext(ctx, `SELECT c.group_id FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE t.id=? FOR UPDATE`, typeID).Scan(&groupID); err != nil {
		return err
	}
	if len(brandIDs) > 0 {
		placeholders := make([]string, len(brandIDs))
		args := make([]any, 0, len(brandIDs)+1)
		args = append(args, groupID)
		for i, id := range brandIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		var count int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_brands WHERE group_id=? AND id IN (`+strings.Join(placeholders, ",")+`)`, args...).Scan(&count); err != nil {
			return err
		}
		if count != len(brandIDs) {
			return ErrInvalid
		}
	}
	selected := make(map[uint64]bool, len(brandIDs))
	for _, id := range brandIDs {
		selected[id] = true
	}
	usedRows, err := tx.QueryContext(ctx, `SELECT asset_brand_id FROM asset_models WHERE asset_type_id=? AND asset_brand_id IS NOT NULL UNION SELECT asset_brand_id FROM assets WHERE asset_type_id=? AND asset_brand_id IS NOT NULL`, typeID, typeID)
	if err != nil {
		return err
	}
	for usedRows.Next() {
		var id uint64
		if err := usedRows.Scan(&id); err != nil {
			usedRows.Close()
			return err
		}
		if !selected[id] {
			usedRows.Close()
			return ErrInUse
		}
	}
	if err := usedRows.Err(); err != nil {
		usedRows.Close()
		return err
	}
	usedRows.Close()
	if _, err := tx.ExecContext(ctx, `DELETE FROM asset_type_brands WHERE asset_type_id=?`, typeID); err != nil {
		return err
	}
	for _, brandID := range brandIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO asset_type_brands(asset_type_id,asset_brand_id) VALUES(?,?)`, typeID, brandID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) Create(ctx context.Context, kind string, x Input) error {
	queries := map[string]string{
		"categories": `INSERT INTO asset_categories(group_id,code,name,description,is_active) VALUES(?,?,?,NULLIF(?,''),?)`,
		"types":      `INSERT INTO asset_types(asset_category_id,code,name,description,is_active) VALUES(?,?,?,NULLIF(?,''),?)`,
		"brands":     `INSERT INTO asset_brands(group_id,code,name,description,is_active) VALUES(?,?,?,NULLIF(?,''),?)`,
		"models":     `INSERT INTO asset_models(asset_type_id,asset_brand_id,code,name,description,is_active) VALUES(?,NULLIF(?,0),?,?,NULLIF(?,''),?)`,
		"variants":   `INSERT INTO asset_variants(asset_model_id,code,name,description,is_active) VALUES(?,?,?,NULLIF(?,''),?)`,
	}
	var args []any
	switch kind {
	case "categories", "brands":
		args = []any{x.GroupID, x.Code, x.Name, x.Description, x.IsActive}
	case "types", "variants":
		args = []any{x.ParentID, x.Code, x.Name, x.Description, x.IsActive}
	case "models":
		args = []any{x.ParentID, x.BrandID, x.Code, x.Name, x.Description, x.IsActive}
	default:
		return ErrInvalid
	}
	_, err := r.db.ExecContext(ctx, queries[kind], args...)
	return err
}
func (r *Repository) Update(ctx context.Context, kind string, id uint64, x Input) error {
	queries := map[string]string{
		"categories": `UPDATE asset_categories SET group_id=?,code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`,
		"types":      `UPDATE asset_types SET asset_category_id=?,code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`,
		"brands":     `UPDATE asset_brands SET group_id=?,code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`,
		"models":     `UPDATE asset_models SET asset_type_id=?,asset_brand_id=NULLIF(?,0),code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`,
		"variants":   `UPDATE asset_variants SET asset_model_id=?,code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`,
	}
	var args []any
	switch kind {
	case "categories", "brands":
		args = []any{x.GroupID, x.Code, x.Name, x.Description, x.IsActive, id}
	case "types", "variants":
		args = []any{x.ParentID, x.Code, x.Name, x.Description, x.IsActive, id}
	case "models":
		args = []any{x.ParentID, x.BrandID, x.Code, x.Name, x.Description, x.IsActive, id}
	default:
		return ErrInvalid
	}
	return affected(r.db.ExecContext(ctx, queries[kind], args...))
}
func (r *Repository) Delete(ctx context.Context, kind string, id uint64) error {
	tables := map[string]string{"categories": "asset_categories", "types": "asset_types", "brands": "asset_brands", "models": "asset_models", "variants": "asset_variants"}
	table, ok := tables[kind]
	if !ok {
		return ErrInvalid
	}
	return affected(r.db.ExecContext(ctx, "DELETE FROM "+table+" WHERE id=?", id))
}
func affected(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

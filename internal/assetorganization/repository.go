package assetorganization

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Groups(ctx context.Context) ([]CompanyGroup, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT g.id,g.code,g.name,COALESCE(g.description,''),COALESCE(g.address,''),COALESCE(g.phone,''),COALESCE(g.email,''),g.is_active,COUNT(c.id) FROM company_groups g LEFT JOIN companies c ON c.group_id=g.id GROUP BY g.id ORDER BY g.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CompanyGroup
	for rows.Next() {
		var x CompanyGroup
		if err := rows.Scan(&x.ID, &x.Code, &x.Name, &x.Description, &x.Address, &x.Phone, &x.Email, &x.IsActive, &x.CompanyCount); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}
func (r *Repository) Companies(ctx context.Context) ([]Company, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id,c.group_id,g.name,c.code,c.name,COALESCE(c.legal_name,''),COALESCE(c.tax_number,''),COALESCE(c.address,''),COALESCE(c.phone,''),COALESCE(c.email,''),c.is_active,COUNT(b.id) FROM companies c JOIN company_groups g ON g.id=c.group_id LEFT JOIN branches b ON b.company_id=c.id GROUP BY c.id ORDER BY g.name,c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Company
	for rows.Next() {
		var x Company
		if err := rows.Scan(&x.ID, &x.GroupID, &x.GroupName, &x.Code, &x.Name, &x.LegalName, &x.TaxNumber, &x.Address, &x.Phone, &x.Email, &x.IsActive, &x.BranchCount); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}
func (r *Repository) Branches(ctx context.Context) ([]Branch, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT b.id,b.company_id,c.name,b.code,b.name,b.branch_type,COALESCE(b.address,''),COALESCE(b.city,''),COALESCE(b.province,''),COALESCE(b.postal_code,''),COALESCE(b.phone,''),b.is_active,COUNT(l.id) FROM branches b JOIN companies c ON c.id=b.company_id LEFT JOIN asset_locations l ON l.branch_id=b.id GROUP BY b.id ORDER BY c.name,b.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Branch
	for rows.Next() {
		var x Branch
		if err := rows.Scan(&x.ID, &x.CompanyID, &x.CompanyName, &x.Code, &x.Name, &x.BranchType, &x.Address, &x.City, &x.Province, &x.PostalCode, &x.Phone, &x.IsActive, &x.LocationCount); err != nil {
			return nil, err
		}
		x.BranchTypeName = branchTypeName(x.BranchType)
		items = append(items, x)
	}
	return items, rows.Err()
}
func (r *Repository) Departments(ctx context.Context) ([]Department, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT d.id,d.group_id,g.name,COALESCE(d.parent_id,0),COALESCE(p.name,''),d.code,d.name,COALESCE(d.description,''),d.is_active,COUNT(ch.id) FROM departments d JOIN company_groups g ON g.id=d.group_id LEFT JOIN departments p ON p.id=d.parent_id LEFT JOIN departments ch ON ch.parent_id=d.id GROUP BY d.id ORDER BY g.name,d.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Department
	for rows.Next() {
		var x Department
		if err := rows.Scan(&x.ID, &x.GroupID, &x.GroupName, &x.ParentID, &x.ParentName, &x.Code, &x.Name, &x.Description, &x.IsActive, &x.ChildCount); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}
func (r *Repository) Locations(ctx context.Context) ([]AssetLocation, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT l.id,l.branch_id,b.name,COALESCE(l.parent_id,0),COALESCE(p.name,''),l.code,l.name,l.location_type,COALESCE(l.description,''),l.is_active,COUNT(ch.id) FROM asset_locations l JOIN branches b ON b.id=l.branch_id LEFT JOIN asset_locations p ON p.id=l.parent_id LEFT JOIN asset_locations ch ON ch.parent_id=l.id GROUP BY l.id ORDER BY b.name,l.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssetLocation
	for rows.Next() {
		var x AssetLocation
		if err := rows.Scan(&x.ID, &x.BranchID, &x.BranchName, &x.ParentID, &x.ParentName, &x.Code, &x.Name, &x.LocationType, &x.Description, &x.IsActive, &x.ChildCount); err != nil {
			return nil, err
		}
		x.LocationTypeName = locationTypeName(x.LocationType)
		items = append(items, x)
	}
	return items, rows.Err()
}

func (r *Repository) CreateGroup(ctx context.Context, x GroupInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO company_groups(code,name,description,address,phone,email,is_active) VALUES(?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),?)`, x.Code, x.Name, x.Description, x.Address, x.Phone, x.Email, x.IsActive)
	return err
}
func (r *Repository) UpdateGroup(ctx context.Context, id uint64, x GroupInput) error {
	return affected(r.db.ExecContext(ctx, `UPDATE company_groups SET code=?,name=?,description=NULLIF(?,''),address=NULLIF(?,''),phone=NULLIF(?,''),email=NULLIF(?,''),is_active=? WHERE id=?`, x.Code, x.Name, x.Description, x.Address, x.Phone, x.Email, x.IsActive, id))
}
func (r *Repository) DeleteGroup(ctx context.Context, id uint64) error {
	return affected(r.db.ExecContext(ctx, `DELETE FROM company_groups WHERE id=?`, id))
}
func (r *Repository) CreateCompany(ctx context.Context, x CompanyInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO companies(group_id,code,name,legal_name,tax_number,address,phone,email,is_active) VALUES(?,?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),?)`, x.GroupID, x.Code, x.Name, x.LegalName, x.TaxNumber, x.Address, x.Phone, x.Email, x.IsActive)
	return err
}
func (r *Repository) UpdateCompany(ctx context.Context, id uint64, x CompanyInput) error {
	return affected(r.db.ExecContext(ctx, `UPDATE companies SET group_id=?,code=?,name=?,legal_name=NULLIF(?,''),tax_number=NULLIF(?,''),address=NULLIF(?,''),phone=NULLIF(?,''),email=NULLIF(?,''),is_active=? WHERE id=?`, x.GroupID, x.Code, x.Name, x.LegalName, x.TaxNumber, x.Address, x.Phone, x.Email, x.IsActive, id))
}
func (r *Repository) DeleteCompany(ctx context.Context, id uint64) error {
	return affected(r.db.ExecContext(ctx, `DELETE FROM companies WHERE id=?`, id))
}
func (r *Repository) CreateBranch(ctx context.Context, x BranchInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO branches(company_id,code,name,branch_type,address,city,province,postal_code,phone,is_active) VALUES(?,?,?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),?)`, x.CompanyID, x.Code, x.Name, x.BranchType, x.Address, x.City, x.Province, x.PostalCode, x.Phone, x.IsActive)
	return err
}
func (r *Repository) UpdateBranch(ctx context.Context, id uint64, x BranchInput) error {
	return affected(r.db.ExecContext(ctx, `UPDATE branches SET company_id=?,code=?,name=?,branch_type=?,address=NULLIF(?,''),city=NULLIF(?,''),province=NULLIF(?,''),postal_code=NULLIF(?,''),phone=NULLIF(?,''),is_active=? WHERE id=?`, x.CompanyID, x.Code, x.Name, x.BranchType, x.Address, x.City, x.Province, x.PostalCode, x.Phone, x.IsActive, id))
}
func (r *Repository) DeleteBranch(ctx context.Context, id uint64) error {
	return affected(r.db.ExecContext(ctx, `DELETE FROM branches WHERE id=?`, id))
}
func (r *Repository) CreateDepartment(ctx context.Context, x DepartmentInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO departments(group_id,parent_id,code,name,description,is_active) VALUES(?,NULLIF(?,0),?,?,NULLIF(?,''),?)`, x.GroupID, x.ParentID, x.Code, x.Name, x.Description, x.IsActive)
	return err
}
func (r *Repository) UpdateDepartment(ctx context.Context, id uint64, x DepartmentInput) error {
	return affected(r.db.ExecContext(ctx, `UPDATE departments SET group_id=?,parent_id=NULLIF(?,0),code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`, x.GroupID, x.ParentID, x.Code, x.Name, x.Description, x.IsActive, id))
}
func (r *Repository) DeleteDepartment(ctx context.Context, id uint64) error {
	return affected(r.db.ExecContext(ctx, `DELETE FROM departments WHERE id=?`, id))
}
func (r *Repository) CreateLocation(ctx context.Context, x LocationInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO asset_locations(branch_id,parent_id,code,name,location_type,description,is_active) VALUES(?,NULLIF(?,0),?,?,?,?,?)`, x.BranchID, x.ParentID, x.Code, x.Name, x.LocationType, nullString(x.Description), x.IsActive)
	return err
}
func (r *Repository) UpdateLocation(ctx context.Context, id uint64, x LocationInput) error {
	return affected(r.db.ExecContext(ctx, `UPDATE asset_locations SET branch_id=?,parent_id=NULLIF(?,0),code=?,name=?,location_type=?,description=NULLIF(?,''),is_active=? WHERE id=?`, x.BranchID, x.ParentID, x.Code, x.Name, x.LocationType, x.Description, x.IsActive, id))
}
func (r *Repository) DeleteLocation(ctx context.Context, id uint64) error {
	return affected(r.db.ExecContext(ctx, `DELETE FROM asset_locations WHERE id=?`, id))
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
func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

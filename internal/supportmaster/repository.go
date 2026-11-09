package supportmaster

import (
	"context"
	"database/sql"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List(ctx context.Context, kind string) ([]map[string]any, error) {
	queries := map[string]string{
		"conditions":            `SELECT id AS ID,code AS Code,name AS Name,COALESCE(description,'') AS Description,is_active AS IsActive FROM asset_conditions ORDER BY name`,
		"statuses":              `SELECT id AS ID,code AS Code,name AS Name,COALESCE(description,'') AS Description,is_active AS IsActive FROM asset_statuses ORDER BY name`,
		"vendors":               `SELECT v.id AS ID,v.code AS Code,v.name AS Name,'' AS Description,v.is_active AS IsActive,v.company_group_id AS CompanyGroupID,g.name AS CompanyGroupName,v.vendor_type AS VendorType,COALESCE(v.tax_number,'') AS TaxNumber,COALESCE(v.address,'') AS Address,COALESCE(v.phone,'') AS Phone,COALESCE(v.email,'') AS Email,COALESCE(v.contact_person,'') AS ContactPerson FROM vendors v JOIN company_groups g ON g.id=v.company_group_id ORDER BY g.name,v.name`,
		"depreciation-methods":  `SELECT id AS ID,code AS Code,name AS Name,COALESCE(description,'') AS Description,is_active AS IsActive FROM depreciation_methods ORDER BY name`,
		"depreciation-rates":    `SELECT id AS ID,code AS Code,name AS Name,'' AS Description,is_active AS IsActive,useful_life_months AS UsefulLifeMonths,FORMAT(annual_rate_percent,2) AS AnnualRatePercent FROM depreciation_rates ORDER BY useful_life_months DESC`,
		"gl-accounts":           `SELECT a.id AS ID,a.code AS Code,a.name AS Name,'' AS Description,a.is_active AS IsActive,a.company_id AS CompanyID,c.name AS CompanyName,COALESCE(a.parent_gl_account_id,0) AS ParentID,COALESCE(p.name,'') AS ParentName,a.account_type AS AccountType,COALESCE(a.external_account_id,'') AS ExternalAccountID FROM gl_accounts a JOIN companies c ON c.id=a.company_id LEFT JOIN gl_accounts p ON p.id=a.parent_gl_account_id ORDER BY c.name,a.code`,
		"depreciation-policies": `SELECT p.id AS ID,'' AS Code,t.name AS Name,'' AS Description,p.is_active AS IsActive,p.company_group_id AS CompanyGroupID,g.name AS CompanyGroupName,COALESCE(p.company_id,0) AS CompanyID,COALESCE(c.name,'Default Group') AS CompanyName,p.asset_type_id AS AssetTypeID,t.name AS AssetTypeName,p.depreciation_method_id AS MethodID,m.name AS MethodName,p.depreciation_rate_id AS DepreciationRateID,r.name AS DepreciationRateName,p.book_type AS BookType,r.useful_life_months AS UsefulLifeMonths,FORMAT(r.annual_rate_percent,2) AS AnnualRatePercent,p.residual_value_percent AS ResidualValuePercent,DATE_FORMAT(p.effective_from,'%Y-%m-%d') AS EffectiveFrom,COALESCE(DATE_FORMAT(p.effective_to,'%Y-%m-%d'),'') AS EffectiveTo FROM asset_depreciation_policies p JOIN company_groups g ON g.id=p.company_group_id LEFT JOIN companies c ON c.id=p.company_id JOIN asset_types t ON t.id=p.asset_type_id JOIN depreciation_methods m ON m.id=p.depreciation_method_id JOIN depreciation_rates r ON r.id=p.depreciation_rate_id ORDER BY g.name,t.name,p.book_type,p.effective_from DESC`,
		"accounting-mappings":   `SELECT m.id AS ID,'' AS Code,CONCAT(t.name,' - ',p.book_type) AS Name,'' AS Description,m.is_active AS IsActive,m.company_id AS CompanyID,c.name AS CompanyName,m.depreciation_policy_id AS DepreciationPolicyID,CONCAT(t.name,' - ',p.book_type,' (',DATE_FORMAT(p.effective_from,'%Y-%m-%d'),')') AS PolicyName,CONCAT(t.name,' - ',p.book_type) AS AccountingGroupName,m.asset_gl_account_id AS AssetGLID,a.name AS AssetGLName,COALESCE(m.accumulated_depreciation_gl_account_id,0) AS AccumulatedGLID,COALESCE(ad.name,'') AS AccumulatedGLName,COALESCE(m.depreciation_expense_gl_account_id,0) AS ExpenseGLID,COALESCE(e.name,'') AS ExpenseGLName FROM asset_accounting_mappings m JOIN companies c ON c.id=m.company_id JOIN asset_depreciation_policies p ON p.id=m.depreciation_policy_id JOIN asset_types t ON t.id=p.asset_type_id JOIN gl_accounts a ON a.id=m.asset_gl_account_id LEFT JOIN gl_accounts ad ON ad.id=m.accumulated_depreciation_gl_account_id LEFT JOIN gl_accounts e ON e.id=m.depreciation_expense_gl_account_id ORDER BY c.name,t.name,p.book_type,p.effective_from DESC`,
	}
	q, ok := queries[kind]
	if !ok {
		return nil, ErrInvalid
	}
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	var out []map[string]any
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		item := map[string]any{}
		for i, col := range cols {
			v := values[i]
			if b, ok := v.([]byte); ok {
				v = string(b)
			}
			item[col] = v
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) Options(ctx context.Context) (map[string][]Option, error) {
	queries := map[string]string{
		"Groups":               `SELECT id,id,0,code,name,'' FROM company_groups WHERE is_active=1 ORDER BY name`,
		"Companies":            `SELECT id,group_id,id,code,name,'' FROM companies WHERE is_active=1 ORDER BY name`,
		"AssetTypes":           `SELECT t.id,c.group_id,0,t.code,t.name,'' FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE t.is_active=1 ORDER BY t.name`,
		"Methods":              `SELECT id,0,0,code,name,'' FROM depreciation_methods WHERE is_active=1 ORDER BY name`,
		"DepreciationRates":    `SELECT id,0,0,code,CONCAT(name,' - ',FORMAT(annual_rate_percent,2),'%'),CAST(useful_life_months AS CHAR) FROM depreciation_rates WHERE is_active=1 ORDER BY useful_life_months DESC`,
		"DepreciationPolicies": `SELECT p.id,p.company_group_id,COALESCE(p.company_id,0),'',CONCAT(t.name,' - ',p.book_type,' (',DATE_FORMAT(p.effective_from,'%Y-%m-%d'),')'),'' FROM asset_depreciation_policies p JOIN asset_types t ON t.id=p.asset_type_id WHERE p.is_active=1 ORDER BY t.name,p.book_type,p.effective_from DESC`,
		"GLAccounts":           `SELECT a.id,c.group_id,a.company_id,a.code,a.name,a.account_type FROM gl_accounts a JOIN companies c ON c.id=a.company_id WHERE a.is_active=1 ORDER BY c.name,a.code`,
	}
	out := map[string][]Option{}
	for key, q := range queries {
		rows, err := r.db.QueryContext(ctx, q)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var x Option
			if err := rows.Scan(&x.ID, &x.GroupID, &x.CompanyID, &x.Code, &x.Name, &x.Extra); err != nil {
				rows.Close()
				return nil, err
			}
			out[key] = append(out[key], x)
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (r *Repository) Save(ctx context.Context, kind string, id uint64, x Input) error {
	var q string
	var a []any
	if id == 0 {
		switch kind {
		case "conditions":
			q = `INSERT INTO asset_conditions(code,name,description,is_active) VALUES(?,?,NULLIF(?,''),?)`
			a = []any{x.Code, x.Name, x.Description, x.IsActive}
		case "statuses":
			q = `INSERT INTO asset_statuses(code,name,description,is_active) VALUES(?,?,NULLIF(?,''),?)`
			a = []any{x.Code, x.Name, x.Description, x.IsActive}
		case "vendors":
			q = `INSERT INTO vendors(company_group_id,code,name,vendor_type,tax_number,address,phone,email,contact_person,is_active) VALUES(?,?,?,?,NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),NULLIF(?,''),?)`
			a = []any{x.CompanyGroupID, x.Code, x.Name, x.VendorType, x.TaxNumber, x.Address, x.Phone, x.Email, x.ContactPerson, x.IsActive}
		case "depreciation-methods":
			q = `INSERT INTO depreciation_methods(code,name,description,is_active) VALUES(?,?,NULLIF(?,''),?)`
			a = []any{x.Code, x.Name, x.Description, x.IsActive}
		case "depreciation-rates":
			q = `INSERT INTO depreciation_rates(code,name,useful_life_months,annual_rate_percent,is_active) VALUES(?,?,?,?,?)`
			a = []any{x.Code, x.Name, x.UsefulLifeMonths, x.AnnualRatePercent, x.IsActive}
		case "gl-accounts":
			q = `INSERT INTO gl_accounts(company_id,parent_gl_account_id,code,name,account_type,external_account_id,is_active) VALUES(?,NULLIF(?,0),?,?,?,NULLIF(?,''),?)`
			a = []any{x.CompanyID, x.ParentID, x.Code, x.Name, x.AccountType, x.ExternalAccountID, x.IsActive}
		case "depreciation-policies":
			q = `INSERT INTO asset_depreciation_policies(company_group_id,company_id,asset_type_id,depreciation_method_id,depreciation_rate_id,book_type,residual_value_percent,effective_from,effective_to,is_active) VALUES(?,NULLIF(?,0),?,?,?,?,?,?,NULLIF(?,''),?)`
			a = []any{x.CompanyGroupID, x.CompanyID, x.AssetTypeID, x.MethodID, x.DepreciationRateID, x.BookType, x.ResidualValuePercent, x.EffectiveFrom, x.EffectiveTo, x.IsActive}
		case "accounting-mappings":
			q = `INSERT INTO asset_accounting_mappings(company_id,depreciation_policy_id,asset_gl_account_id,accumulated_depreciation_gl_account_id,depreciation_expense_gl_account_id,is_active) VALUES(?,?,?,NULLIF(?,0),NULLIF(?,0),?)`
			a = []any{x.CompanyID, x.DepreciationPolicyID, x.AssetGLID, x.AccumulatedGLID, x.ExpenseGLID, x.IsActive}
		default:
			return ErrInvalid
		}
	} else {
		switch kind {
		case "conditions":
			q = `UPDATE asset_conditions SET code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.Code, x.Name, x.Description, x.IsActive, id}
		case "statuses":
			q = `UPDATE asset_statuses SET code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.Code, x.Name, x.Description, x.IsActive, id}
		case "vendors":
			q = `UPDATE vendors SET company_group_id=?,code=?,name=?,vendor_type=?,tax_number=NULLIF(?,''),address=NULLIF(?,''),phone=NULLIF(?,''),email=NULLIF(?,''),contact_person=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.CompanyGroupID, x.Code, x.Name, x.VendorType, x.TaxNumber, x.Address, x.Phone, x.Email, x.ContactPerson, x.IsActive, id}
		case "depreciation-methods":
			q = `UPDATE depreciation_methods SET code=?,name=?,description=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.Code, x.Name, x.Description, x.IsActive, id}
		case "depreciation-rates":
			q = `UPDATE depreciation_rates SET code=?,name=?,useful_life_months=?,annual_rate_percent=?,is_active=? WHERE id=?`
			a = []any{x.Code, x.Name, x.UsefulLifeMonths, x.AnnualRatePercent, x.IsActive, id}
		case "gl-accounts":
			q = `UPDATE gl_accounts SET company_id=?,parent_gl_account_id=NULLIF(?,0),code=?,name=?,account_type=?,external_account_id=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.CompanyID, x.ParentID, x.Code, x.Name, x.AccountType, x.ExternalAccountID, x.IsActive, id}
		case "depreciation-policies":
			q = `UPDATE asset_depreciation_policies SET company_group_id=?,company_id=NULLIF(?,0),asset_type_id=?,depreciation_method_id=?,depreciation_rate_id=?,book_type=?,residual_value_percent=?,effective_from=?,effective_to=NULLIF(?,''),is_active=? WHERE id=?`
			a = []any{x.CompanyGroupID, x.CompanyID, x.AssetTypeID, x.MethodID, x.DepreciationRateID, x.BookType, x.ResidualValuePercent, x.EffectiveFrom, x.EffectiveTo, x.IsActive, id}
		case "accounting-mappings":
			q = `UPDATE asset_accounting_mappings SET company_id=?,depreciation_policy_id=?,asset_gl_account_id=?,accumulated_depreciation_gl_account_id=NULLIF(?,0),depreciation_expense_gl_account_id=NULLIF(?,0),is_active=? WHERE id=?`
			a = []any{x.CompanyID, x.DepreciationPolicyID, x.AssetGLID, x.AccumulatedGLID, x.ExpenseGLID, x.IsActive, id}
		default:
			return ErrInvalid
		}
	}
	res, err := r.db.ExecContext(ctx, q, a...)
	if err != nil {
		return err
	}
	if id > 0 {
		n, _ := res.RowsAffected()
		if n == 0 {
			return sql.ErrNoRows
		}
	}
	return nil
}
func (r *Repository) Delete(ctx context.Context, kind string, id uint64) error {
	tables := map[string]string{"conditions": "asset_conditions", "statuses": "asset_statuses", "vendors": "vendors", "depreciation-methods": "depreciation_methods", "depreciation-rates": "depreciation_rates", "depreciation-policies": "asset_depreciation_policies", "gl-accounts": "gl_accounts", "accounting-mappings": "asset_accounting_mappings"}
	t, ok := tables[kind]
	if !ok {
		return ErrInvalid
	}
	res, err := r.db.ExecContext(ctx, "DELETE FROM "+t+" WHERE id=?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) Validate(ctx context.Context, kind string, id uint64, x Input) error {
	var count int
	switch kind {
	case "depreciation-policies":
		err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_types t JOIN asset_categories ac ON ac.id=t.asset_category_id JOIN depreciation_rates r ON r.id=? AND r.is_active=1 LEFT JOIN companies c ON c.id=? WHERE t.id=? AND ac.group_id=? AND (?=0 OR c.group_id=?)`, x.DepreciationRateID, x.CompanyID, x.AssetTypeID, x.CompanyGroupID, x.CompanyID, x.CompanyGroupID).Scan(&count)
		if err != nil {
			return err
		}
		if count != 1 {
			return ErrRelation
		}
		var overlap int
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM asset_depreciation_policies WHERE id<>? AND company_group_id=? AND COALESCE(company_id,0)=? AND asset_type_id=? AND book_type=? AND effective_from<=COALESCE(NULLIF(?,''),'9999-12-31') AND COALESCE(effective_to,'9999-12-31')>=?`, id, x.CompanyGroupID, x.CompanyID, x.AssetTypeID, x.BookType, x.EffectiveTo, x.EffectiveFrom).Scan(&overlap)
		if err != nil {
			return err
		}
		if overlap > 0 {
			return ErrOverlap
		}
	case "accounting-mappings":
		err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM companies c JOIN asset_depreciation_policies p ON p.id=? AND p.company_group_id=c.group_id AND (p.company_id IS NULL OR p.company_id=c.id) JOIN gl_accounts a ON a.id=? AND a.company_id=c.id AND a.account_type='ASSET' LEFT JOIN gl_accounts d ON d.id=NULLIF(?,0) AND d.company_id=c.id AND d.account_type='CONTRA_ASSET' LEFT JOIN gl_accounts e ON e.id=NULLIF(?,0) AND e.company_id=c.id AND e.account_type='EXPENSE' WHERE c.id=? AND ((?=0 AND ?=0) OR (d.id IS NOT NULL AND e.id IS NOT NULL))`, x.DepreciationPolicyID, x.AssetGLID, x.AccumulatedGLID, x.ExpenseGLID, x.CompanyID, x.AccumulatedGLID, x.ExpenseGLID).Scan(&count)
		if err != nil {
			return err
		}
		if count != 1 {
			return ErrRelation
		}
	case "gl-accounts":
		if x.ParentID > 0 {
			err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gl_accounts WHERE id=? AND company_id=? AND id<>?`, x.ParentID, x.CompanyID, id).Scan(&count)
			if err != nil {
				return err
			}
			if count != 1 {
				return ErrRelation
			}
		}
	}
	return nil
}

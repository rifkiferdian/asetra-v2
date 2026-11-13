package navigation

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"asetra/internal/auth"
)

type Item struct {
	Key      string
	Label    string
	URL      string
	Active   bool
	Open     bool
	Children []Item
}

type Group struct {
	Label string
	Items []Item
}

type definition struct {
	group      string
	key        string
	label      string
	url        string
	permission string
	parent     string
}

var definitions = []definition{
	{group: "Menu utama", key: "dashboard", label: "Dashboard", url: "/dashboard"},
	{group: "Pengelolaan", key: "assets", label: "Data Aset", url: "/assets", permission: "assets.view"},
	{group: "Pengelolaan", key: "asset-depreciations", label: "Penyusutan Aset", url: "/asset-depreciations", permission: "asset_depreciations.view"},
	{group: "Pengelolaan", key: "location-qr-scanner", label: "Scan QR Lokasi", url: "/tools/scan-location-qr", permission: "assets.view"},
	{group: "Pengelolaan", key: "asset-location-report", label: "Report Aset per Lokasi", url: "/reports/assets-by-location", permission: "assets.view"},
	{group: "Organisasi", parent: "asset-organization", key: "companies", label: "Company", url: "/organization/companies", permission: "companies.view"},
	{group: "Organisasi", parent: "asset-organization", key: "branches", label: "Branch", url: "/organization/branches", permission: "branches.view"},
	{group: "Organisasi", parent: "asset-organization", key: "asset-departments", label: "Department", url: "/organization/departments", permission: "asset_departments.view"},
	{group: "Organisasi", parent: "asset-organization", key: "asset-locations", label: "Asset Location", url: "/organization/asset-locations", permission: "asset_locations.view"},
	{group: "Pengelolaan", parent: "access", key: "stores", label: "Master Data Store", url: "/master/stores", permission: "stores.view"},
	{group: "Master Data", parent: "asset-master", key: "asset-categories", label: "Asset Category", url: "/master/assets/categories", permission: "asset_categories.view"},
	{group: "Master Data", parent: "asset-master", key: "asset-types", label: "Asset Type", url: "/master/assets/types", permission: "asset_types.view"},
	{group: "Master Data", parent: "asset-master", key: "asset-brands", label: "Asset Brand", url: "/master/assets/brands", permission: "asset_brands.view"},
	{group: "Master Data", parent: "asset-master", key: "asset-models", label: "Asset Model", url: "/master/assets/models", permission: "asset_models.view"},
	{group: "Master Data", parent: "asset-master", key: "asset-variants", label: "Asset Variant", url: "/master/assets/variants", permission: "asset_variants.view"},
	{group: "Master Data", parent: "operational-master", key: "asset-conditions", label: "Kondisi Aset", url: "/master/support/conditions", permission: "asset_conditions.view"},
	{group: "Master Data", parent: "operational-master", key: "asset-statuses", label: "Status Aset", url: "/master/support/statuses", permission: "asset_statuses.view"},
	{group: "Master Data", parent: "operational-master", key: "vendors", label: "Vendor", url: "/master/support/vendors", permission: "vendors.view"},
	{group: "Master Data", parent: "accounting-master", key: "depreciation-methods", label: "Metode Penyusutan", url: "/master/support/depreciation-methods", permission: "depreciation_methods.view"},
	{group: "Master Data", parent: "accounting-master", key: "depreciation-rates", label: "Tarif Penyusutan", url: "/master/support/depreciation-rates", permission: "depreciation_rates.view"},
	{group: "Master Data", parent: "accounting-master", key: "depreciation-policies", label: "Kebijakan Penyusutan", url: "/master/support/depreciation-policies", permission: "asset_depreciation_policies.view"},
	{group: "Master Data", parent: "accounting-master", key: "gl-accounts", label: "Chart of Accounts (GL)", url: "/master/support/gl-accounts", permission: "gl_accounts.view"},
	{group: "Master Data", parent: "accounting-master", key: "asset-accounting-mappings", label: "Mapping Akun Accounting", url: "/master/support/accounting-mappings", permission: "asset_accounting_mappings.view"},
	{group: "Pengelolaan", parent: "access", key: "users", label: "User", url: "/settings/users", permission: "users.view"},
	{group: "Pengelolaan", parent: "access", key: "roles", label: "Role", url: "/settings/roles", permission: "roles.view"},
	{group: "Pengelolaan", parent: "access", key: "permissions", label: "Permission", url: "/settings/permissions", permission: "permissions.view"},
}

type Service struct{ db *sql.DB }

func NewService(db *sql.DB) *Service { return &Service{db: db} }

func (s *Service) Menus(ctx context.Context, user auth.SessionUser, activePage string) ([]Group, error) {
	allowed := make(map[string]bool)
	if hasRole(user.RoleName, "superadmin") {
		for _, item := range definitions {
			allowed[item.permission] = true
		}
	} else {
		userStoreID, err := strconv.ParseUint(user.UserStoreID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user store session: %w", err)
		}
		rows, err := s.db.QueryContext(ctx, `SELECT DISTINCT p.name
			FROM permissions p
			WHERE EXISTS(SELECT 1 FROM user_store_permissions usp WHERE usp.user_store_id=? AND usp.permission_id=p.id)
			   OR EXISTS(SELECT 1 FROM user_store_roles usr JOIN role_permissions rp ON rp.role_id=usr.role_id WHERE usr.user_store_id=? AND rp.permission_id=p.id)`, userStoreID, userStoreID)
		if err != nil {
			return nil, fmt.Errorf("load navigation permissions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return nil, err
			}
			allowed[name] = true
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	groups := make([]Group, 0, 2)
	indexes := make(map[string]int)
	children := make(map[string][]Item)
	for _, definition := range definitions {
		if definition.permission != "" && !allowed[definition.permission] {
			continue
		}
		item := Item{Key: definition.key, Label: definition.label, URL: definition.url, Active: definition.key == activePage}
		if definition.parent != "" {
			children[definition.parent] = append(children[definition.parent], item)
			continue
		}
		index, exists := indexes[definition.group]
		if !exists {
			index = len(groups)
			indexes[definition.group] = index
			groups = append(groups, Group{Label: definition.group})
		}
		groups[index].Items = append(groups[index].Items, item)
	}
	if accessItems := children["access"]; len(accessItems) > 0 {
		index, exists := indexes["Pengelolaan"]
		if !exists {
			index = len(groups)
			indexes["Pengelolaan"] = index
			groups = append(groups, Group{Label: "Pengelolaan"})
		}
		open := false
		for _, item := range accessItems {
			open = open || item.Active
		}
		groups[index].Items = append(groups[index].Items, Item{Key: "access", Label: "Manajemen Akses", Open: open, Active: open, Children: accessItems})
	}
	if masterItems := children["master"]; len(masterItems) > 0 {
		index, exists := indexes["Master Data"]
		if !exists {
			index = len(groups)
			indexes["Master Data"] = index
			groups = append(groups, Group{Label: "Master Data"})
		}
		open := false
		for _, item := range masterItems {
			open = open || item.Active
		}
		groups[index].Items = append(groups[index].Items, Item{Key: "master", Label: "Master Data", Open: open, Active: open, Children: masterItems})
	}
	if organizationItems := children["asset-organization"]; len(organizationItems) > 0 {
		index, exists := indexes["Organisasi"]
		if !exists {
			index = len(groups)
			indexes["Organisasi"] = index
			groups = append(groups, Group{Label: "Organisasi"})
		}
		open := false
		for _, item := range organizationItems {
			open = open || item.Active
		}
		groups[index].Items = append(groups[index].Items, Item{Key: "asset-organization", Label: "Struktur Organisasi", Open: open, Active: open, Children: organizationItems})
	}
	if assetMasterItems := children["asset-master"]; len(assetMasterItems) > 0 {
		index, exists := indexes["Master Data"]
		if !exists {
			index = len(groups)
			indexes["Master Data"] = index
			groups = append(groups, Group{Label: "Master Data"})
		}
		open := false
		for _, item := range assetMasterItems {
			open = open || item.Active
		}
		groups[index].Items = append(groups[index].Items, Item{Key: "asset-master", Label: "Master Data Asset", Open: open, Active: open, Children: assetMasterItems})
	}
	for _, parent := range []struct{ Key, Label string }{{"operational-master", "Master Operasional"}, {"accounting-master", "Master Accounting"}} {
		items := children[parent.Key]
		if len(items) == 0 {
			continue
		}
		index, exists := indexes["Master Data"]
		if !exists {
			index = len(groups)
			indexes["Master Data"] = index
			groups = append(groups, Group{Label: "Master Data"})
		}
		open := false
		for _, item := range items {
			open = open || item.Active
		}
		groups[index].Items = append(groups[index].Items, Item{Key: parent.Key, Label: parent.Label, Open: open, Active: open, Children: items})
	}
	return groups, nil
}

func hasRole(roleNames, required string) bool {
	for _, role := range strings.Split(roleNames, ",") {
		if strings.TrimSpace(role) == required {
			return true
		}
	}
	return false
}

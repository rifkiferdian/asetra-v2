START TRANSACTION;
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id
WHERE SUBSTRING_INDEX(p.name,'.',1) IN ('asset_conditions','asset_statuses','vendors','depreciation_methods','asset_depreciation_policies','gl_accounts','asset_accounting_groups','asset_accounting_mappings');
DELETE FROM permissions WHERE SUBSTRING_INDEX(name,'.',1) IN ('asset_conditions','asset_statuses','vendors','depreciation_methods','asset_depreciation_policies','gl_accounts','asset_accounting_groups','asset_accounting_mappings');
COMMIT;

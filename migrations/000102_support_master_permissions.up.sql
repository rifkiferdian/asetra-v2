START TRANSACTION;
INSERT INTO permissions(name,guard_name,description)
SELECT CONCAT(resource,'.',action),'web',CONCAT(action,' ',resource)
FROM (
 SELECT 'asset_conditions' resource UNION ALL SELECT 'asset_statuses' UNION ALL SELECT 'vendors'
 UNION ALL SELECT 'depreciation_methods' UNION ALL SELECT 'asset_depreciation_policies'
 UNION ALL SELECT 'gl_accounts' UNION ALL SELECT 'asset_accounting_groups' UNION ALL SELECT 'asset_accounting_mappings'
) resources CROSS JOIN (
 SELECT 'view' action UNION ALL SELECT 'create' UNION ALL SELECT 'update' UNION ALL SELECT 'delete'
) actions
WHERE 1=1
ON DUPLICATE KEY UPDATE description=VALUES(description);
INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web' AND
SUBSTRING_INDEX(p.name,'.',1) IN ('asset_conditions','asset_statuses','vendors','depreciation_methods','asset_depreciation_policies','gl_accounts','asset_accounting_groups','asset_accounting_mappings');
COMMIT;

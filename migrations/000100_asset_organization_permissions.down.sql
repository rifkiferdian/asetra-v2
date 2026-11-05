START TRANSACTION;
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id
WHERE p.name LIKE 'companies.%' OR p.name LIKE 'branches.%' OR p.name LIKE 'asset_departments.%' OR p.name LIKE 'asset_locations.%';
DELETE FROM permissions
WHERE name LIKE 'companies.%' OR name LIKE 'branches.%' OR name LIKE 'asset_departments.%' OR name LIKE 'asset_locations.%';
COMMIT;

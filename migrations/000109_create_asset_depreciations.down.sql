START TRANSACTION;
DELETE usp FROM user_store_permissions usp JOIN permissions p ON p.id=usp.permission_id WHERE p.name LIKE 'asset_depreciations.%';
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id WHERE p.name LIKE 'asset_depreciations.%';
DELETE FROM permissions WHERE name LIKE 'asset_depreciations.%';
COMMIT;
DROP TABLE asset_depreciations;

START TRANSACTION;
DELETE usp FROM user_store_permissions usp JOIN permissions p ON p.id=usp.permission_id WHERE p.name LIKE 'assets.%';
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id WHERE p.name LIKE 'assets.%';
DELETE FROM permissions WHERE name LIKE 'assets.%';
COMMIT;
DROP TABLE IF EXISTS assets;

START TRANSACTION;
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id
WHERE p.name REGEXP '^asset_(categories|types|brands|models|variants)\\.';
DELETE FROM permissions WHERE name REGEXP '^asset_(categories|types|brands|models|variants)\\.';
COMMIT;

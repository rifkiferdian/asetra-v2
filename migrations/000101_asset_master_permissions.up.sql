START TRANSACTION;
INSERT INTO permissions(name,guard_name,description) VALUES
('asset_categories.view','web','Melihat master kategori aset'),('asset_categories.create','web','Membuat kategori aset'),('asset_categories.update','web','Memperbarui kategori aset'),('asset_categories.delete','web','Menghapus kategori aset'),
('asset_types.view','web','Melihat master tipe aset'),('asset_types.create','web','Membuat tipe aset'),('asset_types.update','web','Memperbarui tipe aset'),('asset_types.delete','web','Menghapus tipe aset'),
('asset_brands.view','web','Melihat master brand aset'),('asset_brands.create','web','Membuat brand aset'),('asset_brands.update','web','Memperbarui brand aset'),('asset_brands.delete','web','Menghapus brand aset'),
('asset_models.view','web','Melihat master model aset'),('asset_models.create','web','Membuat model aset'),('asset_models.update','web','Memperbarui model aset'),('asset_models.delete','web','Menghapus model aset'),
('asset_variants.view','web','Melihat master varian aset'),('asset_variants.create','web','Membuat varian aset'),('asset_variants.update','web','Memperbarui varian aset'),('asset_variants.delete','web','Menghapus varian aset')
ON DUPLICATE KEY UPDATE description=VALUES(description);
INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web' AND p.name REGEXP '^asset_(categories|types|brands|models|variants)\\.';
COMMIT;

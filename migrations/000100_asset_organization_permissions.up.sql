-- Permission Sprint 1: Struktur Organisasi Asetra.
START TRANSACTION;

INSERT INTO permissions(name,guard_name,description) VALUES
('companies.view','web','Melihat Company dan Company Group'),
('companies.create','web','Membuat Company dan Company Group'),
('companies.update','web','Memperbarui Company dan Company Group'),
('companies.delete','web','Menghapus Company dan Company Group'),
('branches.view','web','Melihat Branch'),
('branches.create','web','Membuat Branch'),
('branches.update','web','Memperbarui Branch'),
('branches.delete','web','Menghapus Branch'),
('asset_departments.view','web','Melihat Department organisasi aset'),
('asset_departments.create','web','Membuat Department organisasi aset'),
('asset_departments.update','web','Memperbarui Department organisasi aset'),
('asset_departments.delete','web','Menghapus Department organisasi aset'),
('asset_locations.view','web','Melihat Asset Location'),
('asset_locations.create','web','Membuat Asset Location'),
('asset_locations.update','web','Memperbarui Asset Location'),
('asset_locations.delete','web','Menghapus Asset Location')
ON DUPLICATE KEY UPDATE description=VALUES(description);

INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web'
  AND (p.name LIKE 'companies.%' OR p.name LIKE 'branches.%' OR p.name LIKE 'asset_departments.%' OR p.name LIKE 'asset_locations.%');

COMMIT;

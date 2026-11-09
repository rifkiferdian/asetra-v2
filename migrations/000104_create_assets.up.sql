CREATE TABLE IF NOT EXISTS assets (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  company_id BIGINT UNSIGNED NOT NULL,
  branch_id BIGINT UNSIGNED NOT NULL,
  asset_location_id BIGINT UNSIGNED NOT NULL,
  asset_type_id BIGINT UNSIGNED NOT NULL,
  asset_brand_id BIGINT UNSIGNED NULL,
  asset_model_id BIGINT UNSIGNED NULL,
  asset_variant_id BIGINT UNSIGNED NULL,
  condition_id BIGINT UNSIGNED NOT NULL,
  status_id BIGINT UNSIGNED NOT NULL,
  vendor_id BIGINT UNSIGNED NULL,
  asset_number VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  serial_number VARCHAR(120) NULL,
  barcode VARCHAR(120) NULL,
  description TEXT NULL,
  acquisition_date DATE NOT NULL,
  acquisition_cost DECIMAL(18,2) NOT NULL DEFAULT 0,
  warranty_until DATE NULL,
  notes TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_assets_number (asset_number),
  UNIQUE KEY uq_assets_barcode (barcode),
  KEY idx_assets_company (company_id),
  KEY idx_assets_branch (branch_id),
  KEY idx_assets_location (asset_location_id),
  KEY idx_assets_type (asset_type_id),
  KEY idx_assets_status (status_id),
  CONSTRAINT fk_assets_company FOREIGN KEY (company_id) REFERENCES companies(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_location FOREIGN KEY (asset_location_id) REFERENCES asset_locations(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_type FOREIGN KEY (asset_type_id) REFERENCES asset_types(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_brand FOREIGN KEY (asset_brand_id) REFERENCES asset_brands(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_model FOREIGN KEY (asset_model_id) REFERENCES asset_models(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_variant FOREIGN KEY (asset_variant_id) REFERENCES asset_variants(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_condition FOREIGN KEY (condition_id) REFERENCES asset_conditions(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_status FOREIGN KEY (status_id) REFERENCES asset_statuses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_assets_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT chk_assets_acquisition_cost CHECK (acquisition_cost >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

START TRANSACTION;
INSERT INTO permissions(name,guard_name,description) VALUES
('assets.view','web','Melihat daftar dan detail aset'),
('assets.create','web','Membuat data aset'),
('assets.update','web','Memperbarui data aset'),
('assets.delete','web','Menghapus data aset')
ON DUPLICATE KEY UPDATE description=VALUES(description);
INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web' AND p.name LIKE 'assets.%';
COMMIT;

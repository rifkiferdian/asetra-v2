CREATE TABLE asset_depreciations (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  asset_id BIGINT UNSIGNED NOT NULL,
  depreciation_policy_id BIGINT UNSIGNED NOT NULL,
  period_date DATE NOT NULL,
  depreciation_amount DECIMAL(18,2) NOT NULL,
  status ENUM('POSTED','CANCELLED') NOT NULL DEFAULT 'POSTED',
  posted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  posted_by BIGINT UNSIGNED NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_asset_depreciations_asset_period (asset_id,period_date),
  KEY idx_asset_depreciations_period (period_date),
  KEY idx_asset_depreciations_policy (depreciation_policy_id),
  KEY idx_asset_depreciations_posted_by (posted_by),
  CONSTRAINT fk_asset_depreciations_asset FOREIGN KEY (asset_id) REFERENCES assets(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_asset_depreciations_policy FOREIGN KEY (depreciation_policy_id) REFERENCES asset_depreciation_policies(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_asset_depreciations_posted_by FOREIGN KEY (posted_by) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT chk_asset_depreciations_amount CHECK (depreciation_amount > 0),
  CONSTRAINT chk_asset_depreciations_period CHECK (DAY(period_date) = 1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

START TRANSACTION;
INSERT INTO permissions(name,guard_name,description) VALUES
('asset_depreciations.view','web','Melihat posting penyusutan aset'),
('asset_depreciations.create','web','Memposting penyusutan aset')
ON DUPLICATE KEY UPDATE description=VALUES(description);
INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web' AND p.name LIKE 'asset_depreciations.%';
COMMIT;

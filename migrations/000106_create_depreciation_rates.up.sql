CREATE TABLE depreciation_rates (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(40) NOT NULL,
  name VARCHAR(120) NOT NULL,
  useful_life_months SMALLINT UNSIGNED NOT NULL,
  annual_rate_percent DECIMAL(8,4) NOT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_depreciation_rates_code (code),
  UNIQUE KEY uq_depreciation_rates_life (useful_life_months),
  CONSTRAINT chk_depreciation_rates_life CHECK (useful_life_months > 0),
  CONSTRAINT chk_depreciation_rates_percent CHECK (annual_rate_percent > 0 AND annual_rate_percent <= 100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO depreciation_rates(code,name,useful_life_months,annual_rate_percent,is_active) VALUES
('DEP-20Y','20 Tahun',240,5.0000,1),
('DEP-10Y','10 Tahun',120,10.0000,1),
('DEP-8Y','8 Tahun',96,12.5000,1),
('DEP-5Y','5 Tahun',60,20.0000,1),
('DEP-4Y','4 Tahun',48,25.0000,1);

ALTER TABLE asset_depreciation_policies
  ADD COLUMN depreciation_rate_id BIGINT UNSIGNED NULL AFTER depreciation_method_id;

UPDATE asset_depreciation_policies p
JOIN depreciation_rates r ON r.useful_life_months=p.useful_life_months
SET p.depreciation_rate_id=r.id;

ALTER TABLE asset_depreciation_policies
  MODIFY depreciation_rate_id BIGINT UNSIGNED NOT NULL,
  ADD KEY idx_asset_depreciation_policies_rate (depreciation_rate_id),
  ADD CONSTRAINT fk_asset_depreciation_policies_rate
    FOREIGN KEY (depreciation_rate_id) REFERENCES depreciation_rates(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  DROP COLUMN useful_life_months;

START TRANSACTION;
INSERT INTO permissions(name,guard_name,description)
SELECT CONCAT('depreciation_rates.',action),'web',CONCAT(action,' tarif penyusutan')
FROM (SELECT 'view' action UNION ALL SELECT 'create' UNION ALL SELECT 'update' UNION ALL SELECT 'delete') actions
ON DUPLICATE KEY UPDATE description=VALUES(description);
INSERT IGNORE INTO role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM roles r CROSS JOIN permissions p
WHERE r.name='superadmin' AND r.guard_name='web' AND p.name LIKE 'depreciation_rates.%';
COMMIT;

ALTER TABLE asset_depreciation_policies
  ADD COLUMN useful_life_months SMALLINT UNSIGNED NULL AFTER book_type;
UPDATE asset_depreciation_policies p JOIN depreciation_rates r ON r.id=p.depreciation_rate_id
SET p.useful_life_months=r.useful_life_months;
ALTER TABLE asset_depreciation_policies
  MODIFY useful_life_months SMALLINT UNSIGNED NOT NULL,
  DROP FOREIGN KEY fk_asset_depreciation_policies_rate,
  DROP INDEX idx_asset_depreciation_policies_rate,
  DROP COLUMN depreciation_rate_id;
START TRANSACTION;
DELETE usp FROM user_store_permissions usp JOIN permissions p ON p.id=usp.permission_id WHERE p.name LIKE 'depreciation_rates.%';
DELETE rp FROM role_permissions rp JOIN permissions p ON p.id=rp.permission_id WHERE p.name LIKE 'depreciation_rates.%';
DELETE FROM permissions WHERE name LIKE 'depreciation_rates.%';
COMMIT;
DROP TABLE depreciation_rates;

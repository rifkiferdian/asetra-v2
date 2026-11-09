-- Rollback mengembalikan struktur kelompok accounting, tetapi data mapping lama
-- tidak dapat dipulihkan karena migrasi up memang mengosongkannya.
START TRANSACTION;

DELETE FROM asset_accounting_mappings;
ALTER TABLE asset_accounting_mappings
  DROP FOREIGN KEY fk_asset_accounting_mappings_category,
  DROP INDEX uq_asset_accounting_mappings_company_category,
  DROP INDEX idx_asset_accounting_mappings_category;

CREATE TABLE asset_accounting_groups (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  company_group_id BIGINT UNSIGNED NOT NULL,
  code VARCHAR(40) NOT NULL,
  name VARCHAR(150) NOT NULL,
  description TEXT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_asset_accounting_groups_group_code (company_group_id, code),
  CONSTRAINT fk_asset_accounting_groups_company_group
    FOREIGN KEY (company_group_id) REFERENCES company_groups(id)
    ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE asset_accounting_mappings
  ADD COLUMN asset_accounting_group_id BIGINT UNSIGNED NULL AFTER company_id,
  DROP COLUMN asset_category_id;

ALTER TABLE asset_accounting_mappings
  MODIFY asset_accounting_group_id BIGINT UNSIGNED NOT NULL,
  ADD UNIQUE KEY uq_asset_accounting_mappings_company_group (company_id, asset_accounting_group_id),
  ADD CONSTRAINT fk_asset_accounting_mappings_group
    FOREIGN KEY (asset_accounting_group_id) REFERENCES asset_accounting_groups(id)
    ON UPDATE CASCADE ON DELETE RESTRICT;

COMMIT;

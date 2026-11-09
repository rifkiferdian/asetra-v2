START TRANSACTION;
DELETE FROM asset_accounting_mappings;
COMMIT;

ALTER TABLE asset_accounting_mappings
  DROP FOREIGN KEY fk_asset_accounting_mappings_policy,
  DROP INDEX uq_asset_accounting_mappings_company_policy,
  DROP INDEX idx_asset_accounting_mappings_policy,
  DROP COLUMN depreciation_policy_id,
  ADD COLUMN asset_category_id BIGINT UNSIGNED NOT NULL AFTER company_id,
  ADD UNIQUE KEY uq_asset_accounting_mappings_company_category (company_id,asset_category_id),
  ADD KEY idx_asset_accounting_mappings_category (asset_category_id),
  ADD CONSTRAINT fk_asset_accounting_mappings_category
    FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id)
    ON UPDATE CASCADE ON DELETE RESTRICT;

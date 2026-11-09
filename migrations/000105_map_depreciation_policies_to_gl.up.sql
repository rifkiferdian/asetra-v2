-- Mapping GL kini langsung mengacu ke kebijakan penyusutan, bukan kategori aset.
START TRANSACTION;
DELETE FROM asset_accounting_mappings;
COMMIT;

ALTER TABLE asset_accounting_mappings
  DROP FOREIGN KEY fk_asset_accounting_mappings_category,
  DROP INDEX uq_asset_accounting_mappings_company_category,
  DROP INDEX idx_asset_accounting_mappings_category,
  DROP COLUMN asset_category_id,
  ADD COLUMN depreciation_policy_id BIGINT UNSIGNED NOT NULL AFTER company_id,
  ADD UNIQUE KEY uq_asset_accounting_mappings_company_policy (company_id,depreciation_policy_id),
  ADD KEY idx_asset_accounting_mappings_policy (depreciation_policy_id),
  ADD CONSTRAINT fk_asset_accounting_mappings_policy
    FOREIGN KEY (depreciation_policy_id) REFERENCES asset_depreciation_policies(id)
    ON UPDATE CASCADE ON DELETE RESTRICT;

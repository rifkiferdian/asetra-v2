ALTER TABLE asset_accounting_mappings
  DROP FOREIGN KEY fk_aam_gl,
  DROP INDEX idx_aam_gl,
  CHANGE COLUMN gl_account_id asset_gl_account_id BIGINT UNSIGNED NOT NULL,
  ADD COLUMN accumulated_depreciation_gl_account_id BIGINT UNSIGNED NULL AFTER asset_gl_account_id,
  ADD COLUMN depreciation_expense_gl_account_id BIGINT UNSIGNED NULL AFTER accumulated_depreciation_gl_account_id,
  ADD KEY idx_aam_asset_gl (asset_gl_account_id,company_id),
  ADD KEY idx_aam_accum_gl (accumulated_depreciation_gl_account_id,company_id),
  ADD KEY idx_aam_expense_gl (depreciation_expense_gl_account_id,company_id),
  ADD CONSTRAINT fk_aam_asset_gl FOREIGN KEY (asset_gl_account_id,company_id)
    REFERENCES gl_accounts(id,company_id),
  ADD CONSTRAINT fk_aam_accum_gl FOREIGN KEY (accumulated_depreciation_gl_account_id,company_id)
    REFERENCES gl_accounts(id,company_id),
  ADD CONSTRAINT fk_aam_expense_gl FOREIGN KEY (depreciation_expense_gl_account_id,company_id)
    REFERENCES gl_accounts(id,company_id),
  ADD CONSTRAINT chk_aam_depreciation_pair CHECK (
    (accumulated_depreciation_gl_account_id IS NULL AND depreciation_expense_gl_account_id IS NULL)
    OR (accumulated_depreciation_gl_account_id IS NOT NULL AND depreciation_expense_gl_account_id IS NOT NULL)
  );

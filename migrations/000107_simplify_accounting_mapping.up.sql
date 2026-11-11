-- Mapping accounting hanya menghubungkan kebijakan penyusutan dengan satu akun GL.
ALTER TABLE asset_accounting_mappings
  DROP FOREIGN KEY fk_aam_asset_gl,
  DROP FOREIGN KEY fk_aam_accum_gl,
  DROP FOREIGN KEY fk_aam_expense_gl,
  DROP INDEX idx_aam_asset_gl,
  DROP INDEX idx_aam_accum_gl,
  DROP INDEX idx_aam_expense_gl,
  DROP CONSTRAINT chk_aam_depreciation_pair,
  DROP COLUMN accumulated_depreciation_gl_account_id,
  DROP COLUMN depreciation_expense_gl_account_id,
  CHANGE COLUMN asset_gl_account_id gl_account_id BIGINT UNSIGNED NOT NULL,
  ADD KEY idx_aam_gl (gl_account_id,company_id),
  ADD CONSTRAINT fk_aam_gl FOREIGN KEY (gl_account_id,company_id)
    REFERENCES gl_accounts(id,company_id);

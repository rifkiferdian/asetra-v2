CREATE TABLE asset_number_sequences (
  company_id BIGINT UNSIGNED NOT NULL,
  asset_type_id BIGINT UNSIGNED NOT NULL,
  sequence_year SMALLINT UNSIGNED NOT NULL,
  last_number INT UNSIGNED NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (company_id,asset_type_id,sequence_year),
  CONSTRAINT fk_asset_number_sequences_company FOREIGN KEY (company_id) REFERENCES companies(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_asset_number_sequences_type FOREIGN KEY (asset_type_id) REFERENCES asset_types(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT chk_asset_number_sequences_year CHECK (sequence_year BETWEEN 2000 AND 9999)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

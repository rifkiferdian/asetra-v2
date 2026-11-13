CREATE TABLE asset_type_brands (
  asset_type_id BIGINT UNSIGNED NOT NULL,
  asset_brand_id BIGINT UNSIGNED NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (asset_type_id,asset_brand_id),
  KEY idx_asset_type_brands_brand (asset_brand_id,asset_type_id),
  CONSTRAINT fk_asset_type_brands_type FOREIGN KEY (asset_type_id) REFERENCES asset_types(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT fk_asset_type_brands_brand FOREIGN KEY (asset_brand_id) REFERENCES asset_brands(id)
    ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO asset_type_brands(asset_type_id,asset_brand_id)
SELECT asset_type_id,asset_brand_id FROM asset_models WHERE asset_brand_id IS NOT NULL
UNION
SELECT asset_type_id,asset_brand_id FROM assets WHERE asset_brand_id IS NOT NULL;

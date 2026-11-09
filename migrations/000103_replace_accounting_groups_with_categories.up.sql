-- Mapping accounting sekarang langsung menggunakan Asset Category.
-- Mapping lama sengaja dikosongkan karena tidak ada relasi deterministik dari
-- asset_accounting_groups ke asset_categories.
START TRANSACTION;

DELETE FROM asset_accounting_mappings;

SET @mapping_group_fk = (
  SELECT kcu.CONSTRAINT_NAME
  FROM information_schema.KEY_COLUMN_USAGE kcu
  WHERE kcu.CONSTRAINT_SCHEMA = DATABASE()
    AND kcu.TABLE_NAME = 'asset_accounting_mappings'
    AND kcu.COLUMN_NAME = 'asset_accounting_group_id'
    AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
  LIMIT 1
);
SET @drop_mapping_group_fk = IF(
  @mapping_group_fk IS NULL,
  'SELECT 1',
  CONCAT('ALTER TABLE asset_accounting_mappings DROP FOREIGN KEY `', REPLACE(@mapping_group_fk, '`', '``'), '`')
);
PREPARE drop_mapping_group_fk_stmt FROM @drop_mapping_group_fk;
EXECUTE drop_mapping_group_fk_stmt;
DEALLOCATE PREPARE drop_mapping_group_fk_stmt;

SET @has_company_index = (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS s
  WHERE s.TABLE_SCHEMA = DATABASE()
    AND s.TABLE_NAME = 'asset_accounting_mappings'
    AND s.COLUMN_NAME = 'company_id'
    AND s.SEQ_IN_INDEX = 1
    AND s.INDEX_NAME NOT IN (
      SELECT s2.INDEX_NAME
      FROM information_schema.STATISTICS s2
      WHERE s2.TABLE_SCHEMA = DATABASE()
        AND s2.TABLE_NAME = 'asset_accounting_mappings'
        AND s2.COLUMN_NAME = 'asset_accounting_group_id'
    )
);
SET @add_company_index = IF(
  @has_company_index = 0,
  'ALTER TABLE asset_accounting_mappings ADD KEY idx_asset_accounting_mappings_company (company_id)',
  'SELECT 1'
);
PREPARE add_company_index_stmt FROM @add_company_index;
EXECUTE add_company_index_stmt;
DEALLOCATE PREPARE add_company_index_stmt;

SET @old_mapping_indexes = (
  SELECT GROUP_CONCAT(DISTINCT CONCAT('DROP INDEX `', REPLACE(s.INDEX_NAME, '`', '``'), '`') SEPARATOR ', ')
  FROM information_schema.STATISTICS s
  WHERE s.TABLE_SCHEMA = DATABASE()
    AND s.TABLE_NAME = 'asset_accounting_mappings'
    AND s.COLUMN_NAME = 'asset_accounting_group_id'
    AND s.INDEX_NAME <> 'PRIMARY'
);
SET @drop_old_mapping_indexes = IF(
  @old_mapping_indexes IS NULL,
  'SELECT 1',
  CONCAT('ALTER TABLE asset_accounting_mappings ', @old_mapping_indexes)
);
PREPARE drop_old_mapping_indexes_stmt FROM @drop_old_mapping_indexes;
EXECUTE drop_old_mapping_indexes_stmt;
DEALLOCATE PREPARE drop_old_mapping_indexes_stmt;

SET @has_category_column = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'asset_accounting_mappings' AND COLUMN_NAME = 'asset_category_id'
);
SET @add_category_column = IF(
  @has_category_column = 0,
  'ALTER TABLE asset_accounting_mappings ADD COLUMN asset_category_id BIGINT UNSIGNED NULL AFTER company_id',
  'SELECT 1'
);
PREPARE add_category_column_stmt FROM @add_category_column;
EXECUTE add_category_column_stmt;
DEALLOCATE PREPARE add_category_column_stmt;

SET @has_old_group_column = (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'asset_accounting_mappings' AND COLUMN_NAME = 'asset_accounting_group_id'
);
SET @drop_old_group_column = IF(
  @has_old_group_column = 1,
  'ALTER TABLE asset_accounting_mappings DROP COLUMN asset_accounting_group_id',
  'SELECT 1'
);
PREPARE drop_old_group_column_stmt FROM @drop_old_group_column;
EXECUTE drop_old_group_column_stmt;
DEALLOCATE PREPARE drop_old_group_column_stmt;

ALTER TABLE asset_accounting_mappings
  MODIFY asset_category_id BIGINT UNSIGNED NOT NULL;

SET @has_category_unique = (
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'asset_accounting_mappings'
    AND INDEX_NAME = 'uq_asset_accounting_mappings_company_category'
);
SET @add_category_unique = IF(
  @has_category_unique = 0,
  'ALTER TABLE asset_accounting_mappings ADD UNIQUE KEY uq_asset_accounting_mappings_company_category (company_id, asset_category_id)',
  'SELECT 1'
);
PREPARE add_category_unique_stmt FROM @add_category_unique;
EXECUTE add_category_unique_stmt;
DEALLOCATE PREPARE add_category_unique_stmt;

SET @has_category_index = (
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'asset_accounting_mappings'
    AND INDEX_NAME = 'idx_asset_accounting_mappings_category'
);
SET @add_category_index = IF(
  @has_category_index = 0,
  'ALTER TABLE asset_accounting_mappings ADD KEY idx_asset_accounting_mappings_category (asset_category_id)',
  'SELECT 1'
);
PREPARE add_category_index_stmt FROM @add_category_index;
EXECUTE add_category_index_stmt;
DEALLOCATE PREPARE add_category_index_stmt;

SET @has_category_fk = (
  SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 'asset_accounting_mappings'
    AND CONSTRAINT_NAME = 'fk_asset_accounting_mappings_category'
);
SET @add_category_fk = IF(
  @has_category_fk = 0,
  'ALTER TABLE asset_accounting_mappings ADD CONSTRAINT fk_asset_accounting_mappings_category FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT',
  'SELECT 1'
);
PREPARE add_category_fk_stmt FROM @add_category_fk;
EXECUTE add_category_fk_stmt;
DEALLOCATE PREPARE add_category_fk_stmt;

DELETE usp FROM user_store_permissions usp
JOIN permissions p ON p.id = usp.permission_id
WHERE p.name LIKE 'asset_accounting_groups.%';
DELETE rp FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE p.name LIKE 'asset_accounting_groups.%';
DELETE FROM permissions WHERE name LIKE 'asset_accounting_groups.%';

DROP TABLE IF EXISTS asset_accounting_groups;

COMMIT;

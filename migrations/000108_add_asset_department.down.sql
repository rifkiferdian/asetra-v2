ALTER TABLE assets
  DROP FOREIGN KEY fk_assets_department,
  DROP INDEX idx_assets_department,
  DROP COLUMN department_id;

-- Department menyatakan unit pemilik/pengguna aset; Branch tetap menyatakan penempatan fisik.
ALTER TABLE assets
  ADD COLUMN department_id BIGINT UNSIGNED NULL AFTER asset_location_id,
  ADD KEY idx_assets_department (department_id),
  ADD CONSTRAINT fk_assets_department FOREIGN KEY (department_id)
    REFERENCES departments(id) ON UPDATE CASCADE ON DELETE RESTRICT;

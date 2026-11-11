-- Dummy data registry asset: tepat 2 aset untuk MK5.
-- Bergantung pada seed 000100-000103 dan aman dijalankan ulang.
START TRANSACTION;

SET @branch_id = (SELECT id FROM branches WHERE code='MK5' ORDER BY id LIMIT 1);
SET @company_id = (SELECT company_id FROM branches WHERE id=@branch_id LIMIT 1);
SET @group_id = (
  SELECT c.group_id FROM companies c WHERE c.id=@company_id LIMIT 1
);

SET @location_programmer = (
  SELECT id FROM asset_locations
  WHERE branch_id=@branch_id AND code='PROGRAMMER'
  LIMIT 1
);
SET @location_akuntansi = (
  SELECT id FROM asset_locations
  WHERE branch_id=@branch_id AND code='AKUNTANSI'
  LIMIT 1
);

SET @type_laptop = (
  SELECT t.id FROM asset_types t
  JOIN asset_categories c ON c.id=t.asset_category_id
  WHERE c.group_id=@group_id AND t.code='LAPTOP'
  LIMIT 1
);
SET @type_desk = (
  SELECT t.id FROM asset_types t
  JOIN asset_categories c ON c.id=t.asset_category_id
  WHERE c.group_id=@group_id AND t.code='OFFICE-DESK'
  LIMIT 1
);
SET @brand_lenovo = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='LENOVO' LIMIT 1);
SET @brand_ikea = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='IKEA' LIMIT 1);
SET @model_thinkpad = (SELECT id FROM asset_models WHERE asset_type_id=@type_laptop AND code='THINKPAD-E14' LIMIT 1);
SET @model_bekant = (SELECT id FROM asset_models WHERE asset_type_id=@type_desk AND code='BEKANT-DESK' LIMIT 1);
SET @variant_thinkpad = (SELECT id FROM asset_variants WHERE asset_model_id=@model_thinkpad AND code='I5-16-512' LIMIT 1);
SET @variant_bekant = (SELECT id FROM asset_variants WHERE asset_model_id=@model_bekant AND code='WHITE-140' LIMIT 1);
SET @condition_good = (SELECT id FROM asset_conditions WHERE code='GOOD' LIMIT 1);
SET @status_in_use = (SELECT id FROM asset_statuses WHERE code='IN-USE' LIMIT 1);
SET @status_available = (SELECT id FROM asset_statuses WHERE code='AVAILABLE' LIMIT 1);
SET @vendor_tech = (SELECT id FROM vendors WHERE company_group_id=@group_id AND code='VND-TECH' LIMIT 1);
SET @vendor_general = (SELECT id FROM vendors WHERE company_group_id=@group_id AND code='VND-GENERAL' LIMIT 1);
SET @department_it = (SELECT id FROM departments WHERE group_id=@group_id AND code='IT' LIMIT 1);
SET @department_fin = (SELECT id FROM departments WHERE group_id=@group_id AND code='FIN' LIMIT 1);

INSERT INTO assets(
  company_id,branch_id,asset_location_id,department_id,asset_type_id,asset_brand_id,asset_model_id,
  asset_variant_id,condition_id,status_id,vendor_id,asset_number,name,serial_number,
  barcode,description,acquisition_date,acquisition_cost,warranty_until,notes
) VALUES
(@company_id,@branch_id,@location_programmer,@department_it,@type_laptop,@brand_lenovo,@model_thinkpad,
 @variant_thinkpad,@condition_good,@status_in_use,@vendor_tech,'AST-MK5-0001',
 'Laptop Programmer 01','PF4MK50001','8990000000001','Laptop kerja tim programmer',
 '2026-01-15',14500000.00,'2029-01-15','Dummy asset development'),
(@company_id,@branch_id,@location_akuntansi,@department_fin,@type_desk,@brand_ikea,@model_bekant,
 @variant_bekant,@condition_good,@status_available,@vendor_general,'AST-MK5-0002',
 'Meja Kerja Akuntansi 01','BKT-MK5-0001','8990000000002','Meja kerja untuk ruang akuntansi',
 '2026-02-10',4250000.00,NULL,'Dummy asset development')
ON DUPLICATE KEY UPDATE
  company_id=VALUES(company_id),branch_id=VALUES(branch_id),asset_location_id=VALUES(asset_location_id),department_id=VALUES(department_id),
  asset_type_id=VALUES(asset_type_id),asset_brand_id=VALUES(asset_brand_id),asset_model_id=VALUES(asset_model_id),
  asset_variant_id=VALUES(asset_variant_id),condition_id=VALUES(condition_id),status_id=VALUES(status_id),
  vendor_id=VALUES(vendor_id),name=VALUES(name),serial_number=VALUES(serial_number),
  barcode=VALUES(barcode),description=VALUES(description),acquisition_date=VALUES(acquisition_date),
  acquisition_cost=VALUES(acquisition_cost),warranty_until=VALUES(warranty_until),notes=VALUES(notes);

COMMIT;

-- PC operasional untuk monitoring dan administrasi ruang server MK5.
START TRANSACTION;

SET @branch_id = (SELECT id FROM branches WHERE code='MK5' ORDER BY id LIMIT 1);
SET @company_id = (SELECT company_id FROM branches WHERE id=@branch_id LIMIT 1);
SET @group_id = (SELECT group_id FROM companies WHERE id=@company_id LIMIT 1);
SET @location_server = (
  SELECT id FROM asset_locations WHERE branch_id=@branch_id AND code='SERVER' LIMIT 1
);

INSERT INTO asset_categories(group_id,code,name,description,is_active) VALUES
(@group_id,'ELECTRONIC','Elektronik','Perangkat elektronik dan teknologi informasi',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @category_electronic = (
  SELECT id FROM asset_categories WHERE group_id=@group_id AND code='ELECTRONIC' LIMIT 1
);

INSERT INTO asset_types(asset_category_id,code,name,description,is_active) VALUES
(@category_electronic,'DESKTOP-PC','Desktop PC','Komputer desktop untuk kebutuhan operasional',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_brands(group_id,code,name,description,is_active) VALUES
(@group_id,'DELL','Dell','Perangkat komputer dan server',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @type_desktop = (
  SELECT t.id FROM asset_types t
  JOIN asset_categories c ON c.id=t.asset_category_id
  WHERE c.group_id=@group_id AND t.code='DESKTOP-PC' LIMIT 1
);
SET @brand_dell = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='DELL' LIMIT 1);

INSERT IGNORE INTO asset_type_brands(asset_type_id,asset_brand_id)
VALUES(@type_desktop,@brand_dell);

INSERT INTO asset_models(asset_type_id,asset_brand_id,code,name,description,is_active) VALUES
(@type_desktop,@brand_dell,'OPTIPLEX-7010','OptiPlex 7010','PC monitoring dan administrasi ruang server',1)
ON DUPLICATE KEY UPDATE asset_brand_id=VALUES(asset_brand_id),name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @model_optiplex = (
  SELECT id FROM asset_models WHERE asset_type_id=@type_desktop AND asset_brand_id=@brand_dell AND code='OPTIPLEX-7010' LIMIT 1
);

INSERT INTO asset_variants(asset_model_id,code,name,description,is_active) VALUES
(@model_optiplex,'I7-16-512','i7 16GB 512GB','Intel Core i7, RAM 16 GB, SSD 512 GB',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @variant_optiplex = (
  SELECT id FROM asset_variants WHERE asset_model_id=@model_optiplex AND code='I7-16-512' LIMIT 1
);
SET @condition_good = (SELECT id FROM asset_conditions WHERE code='GOOD' LIMIT 1);
SET @status_in_use = (SELECT id FROM asset_statuses WHERE code='IN-USE' LIMIT 1);

INSERT INTO assets(
  company_id,branch_id,asset_location_id,department_id,asset_type_id,asset_brand_id,asset_model_id,asset_variant_id,
  condition_id,status_id,vendor_id,asset_number,name,serial_number,barcode,description,
  acquisition_date,acquisition_cost,warranty_until,notes
) VALUES(
  @company_id,@branch_id,@location_server,NULL,@type_desktop,@brand_dell,@model_optiplex,@variant_optiplex,
  @condition_good,@status_in_use,NULL,'AST/MK5/SERVER/009','Dell OptiPlex 7010 i7 16GB 512GB',
  'SRV-PC-001','ASTMK5SERVER009','PC untuk monitoring dan administrasi perangkat ruang server',
  '2026-01-05',18500000.00,'2029-01-05','Dummy asset migration MK5 server PC'
)
ON DUPLICATE KEY UPDATE
  asset_location_id=VALUES(asset_location_id),asset_type_id=VALUES(asset_type_id),asset_brand_id=VALUES(asset_brand_id),
  asset_model_id=VALUES(asset_model_id),asset_variant_id=VALUES(asset_variant_id),condition_id=VALUES(condition_id),
  status_id=VALUES(status_id),name=VALUES(name),serial_number=VALUES(serial_number),barcode=VALUES(barcode),
  description=VALUES(description),acquisition_date=VALUES(acquisition_date),acquisition_cost=VALUES(acquisition_cost),
  warranty_until=VALUES(warranty_until),notes=VALUES(notes);

COMMIT;

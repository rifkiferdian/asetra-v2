-- Dummy Master Data Asset.
-- Masing-masing 3 category, type, brand, model, dan variant.
-- Idempotent: aman dijalankan ulang pada database development.
START TRANSACTION;

SET @group_id = COALESCE(
  (SELECT id FROM company_groups WHERE code='MANNA' LIMIT 1),
  (SELECT id FROM company_groups WHERE is_active=1 ORDER BY id LIMIT 1)
);

INSERT INTO asset_categories(group_id,code,name,description,is_active) VALUES
(@group_id,'ELECTRONIC','Elektronik','Perangkat elektronik dan teknologi informasi',1),
(@group_id,'FURNITURE','Furniture','Meja, kursi, lemari, dan perlengkapan ruangan',1),
(@group_id,'VEHICLE','Kendaraan','Kendaraan operasional perusahaan',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @category_electronic = (SELECT id FROM asset_categories WHERE group_id=@group_id AND code='ELECTRONIC' LIMIT 1);
SET @category_furniture = (SELECT id FROM asset_categories WHERE group_id=@group_id AND code='FURNITURE' LIMIT 1);
SET @category_vehicle = (SELECT id FROM asset_categories WHERE group_id=@group_id AND code='VEHICLE' LIMIT 1);

INSERT INTO asset_types(asset_category_id,code,name,description,is_active) VALUES
(@category_electronic,'LAPTOP','Laptop','Komputer portabel untuk kebutuhan kerja',1),
(@category_furniture,'OFFICE-DESK','Meja Kantor','Meja kerja untuk area kantor',1),
(@category_vehicle,'MOTORCYCLE','Sepeda Motor','Sepeda motor untuk kegiatan operasional',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_brands(group_id,code,name,description,is_active) VALUES
(@group_id,'LENOVO','Lenovo','Brand perangkat komputer',1),
(@group_id,'IKEA','IKEA','Brand furniture dan perlengkapan kantor',1),
(@group_id,'HONDA','Honda','Brand kendaraan bermotor',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @type_laptop = (SELECT id FROM asset_types WHERE asset_category_id=@category_electronic AND code='LAPTOP' LIMIT 1);
SET @type_desk = (SELECT id FROM asset_types WHERE asset_category_id=@category_furniture AND code='OFFICE-DESK' LIMIT 1);
SET @type_motorcycle = (SELECT id FROM asset_types WHERE asset_category_id=@category_vehicle AND code='MOTORCYCLE' LIMIT 1);
SET @brand_lenovo = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='LENOVO' LIMIT 1);
SET @brand_ikea = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='IKEA' LIMIT 1);
SET @brand_honda = (SELECT id FROM asset_brands WHERE group_id=@group_id AND code='HONDA' LIMIT 1);

INSERT INTO asset_models(asset_type_id,asset_brand_id,code,name,description,is_active) VALUES
(@type_laptop,@brand_lenovo,'THINKPAD-E14','ThinkPad E14','Laptop bisnis seri ThinkPad E14',1),
(@type_desk,@brand_ikea,'BEKANT-DESK','Bekant Desk','Meja kerja seri Bekant',1),
(@type_motorcycle,@brand_honda,'VARIO-160','Vario 160','Sepeda motor operasional Vario 160',1)
ON DUPLICATE KEY UPDATE asset_brand_id=VALUES(asset_brand_id),name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @model_thinkpad = (SELECT id FROM asset_models WHERE asset_type_id=@type_laptop AND code='THINKPAD-E14' LIMIT 1);
SET @model_bekant = (SELECT id FROM asset_models WHERE asset_type_id=@type_desk AND code='BEKANT-DESK' LIMIT 1);
SET @model_vario = (SELECT id FROM asset_models WHERE asset_type_id=@type_motorcycle AND code='VARIO-160' LIMIT 1);

INSERT INTO asset_variants(asset_model_id,code,name,description,is_active) VALUES
(@model_thinkpad,'I5-16-512','Core i5 / 16 GB / 512 GB','Varian RAM 16 GB dan SSD 512 GB',1),
(@model_bekant,'WHITE-140','Putih 140 cm','Varian warna putih dengan panjang 140 cm',1),
(@model_vario,'CBS-BLACK','CBS Active Black','Varian CBS warna Active Black',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

COMMIT;

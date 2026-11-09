-- Mengganti seluruh Asset Location development dengan hierarki MK5 Godean.
START TRANSACTION;

SET @mk5_branch_id = (
  SELECT id FROM branches WHERE code='MK5' ORDER BY id LIMIT 1
);

-- Hapus child lebih dahulu agar foreign key hierarki tetap terpenuhi.
DELETE FROM asset_locations WHERE parent_id IS NOT NULL;
DELETE FROM asset_locations WHERE parent_id IS NULL;

INSERT INTO asset_locations(branch_id,parent_id,code,name,location_type,description,is_active) VALUES
(@mk5_branch_id,NULL,'TOKO','TOKO','AREA','Area utama toko MK5',1),
(@mk5_branch_id,NULL,'GUDANG','GUDANG','AREA','Area utama gudang MK5',1),
(@mk5_branch_id,NULL,'OFFICE','OFFICE','AREA','Area utama office MK5',1);

SET @toko_id = (SELECT id FROM asset_locations WHERE branch_id=@mk5_branch_id AND code='TOKO' LIMIT 1);
SET @gudang_id = (SELECT id FROM asset_locations WHERE branch_id=@mk5_branch_id AND code='GUDANG' LIMIT 1);
SET @office_id = (SELECT id FROM asset_locations WHERE branch_id=@mk5_branch_id AND code='OFFICE' LIMIT 1);

INSERT INTO asset_locations(branch_id,parent_id,code,name,location_type,description,is_active) VALUES
(@mk5_branch_id,@toko_id,'KASSA-UTARA','KASSA UTARA','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'KASSA-SELATAN','KASSA SELATAN','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'TOILETRIS','TOILETRIS','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'KOSMETIK','KOSMETIK','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'FOOD','FOOD','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'LISTRIK','LISTRIK','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'PECAH-BELAH','PECAH BELAH','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'AREA-PROMO','AREA PROMO','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'MANNA-BAKERY','MANNA BAKERY','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'PARKIR-UTARA','PARKIR UTARA','AREA',NULL,1),
(@mk5_branch_id,@toko_id,'PARKIR-SELATAN','PARKIR SELATAN','AREA',NULL,1),

(@mk5_branch_id,@gudang_id,'RETUR','RETUR','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'PENERIMAAN','PENERIMAAN','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'INPUT','INPUT','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'LLB','LLB','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'FRESHFOOD','FRESHFOOD','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'GUDANG-LT-1','GUDANG LT 1','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'GUDANG-LT-2','GUDANG LT 2','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'PACKING','PACKING','AREA',NULL,1),
(@mk5_branch_id,@gudang_id,'KAS-BESAR','KAS BESAR','AREA',NULL,1),

(@mk5_branch_id,@office_id,'MK-ONLINE','MK ONLINE','AREA',NULL,1),
(@mk5_branch_id,@office_id,'ME','ME','AREA',NULL,1),
(@mk5_branch_id,@office_id,'SERVER','SERVER','AREA',NULL,1),
(@mk5_branch_id,@office_id,'RUANG-TAMU','RUANG TAMU','AREA',NULL,1),
(@mk5_branch_id,@office_id,'PEMBELIAN','PEMBELIAN','AREA',NULL,1),
(@mk5_branch_id,@office_id,'PEMBAYARAN','PEMBAYARAN','AREA',NULL,1),
(@mk5_branch_id,@office_id,'AKUNTANSI','AKUNTANSI','AREA',NULL,1),
(@mk5_branch_id,@office_id,'BANK-DATA','BANK DATA','AREA',NULL,1),
(@mk5_branch_id,@office_id,'SOSMED','SOSMED','AREA',NULL,1),
(@mk5_branch_id,@office_id,'PROGRAMMER','PROGRAMMER','AREA',NULL,1),
(@mk5_branch_id,@office_id,'PERSONALIA','PERSONALIA','AREA',NULL,1),
(@mk5_branch_id,@office_id,'INTERVIEW','INTERVIEW','AREA',NULL,1),
(@mk5_branch_id,@office_id,'KINASIH','KINASIH','AREA',NULL,1),
(@mk5_branch_id,@office_id,'KINASIH-BARU','KINASIH BARU','AREA',NULL,1),
(@mk5_branch_id,@office_id,'ASMEN-AUDIT','ASMEN AUDIT','AREA',NULL,1),
(@mk5_branch_id,@office_id,'HRD-PENGEMBANGAN','HRD PENGEMBANGAN','AREA',NULL,1),
(@mk5_branch_id,@office_id,'MONITORING','MONITORING','AREA',NULL,1);

COMMIT;

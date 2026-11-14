-- Melengkapi dummy aset untuk area office MK5.
-- SERVER: 8 aset infrastruktur; PROGRAMMER, SOSMED, BANK DATA, dan ME: masing-masing 3 aset.
-- Idempotent berdasarkan asset_number dan kode master.
START TRANSACTION;

SET @branch_id = (SELECT id FROM branches WHERE code='MK5' ORDER BY id LIMIT 1);
SET @company_id = (SELECT company_id FROM branches WHERE id=@branch_id LIMIT 1);
SET @group_id = (SELECT group_id FROM companies WHERE id=@company_id LIMIT 1);

INSERT INTO asset_conditions(code,name,description,is_active) VALUES
('GOOD','Baik','Aset berfungsi normal dan layak digunakan',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_statuses(code,name,description,is_active) VALUES
('IN-USE','Digunakan','Aset sedang digunakan atau ditugaskan',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_categories(group_id,code,name,description,is_active) VALUES
(@group_id,'ELECTRONIC','Elektronik','Perangkat elektronik dan teknologi informasi',1),
(@group_id,'TOOLS','Peralatan Teknik','Peralatan kerja dan pengukuran teknis',1),
(@group_id,'INFRASTRUCTURE','Infrastruktur TI','Perangkat infrastruktur jaringan dan pusat data',1),
(@group_id,'HVAC','Tata Udara','Perangkat pendingin dan tata udara',1),
(@group_id,'FURNITURE','Furniture','Furniture dan kabinet operasional',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_types(asset_category_id,code,name,description,is_active)
SELECT c.id,x.code,x.name,x.description,1
FROM (
  SELECT 'ELECTRONIC' category_code,'LAPTOP' code,'Laptop' name,'Komputer portabel untuk kebutuhan kerja' description UNION ALL
  SELECT 'ELECTRONIC','CAMERA','Kamera','Kamera produksi konten digital' UNION ALL
  SELECT 'ELECTRONIC','GIMBAL','Gimbal Kamera','Stabilisator kamera' UNION ALL
  SELECT 'ELECTRONIC','SCANNER','Scanner Dokumen','Pemindai dokumen operasional' UNION ALL
  SELECT 'INFRASTRUCTURE','NAS','Network Attached Storage','Penyimpanan data melalui jaringan' UNION ALL
  SELECT 'TOOLS','MULTIMETER','Multimeter','Alat ukur kelistrikan' UNION ALL
  SELECT 'TOOLS','POWER-DRILL','Bor Listrik','Peralatan bor listrik' UNION ALL
  SELECT 'TOOLS','BLOWER','Blower','Peralatan pembersih dan pemeliharaan' UNION ALL
  SELECT 'INFRASTRUCTURE','RACK-SERVER','Rack Server','Server fisik untuk pusat data' UNION ALL
  SELECT 'INFRASTRUCTURE','NETWORK-SWITCH','Network Switch','Perangkat distribusi jaringan' UNION ALL
  SELECT 'INFRASTRUCTURE','FIREWALL','Firewall','Perangkat keamanan jaringan' UNION ALL
  SELECT 'INFRASTRUCTURE','UPS','UPS','Catu daya tanpa putus' UNION ALL
  SELECT 'FURNITURE','SERVER-RACK','Kabinet Server','Kabinet untuk perangkat pusat data' UNION ALL
  SELECT 'HVAC','PRECISION-AC','AC Server','Pendingin khusus ruang server'
) x
JOIN asset_categories c ON c.group_id=@group_id AND c.code=x.category_code
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_brands(group_id,code,name,description,is_active) VALUES
(@group_id,'LENOVO','Lenovo','Perangkat komputer',1),
(@group_id,'APPLE','Apple','Perangkat komputer dan kreatif',1),
(@group_id,'CANON','Canon','Perangkat kamera',1),
(@group_id,'DJI','DJI','Perangkat produksi konten',1),
(@group_id,'BROTHER','Brother','Perangkat pemindai dokumen',1),
(@group_id,'FUJITSU','Fujitsu','Perangkat pemindai dokumen',1),
(@group_id,'SYNOLOGY','Synology','Perangkat penyimpanan jaringan',1),
(@group_id,'FLUKE','Fluke','Peralatan ukur teknik',1),
(@group_id,'BOSCH','Bosch','Peralatan teknik',1),
(@group_id,'MAKITA','Makita','Peralatan teknik',1),
(@group_id,'DELL','Dell','Perangkat server',1),
(@group_id,'CISCO','Cisco','Perangkat jaringan',1),
(@group_id,'FORTINET','Fortinet','Perangkat keamanan jaringan',1),
(@group_id,'APC','APC','Perangkat daya dan kabinet pusat data',1),
(@group_id,'DAIKIN','Daikin','Perangkat tata udara',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT IGNORE INTO asset_type_brands(asset_type_id,asset_brand_id)
SELECT t.id,b.id FROM (
  SELECT 'LAPTOP' type_code,'LENOVO' brand_code UNION ALL SELECT 'LAPTOP','APPLE' UNION ALL
  SELECT 'CAMERA','CANON' UNION ALL SELECT 'GIMBAL','DJI' UNION ALL
  SELECT 'SCANNER','BROTHER' UNION ALL SELECT 'SCANNER','FUJITSU' UNION ALL
  SELECT 'NAS','SYNOLOGY' UNION ALL SELECT 'MULTIMETER','FLUKE' UNION ALL
  SELECT 'POWER-DRILL','BOSCH' UNION ALL SELECT 'BLOWER','MAKITA' UNION ALL
  SELECT 'RACK-SERVER','DELL' UNION ALL SELECT 'NETWORK-SWITCH','CISCO' UNION ALL
  SELECT 'FIREWALL','FORTINET' UNION ALL SELECT 'UPS','APC' UNION ALL
  SELECT 'SERVER-RACK','APC' UNION ALL SELECT 'PRECISION-AC','DAIKIN'
) x
JOIN asset_types t ON t.code=x.type_code
JOIN asset_categories c ON c.id=t.asset_category_id AND c.group_id=@group_id
JOIN asset_brands b ON b.group_id=@group_id AND b.code=x.brand_code;

INSERT INTO asset_models(asset_type_id,asset_brand_id,code,name,description,is_active)
SELECT t.id,b.id,x.code,x.name,x.description,1 FROM (
  SELECT 'LAPTOP' type_code,'LENOVO' brand_code,'THINKPAD-E14' code,'ThinkPad E14' name,'Laptop bisnis untuk tim programmer' description UNION ALL
  SELECT 'LAPTOP','APPLE','MACBOOK-PRO-14','MacBook Pro 14','Laptop produksi konten' UNION ALL
  SELECT 'CAMERA','CANON','EOS-R10','EOS R10','Kamera mirrorless produksi konten' UNION ALL
  SELECT 'GIMBAL','DJI','RS-3-MINI','RS 3 Mini','Gimbal kamera ringkas' UNION ALL
  SELECT 'SCANNER','BROTHER','ADS-4300N','ADS-4300N','Scanner dokumen jaringan' UNION ALL
  SELECT 'SCANNER','FUJITSU','FI-8170','fi-8170','Scanner dokumen volume tinggi' UNION ALL
  SELECT 'NAS','SYNOLOGY','DS923-PLUS','DS923+','Penyimpanan dokumen terpusat' UNION ALL
  SELECT 'MULTIMETER','FLUKE','FLUKE-117','117','Multimeter digital teknisi' UNION ALL
  SELECT 'POWER-DRILL','BOSCH','GSB-13-RE','GSB 13 RE','Bor listrik teknisi' UNION ALL
  SELECT 'BLOWER','MAKITA','M4000B','M4000B','Blower pemeliharaan' UNION ALL
  SELECT 'RACK-SERVER','DELL','POWEREDGE-R550','PowerEdge R550','Server rack operasional' UNION ALL
  SELECT 'NETWORK-SWITCH','CISCO','CBS350-24P','CBS350-24P','Managed switch pusat data' UNION ALL
  SELECT 'FIREWALL','FORTINET','FORTIGATE-100F','FortiGate 100F','Firewall perimeter jaringan' UNION ALL
  SELECT 'UPS','APC','SRT3000RMXLI','Smart-UPS SRT 3000VA','UPS rack ruang server' UNION ALL
  SELECT 'SERVER-RACK','APC','NETSHELTER-SX-42U','NetShelter SX 42U','Kabinet perangkat ruang server' UNION ALL
  SELECT 'PRECISION-AC','DAIKIN','FTKC35','FTKC35','Pendingin ruang server'
) x
JOIN asset_types t ON t.code=x.type_code
JOIN asset_categories c ON c.id=t.asset_category_id AND c.group_id=@group_id
JOIN asset_brands b ON b.group_id=@group_id AND b.code=x.brand_code
ON DUPLICATE KEY UPDATE asset_brand_id=VALUES(asset_brand_id),name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_variants(asset_model_id,code,name,description,is_active)
SELECT m.id,x.code,x.name,x.description,1 FROM (
  SELECT 'THINKPAD-E14' model_code,'I5-16-512' code,'i5 16GB 512GB' name,'Intel Core i5, RAM 16 GB, SSD 512 GB' description UNION ALL
  SELECT 'MACBOOK-PRO-14','M3-18-512','M3 Pro 18GB 512GB','Apple M3 Pro, memori 18 GB, SSD 512 GB' UNION ALL
  SELECT 'EOS-R10','KIT-18-150','Kit 18-150mm','Kamera dengan lensa kit RF-S 18-150mm' UNION ALL
  SELECT 'RS-3-MINI','BLACK','Black','Warna hitam' UNION ALL
  SELECT 'ADS-4300N','A4-DUPLEX','A4 Duplex','Pemindaian A4 dua sisi' UNION ALL
  SELECT 'FI-8170','A4-DUPLEX','A4 Duplex','Pemindaian A4 dua sisi' UNION ALL
  SELECT 'DS923-PLUS','4X8TB','4x8TB','Empat hard disk 8 TB' UNION ALL
  SELECT 'FLUKE-117','TRUE-RMS','True RMS','Varian pengukuran True RMS' UNION ALL
  SELECT 'GSB-13-RE','13MM','13mm','Chuck 13 mm' UNION ALL
  SELECT 'M4000B','530W','530W','Daya 530 watt' UNION ALL
  SELECT 'POWEREDGE-R550','XEON-64-2TB','Xeon 64GB 2TB','Xeon, RAM 64 GB, storage 2 TB' UNION ALL
  SELECT 'CBS350-24P','24P-POE','24 Port PoE','Dua puluh empat port dengan PoE' UNION ALL
  SELECT 'FORTIGATE-100F','FG-100F','100F','Appliance FortiGate 100F' UNION ALL
  SELECT 'SRT3000RMXLI','3000VA','3000VA','Kapasitas 3000 VA' UNION ALL
  SELECT 'NETSHELTER-SX-42U','42U','42U','Kabinet rack 42U' UNION ALL
  SELECT 'FTKC35','1-5PK','1.5 PK','Kapasitas pendingin 1,5 PK'
) x
JOIN asset_models m ON m.code=x.model_code
JOIN asset_types t ON t.id=m.asset_type_id
JOIN asset_categories c ON c.id=t.asset_category_id AND c.group_id=@group_id
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

DROP TEMPORARY TABLE IF EXISTS tmp_mk5_assets;
CREATE TEMPORARY TABLE tmp_mk5_assets (
  location_code VARCHAR(30), type_code VARCHAR(50), brand_code VARCHAR(50), model_code VARCHAR(50), variant_code VARCHAR(50),
  asset_number VARCHAR(50), serial_number VARCHAR(120), barcode VARCHAR(120), description VARCHAR(255),
  acquisition_date DATE, acquisition_cost DECIMAL(18,2)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO tmp_mk5_assets VALUES
('PROGRAMMER','LAPTOP','LENOVO','THINKPAD-E14','I5-16-512','AST/MK5/PROGRAMMER/001','PRG-E14-001','ASTMK5PROGRAMMER001','Laptop kerja programmer 1','2026-01-15',14500000),
('PROGRAMMER','LAPTOP','LENOVO','THINKPAD-E14','I5-16-512','AST/MK5/PROGRAMMER/002','PRG-E14-002','ASTMK5PROGRAMMER002','Laptop kerja programmer 2','2026-01-15',14500000),
('PROGRAMMER','LAPTOP','LENOVO','THINKPAD-E14','I5-16-512','AST/MK5/PROGRAMMER/003','PRG-E14-003','ASTMK5PROGRAMMER003','Laptop kerja programmer 3','2026-01-15',14500000),
('SOSMED','LAPTOP','APPLE','MACBOOK-PRO-14','M3-18-512','AST/MK5/SOSMED/001','SMD-MBP-001','ASTMK5SOSMED001','Laptop editing konten sosial media','2026-02-01',32000000),
('SOSMED','CAMERA','CANON','EOS-R10','KIT-18-150','AST/MK5/SOSMED/002','SMD-R10-001','ASTMK5SOSMED002','Kamera produksi konten sosial media','2026-02-01',18500000),
('SOSMED','GIMBAL','DJI','RS-3-MINI','BLACK','AST/MK5/SOSMED/003','SMD-RS3-001','ASTMK5SOSMED003','Gimbal produksi video sosial media','2026-02-01',5500000),
('BANK-DATA','SCANNER','BROTHER','ADS-4300N','A4-DUPLEX','AST/MK5/BANKDATA/001','BD-ADS-001','ASTMK5BANKDATA001','Scanner dokumen jaringan','2026-02-10',8500000),
('BANK-DATA','SCANNER','FUJITSU','FI-8170','A4-DUPLEX','AST/MK5/BANKDATA/002','BD-FI-001','ASTMK5BANKDATA002','Scanner arsip volume tinggi','2026-02-10',14500000),
('BANK-DATA','NAS','SYNOLOGY','DS923-PLUS','4X8TB','AST/MK5/BANKDATA/003','BD-NAS-001','ASTMK5BANKDATA003','Penyimpanan arsip digital','2026-02-10',28000000),
('ME','MULTIMETER','FLUKE','FLUKE-117','TRUE-RMS','AST/MK5/ME/001','ME-FLK-001','ASTMK5ME001','Multimeter teknisi ME','2026-03-01',5500000),
('ME','POWER-DRILL','BOSCH','GSB-13-RE','13MM','AST/MK5/ME/002','ME-BSH-001','ASTMK5ME002','Bor listrik teknisi ME','2026-03-01',1800000),
('ME','BLOWER','MAKITA','M4000B','530W','AST/MK5/ME/003','ME-MKT-001','ASTMK5ME003','Blower pemeliharaan teknisi ME','2026-03-01',1600000),
('SERVER','SERVER-RACK','APC','NETSHELTER-SX-42U','42U','AST/MK5/SERVER/001','SRV-RACK-001','ASTMK5SERVER001','Kabinet utama ruang server','2026-01-05',24000000),
('SERVER','RACK-SERVER','DELL','POWEREDGE-R550','XEON-64-2TB','AST/MK5/SERVER/002','SRV-DELL-001','ASTMK5SERVER002','Server aplikasi utama','2026-01-05',85000000),
('SERVER','RACK-SERVER','DELL','POWEREDGE-R550','XEON-64-2TB','AST/MK5/SERVER/003','SRV-DELL-002','ASTMK5SERVER003','Server aplikasi cadangan','2026-01-05',85000000),
('SERVER','NETWORK-SWITCH','CISCO','CBS350-24P','24P-POE','AST/MK5/SERVER/004','SRV-CSC-001','ASTMK5SERVER004','Managed switch utama','2026-01-05',12500000),
('SERVER','FIREWALL','FORTINET','FORTIGATE-100F','FG-100F','AST/MK5/SERVER/005','SRV-FTG-001','ASTMK5SERVER005','Firewall perimeter jaringan','2026-01-05',48000000),
('SERVER','UPS','APC','SRT3000RMXLI','3000VA','AST/MK5/SERVER/006','SRV-UPS-001','ASTMK5SERVER006','UPS perangkat ruang server','2026-01-05',35000000),
('SERVER','NAS','SYNOLOGY','DS923-PLUS','4X8TB','AST/MK5/SERVER/007','SRV-NAS-001','ASTMK5SERVER007','Penyimpanan backup server','2026-01-05',28000000),
('SERVER','PRECISION-AC','DAIKIN','FTKC35','1-5PK','AST/MK5/SERVER/008','SRV-AC-001','ASTMK5SERVER008','Pendingin khusus ruang server','2026-01-05',12000000);

SET @condition_good = (SELECT id FROM asset_conditions WHERE code='GOOD' LIMIT 1);
SET @status_in_use = (SELECT id FROM asset_statuses WHERE code='IN-USE' LIMIT 1);

INSERT INTO assets(
  company_id,branch_id,asset_location_id,department_id,asset_type_id,asset_brand_id,asset_model_id,asset_variant_id,
  condition_id,status_id,vendor_id,asset_number,name,serial_number,barcode,description,acquisition_date,acquisition_cost,warranty_until,notes
)
SELECT @company_id,@branch_id,l.id,NULL,t.id,b.id,m.id,v.id,@condition_good,@status_in_use,NULL,
  x.asset_number,CONCAT_WS(' ',b.name,m.name,v.name),x.serial_number,x.barcode,x.description,
  x.acquisition_date,x.acquisition_cost,DATE_ADD(x.acquisition_date,INTERVAL 3 YEAR),'Dummy asset migration MK5 office'
FROM tmp_mk5_assets x
JOIN asset_locations l ON l.branch_id=@branch_id AND l.code=x.location_code
JOIN asset_types t ON t.code=x.type_code
JOIN asset_categories c ON c.id=t.asset_category_id AND c.group_id=@group_id
JOIN asset_brands b ON b.group_id=@group_id AND b.code=x.brand_code
JOIN asset_models m ON m.asset_type_id=t.id AND m.asset_brand_id=b.id AND m.code=x.model_code
JOIN asset_variants v ON v.asset_model_id=m.id AND v.code=x.variant_code
ON DUPLICATE KEY UPDATE
  asset_location_id=VALUES(asset_location_id),asset_type_id=VALUES(asset_type_id),asset_brand_id=VALUES(asset_brand_id),
  asset_model_id=VALUES(asset_model_id),asset_variant_id=VALUES(asset_variant_id),condition_id=VALUES(condition_id),
  status_id=VALUES(status_id),name=VALUES(name),serial_number=VALUES(serial_number),barcode=VALUES(barcode),
  description=VALUES(description),acquisition_date=VALUES(acquisition_date),acquisition_cost=VALUES(acquisition_cost),
  warranty_until=VALUES(warranty_until),notes=VALUES(notes);

DROP TEMPORARY TABLE tmp_mk5_assets;
COMMIT;

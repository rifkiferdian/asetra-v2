-- Dummy Sprint 1: Struktur Organisasi Manna Kampus.
-- Idempotent: aman dijalankan ulang pada database development.
START TRANSACTION;

INSERT INTO company_groups(code,name,description,address,phone,email,is_active)
VALUES('MANNA','Manna Kampus','Grup usaha retail Manna Kampus',NULL,NULL,NULL,1)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

SET @group_id = (SELECT id FROM company_groups WHERE code='MANNA' LIMIT 1);

INSERT INTO companies(group_id,code,name,legal_name,tax_number,address,phone,email,is_active)
VALUES(@group_id,'MK','Manna Kampus',NULL,NULL,NULL,NULL,NULL,1)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),is_active=VALUES(is_active);

SET @company_id = (SELECT id FROM companies WHERE group_id=@group_id AND code='MK' LIMIT 1);

INSERT INTO branches(company_id,code,name,branch_type,address,city,province,postal_code,phone,is_active) VALUES
(@company_id,'STORE-01','Store 1','STORE',NULL,NULL,NULL,NULL,NULL,1),
(@company_id,'STORE-02','Store 2','STORE',NULL,NULL,NULL,NULL,NULL,1)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),branch_type=VALUES(branch_type),is_active=VALUES(is_active);

INSERT INTO departments(group_id,code,name,description,parent_id,is_active) VALUES
(@group_id,'OPS','Operasional Store','Mengelola kegiatan operasional harian Store',NULL,1),
(@group_id,'MD','Merchandising','Mengelola kategori, produk, harga, dan display',NULL,1),
(@group_id,'PUR','Purchasing','Mengelola pembelian barang dan hubungan pemasok',NULL,1),
(@group_id,'WH','Warehouse & Inventory','Mengelola gudang, stok, dan distribusi barang',NULL,1),
(@group_id,'FIN','Finance & Accounting','Mengelola keuangan, kas, pajak, dan akuntansi',NULL,1),
(@group_id,'HR','Human Resources','Mengelola sumber daya manusia dan administrasi karyawan',NULL,1),
(@group_id,'IT','Information Technology','Mengelola sistem, perangkat, jaringan, dan dukungan teknologi',NULL,1),
(@group_id,'MKT','Marketing','Mengelola promosi, komunikasi, dan program pelanggan',NULL,1),
(@group_id,'GA','General Affairs','Mengelola fasilitas, perizinan, dan kebutuhan umum',NULL,1),
(@group_id,'SEC','Security & Loss Prevention','Mengelola keamanan dan pencegahan kehilangan',NULL,1)
ON DUPLICATE KEY UPDATE
  name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

COMMIT;

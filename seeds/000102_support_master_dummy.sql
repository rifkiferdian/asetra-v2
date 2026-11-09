-- Dummy Master Operasional & Accounting Sprint 3.
-- Masing-masing 3 data per tabel. Idempotent untuk database development.
START TRANSACTION;

SET @group_id = COALESCE(
  (SELECT id FROM company_groups WHERE code='MANNA' LIMIT 1),
  (SELECT id FROM company_groups WHERE is_active=1 ORDER BY id LIMIT 1)
);
SET @company_id = (SELECT id FROM companies WHERE group_id=@group_id AND is_active=1 ORDER BY id LIMIT 1);

INSERT INTO asset_conditions(code,name,description,is_active) VALUES
('GOOD','Baik','Aset berfungsi normal dan layak digunakan',1),
('MINOR-DAMAGE','Rusak Ringan','Aset masih dapat digunakan dengan perbaikan ringan',1),
('MAJOR-DAMAGE','Rusak Berat','Aset memerlukan perbaikan besar atau penggantian',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO asset_statuses(code,name,description,is_active) VALUES
('AVAILABLE','Tersedia','Aset tersedia dan siap digunakan',1),
('IN-USE','Digunakan','Aset sedang digunakan atau ditugaskan',1),
('MAINTENANCE','Perawatan','Aset sedang dalam proses perawatan',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO vendors(company_group_id,code,name,vendor_type,tax_number,address,phone,email,contact_person,is_active) VALUES
(@group_id,'VND-TECH','PT Teknologi Nusantara','SUPPLIER','01.111.111.1-001.000','Jl. Teknologi No. 10, Yogyakarta','0274000001','sales@teknologi.example.invalid','Budi Santoso',1),
(@group_id,'VND-SERVICE','CV Mitra Servis Utama','SERVICE','02.222.222.2-002.000','Jl. Perawatan No. 20, Yogyakarta','0274000002','service@mitra.example.invalid','Sari Wulandari',1),
(@group_id,'VND-GENERAL','PT Sumber Asset Indonesia','BOTH','03.333.333.3-003.000','Jl. Industri No. 30, Yogyakarta','0274000003','contact@sumberasset.example.invalid','Andi Pratama',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),vendor_type=VALUES(vendor_type),tax_number=VALUES(tax_number),address=VALUES(address),phone=VALUES(phone),email=VALUES(email),contact_person=VALUES(contact_person),is_active=VALUES(is_active);

INSERT INTO depreciation_methods(code,name,description,is_active) VALUES
('SL','Garis Lurus','Penyusutan dengan nilai beban tetap setiap periode',1),
('DB','Saldo Menurun','Penyusutan berdasarkan persentase nilai buku',1),
('UOP','Unit Produksi','Penyusutan berdasarkan jumlah unit produksi atau pemakaian',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),description=VALUES(description),is_active=VALUES(is_active);

INSERT INTO gl_accounts(company_id,parent_gl_account_id,code,name,account_type,external_account_id,is_active) VALUES
(@company_id,NULL,'110100','Aset Tetap','ASSET','EXT-ASSET-110100',1),
(@company_id,NULL,'119100','Akumulasi Penyusutan Aset','CONTRA_ASSET','EXT-ACCUM-119100',1),
(@company_id,NULL,'610100','Beban Penyusutan Aset','EXPENSE','EXT-EXP-610100',1)
ON DUPLICATE KEY UPDATE name=VALUES(name),account_type=VALUES(account_type),external_account_id=VALUES(external_account_id),is_active=VALUES(is_active);

SET @method_sl = (SELECT id FROM depreciation_methods WHERE code='SL' LIMIT 1);
SET @type_1 = (SELECT t.id FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE c.group_id=@group_id ORDER BY t.id LIMIT 1 OFFSET 0);
SET @type_2 = (SELECT t.id FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE c.group_id=@group_id ORDER BY t.id LIMIT 1 OFFSET 1);
SET @type_3 = (SELECT t.id FROM asset_types t JOIN asset_categories c ON c.id=t.asset_category_id WHERE c.group_id=@group_id ORDER BY t.id LIMIT 1 OFFSET 2);
SET @rate_4y = (SELECT id FROM depreciation_rates WHERE useful_life_months=48 LIMIT 1);
SET @rate_5y = (SELECT id FROM depreciation_rates WHERE useful_life_months=60 LIMIT 1);
SET @rate_8y = (SELECT id FROM depreciation_rates WHERE useful_life_months=96 LIMIT 1);

INSERT INTO asset_depreciation_policies(company_group_id,company_id,asset_type_id,depreciation_method_id,depreciation_rate_id,book_type,residual_value_percent,effective_from,effective_to,is_active)
SELECT @group_id,NULL,@type_1,@method_sl,@rate_4y,'ACCOUNTING',5.00,'2026-01-01',NULL,1 FROM DUAL
WHERE @type_1 IS NOT NULL AND NOT EXISTS(SELECT 1 FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_1 AND book_type='ACCOUNTING' AND effective_from='2026-01-01');
INSERT INTO asset_depreciation_policies(company_group_id,company_id,asset_type_id,depreciation_method_id,depreciation_rate_id,book_type,residual_value_percent,effective_from,effective_to,is_active)
SELECT @group_id,NULL,@type_2,@method_sl,@rate_5y,'ACCOUNTING',5.00,'2026-01-01',NULL,1 FROM DUAL
WHERE @type_2 IS NOT NULL AND NOT EXISTS(SELECT 1 FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_2 AND book_type='ACCOUNTING' AND effective_from='2026-01-01');
INSERT INTO asset_depreciation_policies(company_group_id,company_id,asset_type_id,depreciation_method_id,depreciation_rate_id,book_type,residual_value_percent,effective_from,effective_to,is_active)
SELECT @group_id,NULL,@type_3,@method_sl,@rate_8y,'ACCOUNTING',10.00,'2026-01-01',NULL,1 FROM DUAL
WHERE @type_3 IS NOT NULL AND NOT EXISTS(SELECT 1 FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_3 AND book_type='ACCOUNTING' AND effective_from='2026-01-01');

SET @gl_asset = (SELECT id FROM gl_accounts WHERE company_id=@company_id AND code='110100' LIMIT 1);
SET @gl_accumulated = (SELECT id FROM gl_accounts WHERE company_id=@company_id AND code='119100' LIMIT 1);
SET @gl_expense = (SELECT id FROM gl_accounts WHERE company_id=@company_id AND code='610100' LIMIT 1);
SET @mapping_policy_1 = (SELECT id FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_1 AND book_type='ACCOUNTING' ORDER BY effective_from DESC LIMIT 1);
SET @mapping_policy_2 = (SELECT id FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_2 AND book_type='ACCOUNTING' ORDER BY effective_from DESC LIMIT 1);
SET @mapping_policy_3 = (SELECT id FROM asset_depreciation_policies WHERE company_group_id=@group_id AND company_id IS NULL AND asset_type_id=@type_3 AND book_type='ACCOUNTING' ORDER BY effective_from DESC LIMIT 1);

INSERT INTO asset_accounting_mappings(company_id,depreciation_policy_id,asset_gl_account_id,accumulated_depreciation_gl_account_id,depreciation_expense_gl_account_id,is_active) VALUES
(@company_id,@mapping_policy_1,@gl_asset,@gl_accumulated,@gl_expense,1),
(@company_id,@mapping_policy_2,@gl_asset,@gl_accumulated,@gl_expense,1),
(@company_id,@mapping_policy_3,@gl_asset,@gl_accumulated,@gl_expense,1)
ON DUPLICATE KEY UPDATE asset_gl_account_id=VALUES(asset_gl_account_id),accumulated_depreciation_gl_account_id=VALUES(accumulated_depreciation_gl_account_id),depreciation_expense_gl_account_id=VALUES(depreciation_expense_gl_account_id),is_active=VALUES(is_active);

COMMIT;

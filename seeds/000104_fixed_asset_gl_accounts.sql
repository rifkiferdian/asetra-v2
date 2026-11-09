-- Akun GL aset tetap untuk seluruh Company aktif.
START TRANSACTION;

INSERT INTO gl_accounts(company_id,parent_gl_account_id,code,name,account_type,external_account_id,is_active)
SELECT c.id,NULL,accounts.code,accounts.name,'ASSET',NULL,1
FROM companies c
CROSS JOIN (
  SELECT '110001' code,'Tanah' name
  UNION ALL SELECT '110002','Gedung'
  UNION ALL SELECT '110003','Emplasement'
  UNION ALL SELECT '110004','Kendaraan'
  UNION ALL SELECT '110005','Mesin & Perlengkapan'
  UNION ALL SELECT '110006','Peralatan Kantor'
  UNION ALL SELECT '110007','Peralatan Toko'
) accounts
WHERE c.code IN ('NAYAN','GODEAN')
ON DUPLICATE KEY UPDATE
  name=VALUES(name),account_type=VALUES(account_type),is_active=VALUES(is_active);

COMMIT;

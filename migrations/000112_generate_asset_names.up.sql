ALTER TABLE assets MODIFY COLUMN name VARCHAR(650) NOT NULL;

UPDATE assets a
JOIN asset_types t ON t.id=a.asset_type_id
LEFT JOIN asset_brands b ON b.id=a.asset_brand_id
LEFT JOIN asset_models m ON m.id=a.asset_model_id
LEFT JOIN asset_variants v ON v.id=a.asset_variant_id
SET a.name=COALESCE(NULLIF(CONCAT_WS(' ',NULLIF(b.name,''),NULLIF(m.name,''),NULLIF(v.name,'')),''),t.name);

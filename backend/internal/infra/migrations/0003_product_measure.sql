ALTER TABLE products ADD COLUMN measure_unit TEXT NOT NULL DEFAULT 'unidade';
ALTER TABLE products ADD COLUMN package_amount REAL NOT NULL DEFAULT 1;

UPDATE products
SET measure_unit = CASE
  WHEN LOWER(TRIM(COALESCE(unit, ''))) IN ('kg', 'quilo', 'quilos') THEN 'kg'
  WHEN LOWER(TRIM(COALESCE(unit, ''))) = 'g' THEN 'g'
  WHEN LOWER(TRIM(COALESCE(unit, ''))) IN ('l', 'lt', 'litro', 'litros') THEN 'l'
  WHEN LOWER(TRIM(COALESCE(unit, ''))) = 'ml' THEN 'ml'
  ELSE 'unidade'
END
WHERE COALESCE(measure_unit, '') = '' OR measure_unit = 'unidade';

UPDATE products
SET package_amount = 1
WHERE package_amount <= 0;

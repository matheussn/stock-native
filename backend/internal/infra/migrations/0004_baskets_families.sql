CREATE TABLE IF NOT EXISTS basket_templates (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  active INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS basket_template_items (
  id TEXT PRIMARY KEY,
  basket_template_id TEXT NOT NULL,
  product_id TEXT NOT NULL,
  quantity_per_basket REAL NOT NULL CHECK(quantity_per_basket > 0),
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  FOREIGN KEY(basket_template_id) REFERENCES basket_templates(id),
  FOREIGN KEY(product_id) REFERENCES products(id),
  UNIQUE(basket_template_id, product_id)
);

CREATE TABLE IF NOT EXISTS families (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  cestas_per_period REAL NOT NULL CHECK(cestas_per_period > 0),
  active INTEGER NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_basket_template_items_template_product ON basket_template_items(basket_template_id, product_id);
CREATE INDEX IF NOT EXISTS idx_families_active ON families(active);

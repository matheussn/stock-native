CREATE INDEX IF NOT EXISTS idx_movements_product_date ON movements(product_id, movement_date);
CREATE INDEX IF NOT EXISTS idx_movements_type_date ON movements(type, movement_date);
CREATE INDEX IF NOT EXISTS idx_movements_origin ON movements(origin_id);
CREATE INDEX IF NOT EXISTS idx_movements_destination ON movements(destination_id);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
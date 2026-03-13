CREATE INDEX IF NOT EXISTS idx_product_variation_product_active_base
    ON product_variation(product_id, is_active, base_quantity DESC, id ASC);

CREATE INDEX IF NOT EXISTS idx_family_assistential_work_active
    ON family(assistential_work_id, is_active, id ASC);

CREATE INDEX IF NOT EXISTS idx_family_group_assignment_family_started
    ON family_group_assignment(family_id, started_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_family_group_assignment_active_family
    ON family_group_assignment(family_id)
    WHERE ended_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_movement_created
    ON movement(created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_movement_type_work_created
    ON movement(type, assistential_work_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_movement_work_created
    ON movement(assistential_work_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_movement_product_item_movement
    ON movement_product_item(movement_id);

CREATE INDEX IF NOT EXISTS idx_movement_group_item_movement
    ON movement_group_item(movement_id);

CREATE INDEX IF NOT EXISTS idx_movement_group_item_resolution_group_item
    ON movement_group_item_resolution(movement_group_item_id);

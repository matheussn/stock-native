CREATE TABLE assistential_work (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE product (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    base_unit   TEXT NOT NULL,
    description TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE product_variation (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id    INTEGER NOT NULL REFERENCES product(id),
    description   TEXT,
    base_quantity INTEGER NOT NULL,
    current_stock INTEGER NOT NULL DEFAULT 0,
    is_active     INTEGER NOT NULL DEFAULT 1,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE product_group (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE product_group_item (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    product_group_id INTEGER NOT NULL REFERENCES product_group(id),
    product_id       INTEGER NOT NULL REFERENCES product(id),
    base_quantity    INTEGER NOT NULL,
    UNIQUE (product_group_id, product_id)
);

CREATE TABLE family (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    assistential_work_id INTEGER NOT NULL REFERENCES assistential_work(id),
    name                 TEXT NOT NULL,
    member_count         INTEGER NOT NULL DEFAULT 1,
    address              TEXT,
    contact              TEXT,
    is_active            INTEGER NOT NULL DEFAULT 1,
    created_at           TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE family_group_assignment (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    family_id        INTEGER NOT NULL REFERENCES family(id),
    product_group_id INTEGER NOT NULL REFERENCES product_group(id),
    started_at       TEXT NOT NULL DEFAULT (datetime('now')),
    ended_at         TEXT
);

CREATE TABLE movement (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    assistential_work_id INTEGER NOT NULL REFERENCES assistential_work(id),
    type                 TEXT NOT NULL CHECK (type IN ('in', 'out')),
    notes                TEXT,
    created_at           TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE movement_product_item (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    movement_id          INTEGER NOT NULL REFERENCES movement(id),
    product_variation_id INTEGER NOT NULL REFERENCES product_variation(id),
    quantity             INTEGER NOT NULL CHECK (quantity > 0)
);

CREATE TABLE movement_group_item (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    movement_id      INTEGER NOT NULL REFERENCES movement(id),
    product_group_id INTEGER NOT NULL REFERENCES product_group(id),
    quantity         INTEGER NOT NULL CHECK (quantity > 0)
);

CREATE TABLE movement_group_item_resolution (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    movement_group_item_id INTEGER NOT NULL REFERENCES movement_group_item(id),
    product_variation_id   INTEGER NOT NULL REFERENCES product_variation(id),
    quantity               INTEGER NOT NULL CHECK (quantity > 0)
);

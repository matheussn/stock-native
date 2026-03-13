# Spiritual Center Ã¢â‚¬â€ Inventory System
## Database Documentation

> **Engine:** SQLite  
> **Convention:** All tables and columns in English. Booleans as `INTEGER` (1/0). Dates as `TEXT` in `datetime('now')` format.

---

## Overview

This database supports an inventory system for a spiritual center. It handles product registration with variations, product groups (baskets), assistential works, family management, and stock movements (in/out).

The stock balance lives in `product_variation.current_stock` and must be updated by the **application layer** Ã¢â‚¬â€ never via triggers. Every stock-affecting operation must run inside a **transaction** to guarantee consistency.

## Migration Policy

- Migration `v1` (`initial_schema`) is the frozen baseline for this project.
- After the first release, **never edit an applied migration**.
- Any schema change must be added as a **new versioned migration** (`v2`, `v3`, ...).
- Migration files live in `internal/db/migrations/*.sql` and are embedded into the application binary.
- The application validates frozen migration checksums on startup to catch accidental edits to the baseline.

## Backup and Restore Policy

- Backups are created as standalone SQLite `.db` files outside the active database path.
- Restore always replaces the active `stock.db` with the selected backup file.
- Before restore, the current active database is preserved as `stock.before-restore-YYYYMMDD-HHMMSS.db`.
- Restored database files are validated and migrated before going back into active use.

## Connection Policy

- The application opens SQLite with `foreign_keys = ON`.
- The active database runs with `journal_mode = WAL`.
- The connection sets `busy_timeout = 5000` to reduce transient lock failures.
- The connection sets `synchronous = NORMAL` as the durability/performance baseline for local desktop usage.

## Runtime Indexes

Migration `v2` (`add_indexes`) adds the secondary indexes used by the current application queries:

- `idx_product_variation_product_active_base` on `product_variation(product_id, is_active, base_quantity DESC, id ASC)`
- `idx_family_assistential_work_active` on `family(assistential_work_id, is_active, id ASC)`
- `idx_family_group_assignment_family_started` on `family_group_assignment(family_id, started_at DESC, id DESC)`
- `idx_family_group_assignment_active_family` on `family_group_assignment(family_id) WHERE ended_at IS NULL`
- `idx_movement_created` on `movement(created_at DESC, id DESC)`
- `idx_movement_type_work_created` on `movement(type, assistential_work_id, created_at DESC, id DESC)`
- `idx_movement_work_created` on `movement(assistential_work_id, created_at DESC, id DESC)`
- `idx_movement_product_item_movement` on `movement_product_item(movement_id)`
- `idx_movement_group_item_movement` on `movement_group_item(movement_id)`
- `idx_movement_group_item_resolution_group_item` on `movement_group_item_resolution(movement_group_item_id)`

---

## Entity Relationship Summary

```
assistential_work
    Ã¢â€Å“Ã¢â€â‚¬Ã¢â€â‚¬ movement (1:N)
    Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ family (1:N)
            Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ family_group_assignment (1:N)
                        Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_group

product
    Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_variation (1:N)  Ã¢â€ Â actual stock lives here

product_group
    Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_group_item (1:N)
                Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product (generic reference, no variation)

movement
    Ã¢â€Å“Ã¢â€â‚¬Ã¢â€â‚¬ movement_product_item (1:N)
    Ã¢â€â€š           Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_variation
    Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ movement_group_item (1:N)
                Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_group
                Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ movement_group_item_resolution (1:N)
                            Ã¢â€â€Ã¢â€â‚¬Ã¢â€â‚¬ product_variation
```

---

## Tables

### `assistential_work`
Represents the assistential works carried out by the spiritual center (ex: Fraternity, Breakfast Program).

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `name` | TEXT NOT NULL | Work name |
| `description` | TEXT | Optional description |
| `is_active` | INTEGER | 1 = active, 0 = inactive. Default: 1 |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

---

### `product`
Generic product registration, without variations (ex: Rice, Beans, Milk).

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `name` | TEXT NOT NULL | Product name (ex: "Rice") |
| `base_unit` | TEXT NOT NULL | Unit used across all variations: `g`, `ml`, or `un` |
| `description` | TEXT | Optional description |
| `is_active` | INTEGER | 1 = active, 0 = inactive. Default: 1 |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

> **Note:** `base_unit` defines the unit for `base_quantity` in both `product_variation` and `product_group_item`. All quantities for the same product must use the same unit.

---

### `product_variation`
Specific variations of a product (ex: Rice 1kg, Rice 5kg). **This is where actual stock is tracked.**

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `product_id` | INTEGER FK | References `product(id)` |
| `description` | TEXT | Optional variation label (ex: "Rice 5kg") |
| `base_quantity` | INTEGER NOT NULL | Quantity in `base_unit` (ex: 5000 for 5kg in grams) |
| `current_stock` | INTEGER NOT NULL | Current unit count of this variation. Default: 0 |
| `is_active` | INTEGER | 1 = active, 0 = inactive. Default: 1 |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

> **Important:** `current_stock` must only be updated by the application inside a transaction alongside the movement registration.

---

### `product_group`
Product groups, commonly known as baskets (ex: Small Basket, Large Basket).

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `name` | TEXT NOT NULL | Group name (ex: "Small Basket") |
| `description` | TEXT | Optional description |
| `is_active` | INTEGER | 1 = active, 0 = inactive. Default: 1 |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

---

### `product_group_item`
Composition of a product group. References the **generic product** (not a variation), with a quantity in the product's `base_unit`.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `product_group_id` | INTEGER FK | References `product_group(id)` |
| `product_id` | INTEGER FK | References `product(id)` Ã¢â‚¬â€ generic, no variation |
| `base_quantity` | INTEGER NOT NULL | Amount needed in `base_unit` (ex: 3000 = 3kg of rice) |

> **Constraint:** `UNIQUE (product_group_id, product_id)` Ã¢â‚¬â€ the same product cannot appear twice in the same group. Adjust the quantity via `base_quantity` instead.

> **Relationship:** Many-to-many between `product_group` and `product`. The same product can belong to multiple groups, and the same group can contain multiple products.

---

### `family`
Families assisted by the center. Each family belongs to exactly one assistential work, which never changes.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `assistential_work_id` | INTEGER FK | References `assistential_work(id)` |
| `name` | TEXT NOT NULL | Family identification name |
| `member_count` | INTEGER NOT NULL | Number of family members. Default: 1 |
| `address` | TEXT | Optional address |
| `contact` | TEXT | Optional contact info |
| `is_active` | INTEGER | 1 = active, 0 = inactive. Default: 1 |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

---

### `family_group_assignment`
Historical record of which product group (basket) a family receives. A new record is inserted on every change. Multiple active records (`ended_at IS NULL`) mean the family currently receives more than one basket.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `family_id` | INTEGER FK | References `family(id)` |
| `product_group_id` | INTEGER FK | References `product_group(id)` |
| `started_at` | TEXT NOT NULL | Assignment start date. Default: `datetime('now')` |
| `ended_at` | TEXT | Assignment end date. NULL = currently active |

> **Query pattern Ã¢â‚¬â€ current assignment:** `WHERE ended_at IS NULL`  
> **Query pattern Ã¢â‚¬â€ assignment at a given date:** `WHERE started_at <= :date AND (ended_at IS NULL OR ended_at > :date)`

---

### `movement`
Header of a stock movement. Each record represents a single in or out operation, linked to an assistential work.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `assistential_work_id` | INTEGER FK | References `assistential_work(id)` |
| `type` | TEXT NOT NULL | `'in'` for stock entry, `'out'` for stock exit |
| `notes` | TEXT | Optional notes |
| `created_at` | TEXT | Creation timestamp. Default: `datetime('now')` |

> **Constraint:** `CHECK (type IN ('in', 'out'))`

---

### `movement_product_item`
Line items of a movement for direct product variations (ex: 3 units of Milk 1L). A single movement can have multiple items.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `movement_id` | INTEGER FK | References `movement(id)` |
| `product_variation_id` | INTEGER FK | References `product_variation(id)` |
| `quantity` | INTEGER NOT NULL | Number of units. Must be > 0 |

---

### `movement_group_item`
Line items of a movement for product groups/baskets (ex: 4 Small Baskets). A single movement can mix product items and group items freely.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `movement_id` | INTEGER FK | References `movement(id)` |
| `product_group_id` | INTEGER FK | References `product_group(id)` |
| `quantity` | INTEGER NOT NULL | Number of baskets. Must be > 0 |

---

### `movement_group_item_resolution`
Records how each basket movement was resolved into actual physical stock variations. Provides full traceability between a basket movement and the specific product variations consumed or received.

| Column | Type | Description |
|---|---|---|
| `id` | INTEGER PK | Auto-incremented identifier |
| `movement_group_item_id` | INTEGER FK | References `movement_group_item(id)` |
| `product_variation_id` | INTEGER FK | References `product_variation(id)` |
| `quantity` | INTEGER NOT NULL | Number of units of this variation used. Must be > 0 |

> **Example:** A movement of 4 Small Baskets, each requiring 3kg of rice (12kg total), might resolve to: 2Ãƒâ€” Rice 5kg + 2Ãƒâ€” Rice 1kg. Two resolution records are inserted for that `movement_group_item`.

> **Traceability query:** To find which physical variations were consumed by a basket movement, join `movement_group_item` Ã¢â€ â€™ `movement_group_item_resolution` Ã¢â€ â€™ `product_variation`.

---

## Key Business Rules

**Stock update flow**
The application must follow this sequence inside a single transaction when registering a movement:
1. Insert the `movement` record
2. Insert all `movement_product_item` records
3. Insert all `movement_group_item` records
4. For each group item, resolve which `product_variation`s cover the required quantities and insert into `movement_group_item_resolution`
5. Update `current_stock` on all affected `product_variation` records (increment for `'in'`, decrement for `'out'`)
6. Commit Ã¢â‚¬â€ or rollback everything on any failure

**Basket resolution logic**
When a basket is moved, the application must:
- Read the basket's `product_group_item` list
- Multiply each `base_quantity` by the number of baskets
- Decide which `product_variation`s of each generic product to use (e.g. largest variation first)
- Validate that stock is sufficient before committing

**Monthly demand projection**
To check whether current stock covers the month's demand:
- Fetch all active `family_group_assignment` records (`ended_at IS NULL`)
- Expand each assignment into its basket's `product_group_item` list
- Aggregate total required `base_quantity` per generic `product`
- Distribute across available `product_variation`s and compare against `current_stock`

---

## Full SQL Script

```sql
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
```

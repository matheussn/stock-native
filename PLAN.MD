# Spiritual Center — Inventory System
## Development Planning

> **Stack:** Go + Wails + SQLite (file-based)  
> **Runtime:** Fully local, no internet required  
> **Architecture:** Wails exposes Go backend functions directly to the frontend via bindings. The frontend calls Go functions, Go handles all business logic and database access.

---

## Project Structure

```
/
├── main.go                  # Wails app entry point
├── app.go                   # Wails context and app lifecycle
├── wails.json               # Wails config
│
├── internal/
│   ├── db/
│   │   ├── connection.go    # SQLite connection and setup
│   │   └── migrations.go    # Schema creation and versioning
│   │
│   ├── models/              # Structs matching database tables
│   │   ├── assistential_work.go
│   │   ├── product.go
│   │   ├── product_variation.go
│   │   ├── product_group.go
│   │   ├── product_group_item.go
│   │   ├── family.go
│   │   ├── family_group_assignment.go
│   │   └── movement.go
│   │
│   └── services/            # Business logic layer
│       ├── assistential_work_service.go
│       ├── product_service.go
│       ├── product_group_service.go
│       ├── family_service.go
│       ├── movement_service.go
│       └── stock_service.go
│
└── frontend/
    ├── src/
    │   ├── pages/
    │   ├── components/
    │   └── wailsjs/         # Auto-generated Go bindings
    └── ...
```

---

## Phases

---

### Phase 1 — Foundation

**Goal:** Project up and running with database ready.

- [ ] Initialize Wails project
- [ ] Configure SQLite dependency (`modernc.org/sqlite` — pure Go, no CGO required)
- [ ] Implement database connection with file path pointing to local user data directory
- [ ] Implement migrations system — run schema creation on first launch, versioned for future changes
- [ ] Create all models (structs) matching the database tables
- [ ] Validate that Wails window opens and connects to the database correctly

**Deliverable:** App launches, creates the `.db` file on first run, and applies the full schema.

---

### Phase 2 — Assistential Works

**Goal:** First fully working CRUD, serves as the pattern for all others.

- [ ] `AssistentialWorkService`: Create, List, Update, soft delete (toggle `is_active`)
- [ ] Expose service methods as Wails bindings
- [ ] Frontend: list screen with active/inactive filter
- [ ] Frontend: create and edit form
- [ ] Frontend: confirm before deactivating

**Deliverable:** User can manage assistential works end to end.

---

### Phase 3 — Products and Variations

**Goal:** Register generic products and their physical variations with stock tracking.

- [ ] `ProductService`: Create, List, Update, soft delete for `product`
- [ ] `ProductVariationService`: Create, List by product, Update, soft delete for `product_variation`
- [ ] Frontend: product list with expandable variations per product
- [ ] Frontend: form to create product and add variations inline
- [ ] Frontend: display `current_stock` for each variation

**Deliverable:** User can manage products and their variations. Stock levels are visible.

---

### Phase 4 — Product Groups (Baskets)

**Goal:** Register baskets composed of generic products with quantities.

- [ ] `ProductGroupService`: Create, List, Update, soft delete for `product_group`
- [ ] Manage `product_group_item` composition (add, update quantity, remove items)
- [ ] Frontend: basket list
- [ ] Frontend: basket detail screen showing composition (product name + quantity in readable unit)
- [ ] Frontend: form to create basket and manage its items

**Deliverable:** User can create and manage baskets and their composition.

---

### Phase 5 — Families

**Goal:** Register families and link them to their baskets.

- [ ] `FamilyService`: Create, List (with filter by assistential work), Update, soft delete
- [ ] Manage `family_group_assignment`: create new assignment, close current assignment (`ended_at`), list history
- [ ] Frontend: family list with filter by assistential work
- [ ] Frontend: family detail showing current basket assignment(s) and history
- [ ] Frontend: form to create family and assign initial basket
- [ ] Frontend: action to change basket (closes current, opens new assignment)

**Deliverable:** User can manage families and track their basket assignment history.

---

### Phase 6 — Stock Movements

**Goal:** Core of the system. Register in/out movements with product items and basket items in the same operation.

**Service — `MovementService`:**
- [ ] Begin transaction
- [ ] Insert `movement` header
- [ ] Insert `movement_product_item` lines
- [ ] Insert `movement_group_item` lines
- [ ] For each group item: resolve basket composition into variations (see resolution logic below)
- [ ] Insert `movement_group_item_resolution` records
- [ ] Update `current_stock` on all affected `product_variation` records
- [ ] Commit or rollback

**Basket resolution logic (application layer):**
- Read basket's `product_group_item` list
- Multiply each `base_quantity` by the number of baskets requested
- For each product, fetch available variations ordered by `base_quantity DESC` (largest first)
- Greedily fill the required quantity using available stock
- If stock is insufficient for any product, abort the entire transaction and return a descriptive error

**Frontend:**
- [ ] Movement creation screen with two sections: direct products and baskets
- [ ] Search/select for product variations with current stock shown
- [ ] Search/select for baskets with quantity input
- [ ] Stock validation feedback before confirming (show which products would be insufficient)
- [ ] Movement history list with filter by type (in/out) and by assistential work
- [ ] Movement detail showing all items and basket resolutions

**Deliverable:** User can register complete in/out movements. Stock is updated automatically and atomically.

---

### Phase 7 — Stock Dashboard

**Goal:** Give visibility into current stock health and monthly demand coverage.

- [ ] `StockService`:
  - [ ] List all variations with current stock
  - [ ] Monthly demand projection: aggregate all active `family_group_assignment` → expand basket composition → total required per product
  - [ ] Stock coverage check: compare projected demand against current stock and return shortfall per product/variation

- [ ] Frontend: stock overview screen
  - [ ] Cards or table showing each variation with stock level and status (ok / low / critical)
  - [ ] Monthly demand panel: shows whether current stock covers all active family assignments
  - [ ] Highlight products with insufficient stock for the month, showing the gap

**Deliverable:** User can see at a glance if the stock is healthy and if the month's demand can be met.

---

### Phase 8 — Polish and Distribution

**Goal:** Production-ready app for local distribution.

- [ ] App icon and Wails window configuration (title, size, min size, resizable)
- [ ] Global error handling — friendly messages for constraint violations, insufficient stock, etc.
- [ ] Confirmation dialogs for destructive actions (deactivate, delete)
- [ ] Empty states for all list screens
- [ ] Loading states for async operations
- [ ] Build pipeline: `wails build` generating a single executable for Windows
- [ ] Test full flow end to end on a clean machine (no dev tools installed)
- [ ] Document how to run/distribute the `.exe` and where the `.db` file is stored

**Deliverable:** Distributable executable ready to hand off to the spiritual center.

---

## Technical Decisions

| Topic | Decision |
|---|---|
| SQLite driver | `modernc.org/sqlite` — pure Go, no CGO, no external DLL needed |
| DB file location | OS user data directory (e.g. `%APPDATA%` on Windows) via `os.UserConfigDir()` |
| Migrations | Simple versioned script executed on startup — no ORM |
| Data access | Raw `database/sql` queries — no ORM, full control |
| Business logic | Service layer in Go, never in the frontend |
| Stock update | Application-managed inside transactions, never via triggers |
| Frontend | Wails default (React or Svelte — team's choice) |
| Build target | Windows primary (single `.exe`), macOS optional |

---

## Rules for the Agent

- **Never update `current_stock` outside of a movement transaction.** All stock changes go through `MovementService`.
- **Never expose raw SQL to the frontend.** The frontend only calls Go service methods via Wails bindings.
- **All destructive operations are soft deletes** via `is_active = 0`. No `DELETE` statements on core tables.
- **Basket resolution must validate stock before committing.** Return clear errors per product if stock is insufficient.
- **All service methods that write data must use transactions,** even if only a single table is affected — for consistency.
- **`base_unit` is the contract.** All `base_quantity` values across `product_variation`, `product_group_item`, and `family_group_assignment` use the unit defined in the parent `product`. Never mix units.

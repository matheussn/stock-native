package services

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"stock/internal/db"

	_ "modernc.org/sqlite"
)

func TestProductGroupServiceCRUDAndItems(t *testing.T) {
	conn := openProductGroupTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	productSvc := NewProductService(conn)
	groupSvc := NewProductGroupService(conn)
	ctx := context.Background()

	p1, err := productSvc.Create(ctx, "Arroz", "kg", "")
	if err != nil {
		t.Fatalf("create first product: %v", err)
	}
	p2, err := productSvc.Create(ctx, "Leite", "ml", "")
	if err != nil {
		t.Fatalf("create second product: %v", err)
	}

	group, err := groupSvc.Create(ctx, "Cesta Basica", "Modelo mensal")
	if err != nil {
		t.Fatalf("create product group: %v", err)
	}

	groups, err := groupSvc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active groups: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 active group, got %d", len(groups))
	}

	updated, err := groupSvc.Update(ctx, group.ID, "Cesta Solidaria", "Atualizada")
	if err != nil {
		t.Fatalf("update product group: %v", err)
	}
	if updated.Name != "Cesta Solidaria" {
		t.Fatalf("unexpected updated name: %q", updated.Name)
	}

	item1, err := groupSvc.UpsertItem(ctx, group.ID, p1.ID, 5)
	if err != nil {
		t.Fatalf("upsert first item: %v", err)
	}
	if item1.BaseQuantity != 5 {
		t.Fatalf("expected first item quantity to be 5, got %d", item1.BaseQuantity)
	}

	_, err = groupSvc.UpsertItem(ctx, group.ID, p2.ID, 1000)
	if err != nil {
		t.Fatalf("upsert second item: %v", err)
	}

	item1, err = groupSvc.UpsertItem(ctx, group.ID, p1.ID, 7)
	if err != nil {
		t.Fatalf("upsert first item quantity update: %v", err)
	}
	if item1.BaseQuantity != 7 {
		t.Fatalf("expected first item quantity to be 7 after update, got %d", item1.BaseQuantity)
	}

	items, err := groupSvc.ListItems(ctx, group.ID)
	if err != nil {
		t.Fatalf("list group items: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 group items, got %d", len(items))
	}

	if err := groupSvc.RemoveItem(ctx, group.ID, p2.ID); err != nil {
		t.Fatalf("remove second item: %v", err)
	}

	items, err = groupSvc.ListItems(ctx, group.ID)
	if err != nil {
		t.Fatalf("list group items after remove: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 group item after remove, got %d", len(items))
	}

	if err := groupSvc.SetActive(ctx, group.ID, false); err != nil {
		t.Fatalf("deactivate product group: %v", err)
	}

	groups, err = groupSvc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active groups after deactivate: %v", err)
	}
	if len(groups) != 0 {
		t.Fatalf("expected 0 active groups after deactivate, got %d", len(groups))
	}

	groups, err = groupSvc.List(ctx, true)
	if err != nil {
		t.Fatalf("list groups with inactive: %v", err)
	}
	if len(groups) != 1 || groups[0].IsActive != 0 {
		t.Fatalf("expected one inactive group after deactivate")
	}
}

func openProductGroupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	dsn := fmt.Sprintf("file:%s", dbPath)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open sqlite test db: %v", err)
	}

	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("enable foreign keys in test db: %v", err)
	}

	return conn
}

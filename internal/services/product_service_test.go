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

func TestProductServiceCRUD(t *testing.T) {
	conn := openProductTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	svc := NewProductService(conn)
	ctx := context.Background()

	first, err := svc.Create(ctx, "Arroz", "g", "Arroz branco")
	if err != nil {
		t.Fatalf("create first product: %v", err)
	}

	second, err := svc.Create(ctx, "Leite", "ml", "")
	if err != nil {
		t.Fatalf("create second product: %v", err)
	}

	activeOnly, err := svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active products: %v", err)
	}
	if len(activeOnly) != 2 {
		t.Fatalf("expected 2 active products, got %d", len(activeOnly))
	}

	updated, err := svc.Update(ctx, second.ID, "Leite Integral", "ml", "Caixa 1L")
	if err != nil {
		t.Fatalf("update product: %v", err)
	}
	if updated.Name != "Leite Integral" {
		t.Fatalf("expected updated name to be 'Leite Integral', got %q", updated.Name)
	}

	if err := svc.SetActive(ctx, first.ID, false); err != nil {
		t.Fatalf("deactivate product: %v", err)
	}

	activeOnly, err = svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active products after deactivate: %v", err)
	}
	if len(activeOnly) != 1 {
		t.Fatalf("expected 1 active product after deactivate, got %d", len(activeOnly))
	}

	withInactive, err := svc.List(ctx, true)
	if err != nil {
		t.Fatalf("list products with inactive: %v", err)
	}
	if len(withInactive) != 2 {
		t.Fatalf("expected 2 products including inactive, got %d", len(withInactive))
	}
	if withInactive[0].IsActive != 0 {
		t.Fatalf("expected first product to be inactive")
	}
}

func openProductTestDB(t *testing.T) *sql.DB {
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

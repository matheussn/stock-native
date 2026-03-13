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

func TestProductVariationServiceCRUD(t *testing.T) {
	conn := openProductVariationTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	ctx := context.Background()

	product, err := productSvc.Create(ctx, "Arroz", "g", "")
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	first, err := variationSvc.Create(ctx, product.ID, "", 1000)
	if err != nil {
		t.Fatalf("create first variation: %v", err)
	}
	if first.Description != "" {
		t.Fatalf("expected empty description, got %q", first.Description)
	}

	_, err = variationSvc.Create(ctx, product.ID, "Pacote 5kg", 5000)
	if err != nil {
		t.Fatalf("create second variation: %v", err)
	}

	activeOnly, err := variationSvc.ListByProduct(ctx, product.ID, false)
	if err != nil {
		t.Fatalf("list active variations: %v", err)
	}
	if len(activeOnly) != 2 {
		t.Fatalf("expected 2 active variations, got %d", len(activeOnly))
	}

	updated, err := variationSvc.Update(ctx, first.ID, "", 1000)
	if err != nil {
		t.Fatalf("update variation: %v", err)
	}
	if updated.Description != "" {
		t.Fatalf("expected empty updated description, got %q", updated.Description)
	}

	if err := variationSvc.SetActive(ctx, first.ID, false); err != nil {
		t.Fatalf("deactivate variation: %v", err)
	}

	activeOnly, err = variationSvc.ListByProduct(ctx, product.ID, false)
	if err != nil {
		t.Fatalf("list active variations after deactivate: %v", err)
	}
	if len(activeOnly) != 1 {
		t.Fatalf("expected 1 active variation after deactivate, got %d", len(activeOnly))
	}

	withInactive, err := variationSvc.ListByProduct(ctx, product.ID, true)
	if err != nil {
		t.Fatalf("list variations with inactive: %v", err)
	}
	if len(withInactive) != 2 {
		t.Fatalf("expected 2 variations including inactive, got %d", len(withInactive))
	}
	if withInactive[1].IsActive != 0 {
		t.Fatalf("expected smaller variation to be inactive")
	}
}

func openProductVariationTestDB(t *testing.T) *sql.DB {
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

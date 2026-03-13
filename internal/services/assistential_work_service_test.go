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

func TestAssistentialWorkServiceCRUD(t *testing.T) {
	conn := openTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	svc := NewAssistentialWorkService(conn)
	ctx := context.Background()

	first, err := svc.Create(ctx, "Fraternidade", "Atendimento semanal")
	if err != nil {
		t.Fatalf("create first assistential work: %v", err)
	}

	second, err := svc.Create(ctx, "Sopa Solidaria", "")
	if err != nil {
		t.Fatalf("create second assistential work: %v", err)
	}

	activeOnly, err := svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active assistential works: %v", err)
	}
	if len(activeOnly) != 2 {
		t.Fatalf("expected 2 active assistential works, got %d", len(activeOnly))
	}

	updated, err := svc.Update(ctx, second.ID, "Sopa Fraterna", "Distribuicao de sopa")
	if err != nil {
		t.Fatalf("update assistential work: %v", err)
	}
	if updated.Name != "Sopa Fraterna" {
		t.Fatalf("expected updated name to be 'Sopa Fraterna', got %q", updated.Name)
	}

	if err := svc.SetActive(ctx, first.ID, false); err != nil {
		t.Fatalf("deactivate assistential work: %v", err)
	}

	activeOnly, err = svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active assistential works after deactivate: %v", err)
	}
	if len(activeOnly) != 1 {
		t.Fatalf("expected 1 active assistential work after deactivate, got %d", len(activeOnly))
	}

	withInactive, err := svc.List(ctx, true)
	if err != nil {
		t.Fatalf("list assistential works with inactive: %v", err)
	}
	if len(withInactive) != 2 {
		t.Fatalf("expected 2 assistential works including inactive, got %d", len(withInactive))
	}
	if withInactive[0].IsActive != 0 {
		t.Fatalf("expected first assistential work to be inactive")
	}
}

func openTestDB(t *testing.T) *sql.DB {
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

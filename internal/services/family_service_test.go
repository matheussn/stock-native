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

func TestFamilyServiceCRUDAndFilter(t *testing.T) {
	conn := openFamilyTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	familySvc := NewFamilyService(conn)
	ctx := context.Background()

	workA, err := workSvc.Create(ctx, "Fraternidade", "")
	if err != nil {
		t.Fatalf("create assistential work A: %v", err)
	}
	workB, err := workSvc.Create(ctx, "Sopa", "")
	if err != nil {
		t.Fatalf("create assistential work B: %v", err)
	}

	first, err := familySvc.Create(ctx, workA.ID, "Familia Silva", 4, "Rua 1", "9999-1111")
	if err != nil {
		t.Fatalf("create first family: %v", err)
	}
	second, err := familySvc.Create(ctx, workB.ID, "Familia Souza", 3, "", "")
	if err != nil {
		t.Fatalf("create second family: %v", err)
	}

	filteredA, err := familySvc.List(ctx, false, workA.ID)
	if err != nil {
		t.Fatalf("list families by work A: %v", err)
	}
	if len(filteredA) != 1 || filteredA[0].ID != first.ID {
		t.Fatalf("expected 1 family for work A")
	}

	updated, err := familySvc.Update(ctx, second.ID, "Familia Souza Lima", 5, "Rua 2", "9999-2222")
	if err != nil {
		t.Fatalf("update second family: %v", err)
	}
	if updated.MemberCount != 5 {
		t.Fatalf("expected member_count = 5 after update")
	}

	if err := familySvc.SetActive(ctx, first.ID, false); err != nil {
		t.Fatalf("deactivate first family: %v", err)
	}

	activeOnly, err := familySvc.List(ctx, false, 0)
	if err != nil {
		t.Fatalf("list active families: %v", err)
	}
	if len(activeOnly) != 1 || activeOnly[0].ID != second.ID {
		t.Fatalf("expected only second family active")
	}

	withInactive, err := familySvc.List(ctx, true, 0)
	if err != nil {
		t.Fatalf("list families including inactive: %v", err)
	}
	if len(withInactive) != 2 {
		t.Fatalf("expected 2 families including inactive")
	}
}

func openFamilyTestDB(t *testing.T) *sql.DB {
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

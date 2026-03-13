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

func TestFamilyGroupAssignmentAssignAndHistory(t *testing.T) {
	conn := openFamilyAssignmentTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	familySvc := NewFamilyService(conn)
	groupSvc := NewProductGroupService(conn)
	assignSvc := NewFamilyGroupAssignmentService(conn)
	ctx := context.Background()

	work, err := workSvc.Create(ctx, "Fraternidade", "")
	if err != nil {
		t.Fatalf("create assistential work: %v", err)
	}
	family, err := familySvc.Create(ctx, work.ID, "Familia Alves", 3, "", "")
	if err != nil {
		t.Fatalf("create family: %v", err)
	}
	groupA, err := groupSvc.Create(ctx, "Cesta A", "")
	if err != nil {
		t.Fatalf("create group A: %v", err)
	}
	groupB, err := groupSvc.Create(ctx, "Cesta B", "")
	if err != nil {
		t.Fatalf("create group B: %v", err)
	}

	first, err := assignSvc.Assign(ctx, family.ID, groupA.ID)
	if err != nil {
		t.Fatalf("assign first group: %v", err)
	}
	if first.EndedAt != "" {
		t.Fatalf("expected first assignment to be active")
	}

	second, err := assignSvc.Assign(ctx, family.ID, groupB.ID)
	if err != nil {
		t.Fatalf("assign second group: %v", err)
	}

	assignments, err := assignSvc.ListByFamily(ctx, family.ID)
	if err != nil {
		t.Fatalf("list assignments: %v", err)
	}
	if len(assignments) != 2 {
		t.Fatalf("expected 2 assignments in history, got %d", len(assignments))
	}
	if assignments[0].ID != second.ID || assignments[0].EndedAt != "" {
		t.Fatalf("expected latest assignment active and first in list")
	}
	if assignments[1].ID != first.ID || assignments[1].EndedAt == "" {
		t.Fatalf("expected first assignment to be closed")
	}
}

func openFamilyAssignmentTestDB(t *testing.T) *sql.DB {
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

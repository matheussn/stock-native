package services

import (
	"context"
	"testing"

	"stock/internal/db"
)

func TestInstitutionServiceCRUD(t *testing.T) {
	conn := openTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	svc := NewInstitutionService(conn)
	ctx := context.Background()

	first, err := svc.Create(ctx, "Casa Esperanca", "Rua A", "", "Maria", "11999999999")
	if err != nil {
		t.Fatalf("create first institution: %v", err)
	}

	second, err := svc.Create(ctx, "Lar da Paz", "", "12345678000199", "", "")
	if err != nil {
		t.Fatalf("create second institution: %v", err)
	}

	activeOnly, err := svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active institutions: %v", err)
	}
	if len(activeOnly) != 2 {
		t.Fatalf("expected 2 active institutions, got %d", len(activeOnly))
	}

	updated, err := svc.Update(ctx, second.ID, "Lar da Paz", "Rua B", "12345678000199", "Joao", "11888888888")
	if err != nil {
		t.Fatalf("update institution: %v", err)
	}
	if updated.Address != "Rua B" {
		t.Fatalf("expected updated address to be Rua B, got %q", updated.Address)
	}
	if updated.ResponsibleName != "Joao" {
		t.Fatalf("expected updated responsible name to be Joao, got %q", updated.ResponsibleName)
	}

	if err := svc.SetActive(ctx, first.ID, false); err != nil {
		t.Fatalf("deactivate institution: %v", err)
	}

	activeOnly, err = svc.List(ctx, false)
	if err != nil {
		t.Fatalf("list active institutions after deactivate: %v", err)
	}
	if len(activeOnly) != 1 {
		t.Fatalf("expected 1 active institution after deactivate, got %d", len(activeOnly))
	}

	withInactive, err := svc.List(ctx, true)
	if err != nil {
		t.Fatalf("list institutions with inactive: %v", err)
	}
	if len(withInactive) != 2 {
		t.Fatalf("expected 2 institutions including inactive, got %d", len(withInactive))
	}
	if withInactive[0].IsActive != 0 {
		t.Fatalf("expected first institution to be inactive")
	}
}

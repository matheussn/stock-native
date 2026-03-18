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

func TestMovementServiceCreateWithProductsAndGroups(t *testing.T) {
	conn := openMovementTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	institutionSvc := NewInstitutionService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	groupSvc := NewProductGroupService(conn)
	movementSvc := NewMovementService(conn)
	ctx := context.Background()

	work, err := workSvc.Create(ctx, "Fraternidade", "")
	if err != nil {
		t.Fatalf("create work: %v", err)
	}
	institution, err := institutionSvc.Create(ctx, "Casa Esperanca", "", "", "", "")
	if err != nil {
		t.Fatalf("create institution: %v", err)
	}

	rice, err := productSvc.Create(ctx, "Arroz", "g", "")
	if err != nil {
		t.Fatalf("create rice product: %v", err)
	}
	beans, err := productSvc.Create(ctx, "Feijao", "g", "")
	if err != nil {
		t.Fatalf("create beans product: %v", err)
	}

	rice5kg, err := variationSvc.Create(ctx, rice.ID, "Pacote 5kg", 5000)
	if err != nil {
		t.Fatalf("create rice 5kg variation: %v", err)
	}
	rice1kg, err := variationSvc.Create(ctx, rice.ID, "Pacote 1kg", 1000)
	if err != nil {
		t.Fatalf("create rice 1kg variation: %v", err)
	}
	beans1kg, err := variationSvc.Create(ctx, beans.ID, "Pacote 1kg", 1000)
	if err != nil {
		t.Fatalf("create beans 1kg variation: %v", err)
	}

	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 2 WHERE id = ?`, rice5kg.ID); err != nil {
		t.Fatalf("seed stock for rice5kg: %v", err)
	}
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 3 WHERE id = ?`, rice1kg.ID); err != nil {
		t.Fatalf("seed stock for rice1kg: %v", err)
	}
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 5 WHERE id = ?`, beans1kg.ID); err != nil {
		t.Fatalf("seed stock for beans1kg: %v", err)
	}

	group, err := groupSvc.Create(ctx, "Cesta Basica", "")
	if err != nil {
		t.Fatalf("create group: %v", err)
	}
	if _, err := groupSvc.UpsertItem(ctx, group.ID, rice.ID, 3000); err != nil {
		t.Fatalf("upsert rice in group: %v", err)
	}
	if _, err := groupSvc.UpsertItem(ctx, group.ID, beans.ID, 2000); err != nil {
		t.Fatalf("upsert beans in group: %v", err)
	}

	movement, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		InstitutionID:      institution.ID,
		Type:               "out",
		Notes:              "Distribuicao mensal",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice1kg.ID, Quantity: 1},
		},
		GroupItems: []MovementGroupItemInput{
			{ProductGroupID: group.ID, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("create movement: %v", err)
	}
	if movement.ID <= 0 {
		t.Fatalf("expected valid movement id")
	}
	if movement.InstitutionID != institution.ID {
		t.Fatalf("expected institution id %d, got %d", institution.ID, movement.InstitutionID)
	}

	assertStock(t, conn, rice5kg.ID, 1)
	assertStock(t, conn, rice1kg.ID, 1)
	assertStock(t, conn, beans1kg.ID, 1)

	var storedInstitutionID sql.NullInt64
	if err := conn.QueryRow(`SELECT institution_id FROM movement WHERE id = ?`, movement.ID).Scan(&storedInstitutionID); err != nil {
		t.Fatalf("read movement institution: %v", err)
	}
	if !storedInstitutionID.Valid || storedInstitutionID.Int64 != institution.ID {
		t.Fatalf("expected stored institution id %d, got %+v", institution.ID, storedInstitutionID)
	}

	var resolutionCount int64
	if err := conn.QueryRow(`SELECT COUNT(*) FROM movement_group_item_resolution`).Scan(&resolutionCount); err != nil {
		t.Fatalf("count movement group resolutions: %v", err)
	}
	if resolutionCount != 3 {
		t.Fatalf("expected 3 group resolution records, got %d", resolutionCount)
	}
}

func TestMovementServiceInsufficientStockRollsBack(t *testing.T) {
	conn := openMovementTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	groupSvc := NewProductGroupService(conn)
	movementSvc := NewMovementService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Fraternidade", "")
	rice, _ := productSvc.Create(ctx, "Arroz", "g", "")
	rice1kg, _ := variationSvc.Create(ctx, rice.ID, "Pacote 1kg", 1000)
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 1 WHERE id = ?`, rice1kg.ID); err != nil {
		t.Fatalf("seed stock: %v", err)
	}

	group, _ := groupSvc.Create(ctx, "Cesta Basica", "")
	if _, err := groupSvc.UpsertItem(ctx, group.ID, rice.ID, 3000); err != nil {
		t.Fatalf("upsert rice in group: %v", err)
	}

	_, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		Type:               "out",
		GroupItems: []MovementGroupItemInput{
			{ProductGroupID: group.ID, Quantity: 2},
		},
	})
	if err == nil {
		t.Fatalf("expected insufficient stock error")
	}

	assertStock(t, conn, rice1kg.ID, 1)

	var movementCount int64
	if err := conn.QueryRow(`SELECT COUNT(*) FROM movement`).Scan(&movementCount); err != nil {
		t.Fatalf("count movement rows: %v", err)
	}
	if movementCount != 0 {
		t.Fatalf("expected movement rollback on error")
	}
}

func TestMovementServiceRespectsSharedStockAcrossMultipleGroups(t *testing.T) {
	conn := openMovementTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	groupSvc := NewProductGroupService(conn)
	movementSvc := NewMovementService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Visita Fraterna", "")

	rice, _ := productSvc.Create(ctx, "Arroz", "kg", "")
	beans, _ := productSvc.Create(ctx, "Feijao", "kg", "")

	rice2kg, _ := variationSvc.Create(ctx, rice.ID, "2Kg", 2)
	rice1kg, _ := variationSvc.Create(ctx, rice.ID, "1Kg", 1)
	beans1kg, _ := variationSvc.Create(ctx, beans.ID, "1Kg", 1)

	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 4 WHERE id = ?`, rice2kg.ID); err != nil {
		t.Fatalf("seed rice2kg stock: %v", err)
	}
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 5 WHERE id = ?`, rice1kg.ID); err != nil {
		t.Fatalf("seed rice1kg stock: %v", err)
	}
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 1 WHERE id = ?`, beans1kg.ID); err != nil {
		t.Fatalf("seed beans1kg stock: %v", err)
	}

	groupM, _ := groupSvc.Create(ctx, "Cesta M", "")
	groupG, _ := groupSvc.Create(ctx, "Cesta G", "")
	if _, err := groupSvc.UpsertItem(ctx, groupM.ID, rice.ID, 2); err != nil {
		t.Fatalf("upsert rice in group M: %v", err)
	}
	if _, err := groupSvc.UpsertItem(ctx, groupM.ID, beans.ID, 1); err != nil {
		t.Fatalf("upsert beans in group M: %v", err)
	}
	if _, err := groupSvc.UpsertItem(ctx, groupG.ID, rice.ID, 5); err != nil {
		t.Fatalf("upsert rice in group G: %v", err)
	}

	_, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		Type:               "out",
		GroupItems: []MovementGroupItemInput{
			{ProductGroupID: groupM.ID, Quantity: 1},
			{ProductGroupID: groupG.ID, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("movement should succeed with shared stock across groups: %v", err)
	}

	assertStock(t, conn, rice2kg.ID, 0)
	assertStock(t, conn, rice1kg.ID, 1)
	assertStock(t, conn, beans1kg.ID, 0)
}

func TestMovementServiceIgnoresInstitutionForEntry(t *testing.T) {
	conn := openMovementTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	institutionSvc := NewInstitutionService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	movementSvc := NewMovementService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Recepcao", "")
	institution, _ := institutionSvc.Create(ctx, "Lar da Luz", "", "", "", "")
	product, _ := productSvc.Create(ctx, "Arroz", "g", "")
	variation, _ := variationSvc.Create(ctx, product.ID, "Pacote 1kg", 1000)

	movement, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		InstitutionID:      institution.ID,
		Type:               "in",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: variation.ID, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("create entry movement: %v", err)
	}
	if movement.InstitutionID != 0 {
		t.Fatalf("expected entry movement to ignore institution, got %d", movement.InstitutionID)
	}

	var storedInstitutionID sql.NullInt64
	if err := conn.QueryRow(`SELECT institution_id FROM movement WHERE id = ?`, movement.ID).Scan(&storedInstitutionID); err != nil {
		t.Fatalf("read movement institution: %v", err)
	}
	if storedInstitutionID.Valid {
		t.Fatalf("expected null institution for entry movement, got %+v", storedInstitutionID)
	}
}

func TestMovementServiceRejectsInactiveOrMissingInstitutionOnExit(t *testing.T) {
	conn := openMovementTestDB(t)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	institutionSvc := NewInstitutionService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	movementSvc := NewMovementService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Fraternidade", "")
	institution, _ := institutionSvc.Create(ctx, "Casa de Apoio", "", "", "", "")
	product, _ := productSvc.Create(ctx, "Feijao", "g", "")
	variation, _ := variationSvc.Create(ctx, product.ID, "Pacote 1kg", 1000)
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 2 WHERE id = ?`, variation.ID); err != nil {
		t.Fatalf("seed stock: %v", err)
	}

	if err := institutionSvc.SetActive(ctx, institution.ID, false); err != nil {
		t.Fatalf("deactivate institution: %v", err)
	}

	_, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		InstitutionID:      institution.ID,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: variation.ID, Quantity: 1},
		},
	})
	if err == nil {
		t.Fatal("expected inactive institution error")
	}

	_, err = movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		InstitutionID:      9999,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: variation.ID, Quantity: 1},
		},
	})
	if err == nil {
		t.Fatal("expected missing institution error")
	}
}

func assertStock(t *testing.T, conn *sql.DB, variationID, expected int64) {
	t.Helper()
	var actual int64
	if err := conn.QueryRow(`SELECT current_stock FROM product_variation WHERE id = ?`, variationID).Scan(&actual); err != nil {
		t.Fatalf("read stock: %v", err)
	}
	if actual != expected {
		t.Fatalf("unexpected stock for variation %d: got %d, want %d", variationID, actual, expected)
	}
}

func openMovementTestDB(t *testing.T) *sql.DB {
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

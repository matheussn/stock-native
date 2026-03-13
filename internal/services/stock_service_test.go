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

func TestStockServiceProjectionAndCoverage(t *testing.T) {
	conn := openStockTestDB(t)
	t.Cleanup(func() { _ = conn.Close() })

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	familySvc := NewFamilyService(conn)
	groupSvc := NewProductGroupService(conn)
	assignSvc := NewFamilyGroupAssignmentService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	stockSvc := NewStockService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Fraternidade", "")
	family, _ := familySvc.Create(ctx, work.ID, "Familia 1", 3, "", "")

	rice, _ := productSvc.Create(ctx, "Arroz", "kg", "")
	beans, _ := productSvc.Create(ctx, "Feijao", "kg", "")

	riceVar, _ := variationSvc.Create(ctx, rice.ID, "Pacote 5kg", 5)
	beansVar, _ := variationSvc.Create(ctx, beans.ID, "Pacote 1kg", 1)
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 3 WHERE id = ?`, riceVar.ID); err != nil {
		t.Fatalf("seed rice stock: %v", err)
	}
	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 2 WHERE id = ?`, beansVar.ID); err != nil {
		t.Fatalf("seed beans stock: %v", err)
	}

	group, _ := groupSvc.Create(ctx, "Cesta A", "")
	groupSvc.UpsertItem(ctx, group.ID, rice.ID, 5)
	groupSvc.UpsertItem(ctx, group.ID, beans.ID, 3)
	assignSvc.Assign(ctx, family.ID, group.ID)

	status, err := stockSvc.ListVariationStatus(ctx, false)
	if err != nil {
		t.Fatalf("list variation status: %v", err)
	}
	if len(status) != 2 {
		t.Fatalf("expected 2 variation status rows, got %d", len(status))
	}

	projection, err := stockSvc.GetMonthlyDemandProjection(ctx)
	if err != nil {
		t.Fatalf("monthly demand projection: %v", err)
	}
	if len(projection) != 2 {
		t.Fatalf("expected 2 projected products, got %d", len(projection))
	}

	coverage, err := stockSvc.GetCoverageCheck(ctx)
	if err != nil {
		t.Fatalf("coverage check: %v", err)
	}
	if len(coverage) != 2 {
		t.Fatalf("expected 2 coverage entries, got %d", len(coverage))
	}

	var beansCoverage *StockCoverageItem
	for i := range coverage {
		if coverage[i].ProductID == beans.ID {
			beansCoverage = &coverage[i]
		}
	}
	if beansCoverage == nil {
		t.Fatalf("missing beans coverage")
	}
	if beansCoverage.IsCovered {
		t.Fatalf("expected beans not covered")
	}
	if beansCoverage.ShortfallBaseQuantity <= 0 {
		t.Fatalf("expected positive shortfall for beans")
	}
}

func TestStockServiceAssistentialWorkOutflowKg(t *testing.T) {
	conn := openStockTestDB(t)
	t.Cleanup(func() { _ = conn.Close() })

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	groupSvc := NewProductGroupService(conn)
	movementSvc := NewMovementService(conn)
	stockSvc := NewStockService(conn)
	ctx := context.Background()

	workA, _ := workSvc.Create(ctx, "Fraternidade", "")
	workB, _ := workSvc.Create(ctx, "Visita Fraterna", "")

	rice, _ := productSvc.Create(ctx, "Arroz", "kg", "")
	sugar, _ := productSvc.Create(ctx, "Acucar", "g", "")
	milk, _ := productSvc.Create(ctx, "Leite", "L", "")

	rice5kg, _ := variationSvc.Create(ctx, rice.ID, "Pacote 5kg", 5)
	sugar1kg, _ := variationSvc.Create(ctx, sugar.ID, "Pacote 1kg", 1000)
	milk1L, _ := variationSvc.Create(ctx, milk.ID, "Caixa 1L", 1)

	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 10 WHERE id IN (?, ?, ?)`, rice5kg.ID, sugar1kg.ID, milk1L.ID); err != nil {
		t.Fatalf("seed variation stock: %v", err)
	}

	group, _ := groupSvc.Create(ctx, "Cesta Doce", "")
	if _, err := groupSvc.UpsertItem(ctx, group.ID, sugar.ID, 1000); err != nil {
		t.Fatalf("upsert group item: %v", err)
	}

	if _, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: workA.ID,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice5kg.ID, Quantity: 2},
			{ProductVariationID: milk1L.ID, Quantity: 3},
		},
		GroupItems: []MovementGroupItemInput{
			{ProductGroupID: group.ID, Quantity: 1},
		},
	}); err != nil {
		t.Fatalf("create movement for work A: %v", err)
	}

	if _, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: workB.ID,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice5kg.ID, Quantity: 1},
		},
	}); err != nil {
		t.Fatalf("create movement for work B: %v", err)
	}

	outflow, err := stockSvc.GetAssistentialWorkOutflowKg(ctx)
	if err != nil {
		t.Fatalf("assistential work outflow kg: %v", err)
	}
	if len(outflow) != 2 {
		t.Fatalf("expected 2 assistential work rows, got %d", len(outflow))
	}

	if outflow[0].AssistentialWorkID != workA.ID {
		t.Fatalf("expected work A first, got %d", outflow[0].AssistentialWorkID)
	}
	if outflow[0].OutputKg != 11 {
		t.Fatalf("expected 11kg for work A, got %.3f", outflow[0].OutputKg)
	}
	if outflow[1].AssistentialWorkID != workB.ID {
		t.Fatalf("expected work B second, got %d", outflow[1].AssistentialWorkID)
	}
	if outflow[1].OutputKg != 5 {
		t.Fatalf("expected 5kg for work B, got %.3f", outflow[1].OutputKg)
	}
}

func TestStockServiceMonthlyMovementFlowKg(t *testing.T) {
	conn := openStockTestDB(t)
	t.Cleanup(func() { _ = conn.Close() })

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	workSvc := NewAssistentialWorkService(conn)
	productSvc := NewProductService(conn)
	variationSvc := NewProductVariationService(conn)
	groupSvc := NewProductGroupService(conn)
	movementSvc := NewMovementService(conn)
	stockSvc := NewStockService(conn)
	ctx := context.Background()

	work, _ := workSvc.Create(ctx, "Fraternidade", "")

	rice, _ := productSvc.Create(ctx, "Arroz", "kg", "")
	sugar, _ := productSvc.Create(ctx, "Acucar", "g", "")

	rice5kg, _ := variationSvc.Create(ctx, rice.ID, "Pacote 5kg", 5)
	sugar1kg, _ := variationSvc.Create(ctx, sugar.ID, "Pacote 1kg", 1000)

	if _, err := conn.Exec(`UPDATE product_variation SET current_stock = 20 WHERE id IN (?, ?)`, rice5kg.ID, sugar1kg.ID); err != nil {
		t.Fatalf("seed variation stock: %v", err)
	}

	group, _ := groupSvc.Create(ctx, "Cesta Doce", "")
	if _, err := groupSvc.UpsertItem(ctx, group.ID, sugar.ID, 1000); err != nil {
		t.Fatalf("upsert group item: %v", err)
	}

	januaryOut, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice5kg.ID, Quantity: 1},
		},
		GroupItems: []MovementGroupItemInput{
			{ProductGroupID: group.ID, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("create january out movement: %v", err)
	}

	januaryIn, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		Type:               "in",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice5kg.ID, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("create january in movement: %v", err)
	}

	februaryOut, err := movementSvc.Create(ctx, CreateMovementInput{
		AssistentialWorkID: work.ID,
		Type:               "out",
		ProductItems: []MovementProductItemInput{
			{ProductVariationID: rice5kg.ID, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("create february out movement: %v", err)
	}

	if _, err := conn.Exec(`UPDATE movement SET created_at = '2026-01-10 12:00:00' WHERE id = ?`, januaryOut.ID); err != nil {
		t.Fatalf("update january out date: %v", err)
	}
	if _, err := conn.Exec(`UPDATE movement SET created_at = '2026-01-20 12:00:00' WHERE id = ?`, januaryIn.ID); err != nil {
		t.Fatalf("update january in date: %v", err)
	}
	if _, err := conn.Exec(`UPDATE movement SET created_at = '2026-02-05 12:00:00' WHERE id = ?`, februaryOut.ID); err != nil {
		t.Fatalf("update february out date: %v", err)
	}

	flow, err := stockSvc.GetMonthlyMovementFlowKg(ctx)
	if err != nil {
		t.Fatalf("monthly movement flow kg: %v", err)
	}
	if len(flow) != 2 {
		t.Fatalf("expected 2 monthly flow rows, got %d", len(flow))
	}

	if flow[0].MonthKey != "2026-01" || flow[0].InputKg != 10 || flow[0].OutputKg != 6 {
		t.Fatalf("unexpected january row: %+v", flow[0])
	}
	if flow[1].MonthKey != "2026-02" || flow[1].InputKg != 0 || flow[1].OutputKg != 5 {
		t.Fatalf("unexpected february row: %+v", flow[1])
	}
}

func openStockTestDB(t *testing.T) *sql.DB {
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

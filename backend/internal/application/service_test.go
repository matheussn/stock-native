package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"stock/backend/internal/application"
	"stock/backend/internal/domain"
	"stock/backend/internal/infra/sqlite"
)

func setupService(t *testing.T) (*application.Service, func()) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	manager, err := sqlite.NewDBManager(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	repo := sqlite.NewRepository(manager.DB())
	service := application.NewService(repo, manager)
	cleanup := func() {
		_ = manager.Close()
	}
	return service, cleanup
}

func TestCreateExitMovement_InsufficientStock(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Alimentos"})
	product, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Arroz",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	destination, _ := svc.CreateDestination(ctx, application.CreateDestinationInput{Name: "Kits"})

	_, err := svc.CreateExitMovement(ctx, application.CreateMovementInput{
		ProductID:    product.ID,
		Quantity:     "1.000",
		MovementDate: time.Now(),
		SourceID:     destination.ID,
	})
	if err == nil {
		t.Fatal("expected insufficient stock error")
	}
	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != domain.ErrorInsufficientStock {
		t.Fatalf("expected %s, got %s", domain.ErrorInsufficientStock, appErr.Code)
	}
}

func TestCreateEntryAndExitMovement_UpdatesStock(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Higiene"})
	product, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Sabonete",
		CategoryID:    category.ID,
		MeasureUnit:   "unidade",
		PackageAmount: "1",
	})
	origin, _ := svc.CreateOrigin(ctx, application.CreateOriginInput{Name: "Doacao"})
	destination, _ := svc.CreateDestination(ctx, application.CreateDestinationInput{Name: "Familias"})

	if _, err := svc.CreateEntryMovement(ctx, application.CreateMovementInput{
		ProductID:    product.ID,
		Quantity:     "10.500",
		MovementDate: time.Now(),
		SourceID:     origin.ID,
	}); err != nil {
		t.Fatalf("entry failed: %v", err)
	}

	if _, err := svc.CreateExitMovement(ctx, application.CreateMovementInput{
		ProductID:    product.ID,
		Quantity:     "2.250",
		MovementDate: time.Now(),
		SourceID:     destination.ID,
	}); err != nil {
		t.Fatalf("exit failed: %v", err)
	}

	stock, err := svc.GetCurrentStock(ctx, application.CurrentStockFilter{})
	if err != nil {
		t.Fatalf("get stock failed: %v", err)
	}
	if len(stock) != 1 {
		t.Fatalf("expected 1 stock item, got %d", len(stock))
	}

	if stock[0].CurrentStock != 8.25 {
		t.Fatalf("expected stock 8.25, got %v", stock[0].CurrentStock)
	}
}

func TestCreateEntryMovementsBatch_UpdatesStock(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Mercearia"})
	rice, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Arroz",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	beans, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Feijao",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	origin, _ := svc.CreateOrigin(ctx, application.CreateOriginInput{Name: "Doacao"})

	if _, err := svc.CreateEntryMovementsBatch(ctx, application.CreateMovementBatchInput{
		MovementDate: time.Now(),
		SourceID:     origin.ID,
		Items: []application.MovementBatchItemInput{
			{ProductID: rice.ID, Quantity: "5.000"},
			{ProductID: beans.ID, Quantity: "3.500"},
		},
	}); err != nil {
		t.Fatalf("entry batch failed: %v", err)
	}

	stock, err := svc.GetCurrentStock(ctx, application.CurrentStockFilter{})
	if err != nil {
		t.Fatalf("get stock failed: %v", err)
	}

	if len(stock) != 2 {
		t.Fatalf("expected 2 stock items, got %d", len(stock))
	}
}

func TestCreateExitMovementsBatch_InsufficientStockInBatch(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Limpeza"})
	product, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Detergente",
		CategoryID:    category.ID,
		MeasureUnit:   "unidade",
		PackageAmount: "1",
	})
	origin, _ := svc.CreateOrigin(ctx, application.CreateOriginInput{Name: "Compra"})
	destination, _ := svc.CreateDestination(ctx, application.CreateDestinationInput{Name: "Entrega"})

	if _, err := svc.CreateEntryMovement(ctx, application.CreateMovementInput{
		ProductID:    product.ID,
		Quantity:     "5.000",
		MovementDate: time.Now(),
		SourceID:     origin.ID,
	}); err != nil {
		t.Fatalf("entry failed: %v", err)
	}

	_, err := svc.CreateExitMovementsBatch(ctx, application.CreateMovementBatchInput{
		MovementDate: time.Now(),
		SourceID:     destination.ID,
		Items: []application.MovementBatchItemInput{
			{ProductID: product.ID, Quantity: "3.000"},
			{ProductID: product.ID, Quantity: "3.000"},
		},
	})
	if err == nil {
		t.Fatal("expected insufficient stock error")
	}
	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != domain.ErrorInsufficientStock {
		t.Fatalf("expected %s, got %s", domain.ErrorInsufficientStock, appErr.Code)
	}
}

func TestCreateBasketEntryAndExit_WorksWithSubstitution(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Cesta"})
	rice, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Arroz",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	beans, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Feijao",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "1",
	})
	pasta, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Macarrao",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "1",
	})
	origin, _ := svc.CreateOrigin(ctx, application.CreateOriginInput{Name: "Doacao"})
	destination, _ := svc.CreateDestination(ctx, application.CreateDestinationInput{Name: "Entrega"})

	template, err := svc.CreateBasketTemplate(ctx, application.CreateBasketTemplateInput{
		Name: "Cesta Mensal",
		Items: []application.BasketTemplateItemInput{
			{ProductID: rice.ID, QuantityPerBasket: "1.000"},
			{ProductID: beans.ID, QuantityPerBasket: "1.000"},
		},
	})
	if err != nil {
		t.Fatalf("create basket template failed: %v", err)
	}

	if _, err := svc.CreateBasketEntry(ctx, application.CreateBasketMovementInput{
		BasketTemplateID: template.ID,
		BasketsCount:     "3.000",
		SourceID:         origin.ID,
		MovementDate:     time.Now(),
	}); err != nil {
		t.Fatalf("create basket entry failed: %v", err)
	}
	if _, err := svc.CreateEntryMovement(ctx, application.CreateMovementInput{
		ProductID:    pasta.ID,
		Quantity:     "1.000",
		MovementDate: time.Now(),
		SourceID:     origin.ID,
	}); err != nil {
		t.Fatalf("create pasta stock failed: %v", err)
	}

	movements, err := svc.CreateBasketExit(ctx, application.CreateBasketMovementInput{
		BasketTemplateID: template.ID,
		BasketsCount:     "1.000",
		SourceID:         destination.ID,
		MovementDate:     time.Now(),
		Substitutions: []application.BasketSubstitutionInput{
			{OriginalProductID: beans.ID, ProductID: pasta.ID, QuantityPerBasket: "0.500"},
		},
	})
	if err != nil {
		t.Fatalf("create basket exit failed: %v", err)
	}
	if len(movements) != 2 {
		t.Fatalf("expected 2 generated movements, got %d", len(movements))
	}

	history, err := svc.ListMovements(ctx, application.ListMovementsFilter{Type: "saida"})
	if err != nil {
		t.Fatalf("list movements failed: %v", err)
	}
	if len(history) == 0 {
		t.Fatal("expected exit movements")
	}
	if !strings.Contains(history[0].Note, "substitutions") {
		t.Fatalf("expected substitution metadata in note, got %q", history[0].Note)
	}
}

func TestCreateBasketExit_InsufficientStock(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Cesta"})
	product, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Arroz",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	destination, _ := svc.CreateDestination(ctx, application.CreateDestinationInput{Name: "Entrega"})

	template, err := svc.CreateBasketTemplate(ctx, application.CreateBasketTemplateInput{
		Name: "Cesta Mensal",
		Items: []application.BasketTemplateItemInput{
			{ProductID: product.ID, QuantityPerBasket: "1.000"},
		},
	})
	if err != nil {
		t.Fatalf("create basket template failed: %v", err)
	}

	_, err = svc.CreateBasketExit(ctx, application.CreateBasketMovementInput{
		BasketTemplateID: template.ID,
		BasketsCount:     "1.000",
		SourceID:         destination.ID,
		MovementDate:     time.Now(),
	})
	if err == nil {
		t.Fatal("expected insufficient stock error")
	}
	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != domain.ErrorInsufficientStock {
		t.Fatalf("expected %s, got %s", domain.ErrorInsufficientStock, appErr.Code)
	}
}

func TestGetStockPlanning_UsesFamiliesDemand(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()
	ctx := context.Background()

	category, _ := svc.CreateCategory(ctx, application.CreateCategoryInput{Name: "Cesta"})
	rice, _ := svc.CreateProduct(ctx, application.CreateProductInput{
		Name:          "Arroz",
		CategoryID:    category.ID,
		MeasureUnit:   "kg",
		PackageAmount: "5",
	})
	origin, _ := svc.CreateOrigin(ctx, application.CreateOriginInput{Name: "Doacao"})

	template, err := svc.CreateBasketTemplate(ctx, application.CreateBasketTemplateInput{
		Name: "Cesta Mensal",
		Items: []application.BasketTemplateItemInput{
			{ProductID: rice.ID, QuantityPerBasket: "2.000"},
		},
	})
	if err != nil {
		t.Fatalf("create basket template failed: %v", err)
	}

	if _, err := svc.CreateFamily(ctx, application.CreateFamilyInput{Name: "Familia A", CestasPerPeriod: "1.000"}); err != nil {
		t.Fatalf("create family A failed: %v", err)
	}
	if _, err := svc.CreateFamily(ctx, application.CreateFamilyInput{Name: "Familia B", CestasPerPeriod: "2.000"}); err != nil {
		t.Fatalf("create family B failed: %v", err)
	}

	if _, err := svc.CreateEntryMovement(ctx, application.CreateMovementInput{
		ProductID:    rice.ID,
		Quantity:     "4.000",
		MovementDate: time.Now(),
		SourceID:     origin.ID,
	}); err != nil {
		t.Fatalf("create entry failed: %v", err)
	}

	plan, err := svc.GetStockPlanning(ctx, application.StockPlanningFilter{
		BasketTemplateID: template.ID,
		HorizonDays:      30,
	})
	if err != nil {
		t.Fatalf("get stock planning failed: %v", err)
	}
	if plan.RequiredBaskets != 3 {
		t.Fatalf("expected required baskets 3, got %v", plan.RequiredBaskets)
	}
	if len(plan.Items) != 1 {
		t.Fatalf("expected 1 planning item, got %d", len(plan.Items))
	}
	if plan.Items[0].RequiredStock != 6 {
		t.Fatalf("expected required stock 6, got %v", plan.Items[0].RequiredStock)
	}
	if plan.Items[0].ProjectedBalance != -2 {
		t.Fatalf("expected projected balance -2, got %v", plan.Items[0].ProjectedBalance)
	}
	if plan.Items[0].Status != domain.StockPlanningAttention {
		t.Fatalf("expected status %s, got %s", domain.StockPlanningAttention, plan.Items[0].Status)
	}
}

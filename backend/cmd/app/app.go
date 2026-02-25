package app

import (
	"context"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"stock/backend/internal/application"
	"stock/backend/internal/domain"
	"stock/backend/internal/infra/sqlite"
)

type App struct {
	ctx context.Context
	svc *application.Service
	db  *sqlite.DBManager
}

func New() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	dbPath := filepath.Join("data", "app.db")

	manager, err := sqlite.NewDBManager(dbPath)
	if err != nil {
		runtime.LogError(ctx, err.Error())
		panic(err)
	}
	a.db = manager
	repo := sqlite.NewRepository(manager.DB())
	a.svc = application.NewService(repo, manager)
}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
	_ = ctx
}

func (a *App) HealthCheck() (map[string]string, error) {
	return a.svc.HealthCheck(), nil
}

func (a *App) ListCategories() ([]domain.Category, error) {
	return a.svc.ListCategories(a.ctx)
}

func (a *App) CreateCategory(input application.CreateCategoryInput) (*domain.Category, error) {
	return a.svc.CreateCategory(a.ctx, input)
}

func (a *App) UpdateCategory(input application.UpdateCategoryInput) (*domain.Category, error) {
	return a.svc.UpdateCategory(a.ctx, input)
}

func (a *App) DeactivateCategory(id string) error {
	return a.svc.DeactivateCategory(a.ctx, id)
}

func (a *App) ListProducts() ([]domain.Product, error) {
	return a.svc.ListProducts(a.ctx)
}

func (a *App) CreateProduct(input application.CreateProductInput) (*domain.Product, error) {
	return a.svc.CreateProduct(a.ctx, input)
}

func (a *App) UpdateProduct(input application.UpdateProductInput) (*domain.Product, error) {
	return a.svc.UpdateProduct(a.ctx, input)
}

func (a *App) DeactivateProduct(id string) error {
	return a.svc.DeactivateProduct(a.ctx, id)
}

func (a *App) ListOrigins() ([]domain.Origin, error) {
	return a.svc.ListOrigins(a.ctx)
}

func (a *App) CreateOrigin(input application.CreateOriginInput) (*domain.Origin, error) {
	return a.svc.CreateOrigin(a.ctx, input)
}

func (a *App) UpdateOrigin(input application.UpdateOriginInput) (*domain.Origin, error) {
	return a.svc.UpdateOrigin(a.ctx, input)
}

func (a *App) DeactivateOrigin(id string) error {
	return a.svc.DeactivateOrigin(a.ctx, id)
}

func (a *App) ListDestinations() ([]domain.Destination, error) {
	return a.svc.ListDestinations(a.ctx)
}

func (a *App) CreateDestination(input application.CreateDestinationInput) (*domain.Destination, error) {
	return a.svc.CreateDestination(a.ctx, input)
}

func (a *App) UpdateDestination(input application.UpdateDestinationInput) (*domain.Destination, error) {
	return a.svc.UpdateDestination(a.ctx, input)
}

func (a *App) DeactivateDestination(id string) error {
	return a.svc.DeactivateDestination(a.ctx, id)
}

func (a *App) ListBasketTemplates() ([]domain.BasketTemplate, error) {
	return a.svc.ListBasketTemplates(a.ctx)
}

func (a *App) CreateBasketTemplate(input application.CreateBasketTemplateInput) (*domain.BasketTemplate, error) {
	return a.svc.CreateBasketTemplate(a.ctx, input)
}

func (a *App) UpdateBasketTemplate(input application.UpdateBasketTemplateInput) (*domain.BasketTemplate, error) {
	return a.svc.UpdateBasketTemplate(a.ctx, input)
}

func (a *App) DeactivateBasketTemplate(id string) error {
	return a.svc.DeactivateBasketTemplate(a.ctx, id)
}

func (a *App) ListFamilies() ([]domain.Family, error) {
	return a.svc.ListFamilies(a.ctx)
}

func (a *App) CreateFamily(input application.CreateFamilyInput) (*domain.Family, error) {
	return a.svc.CreateFamily(a.ctx, input)
}

func (a *App) UpdateFamily(input application.UpdateFamilyInput) (*domain.Family, error) {
	return a.svc.UpdateFamily(a.ctx, input)
}

func (a *App) DeactivateFamily(id string) error {
	return a.svc.DeactivateFamily(a.ctx, id)
}

func (a *App) CreateEntryMovement(input application.CreateMovementInput) (*domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateEntryMovement(a.ctx, input)
}

func (a *App) CreateEntryMovementsBatch(input application.CreateMovementBatchInput) ([]domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateEntryMovementsBatch(a.ctx, input)
}

func (a *App) CreateExitMovement(input application.CreateMovementInput) (*domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateExitMovement(a.ctx, input)
}

func (a *App) CreateExitMovementsBatch(input application.CreateMovementBatchInput) ([]domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateExitMovementsBatch(a.ctx, input)
}

func (a *App) CreateBasketEntry(input application.CreateBasketMovementInput) ([]domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateBasketEntry(a.ctx, input)
}

func (a *App) CreateBasketExit(input application.CreateBasketMovementInput) ([]domain.Movement, error) {
	if input.MovementDate.IsZero() {
		input.MovementDate = time.Now()
	}
	return a.svc.CreateBasketExit(a.ctx, input)
}

func (a *App) ListMovements(filter application.ListMovementsFilter) ([]domain.Movement, error) {
	return a.svc.ListMovements(a.ctx, filter)
}

func (a *App) GetCurrentStock(filter application.CurrentStockFilter) ([]domain.StockItem, error) {
	return a.svc.GetCurrentStock(a.ctx, filter)
}

func (a *App) GetStockPlanning(filter application.StockPlanningFilter) (*domain.StockPlanning, error) {
	if filter.HorizonDays <= 0 {
		filter.HorizonDays = 30
	}
	return a.svc.GetStockPlanning(a.ctx, filter)
}

func (a *App) GetReports(filter application.ReportsFilter) (*domain.Reports, error) {
	if filter.DateFrom.IsZero() {
		filter.DateFrom = time.Now().AddDate(0, 0, -30)
	}
	if filter.DateTo.IsZero() {
		filter.DateTo = time.Now()
	}
	return a.svc.GetReports(a.ctx, filter)
}

func (a *App) ExportBackup(destination string) error {
	return a.svc.ExportBackup(a.ctx, destination)
}

func (a *App) ImportBackup(input application.ImportBackupInput) error {
	return a.svc.ImportBackup(a.ctx, input)
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"stock/internal/db"
	"stock/internal/models"
	"stock/internal/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	db            *sql.DB
	dbPath        string
	logPath       string
	logger        *operationalLogger
	startupStatus AppStartupStatus
	mu            sync.RWMutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.initializeOperationalLogger(); err != nil && ctx != nil {
		runtime.LogErrorf(ctx, "failed to initialize operational logger: %v", err)
	}

	a.logInfof("application startup initiated")
	if err := a.initializeDatabase(); err != nil {
		a.logErrorf("database initialization failed during startup: %v", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.logInfof("application shutdown started")

	a.mu.Lock()
	currentDB := a.db
	a.db = nil
	a.mu.Unlock()

	if currentDB != nil {
		if err := currentDB.Close(); err != nil {
			a.logErrorf("failed to close database connection: %v", err)
		} else {
			a.logInfof("database connection closed")
		}
	}

	a.logInfof("application shutdown completed")
	a.closeOperationalLogger()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CreateAssistentialWork(name, description string) (models.AssistentialWork, error) {
	work, err := services.NewAssistentialWorkService(a.db).Create(a.requestContext(), name, description)
	if err != nil {
		a.logMutationFailure("create assistential work", err, "name=%q", name)
		return models.AssistentialWork{}, err
	}

	a.logMutationSuccess("assistential work created", "id=%d name=%q", work.ID, work.Name)
	return work, nil
}

func (a *App) ListAssistentialWorks(includeInactive bool) ([]models.AssistentialWork, error) {
	return services.NewAssistentialWorkService(a.db).List(a.requestContext(), includeInactive)
}

func (a *App) UpdateAssistentialWork(id int64, name, description string) (models.AssistentialWork, error) {
	work, err := services.NewAssistentialWorkService(a.db).Update(a.requestContext(), id, name, description)
	if err != nil {
		a.logMutationFailure("update assistential work", err, "id=%d name=%q", id, name)
		return models.AssistentialWork{}, err
	}

	a.logMutationSuccess("assistential work updated", "id=%d name=%q", work.ID, work.Name)
	return work, nil
}

func (a *App) SetAssistentialWorkActive(id int64, isActive bool) error {
	err := services.NewAssistentialWorkService(a.db).SetActive(a.requestContext(), id, isActive)
	if err != nil {
		a.logMutationFailure("toggle assistential work active flag", err, "id=%d active=%t", id, isActive)
		return err
	}

	a.logMutationSuccess("assistential work active flag updated", "id=%d active=%t", id, isActive)
	return nil
}

func (a *App) CreateProduct(name, baseUnit, description string) (models.Product, error) {
	product, err := services.NewProductService(a.db).Create(a.requestContext(), name, baseUnit, description)
	if err != nil {
		a.logMutationFailure("create product", err, "name=%q base_unit=%q", name, baseUnit)
		return models.Product{}, err
	}

	a.logMutationSuccess("product created", "id=%d name=%q base_unit=%q", product.ID, product.Name, product.BaseUnit)
	return product, nil
}

func (a *App) ListProducts(includeInactive bool) ([]models.Product, error) {
	return services.NewProductService(a.db).List(a.requestContext(), includeInactive)
}

func (a *App) UpdateProduct(id int64, name, baseUnit, description string) (models.Product, error) {
	product, err := services.NewProductService(a.db).Update(a.requestContext(), id, name, baseUnit, description)
	if err != nil {
		a.logMutationFailure("update product", err, "id=%d name=%q base_unit=%q", id, name, baseUnit)
		return models.Product{}, err
	}

	a.logMutationSuccess("product updated", "id=%d name=%q base_unit=%q", product.ID, product.Name, product.BaseUnit)
	return product, nil
}

func (a *App) SetProductActive(id int64, isActive bool) error {
	err := services.NewProductService(a.db).SetActive(a.requestContext(), id, isActive)
	if err != nil {
		a.logMutationFailure("toggle product active flag", err, "id=%d active=%t", id, isActive)
		return err
	}

	a.logMutationSuccess("product active flag updated", "id=%d active=%t", id, isActive)
	return nil
}

func (a *App) CreateProductVariation(productID int64, description string, baseQuantity int64) (models.ProductVariation, error) {
	variation, err := services.NewProductVariationService(a.db).Create(a.requestContext(), productID, description, baseQuantity)
	if err != nil {
		a.logMutationFailure("create product variation", err, "product_id=%d base_quantity=%d", productID, baseQuantity)
		return models.ProductVariation{}, err
	}

	a.logMutationSuccess("product variation created", "id=%d product_id=%d base_quantity=%d", variation.ID, variation.ProductID, variation.BaseQuantity)
	return variation, nil
}

func (a *App) ListProductVariationsByProduct(productID int64, includeInactive bool) ([]models.ProductVariation, error) {
	return services.NewProductVariationService(a.db).ListByProduct(a.requestContext(), productID, includeInactive)
}

func (a *App) UpdateProductVariation(id int64, description string, baseQuantity int64) (models.ProductVariation, error) {
	variation, err := services.NewProductVariationService(a.db).Update(a.requestContext(), id, description, baseQuantity)
	if err != nil {
		a.logMutationFailure("update product variation", err, "id=%d base_quantity=%d", id, baseQuantity)
		return models.ProductVariation{}, err
	}

	a.logMutationSuccess("product variation updated", "id=%d product_id=%d base_quantity=%d", variation.ID, variation.ProductID, variation.BaseQuantity)
	return variation, nil
}

func (a *App) SetProductVariationActive(id int64, isActive bool) error {
	err := services.NewProductVariationService(a.db).SetActive(a.requestContext(), id, isActive)
	if err != nil {
		a.logMutationFailure("toggle product variation active flag", err, "id=%d active=%t", id, isActive)
		return err
	}

	a.logMutationSuccess("product variation active flag updated", "id=%d active=%t", id, isActive)
	return nil
}

func (a *App) CreateProductGroup(name, description string) (models.ProductGroup, error) {
	group, err := services.NewProductGroupService(a.db).Create(a.requestContext(), name, description)
	if err != nil {
		a.logMutationFailure("create product group", err, "name=%q", name)
		return models.ProductGroup{}, err
	}

	a.logMutationSuccess("product group created", "id=%d name=%q", group.ID, group.Name)
	return group, nil
}

func (a *App) ListProductGroups(includeInactive bool) ([]models.ProductGroup, error) {
	return services.NewProductGroupService(a.db).List(a.requestContext(), includeInactive)
}

func (a *App) UpdateProductGroup(id int64, name, description string) (models.ProductGroup, error) {
	group, err := services.NewProductGroupService(a.db).Update(a.requestContext(), id, name, description)
	if err != nil {
		a.logMutationFailure("update product group", err, "id=%d name=%q", id, name)
		return models.ProductGroup{}, err
	}

	a.logMutationSuccess("product group updated", "id=%d name=%q", group.ID, group.Name)
	return group, nil
}

func (a *App) SetProductGroupActive(id int64, isActive bool) error {
	err := services.NewProductGroupService(a.db).SetActive(a.requestContext(), id, isActive)
	if err != nil {
		a.logMutationFailure("toggle product group active flag", err, "id=%d active=%t", id, isActive)
		return err
	}

	a.logMutationSuccess("product group active flag updated", "id=%d active=%t", id, isActive)
	return nil
}

func (a *App) ListProductGroupItems(productGroupID int64) ([]models.ProductGroupItem, error) {
	return services.NewProductGroupService(a.db).ListItems(a.requestContext(), productGroupID)
}

func (a *App) UpsertProductGroupItem(productGroupID, productID, baseQuantity int64) (models.ProductGroupItem, error) {
	item, err := services.NewProductGroupService(a.db).UpsertItem(a.requestContext(), productGroupID, productID, baseQuantity)
	if err != nil {
		a.logMutationFailure("upsert product group item", err, "product_group_id=%d product_id=%d base_quantity=%d", productGroupID, productID, baseQuantity)
		return models.ProductGroupItem{}, err
	}

	a.logMutationSuccess("product group item upserted", "id=%d product_group_id=%d product_id=%d base_quantity=%d", item.ID, item.ProductGroupID, item.ProductID, item.BaseQuantity)
	return item, nil
}

func (a *App) RemoveProductGroupItem(productGroupID, productID int64) error {
	err := services.NewProductGroupService(a.db).RemoveItem(a.requestContext(), productGroupID, productID)
	if err != nil {
		a.logMutationFailure("remove product group item", err, "product_group_id=%d product_id=%d", productGroupID, productID)
		return err
	}

	a.logMutationSuccess("product group item removed", "product_group_id=%d product_id=%d", productGroupID, productID)
	return nil
}

func (a *App) CreateFamily(assistentialWorkID int64, name string, memberCount int64, address, contact string) (models.Family, error) {
	family, err := services.NewFamilyService(a.db).Create(a.requestContext(), assistentialWorkID, name, memberCount, address, contact)
	if err != nil {
		a.logMutationFailure("create family", err, "assistential_work_id=%d name=%q", assistentialWorkID, name)
		return models.Family{}, err
	}

	a.logMutationSuccess("family created", "id=%d assistential_work_id=%d name=%q", family.ID, family.AssistentialWorkID, family.Name)
	return family, nil
}

func (a *App) ListFamilies(includeInactive bool, assistentialWorkID int64) ([]models.Family, error) {
	return services.NewFamilyService(a.db).List(a.requestContext(), includeInactive, assistentialWorkID)
}

func (a *App) UpdateFamily(id int64, name string, memberCount int64, address, contact string) (models.Family, error) {
	family, err := services.NewFamilyService(a.db).Update(a.requestContext(), id, name, memberCount, address, contact)
	if err != nil {
		a.logMutationFailure("update family", err, "id=%d name=%q member_count=%d", id, name, memberCount)
		return models.Family{}, err
	}

	a.logMutationSuccess("family updated", "id=%d name=%q member_count=%d", family.ID, family.Name, family.MemberCount)
	return family, nil
}

func (a *App) SetFamilyActive(id int64, isActive bool) error {
	err := services.NewFamilyService(a.db).SetActive(a.requestContext(), id, isActive)
	if err != nil {
		a.logMutationFailure("toggle family active flag", err, "id=%d active=%t", id, isActive)
		return err
	}

	a.logMutationSuccess("family active flag updated", "id=%d active=%t", id, isActive)
	return nil
}

func (a *App) AssignFamilyGroup(familyID, productGroupID int64) (models.FamilyGroupAssignment, error) {
	assignment, err := services.NewFamilyGroupAssignmentService(a.db).Assign(a.requestContext(), familyID, productGroupID)
	if err != nil {
		a.logMutationFailure("assign family group", err, "family_id=%d product_group_id=%d", familyID, productGroupID)
		return models.FamilyGroupAssignment{}, err
	}

	a.logMutationSuccess("family group assigned", "assignment_id=%d family_id=%d product_group_id=%d", assignment.ID, assignment.FamilyID, assignment.ProductGroupID)
	return assignment, nil
}

func (a *App) ListFamilyGroupAssignments(familyID int64) ([]models.FamilyGroupAssignment, error) {
	return services.NewFamilyGroupAssignmentService(a.db).ListByFamily(a.requestContext(), familyID)
}

func (a *App) CreateMovement(input services.CreateMovementInput) (models.Movement, error) {
	movement, err := services.NewMovementService(a.db).Create(a.requestContext(), input)
	if err != nil {
		a.logMutationFailure(
			"create movement",
			err,
			"assistential_work_id=%d type=%q product_items=%d group_items=%d",
			input.AssistentialWorkID,
			input.Type,
			len(input.ProductItems),
			len(input.GroupItems),
		)
		return models.Movement{}, err
	}

	a.logMutationSuccess(
		"movement created",
		"id=%d assistential_work_id=%d type=%q product_items=%d group_items=%d",
		movement.ID,
		movement.AssistentialWorkID,
		movement.Type,
		len(input.ProductItems),
		len(input.GroupItems),
	)
	return movement, nil
}

func (a *App) ListMovements(movementType string, assistentialWorkID int64) ([]models.Movement, error) {
	return services.NewMovementService(a.db).List(a.requestContext(), movementType, assistentialWorkID)
}

func (a *App) ListMovementProductItems(movementID int64) ([]models.MovementProductItem, error) {
	return services.NewMovementService(a.db).ListProductItemsByMovement(a.requestContext(), movementID)
}

func (a *App) ListMovementGroupItems(movementID int64) ([]models.MovementGroupItem, error) {
	return services.NewMovementService(a.db).ListGroupItemsByMovement(a.requestContext(), movementID)
}

func (a *App) ListMovementGroupItemResolutions(movementGroupItemID int64) ([]models.MovementGroupItemResolution, error) {
	return services.NewMovementService(a.db).ListGroupItemResolutions(a.requestContext(), movementGroupItemID)
}

func (a *App) ListVariationStockStatus(includeInactive bool) ([]services.VariationStockStatus, error) {
	return services.NewStockService(a.db).ListVariationStatus(a.requestContext(), includeInactive)
}

func (a *App) GetMonthlyDemandProjection() ([]services.MonthlyDemandItem, error) {
	return services.NewStockService(a.db).GetMonthlyDemandProjection(a.requestContext())
}

func (a *App) GetStockCoverageCheck() ([]services.StockCoverageItem, error) {
	return services.NewStockService(a.db).GetCoverageCheck(a.requestContext())
}

func (a *App) GetAssistentialWorkOutflowKg() ([]services.AssistentialWorkOutflowKgItem, error) {
	return services.NewStockService(a.db).GetAssistentialWorkOutflowKg(a.requestContext())
}

func (a *App) GetMonthlyMovementFlowKg() ([]services.MonthlyMovementFlowKgItem, error) {
	return services.NewStockService(a.db).GetMonthlyMovementFlowKg(a.requestContext())
}

func (a *App) BackupDatabase(destinationPath string) (string, error) {
	if err := a.databaseReadyError(); err != nil {
		a.logMutationFailure("database backup", err, "destination=%q", destinationPath)
		return "", err
	}

	backedUpPath, err := db.Backup(a.currentDatabase(), a.GetDatabasePath(), destinationPath)
	if err != nil {
		a.logMutationFailure("database backup", err, "destination=%q", destinationPath)
		return "", err
	}

	a.logMutationSuccess("database backup created", "path=%q", backedUpPath)
	return backedUpPath, nil
}

func (a *App) RestoreDatabase(sourcePath string) (db.RestoreResult, error) {
	a.logInfof("database restore requested: source=%q", sourcePath)

	activeDBPath, err := a.resolveDatabasePath()
	if err != nil {
		a.logMutationFailure("database restore", err, "source=%q", sourcePath)
		return db.RestoreResult{}, err
	}

	a.mu.Lock()
	currentDB := a.db
	a.db = nil
	a.mu.Unlock()

	if currentDB != nil {
		if err := currentDB.Close(); err != nil {
			a.mu.Lock()
			a.db = currentDB
			a.mu.Unlock()
			a.logMutationFailure("database restore", err, "source=%q", sourcePath)
			return db.RestoreResult{}, fmt.Errorf("close current database before restore: %w", err)
		}
	}

	restoreResult, restoreErr := db.RestoreFile(activeDBPath, sourcePath)
	if restoreErr != nil {
		reopenErr := a.initializeDatabase()
		if reopenErr != nil {
			a.logErrorf("database restore failed and reopen failed: source=%q restore_err=%v reopen_err=%v", sourcePath, restoreErr, reopenErr)
			return db.RestoreResult{}, fmt.Errorf("restore database: %w; failed to reopen database: %v", restoreErr, reopenErr)
		}

		a.logMutationFailure("database restore", restoreErr, "source=%q", sourcePath)
		return db.RestoreResult{}, fmt.Errorf("restore database: %w", restoreErr)
	}

	if err := a.initializeDatabase(); err != nil {
		a.logMutationFailure("database restore reopen", err, "source=%q", sourcePath)
		return db.RestoreResult{}, fmt.Errorf("open restored database: %w", err)
	}

	a.logMutationSuccess("database restored", "source=%q restored_path=%q previous_backup_path=%q", sourcePath, restoreResult.RestoredPath, restoreResult.PreviousBackupPath)

	return restoreResult, nil
}

func (a *App) requestContext() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

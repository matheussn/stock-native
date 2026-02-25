package application

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"stock/backend/internal/domain"
)

type BackupManager interface {
	ExportBackup(ctx context.Context, destination string) error
	ImportBackup(ctx context.Context, source string) error
}

type Service struct {
	repo   Repository
	backup BackupManager
}

func NewService(repo Repository, backup BackupManager) *Service {
	return &Service{repo: repo, backup: backup}
}

func (s *Service) HealthCheck() map[string]string {
	return map[string]string{"status": "ok"}
}

func (s *Service) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) CreateCategory(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.Category{ID: uuid.NewString(), Name: strings.TrimSpace(input.Name), Active: true, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateCategory(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateCategory(ctx context.Context, input UpdateCategoryInput) (*domain.Category, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	item := domain.Category{ID: input.ID, Name: strings.TrimSpace(input.Name), UpdatedAt: time.Now()}
	if err := s.repo.UpdateCategory(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) DeactivateCategory(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateCategory(ctx, id))
}

func (s *Service) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.ListProducts(ctx)
}

func parseThreshold(raw string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || value < 0 {
		return 0, domain.NewAppError(domain.ErrorValidation, "lowStockThreshold must be zero or positive", map[string]string{"field": "lowStockThreshold"})
	}
	return domain.Round3(value), nil
}

func parsePackageAmount(raw string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, domain.NewAppError(domain.ErrorValidation, "packageAmount is required", map[string]string{"field": "packageAmount"})
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || value <= 0 {
		return 0, domain.NewAppError(domain.ErrorValidation, "packageAmount must be greater than zero", map[string]string{"field": "packageAmount"})
	}
	return domain.Round3(value), nil
}

func parseRequiredPositiveDecimal(raw string, field string) (float64, error) {
	value, err := domain.ValidateQuantityString(raw)
	if err != nil {
		return 0, domain.NewAppError(domain.ErrorValidation, field+" must be greater than zero with up to 3 decimal places", map[string]string{"field": field})
	}
	return value, nil
}

func (s *Service) CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.CategoryID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "categoryId is required", map[string]string{"field": "categoryId"})
	}
	if err := domain.ValidateMeasureUnit(input.MeasureUnit); err != nil {
		return nil, err
	}
	packageAmount, err := parsePackageAmount(input.PackageAmount)
	if err != nil {
		return nil, err
	}
	threshold, err := parseThreshold(input.LowStockThreshold)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.Product{
		ID:                uuid.NewString(),
		Name:              strings.TrimSpace(input.Name),
		CategoryID:        input.CategoryID,
		MeasureUnit:       domain.MeasureUnit(strings.ToLower(strings.TrimSpace(input.MeasureUnit))),
		PackageAmount:     packageAmount,
		LowStockThreshold: threshold,
		Active:            true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := s.repo.CreateProduct(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateProduct(ctx context.Context, input UpdateProductInput) (*domain.Product, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.CategoryID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "categoryId is required", map[string]string{"field": "categoryId"})
	}
	if err := domain.ValidateMeasureUnit(input.MeasureUnit); err != nil {
		return nil, err
	}
	packageAmount, err := parsePackageAmount(input.PackageAmount)
	if err != nil {
		return nil, err
	}
	threshold, err := parseThreshold(input.LowStockThreshold)
	if err != nil {
		return nil, err
	}
	item := domain.Product{
		ID:                input.ID,
		Name:              strings.TrimSpace(input.Name),
		CategoryID:        input.CategoryID,
		MeasureUnit:       domain.MeasureUnit(strings.ToLower(strings.TrimSpace(input.MeasureUnit))),
		PackageAmount:     packageAmount,
		LowStockThreshold: threshold,
		UpdatedAt:         time.Now(),
	}
	if err := s.repo.UpdateProduct(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) DeactivateProduct(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateProduct(ctx, id))
}

func (s *Service) ListOrigins(ctx context.Context) ([]domain.Origin, error) {
	return s.repo.ListOrigins(ctx)
}

func (s *Service) CreateOrigin(ctx context.Context, input CreateOriginInput) (*domain.Origin, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.Origin{ID: uuid.NewString(), Name: strings.TrimSpace(input.Name), Active: true, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateOrigin(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateOrigin(ctx context.Context, input UpdateOriginInput) (*domain.Origin, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	item := domain.Origin{ID: input.ID, Name: strings.TrimSpace(input.Name), UpdatedAt: time.Now()}
	if err := s.repo.UpdateOrigin(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) DeactivateOrigin(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateOrigin(ctx, id))
}

func (s *Service) ListDestinations(ctx context.Context) ([]domain.Destination, error) {
	return s.repo.ListDestinations(ctx)
}

func (s *Service) CreateDestination(ctx context.Context, input CreateDestinationInput) (*domain.Destination, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.Destination{ID: uuid.NewString(), Name: strings.TrimSpace(input.Name), Active: true, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.CreateDestination(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateDestination(ctx context.Context, input UpdateDestinationInput) (*domain.Destination, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	item := domain.Destination{ID: input.ID, Name: strings.TrimSpace(input.Name), UpdatedAt: time.Now()}
	if err := s.repo.UpdateDestination(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) DeactivateDestination(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateDestination(ctx, id))
}

func (s *Service) ListBasketTemplates(ctx context.Context) ([]domain.BasketTemplate, error) {
	return s.repo.ListBasketTemplates(ctx)
}

func buildBasketTemplateItems(rawItems []BasketTemplateItemInput) ([]domain.BasketTemplateItem, error) {
	if len(rawItems) == 0 {
		return nil, domain.NewAppError(domain.ErrorValidation, "items are required", map[string]string{"field": "items"})
	}
	items := make([]domain.BasketTemplateItem, 0, len(rawItems))
	seenProductIDs := map[string]struct{}{}
	now := time.Now()
	for index, raw := range rawItems {
		productID := strings.TrimSpace(raw.ProductID)
		if productID == "" {
			return nil, domain.NewAppError(domain.ErrorValidation, "productId is required in items", map[string]any{"index": index, "field": "productId"})
		}
		if _, exists := seenProductIDs[productID]; exists {
			return nil, domain.NewAppError(domain.ErrorValidation, "productId duplicated in items", map[string]any{"index": index, "field": "productId"})
		}
		seenProductIDs[productID] = struct{}{}

		quantity, err := parseRequiredPositiveDecimal(raw.QuantityPerBasket, "quantityPerBasket")
		if err != nil {
			return nil, domain.NewAppError(domain.ErrorValidation, "invalid quantityPerBasket in items", map[string]any{"index": index, "error": err.Error()})
		}
		items = append(items, domain.BasketTemplateItem{
			ID:                uuid.NewString(),
			ProductID:         productID,
			QuantityPerBasket: quantity,
			CreatedAt:         now,
			UpdatedAt:         now,
		})
	}
	return items, nil
}

func (s *Service) CreateBasketTemplate(ctx context.Context, input CreateBasketTemplateInput) (*domain.BasketTemplate, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	items, err := buildBasketTemplateItems(input.Items)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.BasketTemplate{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(input.Name),
		Active:    true,
		Items:     items,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for i := range item.Items {
		item.Items[i].BasketTemplateID = item.ID
	}
	if err := s.repo.CreateBasketTemplate(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateBasketTemplate(ctx context.Context, input UpdateBasketTemplateInput) (*domain.BasketTemplate, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	items, err := buildBasketTemplateItems(input.Items)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.BasketTemplate{
		ID:        input.ID,
		Name:      strings.TrimSpace(input.Name),
		Items:     items,
		UpdatedAt: now,
	}
	for i := range item.Items {
		item.Items[i].BasketTemplateID = item.ID
		item.Items[i].UpdatedAt = now
		item.Items[i].CreatedAt = now
	}
	if err := s.repo.UpdateBasketTemplate(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return s.repo.GetBasketTemplate(ctx, item.ID)
}

func (s *Service) DeactivateBasketTemplate(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateBasketTemplate(ctx, id))
}

func parseCestasPerPeriod(raw string) (float64, error) {
	return parseRequiredPositiveDecimal(raw, "cestasPerPeriod")
}

func (s *Service) ListFamilies(ctx context.Context) ([]domain.Family, error) {
	return s.repo.ListFamilies(ctx)
}

func (s *Service) CreateFamily(ctx context.Context, input CreateFamilyInput) (*domain.Family, error) {
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	cestasPerPeriod, err := parseCestasPerPeriod(input.CestasPerPeriod)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := domain.Family{
		ID:              uuid.NewString(),
		Name:            strings.TrimSpace(input.Name),
		CestasPerPeriod: cestasPerPeriod,
		Active:          true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.repo.CreateFamily(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) UpdateFamily(ctx context.Context, input UpdateFamilyInput) (*domain.Family, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	if err := domain.ValidateName(input.Name, "name"); err != nil {
		return nil, err
	}
	cestasPerPeriod, err := parseCestasPerPeriod(input.CestasPerPeriod)
	if err != nil {
		return nil, err
	}
	item := domain.Family{
		ID:              input.ID,
		Name:            strings.TrimSpace(input.Name),
		CestasPerPeriod: cestasPerPeriod,
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.UpdateFamily(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) DeactivateFamily(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.NewAppError(domain.ErrorValidation, "id is required", map[string]string{"field": "id"})
	}
	return mapRepoError(s.repo.DeactivateFamily(ctx, id))
}

type resolvedBasketSubstitution struct {
	OriginalProductID string  `json:"originalProductId"`
	ProductID         string  `json:"productId"`
	QuantityPerBasket float64 `json:"quantityPerBasket"`
}

func resolveBasketItems(template *domain.BasketTemplate, basketsCount float64, substitutions []BasketSubstitutionInput) ([]MovementBatchItemInput, []resolvedBasketSubstitution, error) {
	if template == nil {
		return nil, nil, domain.NewAppError(domain.ErrorValidation, "basketTemplate is required", nil)
	}
	substitutionByOriginal := map[string]BasketSubstitutionInput{}
	for index, sub := range substitutions {
		originalID := strings.TrimSpace(sub.OriginalProductID)
		productID := strings.TrimSpace(sub.ProductID)
		if originalID == "" || productID == "" {
			return nil, nil, domain.NewAppError(domain.ErrorValidation, "originalProductId and productId are required in substitutions", map[string]any{"index": index})
		}
		if _, exists := substitutionByOriginal[originalID]; exists {
			return nil, nil, domain.NewAppError(domain.ErrorValidation, "duplicated substitution for originalProductId", map[string]any{"index": index, "originalProductId": originalID})
		}
		sub.OriginalProductID = originalID
		sub.ProductID = productID
		substitutionByOriginal[originalID] = sub
	}

	availableOriginals := map[string]domain.BasketTemplateItem{}
	for _, item := range template.Items {
		availableOriginals[item.ProductID] = item
	}

	totalsByProduct := map[string]float64{}
	resolvedSubs := make([]resolvedBasketSubstitution, 0, len(substitutions))

	for _, item := range template.Items {
		targetProductID := item.ProductID
		quantityPerBasket := item.QuantityPerBasket
		if sub, hasSub := substitutionByOriginal[item.ProductID]; hasSub {
			targetProductID = sub.ProductID
			if strings.TrimSpace(sub.QuantityPerBasket) != "" {
				parsedQty, err := parseRequiredPositiveDecimal(sub.QuantityPerBasket, "quantityPerBasket")
				if err != nil {
					return nil, nil, domain.NewAppError(domain.ErrorValidation, "invalid quantityPerBasket in substitutions", map[string]any{"originalProductId": item.ProductID, "error": err.Error()})
				}
				quantityPerBasket = parsedQty
			}
			resolvedSubs = append(resolvedSubs, resolvedBasketSubstitution{
				OriginalProductID: item.ProductID,
				ProductID:         targetProductID,
				QuantityPerBasket: quantityPerBasket,
			})
		}
		totalsByProduct[targetProductID] = domain.Round3(totalsByProduct[targetProductID] + quantityPerBasket*basketsCount)
	}

	for originalID := range substitutionByOriginal {
		if _, ok := availableOriginals[originalID]; !ok {
			return nil, nil, domain.NewAppError(domain.ErrorValidation, "substitution originalProductId not found in basket template", map[string]any{"originalProductId": originalID})
		}
	}

	items := make([]MovementBatchItemInput, 0, len(totalsByProduct))
	for productID, quantity := range totalsByProduct {
		items = append(items, MovementBatchItemInput{
			ProductID: productID,
			Quantity:  strconv.FormatFloat(domain.Round3(quantity), 'f', 3, 64),
		})
	}
	return items, resolvedSubs, nil
}

func buildBasketOperationNote(baseNote string, basketTemplateName string, basketsCount float64, substitutions []resolvedBasketSubstitution) string {
	trimmed := strings.TrimSpace(baseNote)
	payload := map[string]any{
		"basketTemplateName": basketTemplateName,
		"basketsCount":       domain.Round3(basketsCount),
	}
	if len(substitutions) > 0 {
		payload["substitutions"] = substitutions
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return trimmed
	}
	if trimmed == "" {
		return string(raw)
	}
	return trimmed + " | " + string(raw)
}

func (s *Service) CreateBasketEntry(ctx context.Context, input CreateBasketMovementInput) ([]domain.Movement, error) {
	if strings.TrimSpace(input.BasketTemplateID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "basketTemplateId is required", map[string]string{"field": "basketTemplateId"})
	}
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}
	basketsCount, err := parseRequiredPositiveDecimal(input.BasketsCount, "basketsCount")
	if err != nil {
		return nil, err
	}
	template, err := s.repo.GetBasketTemplate(ctx, input.BasketTemplateID)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if template == nil || !template.Active {
		return nil, domain.NewAppError(domain.ErrorValidation, "basket template is invalid", map[string]string{"field": "basketTemplateId"})
	}

	items, resolvedSubs, err := resolveBasketItems(template, basketsCount, input.Substitutions)
	if err != nil {
		return nil, err
	}
	note := buildBasketOperationNote(input.Note, template.Name, basketsCount, resolvedSubs)
	return s.CreateEntryMovementsBatch(ctx, CreateMovementBatchInput{
		MovementDate: input.MovementDate,
		SourceID:     input.SourceID,
		Note:         note,
		Items:        items,
	})
}

func (s *Service) CreateBasketExit(ctx context.Context, input CreateBasketMovementInput) ([]domain.Movement, error) {
	if strings.TrimSpace(input.BasketTemplateID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "basketTemplateId is required", map[string]string{"field": "basketTemplateId"})
	}
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}
	basketsCount, err := parseRequiredPositiveDecimal(input.BasketsCount, "basketsCount")
	if err != nil {
		return nil, err
	}
	template, err := s.repo.GetBasketTemplate(ctx, input.BasketTemplateID)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if template == nil || !template.Active {
		return nil, domain.NewAppError(domain.ErrorValidation, "basket template is invalid", map[string]string{"field": "basketTemplateId"})
	}

	items, resolvedSubs, err := resolveBasketItems(template, basketsCount, input.Substitutions)
	if err != nil {
		return nil, err
	}
	note := buildBasketOperationNote(input.Note, template.Name, basketsCount, resolvedSubs)
	return s.CreateExitMovementsBatch(ctx, CreateMovementBatchInput{
		MovementDate: input.MovementDate,
		SourceID:     input.SourceID,
		Note:         note,
		Items:        items,
	})
}

func planningStatus(current float64, required float64) domain.StockPlanningStatus {
	if current >= required {
		return domain.StockPlanningSufficient
	}
	if current > 0 {
		return domain.StockPlanningAttention
	}
	return domain.StockPlanningCritical
}

func (s *Service) GetStockPlanning(ctx context.Context, filter StockPlanningFilter) (*domain.StockPlanning, error) {
	if strings.TrimSpace(filter.BasketTemplateID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "basketTemplateId is required", map[string]string{"field": "basketTemplateId"})
	}
	if filter.HorizonDays <= 0 {
		filter.HorizonDays = 30
	}
	template, err := s.repo.GetBasketTemplate(ctx, filter.BasketTemplateID)
	if err != nil {
		return nil, mapRepoError(err)
	}
	if template == nil {
		return nil, domain.NewAppError(domain.ErrorNotFound, "basket template not found", nil)
	}
	requiredBasketsPer30Days, err := s.repo.GetActiveFamiliesCestas(ctx)
	if err != nil {
		return nil, mapRepoError(err)
	}
	requiredBaskets := domain.Round3(requiredBasketsPer30Days * float64(filter.HorizonDays) / 30.0)
	currentStock, err := s.repo.GetCurrentStock(ctx, CurrentStockFilter{})
	if err != nil {
		return nil, mapRepoError(err)
	}
	stockByProduct := map[string]domain.StockItem{}
	for _, stock := range currentStock {
		stockByProduct[stock.ProductID] = stock
	}

	items := make([]domain.StockPlanningItem, 0, len(template.Items))
	for _, basketItem := range template.Items {
		requiredStock := domain.Round3(requiredBaskets * basketItem.QuantityPerBasket)
		stock, ok := stockByProduct[basketItem.ProductID]
		current := 0.0
		productName := basketItem.ProductName
		measureUnit := basketItem.MeasureUnit
		if ok {
			current = stock.CurrentStock
			productName = stock.ProductName
			measureUnit = stock.MeasureUnit
		}
		projectedBalance := domain.Round3(current - requiredStock)
		items = append(items, domain.StockPlanningItem{
			ProductID:        basketItem.ProductID,
			ProductName:      productName,
			MeasureUnit:      measureUnit,
			CurrentStock:     domain.Round3(current),
			RequiredStock:    requiredStock,
			ProjectedBalance: projectedBalance,
			Status:           planningStatus(current, requiredStock),
		})
	}

	return &domain.StockPlanning{
		BasketTemplateID: template.ID,
		BasketName:       template.Name,
		HorizonDays:      filter.HorizonDays,
		RequiredBaskets:  requiredBaskets,
		Items:            items,
	}, nil
}

func (s *Service) CreateEntryMovement(ctx context.Context, input CreateMovementInput) (*domain.Movement, error) {
	quantity, err := domain.ValidateQuantityString(input.Quantity)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.ProductID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "productId is required", map[string]string{"field": "productId"})
	}
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}
	source := input.SourceID
	now := time.Now()
	item := domain.Movement{ID: uuid.NewString(), Type: domain.MovementTypeEntry, ProductID: input.ProductID, Quantity: quantity, MovementDate: input.MovementDate, OriginID: &source, Note: strings.TrimSpace(input.Note), CreatedAt: now}
	if err := s.repo.CreateMovementEntry(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) CreateEntryMovementsBatch(ctx context.Context, input CreateMovementBatchInput) ([]domain.Movement, error) {
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}
	if len(input.Items) == 0 {
		return nil, domain.NewAppError(domain.ErrorValidation, "items are required", map[string]string{"field": "items"})
	}

	movements := make([]domain.Movement, 0, len(input.Items))
	now := time.Now()
	note := strings.TrimSpace(input.Note)

	for index, raw := range input.Items {
		quantity, err := domain.ValidateQuantityString(raw.Quantity)
		if err != nil {
			return nil, domain.NewAppError(domain.ErrorValidation, "invalid quantity in items", map[string]any{"index": index, "error": err.Error()})
		}
		if strings.TrimSpace(raw.ProductID) == "" {
			return nil, domain.NewAppError(domain.ErrorValidation, "productId is required in items", map[string]any{"index": index, "field": "productId"})
		}
		source := input.SourceID
		movements = append(movements, domain.Movement{
			ID:           uuid.NewString(),
			Type:         domain.MovementTypeEntry,
			ProductID:    raw.ProductID,
			Quantity:     quantity,
			MovementDate: input.MovementDate,
			OriginID:     &source,
			Note:         note,
			CreatedAt:    now,
		})
	}

	if err := s.repo.CreateMovementEntriesBatch(ctx, movements); err != nil {
		return nil, mapRepoError(err)
	}
	return movements, nil
}

func (s *Service) CreateExitMovement(ctx context.Context, input CreateMovementInput) (*domain.Movement, error) {
	quantity, err := domain.ValidateQuantityString(input.Quantity)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.ProductID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "productId is required", map[string]string{"field": "productId"})
	}
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}

	source := input.SourceID
	now := time.Now()
	item := domain.Movement{ID: uuid.NewString(), Type: domain.MovementTypeExit, ProductID: input.ProductID, Quantity: quantity, MovementDate: input.MovementDate, DestinationID: &source, Note: strings.TrimSpace(input.Note), CreatedAt: now}
	if err := s.repo.CreateMovementExit(ctx, item); err != nil {
		return nil, mapRepoError(err)
	}
	return &item, nil
}

func (s *Service) CreateExitMovementsBatch(ctx context.Context, input CreateMovementBatchInput) ([]domain.Movement, error) {
	if strings.TrimSpace(input.SourceID) == "" {
		return nil, domain.NewAppError(domain.ErrorValidation, "sourceId is required", map[string]string{"field": "sourceId"})
	}
	if dateErr := domain.ValidateDateNotFuture(input.MovementDate); dateErr != nil {
		return nil, dateErr
	}
	if len(input.Items) == 0 {
		return nil, domain.NewAppError(domain.ErrorValidation, "items are required", map[string]string{"field": "items"})
	}

	movements := make([]domain.Movement, 0, len(input.Items))
	now := time.Now()
	note := strings.TrimSpace(input.Note)

	for index, raw := range input.Items {
		quantity, err := domain.ValidateQuantityString(raw.Quantity)
		if err != nil {
			return nil, domain.NewAppError(domain.ErrorValidation, "invalid quantity in items", map[string]any{"index": index, "error": err.Error()})
		}
		if strings.TrimSpace(raw.ProductID) == "" {
			return nil, domain.NewAppError(domain.ErrorValidation, "productId is required in items", map[string]any{"index": index, "field": "productId"})
		}
		source := input.SourceID
		movements = append(movements, domain.Movement{
			ID:            uuid.NewString(),
			Type:          domain.MovementTypeExit,
			ProductID:     raw.ProductID,
			Quantity:      quantity,
			MovementDate:  input.MovementDate,
			DestinationID: &source,
			Note:          note,
			CreatedAt:     now,
		})
	}

	if err := s.repo.CreateMovementExitsBatch(ctx, movements); err != nil {
		return nil, mapRepoError(err)
	}
	return movements, nil
}

func (s *Service) ListMovements(ctx context.Context, filter ListMovementsFilter) ([]domain.Movement, error) {
	if filter.DateFrom != nil && filter.DateTo != nil && filter.DateTo.Before(*filter.DateFrom) {
		return nil, domain.NewAppError(domain.ErrorValidation, "dateTo must be after dateFrom", nil)
	}
	return s.repo.ListMovements(ctx, filter)
}

func (s *Service) GetCurrentStock(ctx context.Context, filter CurrentStockFilter) ([]domain.StockItem, error) {
	return s.repo.GetCurrentStock(ctx, filter)
}

func (s *Service) GetReports(ctx context.Context, filter ReportsFilter) (*domain.Reports, error) {
	if filter.DateTo.Before(filter.DateFrom) {
		return nil, domain.NewAppError(domain.ErrorValidation, "dateTo must be after dateFrom", nil)
	}
	entriesByPeriod, err := s.repo.GetEntriesByPeriod(ctx, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, mapRepoError(err)
	}
	entriesByOrigin, err := s.repo.GetEntriesByOrigin(ctx, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, mapRepoError(err)
	}
	exitsByPeriod, err := s.repo.GetExitsByPeriod(ctx, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, mapRepoError(err)
	}
	exitsByDestination, err := s.repo.GetExitsByDestination(ctx, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, mapRepoError(err)
	}
	history, err := s.repo.GetMovementHistory(ctx, filter.DateFrom, filter.DateTo)
	if err != nil {
		return nil, mapRepoError(err)
	}

	return &domain.Reports{
		EntriesByPeriod:    entriesByPeriod,
		EntriesByOrigin:    entriesByOrigin,
		ExitsByPeriod:      exitsByPeriod,
		ExitsByDestination: exitsByDestination,
		MovementHistory:    history,
	}, nil
}

func (s *Service) ExportBackup(ctx context.Context, destination string) error {
	if strings.TrimSpace(destination) == "" {
		return domain.NewAppError(domain.ErrorValidation, "destination path is required", map[string]string{"field": "destination"})
	}
	if err := s.backup.ExportBackup(ctx, destination); err != nil {
		return domain.NewAppError(domain.ErrorInternal, "failed to export backup", err.Error())
	}
	return nil
}

func (s *Service) ImportBackup(ctx context.Context, input ImportBackupInput) error {
	if strings.TrimSpace(input.FilePath) == "" {
		return domain.NewAppError(domain.ErrorValidation, "filePath is required", map[string]string{"field": "filePath"})
	}
	if !input.ConfirmReplace || !input.ConfirmUnderstand {
		return domain.NewAppError(domain.ErrorValidation, "import requires both confirmations", nil)
	}
	if err := s.backup.ImportBackup(ctx, input.FilePath); err != nil {
		return domain.NewAppError(domain.ErrorInternal, "failed to import backup", err.Error())
	}
	return nil
}

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*domain.AppError); ok {
		return appErr
	}
	return domain.NewAppError(domain.ErrorInternal, "unexpected error", err.Error())
}

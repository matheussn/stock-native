package application

import (
	"context"
	"time"

	"stock/backend/internal/domain"
)

type Repository interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, item domain.Category) error
	UpdateCategory(ctx context.Context, item domain.Category) error
	DeactivateCategory(ctx context.Context, id string) error

	ListProducts(ctx context.Context) ([]domain.Product, error)
	CreateProduct(ctx context.Context, item domain.Product) error
	UpdateProduct(ctx context.Context, item domain.Product) error
	DeactivateProduct(ctx context.Context, id string) error

	ListOrigins(ctx context.Context) ([]domain.Origin, error)
	CreateOrigin(ctx context.Context, item domain.Origin) error
	UpdateOrigin(ctx context.Context, item domain.Origin) error
	DeactivateOrigin(ctx context.Context, id string) error

	ListDestinations(ctx context.Context) ([]domain.Destination, error)
	CreateDestination(ctx context.Context, item domain.Destination) error
	UpdateDestination(ctx context.Context, item domain.Destination) error
	DeactivateDestination(ctx context.Context, id string) error

	CreateMovementEntry(ctx context.Context, item domain.Movement) error
	CreateMovementExit(ctx context.Context, item domain.Movement) error
	CreateMovementEntriesBatch(ctx context.Context, items []domain.Movement) error
	CreateMovementExitsBatch(ctx context.Context, items []domain.Movement) error
	ListMovements(ctx context.Context, filter ListMovementsFilter) ([]domain.Movement, error)

	GetCurrentStock(ctx context.Context, filter CurrentStockFilter) ([]domain.StockItem, error)
	GetCurrentStockByProduct(ctx context.Context, productID string) (float64, error)

	GetEntriesByPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error)
	GetEntriesByOrigin(ctx context.Context, from time.Time, to time.Time) ([]domain.GroupedTotal, error)
	GetExitsByPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error)
	GetExitsByDestination(ctx context.Context, from time.Time, to time.Time) ([]domain.GroupedTotal, error)
	GetMovementHistory(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error)

	ListBasketTemplates(ctx context.Context) ([]domain.BasketTemplate, error)
	GetBasketTemplate(ctx context.Context, id string) (*domain.BasketTemplate, error)
	CreateBasketTemplate(ctx context.Context, item domain.BasketTemplate) error
	UpdateBasketTemplate(ctx context.Context, item domain.BasketTemplate) error
	DeactivateBasketTemplate(ctx context.Context, id string) error

	ListFamilies(ctx context.Context) ([]domain.Family, error)
	CreateFamily(ctx context.Context, item domain.Family) error
	UpdateFamily(ctx context.Context, item domain.Family) error
	DeactivateFamily(ctx context.Context, id string) error
	GetActiveFamiliesCestas(ctx context.Context) (float64, error)
}

package domain

import "time"

type ErrorCode string

const (
	ErrorValidation        ErrorCode = "VALIDATION_ERROR"
	ErrorNotFound          ErrorCode = "NOT_FOUND"
	ErrorInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrorConflict          ErrorCode = "CONFLICT"
	ErrorInternal          ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code ErrorCode, message string, details any) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Product struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	CategoryID        string      `json:"categoryId"`
	MeasureUnit       MeasureUnit `json:"measureUnit"`
	PackageAmount     float64     `json:"packageAmount"`
	LowStockThreshold float64     `json:"lowStockThreshold"`
	Active            bool        `json:"active"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

type MeasureUnit string

const (
	MeasureUnitKg      MeasureUnit = "kg"
	MeasureUnitG       MeasureUnit = "g"
	MeasureUnitL       MeasureUnit = "l"
	MeasureUnitMl      MeasureUnit = "ml"
	MeasureUnitUnidade MeasureUnit = "unidade"
)

type Origin struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Destination struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MovementType string

const (
	MovementTypeEntry MovementType = "entrada"
	MovementTypeExit  MovementType = "saida"
)

type Movement struct {
	ID              string       `json:"id"`
	Type            MovementType `json:"type"`
	ProductID       string       `json:"productId"`
	ProductName     string       `json:"productName,omitempty"`
	Quantity        float64      `json:"quantity"`
	MovementDate    time.Time    `json:"movementDate"`
	OriginID        *string      `json:"originId,omitempty"`
	OriginName      *string      `json:"originName,omitempty"`
	DestinationID   *string      `json:"destinationId,omitempty"`
	DestinationName *string      `json:"destinationName,omitempty"`
	Note            string       `json:"note"`
	CreatedAt       time.Time    `json:"createdAt"`
}

type StockItem struct {
	ProductID         string      `json:"productId"`
	ProductName       string      `json:"productName"`
	CategoryID        string      `json:"categoryId"`
	CategoryName      string      `json:"categoryName"`
	MeasureUnit       MeasureUnit `json:"measureUnit"`
	PackageAmount     float64     `json:"packageAmount"`
	LowStockThreshold float64     `json:"lowStockThreshold"`
	CurrentStock      float64     `json:"currentStock"`
	IsLowStock        bool        `json:"isLowStock"`
}

type GroupedTotal struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

type Reports struct {
	EntriesByPeriod    []Movement     `json:"entriesByPeriod"`
	EntriesByOrigin    []GroupedTotal `json:"entriesByOrigin"`
	ExitsByPeriod      []Movement     `json:"exitsByPeriod"`
	ExitsByDestination []GroupedTotal `json:"exitsByDestination"`
	MovementHistory    []Movement     `json:"movementHistory"`
}

type BasketTemplate struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Active    bool                 `json:"active"`
	Items     []BasketTemplateItem `json:"items"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

type BasketTemplateItem struct {
	ID                string      `json:"id"`
	BasketTemplateID  string      `json:"basketTemplateId"`
	ProductID         string      `json:"productId"`
	ProductName       string      `json:"productName,omitempty"`
	MeasureUnit       MeasureUnit `json:"measureUnit"`
	QuantityPerBasket float64     `json:"quantityPerBasket"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

type Family struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	CestasPerPeriod float64   `json:"cestasPerPeriod"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type StockPlanningStatus string

const (
	StockPlanningSufficient StockPlanningStatus = "suficiente"
	StockPlanningAttention  StockPlanningStatus = "atencao"
	StockPlanningCritical   StockPlanningStatus = "critico"
)

type StockPlanningItem struct {
	ProductID        string              `json:"productId"`
	ProductName      string              `json:"productName"`
	MeasureUnit      MeasureUnit         `json:"measureUnit"`
	CurrentStock     float64             `json:"currentStock"`
	RequiredStock    float64             `json:"requiredStock"`
	ProjectedBalance float64             `json:"projectedBalance"`
	Status           StockPlanningStatus `json:"status"`
}

type StockPlanning struct {
	BasketTemplateID string              `json:"basketTemplateId"`
	BasketName       string              `json:"basketName"`
	HorizonDays      int                 `json:"horizonDays"`
	RequiredBaskets  float64             `json:"requiredBaskets"`
	Items            []StockPlanningItem `json:"items"`
}

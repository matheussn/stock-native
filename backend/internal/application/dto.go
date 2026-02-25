package application

import "time"

type CreateCategoryInput struct {
	Name string `json:"name"`
}

type UpdateCategoryInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateProductInput struct {
	Name              string `json:"name"`
	CategoryID        string `json:"categoryId"`
	MeasureUnit       string `json:"measureUnit"`
	PackageAmount     string `json:"packageAmount"`
	LowStockThreshold string `json:"lowStockThreshold"`
}

type UpdateProductInput struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	CategoryID        string `json:"categoryId"`
	MeasureUnit       string `json:"measureUnit"`
	PackageAmount     string `json:"packageAmount"`
	LowStockThreshold string `json:"lowStockThreshold"`
}

type CreateOriginInput struct {
	Name string `json:"name"`
}

type UpdateOriginInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateDestinationInput struct {
	Name string `json:"name"`
}

type UpdateDestinationInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateMovementInput struct {
	ProductID    string    `json:"productId"`
	Quantity     string    `json:"quantity"`
	MovementDate time.Time `json:"movementDate"`
	SourceID     string    `json:"sourceId"`
	Note         string    `json:"note"`
}

type MovementBatchItemInput struct {
	ProductID string `json:"productId"`
	Quantity  string `json:"quantity"`
}

type CreateMovementBatchInput struct {
	MovementDate time.Time                `json:"movementDate"`
	SourceID     string                   `json:"sourceId"`
	Note         string                   `json:"note"`
	Items        []MovementBatchItemInput `json:"items"`
}

type ListMovementsFilter struct {
	Type      string     `json:"type"`
	ProductID string     `json:"productId"`
	DateFrom  *time.Time `json:"dateFrom"`
	DateTo    *time.Time `json:"dateTo"`
}

type CurrentStockFilter struct {
	CategoryID string `json:"categoryId"`
	Search     string `json:"search"`
}

type ReportsFilter struct {
	DateFrom time.Time `json:"dateFrom"`
	DateTo   time.Time `json:"dateTo"`
}

type ImportBackupInput struct {
	FilePath          string `json:"filePath"`
	ConfirmReplace    bool   `json:"confirmReplace"`
	ConfirmUnderstand bool   `json:"confirmUnderstand"`
}

type BasketTemplateItemInput struct {
	ProductID         string `json:"productId"`
	QuantityPerBasket string `json:"quantityPerBasket"`
}

type CreateBasketTemplateInput struct {
	Name  string                    `json:"name"`
	Items []BasketTemplateItemInput `json:"items"`
}

type UpdateBasketTemplateInput struct {
	ID    string                    `json:"id"`
	Name  string                    `json:"name"`
	Items []BasketTemplateItemInput `json:"items"`
}

type CreateFamilyInput struct {
	Name            string `json:"name"`
	CestasPerPeriod string `json:"cestasPerPeriod"`
}

type UpdateFamilyInput struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CestasPerPeriod string `json:"cestasPerPeriod"`
}

type BasketSubstitutionInput struct {
	OriginalProductID string `json:"originalProductId"`
	ProductID         string `json:"productId"`
	QuantityPerBasket string `json:"quantityPerBasket"`
}

type CreateBasketMovementInput struct {
	BasketTemplateID string                    `json:"basketTemplateId"`
	BasketsCount     string                    `json:"basketsCount"`
	SourceID         string                    `json:"sourceId"`
	MovementDate     time.Time                 `json:"movementDate"`
	Note             string                    `json:"note"`
	Substitutions    []BasketSubstitutionInput `json:"substitutions"`
}

type StockPlanningFilter struct {
	BasketTemplateID string `json:"basketTemplateId"`
	HorizonDays      int    `json:"horizonDays"`
}

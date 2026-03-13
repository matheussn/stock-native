package models

type ProductVariation struct {
	ID           int64  `json:"id"`
	ProductID    int64  `json:"product_id"`
	Description  string `json:"description"`
	BaseQuantity int64  `json:"base_quantity"`
	CurrentStock int64  `json:"current_stock"`
	IsActive     int64  `json:"is_active"`
	CreatedAt    string `json:"created_at"`
}

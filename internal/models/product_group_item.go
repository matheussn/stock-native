package models

type ProductGroupItem struct {
	ID             int64 `json:"id"`
	ProductGroupID int64 `json:"product_group_id"`
	ProductID      int64 `json:"product_id"`
	BaseQuantity   int64 `json:"base_quantity"`
}

package models

type MovementProductItem struct {
	ID                 int64 `json:"id"`
	MovementID         int64 `json:"movement_id"`
	ProductVariationID int64 `json:"product_variation_id"`
	Quantity           int64 `json:"quantity"`
}

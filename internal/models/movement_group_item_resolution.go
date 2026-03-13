package models

type MovementGroupItemResolution struct {
	ID                  int64 `json:"id"`
	MovementGroupItemID int64 `json:"movement_group_item_id"`
	ProductVariationID  int64 `json:"product_variation_id"`
	Quantity            int64 `json:"quantity"`
}

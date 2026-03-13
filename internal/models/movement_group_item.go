package models

type MovementGroupItem struct {
	ID             int64 `json:"id"`
	MovementID     int64 `json:"movement_id"`
	ProductGroupID int64 `json:"product_group_id"`
	Quantity       int64 `json:"quantity"`
}

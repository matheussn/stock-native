package models

type ProductGroup struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    int64  `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

package models

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	BaseUnit    string `json:"base_unit"`
	Description string `json:"description"`
	IsActive    int64  `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

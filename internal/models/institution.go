package models

type Institution struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	CNPJ            string `json:"cnpj"`
	ResponsibleName string `json:"responsible_name"`
	Phone           string `json:"phone"`
	IsActive        int64  `json:"is_active"`
	CreatedAt       string `json:"created_at"`
}

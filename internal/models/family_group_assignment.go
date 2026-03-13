package models

type FamilyGroupAssignment struct {
	ID             int64  `json:"id"`
	FamilyID       int64  `json:"family_id"`
	ProductGroupID int64  `json:"product_group_id"`
	StartedAt      string `json:"started_at"`
	EndedAt        string `json:"ended_at"`
}

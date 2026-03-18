package models

type Movement struct {
	ID                 int64  `json:"id"`
	AssistentialWorkID int64  `json:"assistential_work_id"`
	InstitutionID      int64  `json:"institution_id"`
	Type               string `json:"type"`
	Notes              string `json:"notes"`
	CreatedAt          string `json:"created_at"`
}

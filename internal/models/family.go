package models

type Family struct {
	ID                 int64  `json:"id"`
	AssistentialWorkID int64  `json:"assistential_work_id"`
	Name               string `json:"name"`
	MemberCount        int64  `json:"member_count"`
	Address            string `json:"address"`
	Contact            string `json:"contact"`
	IsActive           int64  `json:"is_active"`
	CreatedAt          string `json:"created_at"`
}

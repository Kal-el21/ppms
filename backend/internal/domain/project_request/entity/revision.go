package entity

import "time"

type ProjectRequestRevision struct {
	ID               uint64 `gorm:"primaryKey" json:"id"`
	ProjectRequestID uint64 `gorm:"not null" json:"project_request_id"`

	RevisionNumber int `gorm:"not null" json:"revision_number"`

	Title           string  `json:"title"`
	Description     string  `json:"description"`
	BusinessGoal    string  `json:"business_goal"`
	ExpectedOutcome string  `json:"expected_outcome"`
	EstimatedBudget float64 `json:"estimated_budget"`

	RevisionReason string `json:"revision_reason"`

	RevisedBy uint64 `gorm:"not null" json:"revised_by"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProjectRequestRevision) TableName() string {
	return "project_request_revisions"
}

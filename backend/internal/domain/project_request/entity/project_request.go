package entity

import "time"

type RequestStatus string

const (
	StatusDraft       RequestStatus = "DRAFT"
	StatusSubmitted   RequestStatus = "SUBMITTED"
	StatusUnderReview RequestStatus = "UNDER_REVIEW"
	StatusApproved    RequestStatus = "APPROVED"
	StatusRejected    RequestStatus = "REJECTED"
	StatusRevised     RequestStatus = "REVISED"
)

type ProjectRequest struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	RequesterID uint64 `gorm:"not null" json:"requester_id"`

	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`

	BusinessGoal    string  `json:"business_goal"`
	ExpectedOutcome string  `json:"expected_outcome"`
	EstimatedBudget float64 `json:"estimated_budget"`

	Status RequestStatus `gorm:"not null;default:DRAFT" json:"status"`

	SubmittedAt *time.Time `json:"submitted_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (ProjectRequest) TableName() string {
	return "project_requests"
}

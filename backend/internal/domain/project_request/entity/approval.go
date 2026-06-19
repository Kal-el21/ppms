package entity

import "time"

type ApprovalAction string

const (
	ActionApproved        ApprovalAction = "APPROVED"
	ActionRejected        ApprovalAction = "REJECTED"
	ActionRequestRevision ApprovalAction = "REQUEST_REVISION"
)

type ProjectRequestApproval struct {
	ID               uint64  `gorm:"primaryKey" json:"id"`
	ProjectRequestID uint64  `gorm:"not null" json:"project_request_id"`
	RevisionID       *uint64 `json:"revision_id,omitempty"`

	ReviewedBy uint64         `gorm:"not null" json:"reviewed_by"`
	Action     ApprovalAction `gorm:"not null" json:"action"`
	Comment    string         `json:"comment"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProjectRequestApproval) TableName() string {
	return "project_request_approvals"
}

package entity

import "time"

type ApprovalWorkflow struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ApprovalWorkflow) TableName() string {
	return "approval_workflows"
}

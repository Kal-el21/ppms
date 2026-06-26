package entity

import "time"

type ApprovalLevel struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	WorkflowID  uint64 `gorm:"not null" json:"workflow_id"`
	LevelOrder  int    `gorm:"not null" json:"level_order"`
	RoleRequired string `gorm:"not null" json:"role_required"`

	CreatedAt time.Time `json:"created_at"`
}

func (ApprovalLevel) TableName() string {
	return "approval_levels"
}

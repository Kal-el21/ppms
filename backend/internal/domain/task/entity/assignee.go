package entity

import "time"

type TaskAssignee struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	TaskID uint64 `gorm:"not null" json:"task_id"`
	UserID uint64 `gorm:"not null" json:"user_id"`

	AssignedBy uint64    `gorm:"not null" json:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at"`

	UnassignedAt *time.Time `json:"unassigned_at,omitempty"`
	UnassignedBy *uint64    `json:"unassigned_by,omitempty"`
}

func (TaskAssignee) TableName() string {
	return "task_assignees"
}

package entity

import "time"

type TaskComment struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	TaskID uint64 `gorm:"not null" json:"task_id"`
	UserID uint64 `gorm:"not null" json:"user_id"`

	Comment string `gorm:"not null" json:"comment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TaskComment) TableName() string {
	return "task_comments"
}

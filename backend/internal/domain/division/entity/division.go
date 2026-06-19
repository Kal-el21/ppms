package entity

import "time"

type Division struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`
}

func (Division) TableName() string {
	return "divisions"
}

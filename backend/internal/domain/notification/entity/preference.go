package entity

import "time"

type NotificationPreference struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`

	Type    string `gorm:"not null" json:"type"`
	Enabled bool   `gorm:"default:true" json:"enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

package entity

import "time"

type SystemRole string

const (
	RoleAdmin  SystemRole = "ADMIN"
	RoleUser   SystemRole = "USER"
	RoleViewer SystemRole = "VIEWER"
)

type User struct {
	ID         uint64  `gorm:"primaryKey" json:"id"`
	DivisionID *uint64 `json:"division_id"`

	FullName     string `gorm:"not null" json:"full_name"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`

	SystemRole SystemRole `gorm:"not null" json:"system_role"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`

	// Tambahkan ke struct User:
	ProfilePhotoURL          *string `json:"profile_photo_url,omitempty"`
	TwoFAEnabled             bool    `gorm:"default:false" json:"two_fa_enabled"`
	EmailNotificationEnabled bool    `gorm:"default:true" json:"email_notification_enabled"`
}

func (User) TableName() string {
	return "users"
}

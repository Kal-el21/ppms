package entity

import "time"

type UserSession struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`

	RefreshTokenHash string `gorm:"not null" json:"-"`
	DeviceInfo       string `json:"device_info"`
	IPAddress        string `json:"ip_address"`

	ExpiresAt     time.Time  `json:"expires_at"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokedReason string     `json:"revoked_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

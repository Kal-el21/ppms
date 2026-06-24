package entity

import "time"

type OTPPurpose string

const (
	OTPLogin         OTPPurpose = "LOGIN"
	OTPPasswordReset OTPPurpose = "PASSWORD_RESET"
)

type OTPToken struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	UserID    uint64     `gorm:"not null" json:"user_id"`
	TokenHash string     `gorm:"not null" json:"-"`
	Purpose   OTPPurpose `gorm:"not null" json:"purpose"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	IPAddress string     `json:"ip_address"`
	CreatedAt time.Time  `json:"created_at"`
}

func (OTPToken) TableName() string {
	return "otp_tokens"
}

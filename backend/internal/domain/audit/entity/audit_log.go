package entity

import "time"

type AuditLog struct {
	ID     uint64  `gorm:"primaryKey" json:"id"`
	UserID *uint64 `json:"user_id"`

	Module string `gorm:"not null" json:"module"`
	Action string `gorm:"not null" json:"action"`

	EntityType string  `json:"entity_type"`
	EntityID   *uint64 `json:"entity_id"`

	OldData string `gorm:"type:json" json:"old_data,omitempty"`
	NewData string `gorm:"type:json" json:"new_data,omitempty"`

	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	RequestID string `json:"request_id"`

	CreatedAt time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

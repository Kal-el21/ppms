package entity

import "time"

type HandoverStatus string

const (
	StatusPending  HandoverStatus = "PENDING"
	StatusReceived HandoverStatus = "RECEIVED"
	StatusReturned HandoverStatus = "RETURNED"
)

type Handover struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	ProjectID uint64 `gorm:"not null" json:"project_id"`

	SenderID         uint64  `gorm:"not null" json:"sender_id"`
	SenderDivisionID *uint64 `json:"sender_division_id"`

	ReceiverID *uint64 `json:"receiver_id"`

	Description string `json:"description"`

	DeliveryDate *time.Time `json:"delivery_date"`
	DeliveryTime *string    `json:"delivery_time"` // disimpan sebagai string "HH:MM:SS" untuk kompatibilitas Postgres TIME

	ReceivedAt *time.Time `json:"received_at"`

	Status HandoverStatus `gorm:"not null;default:PENDING" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (Handover) TableName() string {
	return "handovers"
}

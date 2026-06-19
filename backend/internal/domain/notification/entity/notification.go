package entity

import "time"

type Channel string
type DeliveryStatus string

const (
	ChannelInApp Channel = "IN_APP"
	ChannelEmail Channel = "EMAIL"

	DeliveryPending DeliveryStatus = "PENDING"
	DeliverySent    DeliveryStatus = "SENT"
	DeliveryFailed  DeliveryStatus = "FAILED"
)

type Notification struct {
	ID     uint64 `gorm:"primaryKey" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`

	Type string `gorm:"not null" json:"type"`

	Title   string `gorm:"not null" json:"title"`
	Message string `json:"message"`

	EntityType string  `json:"entity_type"`
	EntityID   *uint64 `json:"entity_id"`

	ActionURL string `json:"action_url"`

	Channel        Channel        `gorm:"not null;default:IN_APP" json:"channel"`
	DeliveryStatus DeliveryStatus `gorm:"not null;default:SENT" json:"delivery_status"`

	IsRead bool `gorm:"default:false" json:"is_read"`

	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

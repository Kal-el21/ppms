package dto

import "time"

type CreateHandoverRequest struct {
	ReceiverID       *uint64    `json:"receiver_id"`
	Description      string     `json:"description" validate:"max=1000"`
	DeliveryDate     *time.Time `json:"delivery_date"`
	DeliveryTime     string     `json:"delivery_time" validate:"omitempty"` // format "HH:MM:SS" atau "HH:MM"
	SenderDivisionID *uint64    `json:"sender_division_id"`
}

type MarkReceivedRequest struct {
	Version int `json:"version" validate:"required"`
}

type MarkReturnedRequest struct {
	Reason  string `json:"reason" validate:"required,min=3,max=500"`
	Version int    `json:"version" validate:"required"`
}

type HandoverResponse struct {
	ID               uint64     `json:"id"`
	ProjectID        uint64     `json:"project_id"`
	SenderID         uint64     `json:"sender_id"`
	SenderDivisionID *uint64    `json:"sender_division_id"`
	ReceiverID       *uint64    `json:"receiver_id"`
	Description      string     `json:"description"`
	DeliveryDate     *time.Time `json:"delivery_date"`
	DeliveryTime     *string    `json:"delivery_time"`
	ReceivedAt       *time.Time `json:"received_at"`
	Status           string     `json:"status"`
	Version          int        `json:"version"`
	CreatedAt        time.Time  `json:"created_at"`
}

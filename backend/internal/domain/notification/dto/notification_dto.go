package dto

import "time"

type NotificationResponse struct {
	ID         uint64    `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	EntityType string    `json:"entity_type"`
	EntityID   *uint64   `json:"entity_id"`
	ActionURL  string    `json:"action_url"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdatePreferenceRequest struct {
	Type    string `json:"type" validate:"required"`
	Enabled bool   `json:"enabled"`
}

type PreferenceResponse struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

package entity

import "time"

type EntityType string

const (
	EntityProjectRequest    EntityType = "PROJECT_REQUEST"
	EntityProject           EntityType = "PROJECT"
	EntityMilestone         EntityType = "MILESTONE"
	EntityTask              EntityType = "TASK"
	EntityBudgetTransaction EntityType = "BUDGET_TRANSACTION"
	EntityHandover          EntityType = "HANDOVER"
)

type Attachment struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	UploadedBy uint64 `gorm:"not null" json:"uploaded_by"`

	EntityType EntityType `gorm:"not null" json:"entity_type"`
	EntityID   uint64     `gorm:"not null" json:"entity_id"`

	Version int `gorm:"default:1" json:"version"`

	FileName     string `gorm:"not null" json:"file_name"`
	OriginalName string `gorm:"not null" json:"original_name"`

	FilePath string `gorm:"not null" json:"file_path"`
	FileSize int64  `gorm:"not null" json:"file_size"`
	MimeType string `gorm:"not null" json:"mime_type"`

	CreatedAt time.Time `json:"created_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`
}

func (Attachment) TableName() string {
	return "attachments"
}

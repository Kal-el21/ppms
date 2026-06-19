package entity

import "time"

type TaskPriority string
type TaskStatus string

const (
	PriorityLow    TaskPriority = "LOW"
	PriorityMedium TaskPriority = "MEDIUM"
	PriorityHigh   TaskPriority = "HIGH"
	PriorityUrgent TaskPriority = "URGENT"

	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
	StatusCancelled  TaskStatus = "CANCELLED"
)

type Task struct {
	ID          uint64  `gorm:"primaryKey" json:"id"`
	ProjectID   uint64  `gorm:"not null" json:"project_id"`
	MilestoneID *uint64 `json:"milestone_id"`

	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`

	Priority TaskPriority `gorm:"not null;default:MEDIUM" json:"priority"`
	Status   TaskStatus   `gorm:"not null;default:TODO" json:"status"`

	Progress int `gorm:"default:0" json:"progress"`

	OrderIndex int `gorm:"default:0" json:"order_index"`

	StartDate *time.Time `json:"start_date"`
	DueDate   *time.Time `json:"due_date"`

	CreatedBy uint64 `gorm:"not null" json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	DeletedBy *uint64    `json:"deleted_by,omitempty"`

	Version int `gorm:"default:1" json:"version"`
}

func (Task) TableName() string {
	return "tasks"
}

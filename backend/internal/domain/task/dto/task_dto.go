package dto

import "time"

type CreateTaskRequest struct {
	MilestoneID *uint64    `json:"milestone_id"`
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	Priority    string     `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH URGENT"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	Priority    string     `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH URGENT"`
	MilestoneID *uint64    `json:"milestone_id"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
	Version     int        `json:"version" validate:"required"`
}

type ChangeTaskStatusRequest struct {
	Status  string `json:"status" validate:"required,oneof=TODO IN_PROGRESS DONE CANCELLED"`
	Version int    `json:"version" validate:"required"`
}

type UpdateProgressRequest struct {
	Progress int `json:"progress" validate:"min=0,max=100"`
	Version  int `json:"version" validate:"required"`
}

type AssignUserRequest struct {
	UserIDs []uint64 `json:"user_ids" validate:"required,min=1"`
}

type CreateCommentRequest struct {
	Comment string `json:"comment" validate:"required,min=1,max=2000"`
}

type TaskResponse struct {
	ID          uint64     `json:"id"`
	ProjectID   uint64     `json:"project_id"`
	MilestoneID *uint64    `json:"milestone_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Version     int        `json:"version"`
	Assignees   []uint64   `json:"assignee_ids,omitempty"`
}

package dto

import "time"

type UpdateProjectRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Version     int        `json:"version" validate:"required"`
}

type ChangeStatusRequest struct {
	Status  string `json:"status" validate:"required,oneof=PLANNED ACTIVE ON_HOLD COMPLETED CANCELLED"`
	Version int    `json:"version" validate:"required"`
}

type AddMemberRequest struct {
	UserID      uint64 `json:"user_id" validate:"required"`
	ProjectRole string `json:"project_role" validate:"required,oneof=PROJECT_MANAGER MEMBER OBSERVER"`
}

type ChangeMemberRoleRequest struct {
	ProjectRole string `json:"project_role" validate:"required,oneof=PROJECT_MANAGER MEMBER OBSERVER"`
}

type ProjectResponse struct {
	ID          uint64     `json:"id"`
	ProjectCode string     `json:"project_code,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type MemberResponse struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	ProjectRole string    `json:"project_role"`
	Status      string    `json:"status"`
	JoinedAt    time.Time `json:"joined_at"`
}

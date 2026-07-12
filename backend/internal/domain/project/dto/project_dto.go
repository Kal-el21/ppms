package dto

import "time"

type UpdateProjectRequest struct {
	Name           string     `json:"name" validate:"required,min=3,max=200"`
	Description    string     `json:"description" validate:"max=2000"`
	Category       string     `json:"category" validate:"max=100"`
	InitiationType *string    `json:"initiation_type" validate:"omitempty,oneof=NEW_INITIATIVE RENEWAL ENHANCEMENT"`
	Priority       string     `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	Notes          string     `json:"notes" validate:"max=2000"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Version        int        `json:"version" validate:"required"`
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

type CreateProjectDirectRequest struct {
	Name           string     `json:"name" validate:"required,min=3,max=200"`
	Description    string     `json:"description" validate:"max=2000"`
	Category       string     `json:"category" validate:"max=100"`
	InitiationType *string    `json:"initiation_type" validate:"omitempty,oneof=NEW_INITIATIVE RENEWAL ENHANCEMENT"`
	Priority       string     `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	Notes          string     `json:"notes" validate:"max=2000"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	BudgetType     string     `json:"budget_type" validate:"omitempty,oneof=CAPEX OPEX"`
	BudgetName     string     `json:"budget_name" validate:"max=200"`
	AllocatedBudget float64  `json:"allocated_budget" validate:"omitempty,gt=0"`
}

type ProjectResponse struct {
	ID               uint64     `json:"id"`
	ProjectRequestID *uint64    `json:"project_request_id,omitempty"`
	ProjectCode      string     `json:"project_code,omitempty"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Category         string     `json:"category"`
	InitiationType   *string    `json:"initiation_type,omitempty"`
	Priority         string     `json:"priority"`
	Notes            string     `json:"notes"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Status           string     `json:"status"`
	Progress         float64    `json:"progress"`
	Health           string     `json:"health"`
	BudgetType       *string    `json:"budget_type,omitempty"`
	BudgetAllocated  float64    `json:"budget_allocated,omitempty"`
	BudgetUsed       float64    `json:"budget_used,omitempty"`
	Version          int        `json:"version"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type ProjectDeadline struct {
	ID            uint64     `json:"id"`
	ProjectCode   string     `json:"project_code"`
	Name          string     `json:"name"`
	EndDate       *time.Time `json:"end_date"`
	Status        string     `json:"status"`
	Progress      float64    `json:"progress"`
	DaysRemaining int        `json:"days_remaining"`
}

type MemberResponse struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	ProjectRole string    `json:"project_role"`
	Status      string    `json:"status"`
	JoinedAt    time.Time `json:"joined_at"`
}

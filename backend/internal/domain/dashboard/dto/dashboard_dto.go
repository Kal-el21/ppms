package dto

import "time"

type DashboardSummary struct {
	TotalProjects     int64            `json:"total_projects"`
	ActiveProjects    int64            `json:"active_projects"`
	OnHoldProjects    int64            `json:"on_hold_projects"`
	CompletedProjects int64            `json:"completed_projects"`
	OverdueProjects   int64            `json:"overdue_projects"`
	TotalTasks        int64            `json:"total_tasks"`
	CompletedTasks    int64            `json:"completed_tasks"`
	PendingRequests   int64            `json:"pending_requests"`
	OverdueTasks      int64            `json:"overdue_tasks"`
	TotalBudgetUsage  float64          `json:"total_budget_usage_percentage"`

	ByStatus       map[string]int64 `json:"by_status"`
	ByInitiation   map[string]int64 `json:"by_initiation"`
	ByPriority     map[string]int64 `json:"by_priority"`
	HealthGreen  int64            `json:"health_green"`
	HealthYellow int64            `json:"health_yellow"`
	HealthRed    int64            `json:"health_red"`

	CapexAllocated float64 `json:"capex_allocated"`
	CapexUsed      float64 `json:"capex_used"`
	OpexAllocated  float64 `json:"opex_allocated"`
	OpexUsed       float64 `json:"opex_used"`

	UpcomingDeadlines []UpcomingDeadline `json:"upcoming_deadlines"`
	RecentActivities  []RecentActivity   `json:"recent_activities"`

	BudgetMaster    []BudgetMasterItem  `json:"budget_master"`
	Absorption      []AbsorptionItem    `json:"absorption"`
}

type BudgetMasterItem struct {
	BudgetName      string  `json:"budget_name"`
	BudgetType      string  `json:"budget_type"`
	Allocated       float64 `json:"allocated"`
	Used            float64 `json:"used"`
	Remaining       float64 `json:"remaining"`
	UsagePercentage float64 `json:"usage_percentage"`
}

type AbsorptionItem struct {
	ProjectCode     string  `json:"project_code"`
	ProjectName     string  `json:"project_name"`
	BudgetType      string  `json:"budget_type"`
	Allocated       float64 `json:"allocated"`
	Used            float64 `json:"used"`
	UsagePercentage float64 `json:"usage_percentage"`
}

type UpcomingDeadline struct {
	ID            uint64     `json:"id"`
	ProjectCode   string     `json:"project_code"`
	Name          string     `json:"name"`
	EndDate       *time.Time `json:"end_date"`
	Status        string     `json:"status"`
	DaysRemaining int        `json:"days_remaining"`
}

type RecentActivity struct {
	Module     string    `json:"module"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *uint64   `json:"entity_id"`
	UserID     *uint64   `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

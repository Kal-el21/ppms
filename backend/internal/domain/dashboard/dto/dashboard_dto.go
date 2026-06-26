package dto

import "time"

type DashboardSummary struct {
	TotalProjects     int64            `json:"total_projects"`
	ActiveProjects    int64            `json:"active_projects"`
	CompletedProjects int64            `json:"completed_projects"`
	TotalTasks        int64            `json:"total_tasks"`
	CompletedTasks    int64            `json:"completed_tasks"`
	PendingRequests   int64            `json:"pending_requests"`
	OverdueTasks      int64            `json:"overdue_tasks"`
	TotalBudgetUsage  float64          `json:"total_budget_usage_percentage"`
	RecentActivities  []RecentActivity `json:"recent_activities"`
}

type RecentActivity struct {
	Module     string    `json:"module"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *uint64   `json:"entity_id"`
	UserID     *uint64   `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

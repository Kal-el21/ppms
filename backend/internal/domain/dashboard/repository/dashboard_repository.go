package repository

import (
	"github.com/Kal-el21/backend/internal/domain/dashboard/dto"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	CountProjectsByStatus(status string, userID uint64, isAdmin bool) (int64, error)
	CountTotalProjects(userID uint64, isAdmin bool) (int64, error)
	CountPendingRequests(userID uint64, isAdmin bool) (int64, error)
	CountOverdueTasks(userID uint64, isAdmin bool) (int64, error)
	CountTotalTasks(userID uint64, isAdmin bool) (int64, error)
	CountCompletedTasks(userID uint64, isAdmin bool) (int64, error)
	GetAverageBudgetUsage(userID uint64, isAdmin bool) (float64, error)
	GetRecentActivities(limit int) ([]dto.RecentActivity, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// scopeProjectAccess membatasi query ke project yang user ikuti jika bukan ADMIN,
// konsisten dengan permission rule "USER hanya lihat project yang diikuti".
func (r *dashboardRepository) scopeProjectAccess(query *gorm.DB, userID uint64, isAdmin bool) *gorm.DB {
	if isAdmin {
		return query
	}
	return query.Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
}

func (r *dashboardRepository) CountProjectsByStatus(status string, userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Model(&struct{}{}).Table("projects").
		Where("projects.status = ? AND projects.deleted_at IS NULL", status)
	query = r.scopeProjectAccess(query, userID, isAdmin)
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalProjects(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("projects").Where("projects.deleted_at IS NULL")
	query = r.scopeProjectAccess(query, userID, isAdmin)
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountPendingRequests(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("project_requests").
		Where("status IN ('SUBMITTED', 'UNDER_REVIEW') AND deleted_at IS NULL")

	if !isAdmin {
		query = query.Where("requester_id = ?", userID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountOverdueTasks(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("tasks").
		Where("tasks.due_date < CURRENT_DATE AND tasks.status NOT IN ('DONE', 'CANCELLED') AND tasks.deleted_at IS NULL")

	if !isAdmin {
		query = query.Joins("JOIN project_members ON project_members.project_id = tasks.project_id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountTotalTasks(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("tasks").
		Where("tasks.deleted_at IS NULL")

	if !isAdmin {
		query = query.Joins("JOIN project_members ON project_members.project_id = tasks.project_id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) CountCompletedTasks(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("tasks").
		Where("tasks.status = 'DONE' AND tasks.deleted_at IS NULL")

	if !isAdmin {
		query = query.Joins("JOIN project_members ON project_members.project_id = tasks.project_id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetAverageBudgetUsage(userID uint64, isAdmin bool) (float64, error) {
	type result struct {
		Allocated float64
		Used      float64
	}

	var res result

	query := r.db.Table("budgets").
		Select("COALESCE(SUM(budgets.allocated_budget), 0) as allocated, COALESCE(SUM(bt.used), 0) as used").
		Joins(`LEFT JOIN (
			SELECT budget_id, SUM(amount) as used 
			FROM budget_transactions 
			GROUP BY budget_id
		) bt ON bt.budget_id = budgets.id`).
		Where("budgets.deleted_at IS NULL")

	if !isAdmin {
		query = query.Joins("JOIN projects ON projects.id = budgets.project_id").
			Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	if err := query.Scan(&res).Error; err != nil {
		return 0, err
	}

	if res.Allocated == 0 {
		return 0, nil
	}

	return (res.Used / res.Allocated) * 100, nil
}

func (r *dashboardRepository) GetRecentActivities(limit int) ([]dto.RecentActivity, error) {
	var activities []dto.RecentActivity

	rows, err := r.db.Table("audit_logs").
		Select("module, action, entity_type, entity_id, user_id, created_at").
		Order("created_at DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a dto.RecentActivity
		if err := rows.Scan(&a.Module, &a.Action, &a.EntityType, &a.EntityID, &a.UserID, &a.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	return activities, nil
}

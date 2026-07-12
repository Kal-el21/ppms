package repository

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/dashboard/dto"
	"gorm.io/gorm"
)

type DashboardDeadlineRow struct {
	ID          uint64     `gorm:"column:id"`
	ProjectCode string     `gorm:"column:project_code"`
	Name        string     `gorm:"column:name"`
	EndDate     *time.Time `gorm:"column:end_date"`
	Status      string     `gorm:"column:status"`
}

type DashboardRepository interface {
	CountProjectsByStatus(status string, userID uint64, isAdmin bool) (int64, error)
	CountTotalProjects(userID uint64, isAdmin bool) (int64, error)
	CountPendingRequests(userID uint64, isAdmin bool) (int64, error)
	CountOverdueTasks(userID uint64, isAdmin bool) (int64, error)
	CountTotalTasks(userID uint64, isAdmin bool) (int64, error)
	CountCompletedTasks(userID uint64, isAdmin bool) (int64, error)
	GetAverageBudgetUsage(userID uint64, isAdmin bool) (float64, error)
	GetRecentActivities(limit int) ([]dto.RecentActivity, error)
	CountProjectsByHealth(health string, userID uint64, isAdmin bool) (int64, error)
	CountProjectsByDimension(dim string, userID uint64, isAdmin bool) (map[string]int64, error)
	GetBudgetTotalsByType(userID uint64, isAdmin bool) (capexAllocated float64, capexUsed float64, opexAllocated float64, opexUsed float64, err error)
	CountOverdueProjects(userID uint64, isAdmin bool) (int64, error)
	GetUpcomingDeadlines(limit int, userID uint64, isAdmin bool) ([]DashboardDeadlineRow, error)
	GetBudgetMaster(userID uint64, isAdmin bool) ([]dto.BudgetMasterItem, error)
	GetAbsorption(userID uint64, isAdmin bool) ([]dto.AbsorptionItem, error)
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

func (r *dashboardRepository) CountProjectsByHealth(health string, userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("projects").
		Where("projects.health = ? AND projects.deleted_at IS NULL", health)
	query = r.scopeProjectAccess(query, userID, isAdmin)
	err := query.Count(&count).Error
	return count, err
}

// CountProjectsByDimension menghitung jumlah project per nilai kolom dimension
// (initiation_type / priority / status / health). Nama kolom di-whitelist untuk
// mencegah SQL injection.
func (r *dashboardRepository) CountProjectsByDimension(dim string, userID uint64, isAdmin bool) (map[string]int64, error) {
	allowed := map[string]bool{
		"initiation_type": true,
		"priority":        true,
		"status":          true,
		"health":          true,
	}
	if !allowed[dim] {
		return map[string]int64{}, nil
	}

	type row struct {
		Key   string `gorm:"column:key"`
		Count int64  `gorm:"column:count"`
	}

	var rows []row
	query := r.db.Table("projects").
		Select(dim + " as key, COUNT(*) as count").
		Where("projects.deleted_at IS NULL")
	query = r.scopeProjectAccess(query, userID, isAdmin)
	if err := query.Group(dim).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, x := range rows {
		result[x.Key] = x.Count
	}
	return result, nil
}

// GetBudgetTotalsByType menjumlahkan allocated dan used budget per tipe CAPEX/OPEX
// di seluruh project yang bisa diakses user.
func (r *dashboardRepository) GetBudgetTotalsByType(userID uint64, isAdmin bool) (capexAllocated, capexUsed, opexAllocated, opexUsed float64, err error) {
	type result struct {
		BudgetType string
		Allocated  float64
		Used       float64
	}

	var results []result
	query := r.db.Table("budgets").
		Select("budgets.budget_type as budget_type, COALESCE(SUM(budgets.allocated_budget), 0) as allocated, COALESCE(SUM(bt.used), 0) as used").
		Joins(`LEFT JOIN (
			SELECT budget_id, SUM(amount) as used
			FROM budget_transactions
			GROUP BY budget_id
		) bt ON bt.budget_id = budgets.id`).
		Where("budgets.deleted_at IS NULL AND budgets.budget_type IS NOT NULL").
		Group("budgets.budget_type")

	if !isAdmin {
		query = query.
			Joins("JOIN projects ON projects.id = budgets.project_id").
			Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	for _, x := range results {
		switch x.BudgetType {
		case "CAPEX":
			capexAllocated = x.Allocated
			capexUsed = x.Used
		case "OPEX":
			opexAllocated = x.Allocated
			opexUsed = x.Used
		}
	}
	return capexAllocated, capexUsed, opexAllocated, opexUsed, nil
}

func (r *dashboardRepository) CountOverdueProjects(userID uint64, isAdmin bool) (int64, error) {
	var count int64
	query := r.db.Table("projects").
		Where("projects.end_date < CURRENT_DATE AND projects.status NOT IN ('COMPLETED', 'CANCELLED') AND projects.deleted_at IS NULL")
	query = r.scopeProjectAccess(query, userID, isAdmin)
	err := query.Count(&count).Error
	return count, err
}

func (r *dashboardRepository) GetUpcomingDeadlines(limit int, userID uint64, isAdmin bool) ([]DashboardDeadlineRow, error) {
	var rows []DashboardDeadlineRow
	query := r.db.Table("projects").
		Select("id, project_code, name, end_date, status").
		Where("projects.end_date >= CURRENT_DATE AND projects.status NOT IN ('COMPLETED', 'CANCELLED') AND projects.deleted_at IS NULL")
	query = r.scopeProjectAccess(query, userID, isAdmin)
	if err := query.Order("projects.end_date ASC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *dashboardRepository) GetBudgetMaster(userID uint64, isAdmin bool) ([]dto.BudgetMasterItem, error) {
	type row struct {
		BudgetName      string  `gorm:"column:budget_name"`
		BudgetType      string  `gorm:"column:budget_type"`
		Allocated       float64 `gorm:"column:allocated"`
		Used            float64 `gorm:"column:used"`
	}

	var rows []row
	query := r.db.Table("budgets").
		Select("budgets.budget_name, budgets.budget_type, COALESCE(SUM(budgets.allocated_budget), 0) as allocated, COALESCE(SUM(bt.used), 0) as used").
		Joins(`LEFT JOIN (
			SELECT budget_id, SUM(amount) as used 
			FROM budget_transactions 
			GROUP BY budget_id
		) bt ON bt.budget_id = budgets.id`).
		Where("budgets.deleted_at IS NULL AND budgets.budget_name IS NOT NULL AND budgets.budget_name != ''").
		Group("budgets.budget_name, budgets.budget_type")

	if !isAdmin {
		query = query.
			Joins("JOIN projects ON projects.id = budgets.project_id").
			Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]dto.BudgetMasterItem, 0, len(rows))
	for _, r := range rows {
		remaining := r.Allocated - r.Used
		var usage float64
		if r.Allocated > 0 {
			usage = (r.Used / r.Allocated) * 100
		}
		result = append(result, dto.BudgetMasterItem{
			BudgetName:      r.BudgetName,
			BudgetType:      r.BudgetType,
			Allocated:       r.Allocated,
			Used:            r.Used,
			Remaining:       remaining,
			UsagePercentage: usage,
		})
	}
	return result, nil
}

func (r *dashboardRepository) GetAbsorption(userID uint64, isAdmin bool) ([]dto.AbsorptionItem, error) {
	type row struct {
		ProjectCode string  `gorm:"column:project_code"`
		ProjectName string  `gorm:"column:project_name"`
		BudgetType  string  `gorm:"column:budget_type"`
		Allocated   float64 `gorm:"column:allocated"`
		Used        float64 `gorm:"column:used"`
	}

	var rows []row
	query := r.db.Table("budgets").
		Select("projects.project_code, projects.name as project_name, budgets.budget_type, COALESCE(SUM(budgets.allocated_budget), 0) as allocated, COALESCE(SUM(bt.used), 0) as used").
		Joins("JOIN projects ON projects.id = budgets.project_id AND projects.deleted_at IS NULL").
		Joins(`LEFT JOIN (
			SELECT budget_id, SUM(amount) as used 
			FROM budget_transactions 
			GROUP BY budget_id
		) bt ON bt.budget_id = budgets.id`).
		Where("budgets.deleted_at IS NULL").
		Group("projects.project_code, projects.name, budgets.budget_type")

	if !isAdmin {
		query = query.
			Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]dto.AbsorptionItem, 0, len(rows))
	for _, r := range rows {
		var usage float64
		if r.Allocated > 0 {
			usage = (r.Used / r.Allocated) * 100
		}
		result = append(result, dto.AbsorptionItem{
			ProjectCode:     r.ProjectCode,
			ProjectName:     r.ProjectName,
			BudgetType:      r.BudgetType,
			Allocated:       r.Allocated,
			Used:            r.Used,
			UsagePercentage: usage,
		})
	}
	return result, nil
}

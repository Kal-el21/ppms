package repository

import (
	"github.com/Kal-el21/backend/internal/domain/reporting/dto"
	"gorm.io/gorm"
)

type ReportingRepository interface {
	GetProjectReport(projectID *uint64) ([]dto.ProjectReportRow, error)
	GetTaskReport(projectID *uint64) ([]dto.TaskReportRow, error)
	GetBudgetReport(projectID *uint64) ([]dto.BudgetReportRow, error)
	GetMilestoneReport(projectID *uint64) ([]dto.MilestoneReportRow, error)
	GetHandoverReport(projectID *uint64) ([]dto.HandoverReportRow, error)
}

type reportingRepository struct {
	db *gorm.DB
}

func NewReportingRepository(db *gorm.DB) ReportingRepository {
	return &reportingRepository{db: db}
}

func (r *reportingRepository) GetProjectReport(projectID *uint64) ([]dto.ProjectReportRow, error) {
	var rows []dto.ProjectReportRow

	query := r.db.Table("projects").
		Select(`projects.id, projects.name, projects.status,
			COALESCE(milestone_avg.avg_progress, 0) as progress,
			COALESCE(to_char(projects.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(to_char(projects.end_date, 'YYYY-MM-DD'), '') as end_date`).
		Joins(`LEFT JOIN (
			SELECT m.project_id, AVG(task_avg.progress) as avg_progress
			FROM milestones m
			LEFT JOIN (
				SELECT milestone_id, AVG(progress) as progress
				FROM tasks
				WHERE status != 'CANCELLED' AND deleted_at IS NULL
				GROUP BY milestone_id
			) task_avg ON task_avg.milestone_id = m.id
			WHERE m.status != 'CANCELLED' AND m.deleted_at IS NULL
			GROUP BY m.project_id
		) milestone_avg ON milestone_avg.project_id = projects.id`).
		Where("projects.deleted_at IS NULL")

	if projectID != nil {
		query = query.Where("projects.id = ?", *projectID)
	}

	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetTaskReport(projectID *uint64) ([]dto.TaskReportRow, error) {
	var rows []dto.TaskReportRow

	query := r.db.Table("tasks").
		Select(`tasks.id, tasks.title, tasks.status, tasks.priority, tasks.progress,
			COALESCE(to_char(tasks.due_date, 'YYYY-MM-DD'), '') as due_date`).
		Where("tasks.deleted_at IS NULL")

	if projectID != nil {
		query = query.Where("tasks.project_id = ?", *projectID)
	}

	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetBudgetReport(projectID *uint64) ([]dto.BudgetReportRow, error) {
	var rows []dto.BudgetReportRow

	query := r.db.Table("budgets").
		Select(`projects.name as project_name, budgets.allocated_budget,
			COALESCE(bt.used, 0) as used_budget,
			budgets.allocated_budget - COALESCE(bt.used, 0) as remaining_budget`).
		Joins("JOIN projects ON projects.id = budgets.project_id").
		Joins(`LEFT JOIN (
			SELECT budget_id, SUM(amount) as used
			FROM budget_transactions
			GROUP BY budget_id
		) bt ON bt.budget_id = budgets.id`).
		Where("budgets.deleted_at IS NULL")

	if projectID != nil {
		query = query.Where("budgets.project_id = ?", *projectID)
	}

	err := query.Scan(&rows).Error
	return rows, err
}

// GetMilestoneReport: progress dihitung sama seperti milestone.service
// (Average task progress, cancelled task dikecualikan), demi konsistensi
// angka antara halaman Project Detail dan Report yang di-export.
func (r *reportingRepository) GetMilestoneReport(projectID *uint64) ([]dto.MilestoneReportRow, error) {
	var rows []dto.MilestoneReportRow

	query := r.db.Table("milestones").
		Select(`projects.name as project_name, milestones.name, milestones.status,
			COALESCE(task_avg.progress, 0) as progress,
			COALESCE(to_char(milestones.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(to_char(milestones.end_date, 'YYYY-MM-DD'), '') as end_date`).
		Joins("JOIN projects ON projects.id = milestones.project_id").
		Joins(`LEFT JOIN (
			SELECT milestone_id, AVG(progress) as progress
			FROM tasks
			WHERE status != 'CANCELLED' AND deleted_at IS NULL
			GROUP BY milestone_id
		) task_avg ON task_avg.milestone_id = milestones.id`).
		Where("milestones.deleted_at IS NULL")

	if projectID != nil {
		query = query.Where("milestones.project_id = ?", *projectID)
	}

	query = query.Order("milestones.order_index ASC")

	err := query.Scan(&rows).Error
	return rows, err
}

// GetHandoverReport: join ke users dua kali (sender & receiver) sehingga
// report menampilkan nama, bukan hanya ID.
func (r *reportingRepository) GetHandoverReport(projectID *uint64) ([]dto.HandoverReportRow, error) {
	var rows []dto.HandoverReportRow

	query := r.db.Table("handovers").
		Select(`projects.name as project_name,
			sender.full_name as sender_name,
			COALESCE(receiver.full_name, '-') as receiver_name,
			handovers.status,
			COALESCE(to_char(handovers.delivery_date, 'YYYY-MM-DD'), '') as delivery_date,
			COALESCE(to_char(handovers.received_at, 'YYYY-MM-DD HH24:MI'), '') as received_at`).
		Joins("JOIN projects ON projects.id = handovers.project_id").
		Joins("JOIN users sender ON sender.id = handovers.sender_id").
		Joins("LEFT JOIN users receiver ON receiver.id = handovers.receiver_id").
		Where("handovers.deleted_at IS NULL")

	if projectID != nil {
		query = query.Where("handovers.project_id = ?", *projectID)
	}

	query = query.Order("handovers.created_at DESC")

	err := query.Scan(&rows).Error
	return rows, err
}

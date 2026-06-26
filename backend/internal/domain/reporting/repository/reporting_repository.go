package repository

import (
	"encoding/json"
	"time"

	"github.com/Kal-el21/backend/internal/domain/reporting/dto"
	"github.com/Kal-el21/backend/internal/domain/reporting/entity"
	"gorm.io/gorm"
)

type ReportingRepository interface {
	GetProjectReport(projectID *uint64) ([]dto.ProjectReportRow, error)
	GetTaskReport(projectID *uint64) ([]dto.TaskReportRow, error)
	GetBudgetReport(projectID *uint64) ([]dto.BudgetReportRow, error)
	GetMilestoneReport(projectID *uint64) ([]dto.MilestoneReportRow, error)
	GetHandoverReport(projectID *uint64) ([]dto.HandoverReportRow, error)
	GetProjectReportLive(projectID *uint64) ([]dto.ProjectReportRow, error)
	GetTaskReportLive(projectID *uint64) ([]dto.TaskReportRow, error)
	GetBudgetReportLive(projectID *uint64) ([]dto.BudgetReportRow, error)
	GetMilestoneReportLive(projectID *uint64) ([]dto.MilestoneReportRow, error)
	GetHandoverReportLive(projectID *uint64) ([]dto.HandoverReportRow, error)
	CreateSnapshot(reportType string, projectID *uint64, data map[string]interface{}) error
}

type reportingRepository struct {
	db *gorm.DB
}

func NewReportingRepository(db *gorm.DB) ReportingRepository {
	return &reportingRepository{db: db}
}

func (r *reportingRepository) getLatestSnapshot(projectID *uint64, reportType string) (*entity.ProjectSnapshot, error) {
	var snapshot entity.ProjectSnapshot
	today := time.Now().Format("2006-01-02")
	query := r.db.Where("snapshot_date = ?", today)
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}
	err := query.Order("created_at DESC").First(&snapshot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}

func unmarshalReportRows(snapshot *entity.ProjectSnapshot, reportType string) (interface{}, error) {
	if snapshot == nil || snapshot.SnapshotData == nil {
		return nil, nil
	}
	rowsJSON, err := json.Marshal(snapshot.SnapshotData["rows"])
	if err != nil {
		return nil, err
	}
	switch reportType {
	case "PROJECT":
		var rows []dto.ProjectReportRow
		if err := json.Unmarshal(rowsJSON, &rows); err != nil {
			return nil, err
		}
		return rows, nil
	case "TASK":
		var rows []dto.TaskReportRow
		if err := json.Unmarshal(rowsJSON, &rows); err != nil {
			return nil, err
		}
		return rows, nil
	case "BUDGET":
		var rows []dto.BudgetReportRow
		if err := json.Unmarshal(rowsJSON, &rows); err != nil {
			return nil, err
		}
		return rows, nil
	case "MILESTONE":
		var rows []dto.MilestoneReportRow
		if err := json.Unmarshal(rowsJSON, &rows); err != nil {
			return nil, err
		}
		return rows, nil
	case "HANDOVER":
		var rows []dto.HandoverReportRow
		if err := json.Unmarshal(rowsJSON, &rows); err != nil {
			return nil, err
		}
		return rows, nil
	default:
		return nil, nil
	}
}

func (r *reportingRepository) GetProjectReport(projectID *uint64) ([]dto.ProjectReportRow, error) {
	snapshot, err := r.getLatestSnapshot(projectID, "PROJECT")
	if err != nil {
		return nil, err
	}
	if snapshot != nil && snapshot.SnapshotData != nil {
		rows, err := unmarshalReportRows(snapshot, "PROJECT")
		if err == nil && rows != nil {
			return rows.([]dto.ProjectReportRow), nil
		}
	}
	rows, err := r.GetProjectReportLive(projectID)
	if err != nil {
		return nil, err
	}
	r.createSnapshotFromRows(projectID, "PROJECT", rows)
	return rows, nil
}

func (r *reportingRepository) GetTaskReport(projectID *uint64) ([]dto.TaskReportRow, error) {
	snapshot, err := r.getLatestSnapshot(projectID, "TASK")
	if err != nil {
		return nil, err
	}
	if snapshot != nil && snapshot.SnapshotData != nil {
		rows, err := unmarshalReportRows(snapshot, "TASK")
		if err == nil && rows != nil {
			return rows.([]dto.TaskReportRow), nil
		}
	}
	rows, err := r.GetTaskReportLive(projectID)
	if err != nil {
		return nil, err
	}
	r.createSnapshotFromRows(projectID, "TASK", rows)
	return rows, nil
}

func (r *reportingRepository) GetBudgetReport(projectID *uint64) ([]dto.BudgetReportRow, error) {
	snapshot, err := r.getLatestSnapshot(projectID, "BUDGET")
	if err != nil {
		return nil, err
	}
	if snapshot != nil && snapshot.SnapshotData != nil {
		rows, err := unmarshalReportRows(snapshot, "BUDGET")
		if err == nil && rows != nil {
			return rows.([]dto.BudgetReportRow), nil
		}
	}
	rows, err := r.GetBudgetReportLive(projectID)
	if err != nil {
		return nil, err
	}
	r.createSnapshotFromRows(projectID, "BUDGET", rows)
	return rows, nil
}

func (r *reportingRepository) GetMilestoneReport(projectID *uint64) ([]dto.MilestoneReportRow, error) {
	snapshot, err := r.getLatestSnapshot(projectID, "MILESTONE")
	if err != nil {
		return nil, err
	}
	if snapshot != nil && snapshot.SnapshotData != nil {
		rows, err := unmarshalReportRows(snapshot, "MILESTONE")
		if err == nil && rows != nil {
			return rows.([]dto.MilestoneReportRow), nil
		}
	}
	rows, err := r.GetMilestoneReportLive(projectID)
	if err != nil {
		return nil, err
	}
	r.createSnapshotFromRows(projectID, "MILESTONE", rows)
	return rows, nil
}

func (r *reportingRepository) GetHandoverReport(projectID *uint64) ([]dto.HandoverReportRow, error) {
	snapshot, err := r.getLatestSnapshot(projectID, "HANDOVER")
	if err != nil {
		return nil, err
	}
	if snapshot != nil && snapshot.SnapshotData != nil {
		rows, err := unmarshalReportRows(snapshot, "HANDOVER")
		if err == nil && rows != nil {
			return rows.([]dto.HandoverReportRow), nil
		}
	}
	rows, err := r.GetHandoverReportLive(projectID)
	if err != nil {
		return nil, err
	}
	r.createSnapshotFromRows(projectID, "HANDOVER", rows)
	return rows, nil
}


func (r *reportingRepository) CreateSnapshot(reportType string, projectID *uint64, data map[string]interface{}) error {
	snapshot := entity.ProjectSnapshot{
		SnapshotDate: time.Now(),
		SnapshotData: data,
	}
	if projectID != nil {
		snapshot.ProjectID = *projectID
	}
	if pid, ok := data["project_id"].(uint64); ok {
		snapshot.ProjectID = pid
	} else if pidF, ok := data["project_id"].(float64); ok {
		snapshot.ProjectID = uint64(pidF)
	}
	if status, ok := data["project_status"].(string); ok {
		snapshot.ProjectStatus = status
	}
	return r.db.Create(&snapshot).Error
}



func (r *reportingRepository) createSnapshotFromRows(projectID *uint64, reportType string, rows interface{}) error {
	snapshotData := map[string]interface{}{
		"report_type": reportType,
		"project_id":  projectID,
		"generated_at": time.Now().Format(time.RFC3339),
		"rows":        rows,
	}
	snapshotJSON, err := json.Marshal(snapshotData)
	if err != nil {
		return err
	}
	var dataMap map[string]interface{}
	if err := json.Unmarshal(snapshotJSON, &dataMap); err != nil {
		return err
	}
	snapshot := entity.ProjectSnapshot{
		SnapshotDate: time.Now(),
		SnapshotData: dataMap,
	}
	if projectID != nil {
		snapshot.ProjectID = *projectID
	}
	if pid, ok := dataMap["project_id"].(uint64); ok {
		snapshot.ProjectID = pid
	} else if pidF, ok := dataMap["project_id"].(float64); ok {
		snapshot.ProjectID = uint64(pidF)
	}
	if status, ok := dataMap["project_status"].(string); ok {
		snapshot.ProjectStatus = status
	}
	return r.db.Create(&snapshot).Error
}


func (r *reportingRepository) GetProjectReportLive(projectID *uint64) ([]dto.ProjectReportRow, error) {
	var rows []dto.ProjectReportRow
	query := r.db.Table("projects").
		Select("projects.id, projects.name, projects.status, COALESCE(milestone_avg.avg_progress, 0) as progress, COALESCE(to_char(projects.start_date, 'YYYY-MM-DD'), '') as start_date, COALESCE(to_char(projects.end_date, 'YYYY-MM-DD'), '') as end_date").
		Joins("LEFT JOIN (SELECT m.project_id, AVG(task_avg.progress) as avg_progress FROM milestones m LEFT JOIN (SELECT milestone_id, AVG(progress) as progress FROM tasks WHERE status != 'CANCELLED' AND deleted_at IS NULL GROUP BY milestone_id) task_avg ON task_avg.milestone_id = m.id WHERE m.status != 'CANCELLED' AND m.deleted_at IS NULL GROUP BY m.project_id) milestone_avg ON milestone_avg.project_id = projects.id").
		Where("projects.deleted_at IS NULL")
	if projectID != nil {
		query = query.Where("projects.id = ?", *projectID)
	}
	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetTaskReportLive(projectID *uint64) ([]dto.TaskReportRow, error) {
	var rows []dto.TaskReportRow
	query := r.db.Table("tasks").
		Select("tasks.id, tasks.title, tasks.status, tasks.priority, tasks.progress, COALESCE(to_char(tasks.due_date, 'YYYY-MM-DD'), '') as due_date").
		Where("tasks.deleted_at IS NULL")
	if projectID != nil {
		query = query.Where("tasks.project_id = ?", *projectID)
	}
	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetBudgetReportLive(projectID *uint64) ([]dto.BudgetReportRow, error) {
	var rows []dto.BudgetReportRow
	query := r.db.Table("budgets").
		Select("projects.name as project_name, budgets.allocated_budget, COALESCE(bt.used, 0) as used_budget, budgets.allocated_budget - COALESCE(bt.used, 0) as remaining_budget").
		Joins("JOIN projects ON projects.id = budgets.project_id").
		Joins("LEFT JOIN (SELECT budget_id, SUM(amount) as used FROM budget_transactions GROUP BY budget_id) bt ON bt.budget_id = budgets.id").
		Where("budgets.deleted_at IS NULL")
	if projectID != nil {
		query = query.Where("budgets.project_id = ?", *projectID)
	}
	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetMilestoneReportLive(projectID *uint64) ([]dto.MilestoneReportRow, error) {
	var rows []dto.MilestoneReportRow
	query := r.db.Table("milestones").
		Select("projects.name as project_name, milestones.name, milestones.status, COALESCE(task_avg.progress, 0) as progress, COALESCE(to_char(milestones.start_date, 'YYYY-MM-DD'), '') as start_date, COALESCE(to_char(milestones.end_date, 'YYYY-MM-DD'), '') as end_date").
		Joins("JOIN projects ON projects.id = milestones.project_id").
		Joins("LEFT JOIN (SELECT milestone_id, AVG(progress) as progress FROM tasks WHERE status != 'CANCELLED' AND deleted_at IS NULL GROUP BY milestone_id) task_avg ON task_avg.milestone_id = milestones.id").
		Where("milestones.deleted_at IS NULL")
	if projectID != nil {
		query = query.Where("milestones.project_id = ?", *projectID)
	}
	query = query.Order("milestones.order_index ASC")
	err := query.Scan(&rows).Error
	return rows, err
}

func (r *reportingRepository) GetHandoverReportLive(projectID *uint64) ([]dto.HandoverReportRow, error) {
	var rows []dto.HandoverReportRow
	query := r.db.Table("handovers").
		Select("projects.name as project_name, sender.full_name as sender_name, COALESCE(receiver.full_name, '-') as receiver_name, handovers.status, COALESCE(to_char(handovers.delivery_date, 'YYYY-MM-DD'), '') as delivery_date, COALESCE(to_char(handovers.received_at, 'YYYY-MM-DD HH24:MI'), '') as received_at").
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

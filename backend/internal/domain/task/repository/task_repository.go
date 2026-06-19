package repository

import (
	"github.com/Kal-el21/backend/internal/domain/task/entity"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByID(id uint64) (*entity.Task, error)
	FindByProjectID(projectID uint64) ([]entity.Task, error)
	FindByMilestoneID(milestoneID uint64) ([]entity.Task, error)
	UpdateWithVersionCheck(task *entity.Task, expectedVersion int) (int64, error)
	UpdateProgressWithVersionCheck(id uint64, progress int, expectedVersion int) (int64, error)
	Delete(id uint64, deletedBy uint64) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *entity.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint64) (*entity.Task, error) {
	var task entity.Task
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByProjectID(projectID uint64) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order("order_index ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) FindByMilestoneID(milestoneID uint64) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.Where("milestone_id = ? AND deleted_at IS NULL", milestoneID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateWithVersionCheck(task *entity.Task, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Task{}).
		Where("id = ? AND version = ?", task.ID, expectedVersion).
		Updates(map[string]interface{}{
			"title":        task.Title,
			"description":  task.Description,
			"priority":     task.Priority,
			"status":       task.Status,
			"milestone_id": task.MilestoneID,
			"start_date":   task.StartDate,
			"due_date":     task.DueDate,
			"version":      gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *taskRepository) UpdateProgressWithVersionCheck(id uint64, progress int, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Task{}).
		Where("id = ? AND version = ?", id, expectedVersion).
		Updates(map[string]interface{}{
			"progress": progress,
			"version":  gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *taskRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Task{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}

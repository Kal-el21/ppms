package repository

import (
	"github.com/Kal-el21/backend/internal/domain/task/entity"
	"gorm.io/gorm"
)

type AssigneeRepository interface {
	Assign(assignee *entity.TaskAssignee) error
	Unassign(taskID, userID, unassignedBy uint64) error
	FindActiveByTaskID(taskID uint64) ([]entity.TaskAssignee, error)
	IsAssigned(taskID, userID uint64) (bool, error)
}

type assigneeRepository struct {
	db *gorm.DB
}

func NewAssigneeRepository(db *gorm.DB) AssigneeRepository {
	return &assigneeRepository{db: db}
}

func (r *assigneeRepository) Assign(assignee *entity.TaskAssignee) error {
	return r.db.Create(assignee).Error
}

func (r *assigneeRepository) Unassign(taskID, userID, unassignedBy uint64) error {
	return r.db.Model(&entity.TaskAssignee{}).
		Where("task_id = ? AND user_id = ? AND unassigned_at IS NULL", taskID, userID).
		Updates(map[string]interface{}{
			"unassigned_at": gorm.Expr("now()"),
			"unassigned_by": unassignedBy,
		}).Error
}

func (r *assigneeRepository) FindActiveByTaskID(taskID uint64) ([]entity.TaskAssignee, error) {
	var assignees []entity.TaskAssignee
	err := r.db.Where("task_id = ? AND unassigned_at IS NULL", taskID).Find(&assignees).Error
	return assignees, err
}

func (r *assigneeRepository) IsAssigned(taskID, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&entity.TaskAssignee{}).
		Where("task_id = ? AND user_id = ? AND unassigned_at IS NULL", taskID, userID).
		Count(&count).Error
	return count > 0, err
}

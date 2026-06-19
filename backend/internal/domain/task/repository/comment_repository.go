package repository

import (
	"github.com/Kal-el21/backend/internal/domain/task/entity"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *entity.TaskComment) error
	FindByTaskID(taskID uint64, page, limit int) ([]entity.TaskComment, int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *entity.TaskComment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByTaskID(taskID uint64, page, limit int) ([]entity.TaskComment, int64, error) {
	var comments []entity.TaskComment
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&entity.TaskComment{}).Where("task_id = ?", taskID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("task_id = ?", taskID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error

	return comments, total, err
}

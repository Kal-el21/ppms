package repository

import (
	"github.com/Kal-el21/backend/internal/domain/approval/entity"
	"gorm.io/gorm"
)

type ApprovalLevelRepository interface {
	Create(level *entity.ApprovalLevel) error
	GetByWorkflowID(workflowID uint64) ([]entity.ApprovalLevel, error)
	GetByID(id uint64) (*entity.ApprovalLevel, error)
	Delete(id uint64) error
}

type approvalLevelRepository struct {
	db *gorm.DB
}

func NewApprovalLevelRepository(db *gorm.DB) ApprovalLevelRepository {
	return &approvalLevelRepository{db: db}
}

func (r *approvalLevelRepository) Create(level *entity.ApprovalLevel) error {
	return r.db.Create(level).Error
}

func (r *approvalLevelRepository) GetByWorkflowID(workflowID uint64) ([]entity.ApprovalLevel, error) {
	var levels []entity.ApprovalLevel
	err := r.db.Where("workflow_id = ?", workflowID).Order("level_order ASC").Find(&levels).Error
	return levels, err
}

func (r *approvalLevelRepository) GetByID(id uint64) (*entity.ApprovalLevel, error) {
	var level entity.ApprovalLevel
	err := r.db.First(&level, id).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

func (r *approvalLevelRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.ApprovalLevel{}, id).Error
}

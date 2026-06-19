package repository

import (
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	"gorm.io/gorm"
)

type ApprovalRepository interface {
	Create(approval *entity.ProjectRequestApproval) error
	FindByRequestID(requestID uint64) ([]entity.ProjectRequestApproval, error)
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(approval *entity.ProjectRequestApproval) error {
	return r.db.Create(approval).Error
}

func (r *approvalRepository) FindByRequestID(requestID uint64) ([]entity.ProjectRequestApproval, error) {
	var approvals []entity.ProjectRequestApproval
	err := r.db.Where("project_request_id = ?", requestID).
		Order("created_at ASC").
		Find(&approvals).Error
	return approvals, err
}

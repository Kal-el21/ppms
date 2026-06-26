package repository

import (
	"github.com/Kal-el21/backend/internal/domain/approval/entity"
	"gorm.io/gorm"
)

type ApprovalWorkflowRepository interface {
	Create(workflow *entity.ApprovalWorkflow) error
	GetByID(id uint64) (*entity.ApprovalWorkflow, error)
	GetAll() ([]entity.ApprovalWorkflow, error)
	Update(workflow *entity.ApprovalWorkflow) error
}

type approvalWorkflowRepository struct {
	db *gorm.DB
}

func NewApprovalWorkflowRepository(db *gorm.DB) ApprovalWorkflowRepository {
	return &approvalWorkflowRepository{db: db}
}

func (r *approvalWorkflowRepository) Create(workflow *entity.ApprovalWorkflow) error {
	return r.db.Create(workflow).Error
}

func (r *approvalWorkflowRepository) GetByID(id uint64) (*entity.ApprovalWorkflow, error) {
	var workflow entity.ApprovalWorkflow
	err := r.db.First(&workflow, id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *approvalWorkflowRepository) GetAll() ([]entity.ApprovalWorkflow, error) {
	var workflows []entity.ApprovalWorkflow
	err := r.db.Find(&workflows).Error
	return workflows, err
}

func (r *approvalWorkflowRepository) Update(workflow *entity.ApprovalWorkflow) error {
	return r.db.Save(workflow).Error
}

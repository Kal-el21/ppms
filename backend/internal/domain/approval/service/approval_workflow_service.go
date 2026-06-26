package service

import (
	"github.com/Kal-el21/backend/internal/domain/approval/entity"
	"github.com/Kal-el21/backend/internal/domain/approval/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

type ApprovalWorkflowService interface {
	CreateWorkflow(name string) (*entity.ApprovalWorkflow, error)
	GetWorkflow(id uint64) (*entity.ApprovalWorkflow, error)
	ListWorkflows() ([]entity.ApprovalWorkflow, error)
}

type approvalWorkflowService struct {
	repo repository.ApprovalWorkflowRepository
}

func NewApprovalWorkflowService(repo repository.ApprovalWorkflowRepository) ApprovalWorkflowService {
	return &approvalWorkflowService{repo: repo}
}

func (s *approvalWorkflowService) CreateWorkflow(name string) (*entity.ApprovalWorkflow, error) {
	workflow := &entity.ApprovalWorkflow{
		Name: name,
	}
	if err := s.repo.Create(workflow); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, err.Error())
	}
	return workflow, nil
}

func (s *approvalWorkflowService) GetWorkflow(id uint64) (*entity.ApprovalWorkflow, error) {
	workflow, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrValidation, "workflow not found")
	}
	return workflow, nil
}

func (s *approvalWorkflowService) ListWorkflows() ([]entity.ApprovalWorkflow, error) {
	return s.repo.GetAll()
}

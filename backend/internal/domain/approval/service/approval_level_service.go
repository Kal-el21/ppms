package service

import (
	"github.com/Kal-el21/backend/internal/domain/approval/entity"
	"github.com/Kal-el21/backend/internal/domain/approval/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
)

type ApprovalLevelService interface {
	CreateLevel(workflowID uint64, levelOrder int, roleRequired string) (*entity.ApprovalLevel, error)
	GetLevelsByWorkflow(workflowID uint64) ([]entity.ApprovalLevel, error)
}

type approvalLevelService struct {
	repo repository.ApprovalLevelRepository
}

func NewApprovalLevelService(repo repository.ApprovalLevelRepository) ApprovalLevelService {
	return &approvalLevelService{repo: repo}
}

func (s *approvalLevelService) CreateLevel(workflowID uint64, levelOrder int, roleRequired string) (*entity.ApprovalLevel, error) {
	level := &entity.ApprovalLevel{
		WorkflowID:  workflowID,
		LevelOrder:  levelOrder,
		RoleRequired: roleRequired,
	}
	if err := s.repo.Create(level); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, err.Error())
	}
	return level, nil
}

func (s *approvalLevelService) GetLevelsByWorkflow(workflowID uint64) ([]entity.ApprovalLevel, error) {
	return s.repo.GetByWorkflowID(workflowID)
}

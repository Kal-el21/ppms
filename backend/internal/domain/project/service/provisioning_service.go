package service

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/project/entity"
	"github.com/Kal-el21/backend/internal/domain/project/repository"
)

// ProvisioningService bertanggung jawab khusus untuk auto-create project
// dari approved request (FR-05.01). Service lifecycle/member penuh ada di Phase 3.
type ProvisioningService interface {
	CreateFromApprovedRequest(requestID uint64, title string, description string, requesterID uint64, projectManagerID uint64) (*entity.Project, error)
}

type provisioningService struct {
	projectRepo repository.ProjectRepository
	memberRepo  repository.MemberRepository
}

func NewProvisioningService(projectRepo repository.ProjectRepository, memberRepo repository.MemberRepository) ProvisioningService {
	return &provisioningService{projectRepo: projectRepo, memberRepo: memberRepo}
}

func (s *provisioningService) CreateFromApprovedRequest(requestID uint64, title string, description string, requesterID uint64, projectManagerID uint64) (*entity.Project, error) {
	project := &entity.Project{
		ProjectRequestID: &requestID,
		Name:             title,
		Description:      description,
		Status:           entity.ProjectPlanned,
		CreatedBy:        requesterID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Admin memilih PROJECT_MANAGER pertama saat menyetujui request.
	member := &entity.ProjectMember{
		ProjectID:   project.ID,
		UserID:      projectManagerID,
		ProjectRole: entity.RoleProjectManager,
		Status:      entity.MemberActive,
		JoinedAt:    time.Now(),
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return project, nil
}

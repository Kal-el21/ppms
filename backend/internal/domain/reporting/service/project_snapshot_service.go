package service

import (
	"github.com/Kal-el21/backend/internal/domain/reporting/entity"
	"github.com/Kal-el21/backend/internal/domain/reporting/repository"
	"time"
)

type ProjectSnapshotService interface {
	CreateSnapshot(snapshot *entity.ProjectSnapshot) error
	GetSnapshotsByProject(projectID uint64, fromDate, toDate *time.Time) ([]entity.ProjectSnapshot, error)
	GetLatestSnapshot(projectID uint64) (*entity.ProjectSnapshot, error)
}

type projectSnapshotService struct {
	repo repository.ProjectSnapshotRepository
}

func NewProjectSnapshotService(repo repository.ProjectSnapshotRepository) ProjectSnapshotService {
	return &projectSnapshotService{repo: repo}
}

func (s *projectSnapshotService) CreateSnapshot(snapshot *entity.ProjectSnapshot) error {
	return s.repo.Create(snapshot)
}

func (s *projectSnapshotService) GetSnapshotsByProject(projectID uint64, fromDate, toDate *time.Time) ([]entity.ProjectSnapshot, error) {
	return s.repo.GetByProjectID(projectID, fromDate, toDate)
}

func (s *projectSnapshotService) GetLatestSnapshot(projectID uint64) (*entity.ProjectSnapshot, error) {
	return s.repo.GetLatestByProjectID(projectID)
}

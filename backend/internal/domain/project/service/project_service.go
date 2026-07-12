package service

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/project/errors"
	"github.com/Kal-el21/backend/internal/domain/project/repository"
	"github.com/Kal-el21/backend/internal/domain/project/statemachine"
	"github.com/Kal-el21/backend/internal/events"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type ProjectService interface {
	GetByID(id uint64) (*entity.Project, error)
	GetAll(page, limit int, statusFilter string) ([]entity.Project, int64, error)
	GetMyProjects(userID uint64, page, limit int) ([]entity.Project, int64, error)
	List(filter repository.ProjectListFilter) ([]repository.ProjectListRow, int64, error)
	GetDeadlineProjects(window string) ([]dto.ProjectDeadline, error)
	Update(id uint64, req dto.UpdateProjectRequest) (*entity.Project, error)
	ChangeStatus(id uint64, req dto.ChangeStatusRequest) (*entity.Project, error)
	CalculateProgress(id uint64) (float64, error)
	CalculateHealth(project *entity.Project, progress float64) string
	RecalculateHealth(projectID uint64) error
}

type projectService struct {
	repo          repository.ProjectRepository
	milestoneRepo MilestoneProgressProvider
	eventBus      *events.Bus
}

// MilestoneProgressProvider adalah interface kecil untuk menghindari
// circular dependency antara domain project <-> milestone.
// Diimplementasikan oleh milestone.service di bagian 3.
type MilestoneProgressProvider interface {
	GetAverageProgressByProject(projectID uint64) (float64, error)
}

func NewProjectService(repo repository.ProjectRepository, milestoneRepo MilestoneProgressProvider, eventBus *events.Bus) ProjectService {
	return &projectService{repo: repo, milestoneRepo: milestoneRepo, eventBus: eventBus}
}

func (s *projectService) GetByID(id uint64) (*entity.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

func (s *projectService) GetAll(page, limit int, statusFilter string) ([]entity.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindAll(page, limit, statusFilter)
}

func (s *projectService) GetMyProjects(userID uint64, page, limit int) ([]entity.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindByMemberUserID(userID, page, limit)
}

// List mengambil project portfolio dengan filter & sorting (Phase 5).
// Untuk sort by progress, fetch all terlebih dahulu karena progress
// dihitung dari milestone (bukan kolom stored), lalu sort di memory.
func (s *projectService) List(filter repository.ProjectListFilter) ([]repository.ProjectListRow, int64, error) {
	if filter.Sort != "progress" {
		return s.repo.FindWithFilters(filter)
	}

	allFilter := filter
	allFilter.Page = 1
	allFilter.Limit = 1000
	rows, total, err := s.repo.FindWithFilters(allFilter)
	if err != nil {
		return nil, 0, err
	}

	type rowWithProgress struct {
		row    repository.ProjectListRow
		progress float64
	}
	withProgress := make([]rowWithProgress, 0, len(rows))
	for _, r := range rows {
		progress, _ := s.milestoneRepo.GetAverageProgressByProject(r.ID)
		withProgress = append(withProgress, rowWithProgress{row: r, progress: progress})
	}

	sort.Slice(withProgress, func(i, j int) bool {
		return withProgress[i].progress > withProgress[j].progress
	})

	start := (filter.Page - 1) * filter.Limit
	end := start + filter.Limit
	if start > len(withProgress) {
		start = len(withProgress)
	}
	if end > len(withProgress) {
		end = len(withProgress)
	}
	paged := make([]repository.ProjectListRow, 0, end-start)
	for _, wp := range withProgress[start:end] {
		paged = append(paged, wp.row)
	}
	return paged, total, nil
}

// GetDeadlineProjects mengembalikan project dalam window deadline tertentu
// lengkap dengan progress (dari milestone) dan sisa hari.
func (s *projectService) GetDeadlineProjects(window string) ([]dto.ProjectDeadline, error) {
	rows, err := s.repo.FindDeadlineProjects(window, time.Now())
	if err != nil {
		return nil, err
	}

	result := make([]dto.ProjectDeadline, 0, len(rows))
	for _, row := range rows {
		progress, _ := s.milestoneRepo.GetAverageProgressByProject(row.ID)
		var daysRemaining int
		if row.EndDate != nil {
			daysRemaining = int(startOfDay(time.Now()).Sub(startOfDay(*row.EndDate)).Hours() / 24 * -1)
		}
		status := row.Status
		result = append(result, dto.ProjectDeadline{
			ID:            row.ID,
			ProjectCode:  row.ProjectCode,
			Name:          row.Name,
			EndDate:       row.EndDate,
			Status:        status,
			Progress:      progress,
			DaysRemaining: daysRemaining,
		})
	}
	return result, nil
}

func (s *projectService) Update(id uint64, req dto.UpdateProjectRequest) (*entity.Project, error) {
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	category := strings.TrimSpace(req.Category)
	if len(category) > 100 {
		return nil, apperrors.New(apperrors.ErrValidation, "category is too long")
	}
	notes := strings.TrimSpace(req.Notes)
	if len(notes) > 2000 {
		return nil, apperrors.New(apperrors.ErrValidation, "notes is too long")
	}
	priority, err := normalizeProjectPriority(req.Priority)
	if err != nil {
		return nil, err
	}
	if req.StartDate != nil && req.EndDate != nil && req.EndDate.Before(*req.StartDate) {
		return nil, apperrors.New(apperrors.ErrValidation, "end_date cannot be before start_date")
	}

	project.Name = req.Name
	project.Description = req.Description
	project.Category = category
	if req.InitiationType != nil {
		project.InitiationType = req.InitiationType
	}
	project.Priority = priority
	project.Notes = notes
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate

	rows, err := s.repo.UpdateWithVersionCheck(project, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Recompute health setelah metadata project berubah.
	_ = s.RecalculateHealth(id)

	return updated, nil
}

// RecalculateHealth menghitung ulang kesehatan project berdasarkan status,
// end_date, dan progress milestone terkini, lalu menyimpannya ke kolom `health`.
func (s *projectService) RecalculateHealth(projectID uint64) error {
	project, err := s.repo.FindByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	progress, err := s.milestoneRepo.GetAverageProgressByProject(projectID)
	if err != nil {
		return err
	}

	health := calculateProjectHealth(project, progress, time.Now())
	return s.repo.UpdateHealth(projectID, health)
}

func (s *projectService) ChangeStatus(id uint64, req dto.ChangeStatusRequest) (*entity.Project, error) {
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	newStatus := entity.ProjectStatus(req.Status)

	if err := statemachine.ValidateTransition(project.Status, newStatus); err != nil {
		return nil, err
	}

	project.Status = newStatus

	rows, err := s.repo.UpdateWithVersionCheck(project, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if newStatus == entity.ProjectCompleted {
		s.eventBus.Publish(events.Event{
			Name: "project.completed",
			Data: map[string]interface{}{
				"project_id": updated.ID,
			},
		})
	}

	// Recompute health setelah status berubah (status sangat memengaruhi health).
	_ = s.RecalculateHealth(id)

	return updated, nil
}

// CalculateProgress: Project Progress = Average(Milestone Progress) per SDD section 12.
func (s *projectService) CalculateProgress(id uint64) (float64, error) {
	return s.milestoneRepo.GetAverageProgressByProject(id)
}

func (s *projectService) CalculateHealth(project *entity.Project, progress float64) string {
	return calculateProjectHealth(project, progress, time.Now())
}

func calculateProjectHealth(project *entity.Project, progress float64, now time.Time) string {
	if project == nil {
		return "YELLOW"
	}
	if project.Status == entity.ProjectCompleted {
		return "GREEN"
	}
	if project.Status == entity.ProjectCancelled {
		return "RED"
	}
	if project.Status == entity.ProjectOnHold {
		return "YELLOW"
	}
	if project.EndDate != nil && project.EndDate.Before(startOfDay(now)) && progress < 100 {
		return "RED"
	}
	if project.EndDate != nil && project.EndDate.Before(startOfDay(now).AddDate(0, 0, 7)) && progress < 80 {
		return "YELLOW"
	}
	if project.Status == entity.ProjectActive && progress < 30 {
		return "YELLOW"
	}
	return "GREEN"
}

func startOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func normalizeProjectPriority(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	normalized = strings.ToUpper(normalized)
	normalized = strings.NewReplacer(" ", "_", "-", "_").Replace(normalized)
	if normalized == "" {
		return "MEDIUM", nil
	}
	switch normalized {
	case "LOW", "MEDIUM", "HIGH", "URGENT":
		return normalized, nil
	default:
		return "", apperrors.New(apperrors.ErrValidation, "priority is invalid")
	}
}

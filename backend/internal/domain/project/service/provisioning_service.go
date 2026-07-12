package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	budgetentity "github.com/Kal-el21/backend/internal/domain/budget/entity"
	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	requestentity "github.com/Kal-el21/backend/internal/domain/project_request/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

// ProvisioningService bertanggung jawab khusus untuk auto-create project
// dari approved request (FR-05.01).
type ProvisioningService interface {
	CreateFromApprovedRequest(request *requestentity.ProjectRequest, projectManagerID uint64) (*entity.Project, error)
	CreateDirect(createdBy uint64, req dto.CreateProjectDirectRequest) (*entity.Project, error)
}

type provisioningService struct {
	db *gorm.DB
}

func NewProvisioningService(db *gorm.DB) ProvisioningService {
	return &provisioningService{db: db}
}

func (s *provisioningService) CreateFromApprovedRequest(request *requestentity.ProjectRequest, projectManagerID uint64) (*entity.Project, error) {
	if request == nil {
		return nil, apperrors.New(apperrors.ErrValidation, "approved request is required")
	}

	var project *entity.Project
	err := s.db.Transaction(func(tx *gorm.DB) error {
		projectCode, err := generateProjectCode(tx, time.Now().Year())
		if err != nil {
			return err
		}

		project = &entity.Project{
			ProjectRequestID: &request.ID,
			ProjectCode:      projectCode,
			Name:             request.Title,
			Description:      request.Description,
			Category:         strings.TrimSpace(request.Category),
			InitiationType:   stringPtrFromInitiation(request.InitiationType),
			Priority:         projectPriority(request.Priority),
			Notes:            strings.TrimSpace(request.Notes),
			StartDate:        request.ProposedStartDate,
			EndDate:          request.ProposedEndDate,
			Status:           entity.ProjectPlanned,
			CreatedBy:        request.RequesterID,
		}

		if err := tx.Create(project).Error; err != nil {
			return err
		}

		// Admin memilih PROJECT_MANAGER pertama saat menyetujui request.
		member := &entity.ProjectMember{
			ProjectID:   project.ID,
			UserID:      projectManagerID,
			ProjectRole: entity.RoleProjectManager,
			Status:      entity.MemberActive,
			JoinedAt:    time.Now(),
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		if request.EstimatedBudget > 0 {
			budget := &budgetentity.Budget{
				ProjectID:       project.ID,
				BudgetType:      stringPtrFromBudgetType(request.BudgetType),
				BudgetName:      strings.TrimSpace(request.BudgetName),
				AllocatedBudget: request.EstimatedBudget,
				CreatedBy:       request.RequesterID,
			}
			if err := tx.Create(budget).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *provisioningService) CreateDirect(createdBy uint64, req dto.CreateProjectDirectRequest) (*entity.Project, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New(apperrors.ErrValidation, "name is required")
	}

	var project *entity.Project
	err := s.db.Transaction(func(tx *gorm.DB) error {
		projectCode, err := generateProjectCode(tx, time.Now().Year())
		if err != nil {
			return err
		}

		category := strings.TrimSpace(req.Category)
		notes := strings.TrimSpace(req.Notes)
		priority, err := normalizeProjectPriority(req.Priority)
		if err != nil {
			return err
		}
		budgetName := strings.TrimSpace(req.BudgetName)

		var initiationType *string
		if req.InitiationType != nil {
			val := *req.InitiationType
			initiationType = &val
		}

		project = &entity.Project{
			ProjectCode:      projectCode,
			Name:             name,
			Description:      req.Description,
			Category:         category,
			InitiationType:   initiationType,
			Priority:         priority,
			Notes:            notes,
			StartDate:        req.StartDate,
			EndDate:          req.EndDate,
			Status:           entity.ProjectPlanned,
			CreatedBy:        createdBy,
		}

		if err := tx.Create(project).Error; err != nil {
			return err
		}

		member := &entity.ProjectMember{
			ProjectID:   project.ID,
			UserID:      createdBy,
			ProjectRole: entity.RoleProjectManager,
			Status:      entity.MemberActive,
			JoinedAt:    time.Now(),
		}

		if err := tx.Create(member).Error; err != nil {
			return err
		}

		if req.AllocatedBudget > 0 {
			budgetType := normalizeBudgetTypeDirect(req.BudgetType)
			budget := &budgetentity.Budget{
				ProjectID:       project.ID,
				BudgetType:      budgetType,
				BudgetName:      budgetName,
				AllocatedBudget: req.AllocatedBudget,
				CreatedBy:       createdBy,
			}
			if err := tx.Create(budget).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return project, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func normalizeBudgetTypeDirect(value string) *string {
	normalized := strings.TrimSpace(value)
	normalized = strings.ToUpper(normalized)
	if normalized == "" {
		return nil
	}
	if normalized != "CAPEX" && normalized != "OPEX" {
		return nil
	}
	return &normalized
}

func generateProjectCode(tx *gorm.DB, year int) (string, error) {
	prefix := fmt.Sprintf("PRJ-%d-", year)
	var lastCode string

	if err := tx.Model(&entity.Project{}).
		Select("project_code").
		Where("project_code LIKE ?", prefix+"%").
		Order("project_code DESC").
		Limit(1).
		Scan(&lastCode).Error; err != nil {
		return "", err
	}

	next := 1
	if lastCode != "" {
		if parsed, err := strconv.Atoi(strings.TrimPrefix(lastCode, prefix)); err == nil {
			next = parsed + 1
		}
	}

	for attempt := 0; attempt < 100; attempt++ {
		code := fmt.Sprintf("%s%04d", prefix, next+attempt)
		var count int64
		if err := tx.Model(&entity.Project{}).Where("project_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}

	return "", apperrors.New(apperrors.ErrConflict, "unable to generate unique project code")
}

func stringPtrFromInitiation(value *requestentity.InitiationType) *string {
	if value == nil {
		return nil
	}
	result := string(*value)
	return &result
}

func stringPtrFromBudgetType(value *requestentity.BudgetType) *string {
	if value == nil {
		return nil
	}
	result := string(*value)
	return &result
}

func projectPriority(value requestentity.RequestPriority) string {
	priority := strings.TrimSpace(string(value))
	if priority == "" {
		return string(requestentity.PriorityMedium)
	}
	return priority
}

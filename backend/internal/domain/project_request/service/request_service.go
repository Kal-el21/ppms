package service

import (
	"errors"
	"strings"
	"time"

	projectservice "github.com/Kal-el21/backend/internal/domain/project/service"
	"github.com/Kal-el21/backend/internal/domain/project_request/dto"
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/project_request/errors"
	"github.com/Kal-el21/backend/internal/domain/project_request/repository"
	"github.com/Kal-el21/backend/internal/domain/project_request/statemachine"
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
	userrepository "github.com/Kal-el21/backend/internal/domain/user/repository"
	"github.com/Kal-el21/backend/internal/events"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type RequestService interface {
	CreateDraft(requesterID uint64, req dto.CreateDraftRequest) (*entity.ProjectRequest, error)
	UpdateDraft(id uint64, requesterID uint64, req dto.UpdateDraftRequest) (*entity.ProjectRequest, error)
	Submit(id uint64, requesterID uint64) (*entity.ProjectRequest, error)
	Review(id uint64, reviewerID uint64, req dto.ReviewRequest, isAdmin bool) (*entity.ProjectRequest, error)
	Revise(id uint64, requesterID uint64, req dto.ReviseRequest) (*entity.ProjectRequest, error)
	GetByID(id uint64, requesterID uint64, isAdmin bool) (*entity.ProjectRequest, error)
	GetOwnRequests(requesterID uint64, page, limit int) ([]entity.ProjectRequest, int64, error)
	GetAllRequests(page, limit int, statusFilter string) ([]entity.ProjectRequest, int64, error)
	GetRevisionHistory(id uint64) ([]entity.ProjectRequestRevision, error)
	GetApprovalHistory(id uint64) ([]entity.ProjectRequestApproval, error)
	DeleteDraft(id uint64, requesterID uint64) error
}

type requestService struct {
	requestRepo     repository.RequestRepository
	revisionRepo    repository.RevisionRepository
	approvalRepo    repository.ApprovalRepository
	userRepo        userrepository.UserRepository
	provisioningSvc projectservice.ProvisioningService
	eventBus        *events.Bus
}

func NewRequestService(
	requestRepo repository.RequestRepository,
	revisionRepo repository.RevisionRepository,
	approvalRepo repository.ApprovalRepository,
	userRepo userrepository.UserRepository,
	provisioningSvc projectservice.ProvisioningService,
	eventBus *events.Bus,
) RequestService {
	return &requestService{
		requestRepo:     requestRepo,
		revisionRepo:    revisionRepo,
		approvalRepo:    approvalRepo,
		userRepo:        userRepo,
		provisioningSvc: provisioningSvc,
		eventBus:        eventBus,
	}
}

var (
	allowedInitiationTypes = map[string]struct{}{
		string(entity.InitiationNewInitiative): {},
		string(entity.InitiationRenewal):       {},
		string(entity.InitiationEnhancement):   {},
	}
	allowedPriorities = map[string]struct{}{
		string(entity.PriorityLow):    {},
		string(entity.PriorityMedium): {},
		string(entity.PriorityHigh):   {},
		string(entity.PriorityUrgent): {},
	}
	allowedBudgetTypes = map[string]struct{}{
		string(entity.BudgetTypeCapex): {},
		string(entity.BudgetTypeOpex):  {},
	}
)

func cleanEnum(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToUpper(value)
	value = strings.NewReplacer(" ", "_", "-", "_").Replace(value)
	return value
}

func normalizeOptionalEnum(value string, allowed map[string]struct{}, fieldName string) (string, error) {
	normalized := cleanEnum(value)
	if normalized == "" {
		return "", nil
	}
	if _, ok := allowed[normalized]; !ok {
		return "", apperrors.New(apperrors.ErrValidation, fieldName+" is invalid")
	}
	return normalized, nil
}

func normalizeEnumWithDefault(value string, allowed map[string]struct{}, defaultValue string, fieldName string) (string, error) {
	normalized := cleanEnum(value)
	if normalized == "" {
		normalized = defaultValue
	}
	if _, ok := allowed[normalized]; !ok {
		return "", apperrors.New(apperrors.ErrValidation, fieldName+" is invalid")
	}
	return normalized, nil
}

func trimString(value string, maxLength int, fieldName string) (string, error) {
	value = strings.TrimSpace(value)
	if maxLength > 0 && len(value) > maxLength {
		return "", apperrors.New(apperrors.ErrValidation, fieldName+" is too long")
	}
	return value, nil
}

func applyRequestMetadata(request *entity.ProjectRequest, metadata dto.ProjectRequestMetadata) error {
	category, err := trimString(metadata.Category, 100, "category")
	if err != nil {
		return err
	}
	budgetName, err := trimString(metadata.BudgetName, 200, "budget_name")
	if err != nil {
		return err
	}
	notes, err := trimString(metadata.Notes, 2000, "notes")
	if err != nil {
		return err
	}
	initiationType, err := normalizeOptionalEnum(metadata.InitiationType, allowedInitiationTypes, "initiation_type")
	if err != nil {
		return err
	}
	priority, err := normalizeEnumWithDefault(metadata.Priority, allowedPriorities, string(entity.PriorityMedium), "priority")
	if err != nil {
		return err
	}
	budgetType, err := normalizeOptionalEnum(metadata.BudgetType, allowedBudgetTypes, "budget_type")
	if err != nil {
		return err
	}

	if metadata.ProposedStartDate != nil && metadata.ProposedEndDate != nil &&
		metadata.ProposedEndDate.Before(*metadata.ProposedStartDate) {
		return apperrors.New(apperrors.ErrValidation, "proposed_end_date cannot be before proposed_start_date")
	}

	request.Category = category
	if initiationType == "" {
		request.InitiationType = nil
	} else {
		value := entity.InitiationType(initiationType)
		request.InitiationType = &value
	}
	request.Priority = entity.RequestPriority(priority)
	request.ProposedStartDate = metadata.ProposedStartDate
	request.ProposedEndDate = metadata.ProposedEndDate
	if budgetType == "" {
		request.BudgetType = nil
	} else {
		value := entity.BudgetType(budgetType)
		request.BudgetType = &value
	}
	request.BudgetName = budgetName
	request.Notes = notes

	return nil
}

func (s *requestService) CreateDraft(requesterID uint64, req dto.CreateDraftRequest) (*entity.ProjectRequest, error) {
	request := &entity.ProjectRequest{
		RequesterID:     requesterID,
		Title:           req.Title,
		Description:     req.Description,
		BusinessGoal:    req.BusinessGoal,
		ExpectedOutcome: req.ExpectedOutcome,
		EstimatedBudget: req.EstimatedBudget,
		Status:          entity.StatusDraft,
	}

	if err := applyRequestMetadata(request, req.ProjectRequestMetadata); err != nil {
		return nil, err
	}

	if err := s.requestRepo.Create(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *requestService) getOwnedRequest(id uint64, requesterID uint64) (*entity.ProjectRequest, error) {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrRequestNotFound
		}
		return nil, err
	}

	// Ownership Check (Permission Evaluation Order #2)
	if request.RequesterID != requesterID {
		return nil, domainerrors.ErrNotOwner
	}

	return request, nil
}

func (s *requestService) UpdateDraft(id uint64, requesterID uint64, req dto.UpdateDraftRequest) (*entity.ProjectRequest, error) {
	request, err := s.getOwnedRequest(id, requesterID)
	if err != nil {
		return nil, err
	}

	if request.Status != entity.StatusDraft {
		return nil, domainerrors.ErrCannotEditNonDraft
	}

	request.Title = req.Title
	request.Description = req.Description
	request.BusinessGoal = req.BusinessGoal
	request.ExpectedOutcome = req.ExpectedOutcome
	request.EstimatedBudget = req.EstimatedBudget

	if err := applyRequestMetadata(request, req.ProjectRequestMetadata); err != nil {
		return nil, err
	}

	rows, err := s.requestRepo.UpdateWithVersionCheck(request, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.requestRepo.FindByID(id)
}

func (s *requestService) Submit(id uint64, requesterID uint64) (*entity.ProjectRequest, error) {
	request, err := s.getOwnedRequest(id, requesterID)
	if err != nil {
		return nil, err
	}

	// Submit dapat dilakukan dari status DRAFT atau REVISED (FR-04.06 state machine)
	var targetTransitionCheck entity.RequestStatus
	if request.Status == entity.StatusDraft || request.Status == entity.StatusRevised {
		targetTransitionCheck = request.Status
	} else {
		return nil, apperrors.New(apperrors.ErrInvalidStateTransition, "request cannot be submitted from current status")
	}

	if err := statemachine.ValidateTransition(targetTransitionCheck, entity.StatusSubmitted); err != nil {
		return nil, err
	}

	now := time.Now()
	request.Status = entity.StatusSubmitted
	request.SubmittedAt = &now

	rows, err := s.requestRepo.UpdateWithVersionCheck(request, request.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	updated, err := s.requestRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Auto-transition ke UNDER_REVIEW agar masuk antrian review ADMIN
	updated.Status = entity.StatusUnderReview
	if _, err := s.requestRepo.UpdateWithVersionCheck(updated, updated.Version); err != nil {
		return nil, err
	}

	s.eventBus.Publish(events.Event{
		Name: "project.request.submitted",
		Data: map[string]interface{}{
			"request_id":   updated.ID,
			"title":        updated.Title,
			"requester_id": updated.RequesterID,
		},
	})

	return s.requestRepo.FindByID(id)
}

func (s *requestService) Review(id uint64, reviewerID uint64, req dto.ReviewRequest, isAdmin bool) (*entity.ProjectRequest, error) {
	if !isAdmin {
		return nil, apperrors.New(apperrors.ErrInsufficientSystemRole, "only ADMIN can review requests")
	}

	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrRequestNotFound
		}
		return nil, err
	}

	if request.Status != entity.StatusUnderReview && request.Status != entity.StatusRevised {
		return nil, apperrors.New(apperrors.ErrInvalidStateTransition, "only requests under review or revised requests can be reviewed")
	}

	var newStatus entity.RequestStatus
	var eventName string
	now := time.Now()
	var selectedProjectManagerID *uint64

	switch req.Action {
	case "APPROVED":
		projectManagerID, err := s.validateProjectManager(req.ProjectManagerID)
		if err != nil {
			return nil, err
		}
		selectedProjectManagerID = &projectManagerID
		newStatus = entity.StatusApproved
		eventName = "project.request.approved"
		request.ApprovedAt = &now
		request.RejectedAt = nil
	case "REJECTED":
		newStatus = entity.StatusRejected
		eventName = "project.request.rejected"
		request.RejectedAt = &now
		request.ApprovedAt = nil
	case "REQUEST_REVISION":
		newStatus = entity.StatusRevisionRequested
		eventName = "project.request.revision_requested"
		request.ApprovedAt = nil
		request.RejectedAt = nil
	default:
		return nil, apperrors.New(apperrors.ErrValidation, "invalid review action")
	}

	if err := statemachine.ValidateTransition(request.Status, newStatus); err != nil {
		return nil, err
	}

	request.Status = newStatus
	rows, err := s.requestRepo.UpdateWithVersionCheck(request, request.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	approval := &entity.ProjectRequestApproval{
		ProjectRequestID: id,
		ReviewedBy:       reviewerID,
		Action:           entity.ApprovalAction(req.Action),
		Comment:          req.Comment,
		ProjectManagerID: selectedProjectManagerID,
	}

	if err := s.approvalRepo.Create(approval); err != nil {
		return nil, err
	}

	updated, err := s.requestRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// FR-05.01: Auto-create project saat request APPROVED
	if newStatus == entity.StatusApproved {
		createdProject, err := s.provisioningSvc.CreateFromApprovedRequest(updated, *selectedProjectManagerID)
		if err != nil {
			return nil, err
		}
		// Picu perhitungan health awal project yang baru dibuat.
		s.eventBus.Publish(events.Event{
			Name: "project.created",
			Data: map[string]interface{}{"project_id": createdProject.ID},
		})
	}

	eventData := map[string]interface{}{
		"request_id":   updated.ID,
		"title":        updated.Title,
		"requester_id": updated.RequesterID,
		"reviewer_id":  reviewerID,
		"comment":      req.Comment,
	}
	if selectedProjectManagerID != nil {
		eventData["project_manager_id"] = *selectedProjectManagerID
	}

	s.eventBus.Publish(events.Event{
		Name: eventName,
		Data: eventData,
	})

	return updated, nil
}

func (s *requestService) validateProjectManager(projectManagerID *uint64) (uint64, error) {
	if projectManagerID == nil || *projectManagerID == 0 {
		return 0, apperrors.New(apperrors.ErrValidation, "project manager is required when approving a request")
	}

	pm, err := s.userRepo.FindByID(*projectManagerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, apperrors.New(apperrors.ErrValidation, "selected project manager was not found")
		}
		return 0, err
	}
	if !pm.IsActive {
		return 0, apperrors.New(apperrors.ErrValidation, "selected project manager must be an active user")
	}
	if pm.SystemRole == userentity.RoleViewer {
		return 0, apperrors.New(apperrors.ErrValidation, "viewer users cannot be assigned as project manager")
	}

	return pm.ID, nil
}

func (s *requestService) Revise(id uint64, requesterID uint64, req dto.ReviseRequest) (*entity.ProjectRequest, error) {
	request, err := s.getOwnedRequest(id, requesterID)
	if err != nil {
		return nil, err
	}

	transitionFrom := request.Status
	if request.Status == entity.StatusRejected {
		isLegacyRevisionRequest, err := s.isLegacyRevisionRequest(id)
		if err != nil {
			return nil, err
		}
		if isLegacyRevisionRequest {
			transitionFrom = entity.StatusRevisionRequested
		}
	}

	if err := statemachine.ValidateTransition(transitionFrom, entity.StatusRevised); err != nil {
		return nil, err
	}

	// Simpan snapshot versi LAMA ke revision history sebelum overwrite (FR-04.04)
	revisionCount, err := s.revisionRepo.CountByRequestID(id)
	if err != nil {
		return nil, err
	}

	revision := &entity.ProjectRequestRevision{
		ProjectRequestID:  id,
		RevisionNumber:    int(revisionCount) + 1,
		Title:             request.Title,
		Description:       request.Description,
		BusinessGoal:      request.BusinessGoal,
		ExpectedOutcome:   request.ExpectedOutcome,
		EstimatedBudget:   request.EstimatedBudget,
		Category:          request.Category,
		InitiationType:    request.InitiationType,
		Priority:          request.Priority,
		ProposedStartDate: request.ProposedStartDate,
		ProposedEndDate:   request.ProposedEndDate,
		BudgetType:        request.BudgetType,
		BudgetName:        request.BudgetName,
		Notes:             request.Notes,
		RevisionReason:    req.RevisionReason,
		RevisedBy:         requesterID,
	}

	if err := s.revisionRepo.Create(revision); err != nil {
		return nil, err
	}

	// Update request dengan data baru + transisi status ke REVISED
	request.Title = req.Title
	request.Description = req.Description
	request.BusinessGoal = req.BusinessGoal
	request.ExpectedOutcome = req.ExpectedOutcome
	request.EstimatedBudget = req.EstimatedBudget

	if err := applyRequestMetadata(request, req.ProjectRequestMetadata); err != nil {
		return nil, err
	}

	request.CurrentRevision = request.CurrentRevision + 1
	request.Status = entity.StatusRevised

	rows, err := s.requestRepo.UpdateWithVersionCheck(request, request.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	s.eventBus.Publish(events.Event{
		Name: "project.request.revised",
		Data: map[string]interface{}{
			"request_id":   request.ID,
			"requester_id": requesterID,
			"title":        request.Title,
		},
	})

	return s.requestRepo.FindByID(id)
}

func (s *requestService) isLegacyRevisionRequest(id uint64) (bool, error) {
	approvals, err := s.approvalRepo.FindByRequestID(id)
	if err != nil {
		return false, err
	}
	if len(approvals) == 0 {
		return false, nil
	}

	return approvals[len(approvals)-1].Action == entity.ActionRequestRevision, nil
}

func (s *requestService) GetByID(id uint64, requesterID uint64, isAdmin bool) (*entity.ProjectRequest, error) {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrRequestNotFound
		}
		return nil, err
	}

	// ADMIN Override (Permission Evaluation Order #1), else Ownership Check (#2)
	if !isAdmin && request.RequesterID != requesterID {
		return nil, domainerrors.ErrNotOwner
	}

	return request, nil
}

func (s *requestService) GetOwnRequests(requesterID uint64, page, limit int) ([]entity.ProjectRequest, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.requestRepo.FindOwnRequests(requesterID, page, limit)
}

func (s *requestService) GetAllRequests(page, limit int, statusFilter string) ([]entity.ProjectRequest, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.requestRepo.FindAllRequests(page, limit, statusFilter)
}

func (s *requestService) GetRevisionHistory(id uint64) ([]entity.ProjectRequestRevision, error) {
	return s.revisionRepo.FindByRequestID(id)
}

func (s *requestService) GetApprovalHistory(id uint64) ([]entity.ProjectRequestApproval, error) {
	return s.approvalRepo.FindByRequestID(id)
}

func (s *requestService) DeleteDraft(id uint64, requesterID uint64) error {
	request, err := s.getOwnedRequest(id, requesterID)
	if err != nil {
		return err
	}

	if request.Status != entity.StatusDraft {
		return domainerrors.ErrCannotDeleteNonDraft
	}

	return s.requestRepo.Delete(id, requesterID)
}

package service

import (
	"errors"
	"time"

	projectservice "github.com/Kal-el21/backend/internal/domain/project/service"
	"github.com/Kal-el21/backend/internal/domain/project_request/dto"
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/project_request/errors"
	"github.com/Kal-el21/backend/internal/domain/project_request/repository"
	"github.com/Kal-el21/backend/internal/domain/project_request/statemachine"
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
	provisioningSvc projectservice.ProvisioningService
	eventBus        *events.Bus
}

func NewRequestService(
	requestRepo repository.RequestRepository,
	revisionRepo repository.RevisionRepository,
	approvalRepo repository.ApprovalRepository,
	provisioningSvc projectservice.ProvisioningService,
	eventBus *events.Bus,
) RequestService {
	return &requestService{
		requestRepo:     requestRepo,
		revisionRepo:    revisionRepo,
		approvalRepo:    approvalRepo,
		provisioningSvc: provisioningSvc,
		eventBus:        eventBus,
	}
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

	if request.Status != entity.StatusUnderReview {
		return nil, apperrors.New(apperrors.ErrInvalidStateTransition, "only requests under review can be reviewed")
	}

	var newStatus entity.RequestStatus
	var eventName string

	switch req.Action {
	case "APPROVED":
		newStatus = entity.StatusApproved
		eventName = "project.request.approved"
	case "REJECTED", "REQUEST_REVISION":
		newStatus = entity.StatusRejected
		eventName = "project.request.rejected"
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
		if _, err := s.provisioningSvc.CreateFromApprovedRequest(
			updated.ID, updated.Title, updated.Description, updated.RequesterID,
		); err != nil {
			return nil, err
		}
	}

	s.eventBus.Publish(events.Event{
		Name: eventName,
		Data: map[string]interface{}{
			"request_id":   updated.ID,
			"title":        updated.Title,
			"requester_id": updated.RequesterID,
			"reviewer_id":  reviewerID,
			"comment":      req.Comment,
		},
	})

	return updated, nil
}

func (s *requestService) Revise(id uint64, requesterID uint64, req dto.ReviseRequest) (*entity.ProjectRequest, error) {
	request, err := s.getOwnedRequest(id, requesterID)
	if err != nil {
		return nil, err
	}

	if err := statemachine.ValidateTransition(request.Status, entity.StatusRevised); err != nil {
		return nil, err
	}

	// Simpan snapshot versi LAMA ke revision history sebelum overwrite (FR-04.04)
	revisionCount, err := s.revisionRepo.CountByRequestID(id)
	if err != nil {
		return nil, err
	}

	revision := &entity.ProjectRequestRevision{
		ProjectRequestID: id,
		RevisionNumber:   int(revisionCount) + 1,
		Title:            request.Title,
		Description:      request.Description,
		BusinessGoal:     request.BusinessGoal,
		ExpectedOutcome:  request.ExpectedOutcome,
		EstimatedBudget:  request.EstimatedBudget,
		RevisionReason:   req.RevisionReason,
		RevisedBy:        requesterID,
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
	request.Status = entity.StatusRevised

	rows, err := s.requestRepo.UpdateWithVersionCheck(request, request.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.requestRepo.FindByID(id)
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

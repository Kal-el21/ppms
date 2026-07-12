package repository

import (
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	"gorm.io/gorm"
)

type RequestRepository interface {
	Create(req *entity.ProjectRequest) error
	FindByID(id uint64) (*entity.ProjectRequest, error)
	FindOwnRequests(requesterID uint64, page, limit int) ([]entity.ProjectRequest, int64, error)
	FindAllRequests(page, limit int, statusFilter string) ([]entity.ProjectRequest, int64, error)
	UpdateWithVersionCheck(req *entity.ProjectRequest, expectedVersion int) (int64, error)
	Delete(id uint64, deletedBy uint64) error
}

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(req *entity.ProjectRequest) error {
	return r.db.Create(req).Error
}

func (r *requestRepository) FindByID(id uint64) (*entity.ProjectRequest, error) {
	var req entity.ProjectRequest
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *requestRepository) FindOwnRequests(requesterID uint64, page, limit int) ([]entity.ProjectRequest, int64, error) {
	var requests []entity.ProjectRequest
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.ProjectRequest{}).Where("requester_id = ? AND deleted_at IS NULL", requesterID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&requests).Error
	return requests, total, err
}

func (r *requestRepository) FindAllRequests(page, limit int, statusFilter string) ([]entity.ProjectRequest, int64, error) {
	var requests []entity.ProjectRequest
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.ProjectRequest{}).Where("deleted_at IS NULL")
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&requests).Error
	return requests, total, err
}

// UpdateWithVersionCheck implements optimistic locking (FR-17.01).
// Returns rows affected; 0 rows means version mismatch (concurrent edit).
func (r *requestRepository) UpdateWithVersionCheck(req *entity.ProjectRequest, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.ProjectRequest{}).
		Where("id = ? AND version = ?", req.ID, expectedVersion).
		Updates(map[string]interface{}{
			"title":               req.Title,
			"description":         req.Description,
			"business_goal":       req.BusinessGoal,
			"expected_outcome":    req.ExpectedOutcome,
			"estimated_budget":    req.EstimatedBudget,
			"category":            req.Category,
			"initiation_type":     req.InitiationType,
			"priority":            req.Priority,
			"proposed_start_date": req.ProposedStartDate,
			"proposed_end_date":   req.ProposedEndDate,
			"budget_type":         req.BudgetType,
			"budget_name":         req.BudgetName,
			"notes":               req.Notes,
			"status":              req.Status,
			"current_revision":    req.CurrentRevision,
			"submitted_at":        req.SubmittedAt,
			"approved_at":         req.ApprovedAt,
			"rejected_at":         req.RejectedAt,
			"version":             gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *requestRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.ProjectRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}

package repository

import (
	"github.com/Kal-el21/backend/internal/domain/audit/entity"
	"gorm.io/gorm"
)

type AuditFilter struct {
	UserID     *uint64
	Module     string
	EntityType string
	EntityID   *uint64
	Page       int
	Limit      int
}

type AuditRepository interface {
	Create(log *entity.AuditLog) error
	FindWithFilter(filter AuditFilter) ([]entity.AuditLog, int64, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(log *entity.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditRepository) FindWithFilter(filter AuditFilter) ([]entity.AuditLog, int64, error) {
	var logs []entity.AuditLog
	var total int64

	query := r.db.Model(&entity.AuditLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Module != "" {
		query = query.Where("module = ?", filter.Module)
	}
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

package repository

import (
	"github.com/Kal-el21/backend/internal/domain/project_request/entity"
	"gorm.io/gorm"
)

type RevisionRepository interface {
	Create(revision *entity.ProjectRequestRevision) error
	FindByRequestID(requestID uint64) ([]entity.ProjectRequestRevision, error)
	CountByRequestID(requestID uint64) (int64, error)
}

type revisionRepository struct {
	db *gorm.DB
}

func NewRevisionRepository(db *gorm.DB) RevisionRepository {
	return &revisionRepository{db: db}
}

func (r *revisionRepository) Create(revision *entity.ProjectRequestRevision) error {
	return r.db.Create(revision).Error
}

func (r *revisionRepository) FindByRequestID(requestID uint64) ([]entity.ProjectRequestRevision, error) {
	var revisions []entity.ProjectRequestRevision
	err := r.db.Where("project_request_id = ?", requestID).
		Order("revision_number ASC").
		Find(&revisions).Error
	return revisions, err
}

func (r *revisionRepository) CountByRequestID(requestID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&entity.ProjectRequestRevision{}).
		Where("project_request_id = ?", requestID).
		Count(&count).Error
	return count, err
}

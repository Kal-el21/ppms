package repository

import (
	"time"
	"github.com/Kal-el21/backend/internal/domain/reporting/entity"
	"gorm.io/gorm"
)

type ProjectSnapshotRepository interface {
	Create(snapshot *entity.ProjectSnapshot) error
	GetByProjectID(projectID uint64, fromDate, toDate *time.Time) ([]entity.ProjectSnapshot, error)
	GetLatestByProjectID(projectID uint64) (*entity.ProjectSnapshot, error)
}

type projectSnapshotRepository struct {
	db *gorm.DB
}

func NewProjectSnapshotRepository(db *gorm.DB) ProjectSnapshotRepository {
	return &projectSnapshotRepository{db: db}
}

func (r *projectSnapshotRepository) Create(snapshot *entity.ProjectSnapshot) error {
	return r.db.Create(snapshot).Error
}

func (r *projectSnapshotRepository) GetByProjectID(projectID uint64, fromDate, toDate *time.Time) ([]entity.ProjectSnapshot, error) {
	var snapshots []entity.ProjectSnapshot
	query := r.db.Where("project_id = ?", projectID)

	if fromDate != nil {
		query = query.Where("snapshot_date >= ?", *fromDate)
	}
	if toDate != nil {
		query = query.Where("snapshot_date <= ?", *toDate)
	}

	err := query.Order("snapshot_date DESC").Find(&snapshots).Error
	return snapshots, err
}

func (r *projectSnapshotRepository) GetLatestByProjectID(projectID uint64) (*entity.ProjectSnapshot, error) {
	var snapshot entity.ProjectSnapshot
	err := r.db.Where("project_id = ?", projectID).Order("snapshot_date DESC").First(&snapshot).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}

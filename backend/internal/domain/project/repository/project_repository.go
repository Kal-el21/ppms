package repository

import (
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *entity.Project) error
	FindByID(id uint64) (*entity.Project, error)
	FindAll(page, limit int, statusFilter string) ([]entity.Project, int64, error)
	FindByMemberUserID(userID uint64, page, limit int) ([]entity.Project, int64, error)
	UpdateWithVersionCheck(project *entity.Project, expectedVersion int) (int64, error)
	Delete(id uint64, deletedBy uint64) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *entity.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uint64) (*entity.Project, error) {
	var project entity.Project
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindAll(page, limit int, statusFilter string) ([]entity.Project, int64, error) {
	var projects []entity.Project
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.Project{}).Where("deleted_at IS NULL")
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&projects).Error
	return projects, total, err
}

// FindByMemberUserID mengambil semua project di mana user adalah member aktif.
// Digunakan untuk USER/VIEWER yang hanya boleh lihat project yang mereka ikuti.
func (r *projectRepository) FindByMemberUserID(userID uint64, page, limit int) ([]entity.Project, int64, error) {
	var projects []entity.Project
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.Project{}).
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ? AND project_members.status = 'ACTIVE' AND projects.deleted_at IS NULL", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("projects.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) UpdateWithVersionCheck(project *entity.Project, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Project{}).
		Where("id = ? AND version = ?", project.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"start_date":  project.StartDate,
			"end_date":    project.EndDate,
			"status":      project.Status,
			"version":     gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *projectRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}

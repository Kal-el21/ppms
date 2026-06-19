package repository

import (
	"github.com/Kal-el21/backend/internal/domain/division/entity"
	"gorm.io/gorm"
)

type DivisionRepository interface {
	Create(division *entity.Division) error
	FindByID(id uint64) (*entity.Division, error)
	FindByName(name string) (*entity.Division, error)
	FindAll() ([]entity.Division, error)
	Update(division *entity.Division) error
	Delete(id uint64, deletedBy uint64) error
}

type divisionRepository struct {
	db *gorm.DB
}

func NewDivisionRepository(db *gorm.DB) DivisionRepository {
	return &divisionRepository{db: db}
}

func (r *divisionRepository) Create(division *entity.Division) error {
	return r.db.Create(division).Error
}

func (r *divisionRepository) FindByID(id uint64) (*entity.Division, error) {
	var division entity.Division
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&division).Error
	if err != nil {
		return nil, err
	}
	return &division, nil
}

func (r *divisionRepository) FindByName(name string) (*entity.Division, error) {
	var division entity.Division
	err := r.db.Where("name = ? AND deleted_at IS NULL", name).First(&division).Error
	if err != nil {
		return nil, err
	}
	return &division, nil
}

func (r *divisionRepository) FindAll() ([]entity.Division, error) {
	var divisions []entity.Division
	err := r.db.Where("deleted_at IS NULL").Order("name ASC").Find(&divisions).Error
	return divisions, err
}

func (r *divisionRepository) Update(division *entity.Division) error {
	return r.db.Save(division).Error
}

func (r *divisionRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Division{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}

package repository

import (
	"github.com/Kal-el21/backend/internal/domain/user/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint64) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindBySystemRole(role entity.SystemRole) ([]entity.User, error)
	FindAll(page, limit int) ([]entity.User, int64, error)
	Update(user *entity.User) error
	Deactivate(id uint64) error
	Restore(id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindBySystemRole(role entity.SystemRole) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("system_role = ? AND is_active = true AND deleted_at IS NULL", role).
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindAll(page, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&entity.User{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Deactivate(id uint64) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *userRepository) Restore(id uint64) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  true,
			"deleted_at": nil,
			"deleted_by": nil,
		}).Error
}

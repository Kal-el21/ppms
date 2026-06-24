package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/user/dto"
	"github.com/Kal-el21/backend/internal/domain/user/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/user/errors"
	"github.com/Kal-el21/backend/internal/domain/user/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(req dto.CreateUserRequest, bcryptCost int) (*entity.User, error)
	GetByID(id uint64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetAll(page, limit int) ([]entity.User, int64, error)
	Update(id uint64, req dto.UpdateUserRequest) (*entity.User, error)
	AssignRole(id uint64, req dto.AssignRoleRequest) (*entity.User, error)
	Deactivate(id uint64, deletedBy uint64) error
	Restore(id uint64) error
	UpdateProfile(id uint64, req dto.UpdateProfileRequest) (*entity.User, error)
	UpdateProfilePhoto(id uint64, photoURL string) error
	Toggle2FA(id uint64, enabled bool) error
	ToggleEmailNotification(id uint64, enabled bool) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req dto.CreateUserRequest, bcryptCost int) (*entity.User, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, domainerrors.ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		SystemRole:   entity.SystemRole(req.SystemRole),
		DivisionID:   req.DivisionID,
		IsActive:     true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id uint64) (*entity.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByEmail(email string) (*entity.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAll(page, limit int) ([]entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindAll(page, limit)
}

func (s *userService) Update(id uint64, req dto.UpdateUserRequest) (*entity.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName
	user.DivisionID = req.DivisionID
	user.Version++

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) AssignRole(id uint64, req dto.AssignRoleRequest) (*entity.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.SystemRole = entity.SystemRole(req.SystemRole)
	user.Version++

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Deactivate(id uint64, deletedBy uint64) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	user, _ := s.repo.FindByID(id)
	user.IsActive = false
	user.DeletedBy = &deletedBy

	return s.repo.Update(user)
}

func (s *userService) Restore(id uint64) error {
	return s.repo.Restore(id)
}

func (s *userService) UpdateProfile(id uint64, req dto.UpdateProfileRequest) (*entity.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName
	user.Version++

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfilePhoto(id uint64, photoURL string) error {
	return s.repo.UpdateField(id, "profile_photo_url", photoURL)
}

func (s *userService) Toggle2FA(id uint64, enabled bool) error {
	return s.repo.UpdateField(id, "two_fa_enabled", enabled)
}

func (s *userService) ToggleEmailNotification(id uint64, enabled bool) error {
	return s.repo.UpdateField(id, "email_notification_enabled", enabled)
}

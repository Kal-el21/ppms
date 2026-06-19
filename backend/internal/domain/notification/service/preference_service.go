package service

import (
	"github.com/Kal-el21/backend/internal/domain/notification/dto"
	"github.com/Kal-el21/backend/internal/domain/notification/entity"
	"github.com/Kal-el21/backend/internal/domain/notification/repository"
)

// notificationTypes adalah daftar semua tipe notifikasi yang didukung sistem,
// dipakai untuk menampilkan preference list lengkap meski user belum pernah set apa pun.
var notificationTypes = []string{
	"REQUEST_SUBMITTED", "REQUEST_APPROVED", "REQUEST_REJECTED",
	"TASK_ASSIGNED", "TASK_COMPLETED",
	"BUDGET_WARNING", "BUDGET_OVER_LIMIT",
	"HANDOVER_SENT", "HANDOVER_RECEIVED",
}

type PreferenceService interface {
	GetAll(userID uint64) ([]dto.PreferenceResponse, error)
	Update(userID uint64, req dto.UpdatePreferenceRequest) error
}

type preferenceService struct {
	repo repository.PreferenceRepository
}

func NewPreferenceService(repo repository.PreferenceRepository) PreferenceService {
	return &preferenceService{repo: repo}
}

func (s *preferenceService) GetAll(userID uint64) ([]dto.PreferenceResponse, error) {
	existing, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool)
	for _, p := range existing {
		existingMap[p.Type] = p.Enabled
	}

	result := make([]dto.PreferenceResponse, len(notificationTypes))
	for i, t := range notificationTypes {
		enabled, exists := existingMap[t]
		if !exists {
			enabled = true // default enabled
		}
		result[i] = dto.PreferenceResponse{Type: t, Enabled: enabled}
	}

	return result, nil
}

func (s *preferenceService) Update(userID uint64, req dto.UpdatePreferenceRequest) error {
	pref := &entity.NotificationPreference{
		UserID:  userID,
		Type:    req.Type,
		Enabled: req.Enabled,
	}
	return s.repo.Upsert(pref)
}

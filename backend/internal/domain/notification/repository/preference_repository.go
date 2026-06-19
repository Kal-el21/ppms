package repository

import (
	"github.com/Kal-el21/backend/internal/domain/notification/entity"
	"gorm.io/gorm"
)

type PreferenceRepository interface {
	FindByUserAndType(userID uint64, notifType string) (*entity.NotificationPreference, error)
	FindByUserID(userID uint64) ([]entity.NotificationPreference, error)
	Upsert(pref *entity.NotificationPreference) error
	IsEnabled(userID uint64, notifType string) (bool, error)
}

type preferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) PreferenceRepository {
	return &preferenceRepository{db: db}
}

func (r *preferenceRepository) FindByUserAndType(userID uint64, notifType string) (*entity.NotificationPreference, error) {
	var pref entity.NotificationPreference
	err := r.db.Where("user_id = ? AND type = ?", userID, notifType).First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *preferenceRepository) FindByUserID(userID uint64) ([]entity.NotificationPreference, error) {
	var prefs []entity.NotificationPreference
	err := r.db.Where("user_id = ?", userID).Find(&prefs).Error
	return prefs, err
}

func (r *preferenceRepository) Upsert(pref *entity.NotificationPreference) error {
	return r.db.Clauses(
	// ON CONFLICT (user_id, type) DO UPDATE — sesuai unique constraint Phase 0
	// uq_notification_preference
	).Where("user_id = ? AND type = ?", pref.UserID, pref.Type).
		Assign(map[string]interface{}{"enabled": pref.Enabled}).
		FirstOrCreate(pref).Error
}

// IsEnabled: default true jika belum ada preference record sama sekali
// (opt-out model — user mendapat semua notifikasi kecuali secara eksplisit mute).
func (r *preferenceRepository) IsEnabled(userID uint64, notifType string) (bool, error) {
	pref, err := r.FindByUserAndType(userID, notifType)
	if err != nil {
		return true, nil // no record found = default enabled
	}
	return pref.Enabled, nil
}

package service

import (
	"encoding/json"

	"github.com/Kal-el21/backend/internal/domain/audit/entity"
	"github.com/Kal-el21/backend/internal/domain/audit/repository"
)

type AuditService interface {
	Log(params LogParams) error
	Query(filter repository.AuditFilter) ([]entity.AuditLog, int64, error)
}

type LogParams struct {
	UserID     *uint64
	Module     string
	Action     string
	EntityType string
	EntityID   *uint64
	OldData    interface{}
	NewData    interface{}
	IPAddress  string
	UserAgent  string
	RequestID  string
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

var sensitiveFields = []string{
	"password_hash", "password", "refresh_token_hash",
	"access_token", "refresh_token", // jaga-jaga jika ada object yang menyimpan token mentah
	"old_password", "new_password", // dari ChangePasswordRequest
}

func redact(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return string(raw), nil
	}

	for _, field := range sensitiveFields {
		if _, exists := m[field]; exists {
			m[field] = "[REDACTED]"
		}
	}

	redacted, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(redacted), nil
}

func (s *auditService) Log(params LogParams) error {
	oldDataStr, err := redact(params.OldData)
	if err != nil {
		return err
	}

	newDataStr, err := redact(params.NewData)
	if err != nil {
		return err
	}

	log := &entity.AuditLog{
		UserID:     params.UserID,
		Module:     params.Module,
		Action:     params.Action,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		OldData:    oldDataStr,
		NewData:    newDataStr,
		IPAddress:  params.IPAddress,
		UserAgent:  params.UserAgent,
		RequestID:  params.RequestID,
	}

	return s.repo.Create(log)
}

func (s *auditService) Query(filter repository.AuditFilter) ([]entity.AuditLog, int64, error) {
	return s.repo.FindWithFilter(filter)
}

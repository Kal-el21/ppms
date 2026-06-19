package service

import (
	attachmentservice "github.com/Kal-el21/backend/internal/domain/attachment/service"
	"github.com/Kal-el21/backend/internal/domain/project_request/repository"
)

// requestOwnershipAdapter mengimplementasikan attachmentservice.RequestOwnershipChecker
// menggunakan RequestRepository yang sudah ada (Phase 2), tanpa attachment domain
// perlu import project_request secara langsung di level repository/entity.
type requestOwnershipAdapter struct {
	repo repository.RequestRepository
}

func NewRequestOwnershipAdapter(repo repository.RequestRepository) attachmentservice.RequestOwnershipChecker {
	return &requestOwnershipAdapter{repo: repo}
}

func (a *requestOwnershipAdapter) IsRequestOwner(requestID uint64, userID uint64) (bool, error) {
	request, err := a.repo.FindByID(requestID)
	if err != nil {
		return false, err
	}
	return request.RequesterID == userID, nil
}

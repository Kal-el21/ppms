package service

import (
	"errors"
	"time"

	"github.com/Kal-el21/backend/internal/domain/handover/dto"
	"github.com/Kal-el21/backend/internal/domain/handover/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/handover/errors"
	"github.com/Kal-el21/backend/internal/domain/handover/repository"
	"github.com/Kal-el21/backend/internal/events"
	"gorm.io/gorm"
)

type HandoverService interface {
	Create(projectID uint64, senderID uint64, req dto.CreateHandoverRequest) (*entity.Handover, error)
	GetByID(id uint64) (*entity.Handover, error)
	GetByProjectID(projectID uint64, statusFilter string) ([]entity.Handover, error)
	MarkReceived(id uint64, receiverID uint64, req dto.MarkReceivedRequest) (*entity.Handover, error)
	MarkReturned(id uint64, req dto.MarkReturnedRequest) (*entity.Handover, error)
}

type handoverService struct {
	repo     repository.HandoverRepository
	eventBus *events.Bus
}

func NewHandoverService(repo repository.HandoverRepository, eventBus *events.Bus) HandoverService {
	return &handoverService{repo: repo, eventBus: eventBus}
}

func (s *handoverService) Create(projectID uint64, senderID uint64, req dto.CreateHandoverRequest) (*entity.Handover, error) {
	var deliveryTime *string
	if req.DeliveryTime != "" {
		deliveryTime = &req.DeliveryTime
	}

	handover := &entity.Handover{
		ProjectID:        projectID,
		SenderID:         senderID,
		SenderDivisionID: req.SenderDivisionID,
		ReceiverID:       req.ReceiverID,
		Description:      req.Description,
		DeliveryDate:     req.DeliveryDate,
		DeliveryTime:     deliveryTime,
		Status:           entity.StatusPending,
	}

	if err := s.repo.Create(handover); err != nil {
		return nil, err
	}

	s.eventBus.Publish(events.Event{
		Name: "handover.sent",
		Data: map[string]interface{}{
			"handover_id": handover.ID,
			"project_id":  projectID,
			"sender_id":   senderID,
			"receiver_id": req.ReceiverID,
		},
	})

	s.eventBus.Publish(events.Event{
		Name: "handover.created",
		Data: map[string]interface{}{
			"handover_id": handover.ID,
			"project_id":  projectID,
			"sender_id":   senderID,
			"receiver_id": req.ReceiverID,
		},
	})

	return handover, nil
}

func (s *handoverService) GetByID(id uint64) (*entity.Handover, error) {
	handover, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrHandoverNotFound
		}
		return nil, err
	}
	return handover, nil
}

func (s *handoverService) GetByProjectID(projectID uint64, statusFilter string) ([]entity.Handover, error) {
	return s.repo.FindByProjectID(projectID, statusFilter)
}

func (s *handoverService) MarkReceived(id uint64, receiverID uint64, req dto.MarkReceivedRequest) (*entity.Handover, error) {
	handover, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if handover.Status != entity.StatusPending {
		return nil, domainerrors.ErrInvalidStatusChange
	}

	now := time.Now()
	handover.Status = entity.StatusReceived
	handover.ReceivedAt = &now
	handover.ReceiverID = &receiverID

	rows, err := s.repo.UpdateWithVersionCheck(handover, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.eventBus.Publish(events.Event{
		Name: "handover.received",
		Data: map[string]interface{}{
			"handover_id": updated.ID,
			"project_id":  updated.ProjectID,
			"receiver_id": receiverID,
		},
	})

	return updated, nil
}

func (s *handoverService) MarkReturned(id uint64, req dto.MarkReturnedRequest) (*entity.Handover, error) {
	handover, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if handover.Status != entity.StatusPending && handover.Status != entity.StatusReceived {
		return nil, domainerrors.ErrInvalidStatusChange
	}

	handover.Status = entity.StatusReturned
	handover.Description = handover.Description + " | Return reason: " + req.Reason

	rows, err := s.repo.UpdateWithVersionCheck(handover, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

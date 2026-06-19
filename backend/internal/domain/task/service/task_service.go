package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/task/dto"
	"github.com/Kal-el21/backend/internal/domain/task/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/task/errors"
	"github.com/Kal-el21/backend/internal/domain/task/repository"
	"github.com/Kal-el21/backend/internal/domain/task/statemachine"
	"github.com/Kal-el21/backend/internal/events"
	"gorm.io/gorm"
)

type TaskService interface {
	Create(projectID uint64, createdBy uint64, req dto.CreateTaskRequest) (*entity.Task, error)
	GetByID(id uint64) (*entity.Task, error)
	GetByProjectID(projectID uint64) ([]dto.TaskResponse, error)
	Update(id uint64, req dto.UpdateTaskRequest) (*entity.Task, error)
	ChangeStatus(id uint64, userID uint64, isAssignedOnly bool, req dto.ChangeTaskStatusRequest) (*entity.Task, error)
	UpdateProgress(id uint64, userID uint64, isAssignedOnly bool, req dto.UpdateProgressRequest) (*entity.Task, error)
	AssignUsers(taskID uint64, assignedBy uint64, req dto.AssignUserRequest) error
	AddComment(taskID, userID uint64, comment string) error
	GetComments(taskID uint64, page, limit int) ([]entity.TaskComment, int64, error)
	Delete(id uint64, deletedBy uint64) error
	GetAverageProgressByMilestone(milestoneID uint64) (float64, error)
	IsUserAssigned(taskID, userID uint64) (bool, error)
}

type taskService struct {
	repo         repository.TaskRepository
	assigneeRepo repository.AssigneeRepository
	commentRepo  repository.CommentRepository
	eventBus     *events.Bus
}

func NewTaskService(
	repo repository.TaskRepository,
	assigneeRepo repository.AssigneeRepository,
	commentRepo repository.CommentRepository,
	eventBus *events.Bus,
) TaskService {
	return &taskService{repo: repo, assigneeRepo: assigneeRepo, commentRepo: commentRepo, eventBus: eventBus}
}

func (s *taskService) Create(projectID uint64, createdBy uint64, req dto.CreateTaskRequest) (*entity.Task, error) {
	task := &entity.Task{
		ProjectID:   projectID,
		MilestoneID: req.MilestoneID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    entity.TaskPriority(req.Priority),
		Status:      entity.StatusTodo,
		StartDate:   req.StartDate,
		DueDate:     req.DueDate,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	s.eventBus.Publish(events.Event{
		Name: "task.created",
		Data: map[string]interface{}{"task_id": task.ID, "title": task.Title, "project_id": projectID},
	})

	return task, nil
}

func (s *taskService) GetByID(id uint64) (*entity.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetByProjectID(projectID uint64) ([]dto.TaskResponse, error) {
	tasks, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TaskResponse, len(tasks))
	for i, t := range tasks {
		assignees, _ := s.assigneeRepo.FindActiveByTaskID(t.ID)
		assigneeIDs := make([]uint64, len(assignees))
		for j, a := range assignees {
			assigneeIDs[j] = a.UserID
		}

		result[i] = dto.TaskResponse{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			MilestoneID: t.MilestoneID,
			Title:       t.Title,
			Description: t.Description,
			Priority:    string(t.Priority),
			Status:      string(t.Status),
			Progress:    t.Progress,
			StartDate:   t.StartDate,
			DueDate:     t.DueDate,
			Version:     t.Version,
			Assignees:   assigneeIDs,
		}
	}

	return result, nil
}

func (s *taskService) Update(id uint64, req dto.UpdateTaskRequest) (*entity.Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Priority = entity.TaskPriority(req.Priority)
	task.MilestoneID = req.MilestoneID
	task.StartDate = req.StartDate
	task.DueDate = req.DueDate

	rows, err := s.repo.UpdateWithVersionCheck(task, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

// ChangeStatus: PM/ADMIN bisa ubah status apa pun; MEMBER hanya jika assigned (Permission Matrix section 9).
func (s *taskService) ChangeStatus(id uint64, userID uint64, isAssignedOnly bool, req dto.ChangeTaskStatusRequest) (*entity.Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if isAssignedOnly {
		assigned, err := s.assigneeRepo.IsAssigned(id, userID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, domainerrors.ErrNotAssigned
		}
	}

	newStatus := entity.TaskStatus(req.Status)

	if err := statemachine.ValidateTransition(task.Status, newStatus); err != nil {
		return nil, err
	}

	task.Status = newStatus

	rows, err := s.repo.UpdateWithVersionCheck(task, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	if newStatus == entity.StatusDone {
		s.eventBus.Publish(events.Event{
			Name: "task.completed",
			Data: map[string]interface{}{"task_id": task.ID, "title": task.Title},
		})
	}

	return s.repo.FindByID(id)
}

// UpdateProgress: sama seperti ChangeStatus, MEMBER hanya bisa update task yang assigned ke dirinya.
func (s *taskService) UpdateProgress(id uint64, userID uint64, isAssignedOnly bool, req dto.UpdateProgressRequest) (*entity.Task, error) {
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if isAssignedOnly {
		assigned, err := s.assigneeRepo.IsAssigned(id, userID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, domainerrors.ErrNotAssigned
		}
	}

	if req.Progress < 0 || req.Progress > 100 {
		return nil, domainerrors.ErrInvalidProgress
	}

	rows, err := s.repo.UpdateProgressWithVersionCheck(id, req.Progress, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

func (s *taskService) AssignUsers(taskID uint64, assignedBy uint64, req dto.AssignUserRequest) error {
	for _, userID := range req.UserIDs {
		assignee := &entity.TaskAssignee{
			TaskID:     taskID,
			UserID:     userID,
			AssignedBy: assignedBy,
		}
		if err := s.assigneeRepo.Assign(assignee); err != nil {
			return err
		}

		s.eventBus.Publish(events.Event{
			Name: "task.assigned",
			Data: map[string]interface{}{"task_id": taskID, "user_id": userID},
		})
	}
	return nil
}

func (s *taskService) AddComment(taskID, userID uint64, comment string) error {
	c := &entity.TaskComment{
		TaskID:  taskID,
		UserID:  userID,
		Comment: comment,
	}
	return s.commentRepo.Create(c)
}

func (s *taskService) GetComments(taskID uint64, page, limit int) ([]entity.TaskComment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.commentRepo.FindByTaskID(taskID, page, limit)
}

func (s *taskService) Delete(id uint64, deletedBy uint64) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id, deletedBy)
}

// GetAverageProgressByMilestone mengimplementasikan TaskProgressProvider
// yang dibutuhkan milestone.service. Cancelled task tidak dihitung (FR-06.05).
func (s *taskService) GetAverageProgressByMilestone(milestoneID uint64) (float64, error) {
	tasks, err := s.repo.FindByMilestoneID(milestoneID)
	if err != nil {
		return 0, err
	}

	var total float64
	var counted int

	for _, t := range tasks {
		if t.Status == entity.StatusCancelled {
			continue
		}
		total += float64(t.Progress)
		counted++
	}

	if counted == 0 {
		return 0, nil
	}

	return total / float64(counted), nil
}

func (s *taskService) IsUserAssigned(taskID, userID uint64) (bool, error) {
	return s.assigneeRepo.IsAssigned(taskID, userID)
}

package service

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/dashboard/dto"
	"github.com/Kal-el21/backend/internal/domain/dashboard/repository"
)

type DashboardService interface {
	GetSummary(userID uint64, isAdmin bool) (*dto.DashboardSummary, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetSummary(userID uint64, isAdmin bool) (*dto.DashboardSummary, error) {
	totalProjects, err := s.repo.CountTotalProjects(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	active, err := s.repo.CountProjectsByStatus("ACTIVE", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	completed, err := s.repo.CountProjectsByStatus("COMPLETED", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	totalTasks, err := s.repo.CountTotalTasks(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	completedTasks, err := s.repo.CountCompletedTasks(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	pendingRequests, err := s.repo.CountPendingRequests(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	overdueTasks, err := s.repo.CountOverdueTasks(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	onHold, err := s.repo.CountProjectsByStatus("ON_HOLD", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	overdueProjects, err := s.repo.CountOverdueProjects(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	byInitiation, err := s.repo.CountProjectsByDimension("initiation_type", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	byPriority, err := s.repo.CountProjectsByDimension("priority", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	byStatus, err := s.repo.CountProjectsByDimension("status", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	healthGreen, err := s.repo.CountProjectsByHealth("GREEN", userID, isAdmin)
	if err != nil {
		return nil, err
	}
	healthYellow, err := s.repo.CountProjectsByHealth("YELLOW", userID, isAdmin)
	if err != nil {
		return nil, err
	}
	healthRed, err := s.repo.CountProjectsByHealth("RED", userID, isAdmin)
	if err != nil {
		return nil, err
	}

	capexAllocated, capexUsed, opexAllocated, opexUsed, err := s.repo.GetBudgetTotalsByType(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	budgetUsage, err := s.repo.GetAverageBudgetUsage(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	deadlineRows, err := s.repo.GetUpcomingDeadlines(5, userID, isAdmin)
	if err != nil {
		return nil, err
	}
	upcomingDeadlines := make([]dto.UpcomingDeadline, 0, len(deadlineRows))
	for _, d := range deadlineRows {
		days := 0
		if d.EndDate != nil {
			days = int(time.Until(*d.EndDate).Hours() / 24)
		}
		upcomingDeadlines = append(upcomingDeadlines, dto.UpcomingDeadline{
			ID:            d.ID,
			ProjectCode:   d.ProjectCode,
			Name:          d.Name,
			EndDate:       d.EndDate,
			Status:        d.Status,
			DaysRemaining: days,
		})
	}

	budgetMaster, err := s.repo.GetBudgetMaster(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	absorption, err := s.repo.GetAbsorption(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	var recentActivities []dto.RecentActivity
	if isAdmin {
		recentActivities, err = s.repo.GetRecentActivities(10)
		if err != nil {
			return nil, err
		}
	}

	return &dto.DashboardSummary{
		TotalProjects:     totalProjects,
		ActiveProjects:    active,
		OnHoldProjects:    onHold,
		CompletedProjects: completed,
		OverdueProjects:   overdueProjects,
		TotalTasks:        totalTasks,
		CompletedTasks:    completedTasks,
		PendingRequests:   pendingRequests,
		OverdueTasks:      overdueTasks,
		TotalBudgetUsage:  budgetUsage,
		ByStatus:          byStatus,
		ByInitiation:      byInitiation,
		ByPriority:        byPriority,
		HealthGreen:       healthGreen,
		HealthYellow:      healthYellow,
		HealthRed:         healthRed,
		CapexAllocated:    capexAllocated,
		CapexUsed:         capexUsed,
		OpexAllocated:     opexAllocated,
		OpexUsed:          opexUsed,
		UpcomingDeadlines: upcomingDeadlines,
		RecentActivities:  recentActivities,
		BudgetMaster:      budgetMaster,
		Absorption:        absorption,
	}, nil
}

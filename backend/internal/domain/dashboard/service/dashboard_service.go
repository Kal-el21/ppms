package service

import (
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
	total, err := s.repo.CountTotalProjects(userID, isAdmin)
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

	pendingRequests, err := s.repo.CountPendingRequests(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	overdueTasks, err := s.repo.CountOverdueTasks(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	budgetUsage, err := s.repo.GetAverageBudgetUsage(userID, isAdmin)
	if err != nil {
		return nil, err
	}

	var recentActivities []dto.RecentActivity
	if isAdmin {
		// Audit log hanya untuk ADMIN (Permission Matrix section 15: View Audit Logs hanya ADMIN)
		recentActivities, err = s.repo.GetRecentActivities(10)
		if err != nil {
			return nil, err
		}
	}

	return &dto.DashboardSummary{
		TotalProjects:     total,
		ActiveProjects:    active,
		CompletedProjects: completed,
		PendingRequests:   pendingRequests,
		OverdueTasks:      overdueTasks,
		TotalBudgetUsage:  budgetUsage,
		RecentActivities:  recentActivities,
	}, nil
}

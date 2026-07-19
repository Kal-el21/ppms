package service

import (
	"testing"
	"time"

	"github.com/Kal-el21/backend/internal/domain/dashboard/dto"
	"github.com/Kal-el21/backend/internal/domain/dashboard/repository"
)

type fakeDashboardRepo struct {
	totalProjects   int64
	active          int64
	onHold          int64
	completed       int64
	overdueProjects int64
	byStatus        map[string]int64
	byInitiation    map[string]int64
	byPriority      map[string]int64
	health          map[string]int64
	capexAlloc      float64
	capexUsed       float64
	opexAlloc       float64
	opexUsed        float64
	upcoming        []repository.DashboardDeadlineRow
	budgetMaster    []dto.BudgetMasterItem
	absorption      []dto.AbsorptionItem
}

func (f *fakeDashboardRepo) CountProjectsByStatus(status string, userID uint64, isAdmin bool) (int64, error) {
	switch status {
	case "ACTIVE":
		return f.active, nil
	case "ON_HOLD":
		return f.onHold, nil
	case "COMPLETED":
		return f.completed, nil
	}
	return 0, nil
}
func (f *fakeDashboardRepo) CountTotalProjects(userID uint64, isAdmin bool) (int64, error) {
	return f.totalProjects, nil
}
func (f *fakeDashboardRepo) CountPendingRequests(userID uint64, isAdmin bool) (int64, error) {
	return 0, nil
}
func (f *fakeDashboardRepo) CountOverdueTasks(userID uint64, isAdmin bool) (int64, error) {
	return 0, nil
}
func (f *fakeDashboardRepo) CountTotalTasks(userID uint64, isAdmin bool) (int64, error) {
	return 0, nil
}
func (f *fakeDashboardRepo) CountCompletedTasks(userID uint64, isAdmin bool) (int64, error) {
	return 0, nil
}
func (f *fakeDashboardRepo) GetAverageBudgetUsage(userID uint64, isAdmin bool) (float64, error) {
	return 0, nil
}
func (f *fakeDashboardRepo) GetRecentActivities(limit int) ([]dto.RecentActivity, error) {
	return nil, nil
}
func (f *fakeDashboardRepo) CountProjectsByHealth(health string, userID uint64, isAdmin bool) (int64, error) {
	return f.health[health], nil
}
func (f *fakeDashboardRepo) CountProjectsByDimension(dim string, userID uint64, isAdmin bool) (map[string]int64, error) {
	switch dim {
	case "status":
		return f.byStatus, nil
	case "initiation_type":
		return f.byInitiation, nil
	case "priority":
		return f.byPriority, nil
	}
	return nil, nil
}
func (f *fakeDashboardRepo) GetBudgetTotalsByType(userID uint64, isAdmin bool) (float64, float64, float64, float64, error) {
	return f.capexAlloc, f.capexUsed, f.opexAlloc, f.opexUsed, nil
}
func (f *fakeDashboardRepo) CountOverdueProjects(userID uint64, isAdmin bool) (int64, error) {
	return f.overdueProjects, nil
}
func (f *fakeDashboardRepo) GetUpcomingDeadlines(limit int, userID uint64, isAdmin bool) ([]repository.DashboardDeadlineRow, error) {
	return f.upcoming, nil
}

func (f *fakeDashboardRepo) GetBudgetMaster(userID uint64, isAdmin bool) ([]dto.BudgetMasterItem, error) {
	return f.budgetMaster, nil
}

func (f *fakeDashboardRepo) GetAbsorption(userID uint64, isAdmin bool) ([]dto.AbsorptionItem, error) {
	return f.absorption, nil
}

var _ repository.DashboardRepository = (*fakeDashboardRepo)(nil)

func TestDashboardSummaryAggregations(t *testing.T) {
	repo := &fakeDashboardRepo{
		totalProjects:   10,
		active:          4,
		onHold:          2,
		completed:       3,
		overdueProjects: 1,
		byStatus:        map[string]int64{"PLANNED": 1, "ACTIVE": 4, "ON_HOLD": 2, "COMPLETED": 3},
		byPriority:      map[string]int64{"LOW": 1, "MEDIUM": 4, "HIGH": 3, "URGENT": 2},
		health:          map[string]int64{"GREEN": 6, "YELLOW": 3, "RED": 1},
		capexAlloc:      1000,
		capexUsed:       400,
		opexAlloc:       500,
		opexUsed:        100,
		upcoming: []repository.DashboardDeadlineRow{
			{ID: 1, ProjectCode: "PRJ-2026-0001", Name: "A", EndDate: timePtr(time.Now().AddDate(0, 0, 5)), Status: "ACTIVE"},
		},
		budgetMaster: []dto.BudgetMasterItem{
			{BudgetName: "Server", BudgetType: "CAPEX", Allocated: 1000, Used: 400, Remaining: 600, UsagePercentage: 40},
		},
		absorption: []dto.AbsorptionItem{
			{ProjectCode: "PRJ-2026-0001", ProjectName: "A", BudgetType: "CAPEX", Allocated: 1000, Used: 400, UsagePercentage: 40},
		},
	}

	svc := NewDashboardService(repo)
	summary, err := svc.GetSummary(1, true)
	if err != nil {
		t.Fatalf("GetSummary error: %v", err)
	}

	if summary.TotalProjects != 10 || summary.ActiveProjects != 4 || summary.OnHoldProjects != 2 || summary.CompletedProjects != 3 {
		t.Errorf("status counts wrong: %+v", summary)
	}
	if summary.OverdueProjects != 1 {
		t.Errorf("overdue projects = %d, want 1", summary.OverdueProjects)
	}
	if summary.ByPriority["URGENT"] != 2 {
		t.Errorf("urgent priority count = %d, want 2", summary.ByPriority["URGENT"])
	}
	if summary.HealthRed != 1 || summary.HealthGreen != 6 {
		t.Errorf("health distribution wrong: green=%d red=%d", summary.HealthGreen, summary.HealthRed)
	}
	if summary.CapexAllocated != 1000 || summary.OpexAllocated != 500 {
		t.Errorf("budget totals wrong: capex=%v opex=%v", summary.CapexAllocated, summary.OpexAllocated)
	}
	if len(summary.UpcomingDeadlines) != 1 || summary.UpcomingDeadlines[0].ProjectCode != "PRJ-2026-0001" {
		t.Errorf("upcoming deadlines wrong: %+v", summary.UpcomingDeadlines)
	}
	if len(summary.BudgetMaster) != 1 || summary.BudgetMaster[0].BudgetName != "Server" {
		t.Errorf("budget master wrong: %+v", summary.BudgetMaster)
	}
	if len(summary.Absorption) != 1 || summary.Absorption[0].ProjectCode != "PRJ-2026-0001" {
		t.Errorf("absorption wrong: %+v", summary.Absorption)
	}
}

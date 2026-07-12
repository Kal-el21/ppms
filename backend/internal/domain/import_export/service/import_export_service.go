package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	budgetentity "github.com/Kal-el21/backend/internal/domain/budget/entity"
	"github.com/Kal-el21/backend/internal/domain/import_export/dto"
	"github.com/Kal-el21/backend/internal/domain/import_export/entity"
	milestoneentity "github.com/Kal-el21/backend/internal/domain/milestone/entity"
	projectentity "github.com/Kal-el21/backend/internal/domain/project/entity"
	taskentity "github.com/Kal-el21/backend/internal/domain/task/entity"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

// ImportExportService menangani backup (export) dan restore (import) seluruh
// data project dalam format JSON.
type ImportExportService interface {
	ExportAll(userID uint64, isAdmin bool) (*dto.ExportResponse, error)
	ImportData(payload dto.ImportRequest, actorID uint64) (*dto.ImportResult, error)
}

type importExportService struct {
	db *gorm.DB
}

func NewImportExportService(db *gorm.DB) ImportExportService {
	return &importExportService{db: db}
}

// =========================================================================
// EXPORT
// =========================================================================

func (s *importExportService) ExportAll(userID uint64, isAdmin bool) (*dto.ExportResponse, error) {
	query := s.db.Model(&projectentity.Project{}).Where("projects.deleted_at IS NULL")
	if !isAdmin {
		query = query.
			Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	var projects []projectentity.Project
	if err := query.Order("projects.id ASC").Find(&projects).Error; err != nil {
		return nil, err
	}

	exported := make([]entity.ProjectExport, 0, len(projects))
	for _, p := range projects {
		pe, err := s.buildProjectExport(p)
		if err != nil {
			return nil, err
		}
		exported = append(exported, pe)
	}

	return &dto.ExportResponse{
		Version:    entity.SchemaVersion,
		ExportedAt: time.Now(),
		ExportedBy: userID,
		Projects:   exported,
	}, nil
}

func (s *importExportService) buildProjectExport(p projectentity.Project) (entity.ProjectExport, error) {
	pe := entity.ProjectExport{
		ProjectCode:    p.ProjectCode,
		Name:           p.Name,
		Description:    p.Description,
		Category:       p.Category,
		InitiationType: p.InitiationType,
		Priority:       p.Priority,
		Notes:          p.Notes,
		StartDate:      p.StartDate,
		EndDate:        p.EndDate,
		Status:         string(p.Status),
		Health:         p.Health,
	}

	// Members (semua member yang belum di-hard-delete; status ikut diexport).
	var members []projectentity.ProjectMember
	if err := s.db.Where("project_id = ?", p.ID).Find(&members).Error; err != nil {
		return pe, err
	}
	pe.Members = make([]entity.MemberExport, 0, len(members))
	for _, m := range members {
		pe.Members = append(pe.Members, entity.MemberExport{
			UserID:      m.UserID,
			ProjectRole: string(m.ProjectRole),
			Status:      string(m.Status),
		})
	}

	// Milestones.
	var milestones []milestoneentity.Milestone
	if err := s.db.Where("project_id = ? AND deleted_at IS NULL", p.ID).
		Order("order_index ASC").Find(&milestones).Error; err != nil {
		return pe, err
	}
	pe.Milestones = make([]entity.MilestoneExport, 0, len(milestones))
	for _, m := range milestones {
		pe.Milestones = append(pe.Milestones, entity.MilestoneExport{
			RefID:       m.ID,
			Name:        m.Name,
			Description: m.Description,
			OrderIndex:  m.OrderIndex,
			StartDate:   m.StartDate,
			EndDate:     m.EndDate,
			Status:      string(m.Status),
		})
	}

	// Tasks.
	var tasks []taskentity.Task
	if err := s.db.Where("project_id = ? AND deleted_at IS NULL", p.ID).
		Order("order_index ASC").Find(&tasks).Error; err != nil {
		return pe, err
	}
	pe.Tasks = make([]entity.TaskExport, 0, len(tasks))
	for _, t := range tasks {
		pe.Tasks = append(pe.Tasks, entity.TaskExport{
			Title:          t.Title,
			Description:    t.Description,
			Priority:       string(t.Priority),
			Status:         string(t.Status),
			Progress:       t.Progress,
			OrderIndex:     t.OrderIndex,
			MilestoneRefID: t.MilestoneID,
			StartDate:      t.StartDate,
			DueDate:        t.DueDate,
		})
	}

	// Budget + transactions.
	var budget budgetentity.Budget
	err := s.db.Where("project_id = ? AND deleted_at IS NULL", p.ID).First(&budget).Error
	if err == nil {
		be := entity.BudgetExport{
			BudgetType:      budget.BudgetType,
			BudgetName:      budget.BudgetName,
			AllocatedBudget: budget.AllocatedBudget,
		}
		var transactions []budgetentity.BudgetTransaction
		if err := s.db.Where("budget_id = ?", budget.ID).
			Order("transaction_date ASC").Find(&transactions).Error; err != nil {
			return pe, err
		}
		be.Transactions = make([]entity.TransactionExport, 0, len(transactions))
		for _, tr := range transactions {
			var adjType *string
			if tr.AdjustmentType != nil {
				val := string(*tr.AdjustmentType)
				adjType = &val
			}
			be.Transactions = append(be.Transactions, entity.TransactionExport{
				Type:            string(tr.Type),
				AdjustmentType:  adjType,
				Amount:          tr.Amount,
				Reason:          tr.Reason,
				Description:     tr.Description,
				TransactionDate: tr.TransactionDate,
			})
		}
		pe.Budget = &be
	} else if err != gorm.ErrRecordNotFound {
		return pe, err
	}

	return pe, nil
}

// =========================================================================
// IMPORT
// =========================================================================

func (s *importExportService) ImportData(payload dto.ImportRequest, actorID uint64) (*dto.ImportResult, error) {
	if len(payload.Projects) == 0 {
		return nil, apperrors.New(apperrors.ErrValidation, "no projects found in import payload")
	}

	result := &dto.ImportResult{
		TotalProjects:      len(payload.Projects),
		Errors:             []string{},
		ImportedProjectIDs: []uint64{},
	}

	for idx, pe := range payload.Projects {
		newID, err := s.importSingleProject(pe, actorID)
		if err != nil {
			result.Skipped++
			label := strings.TrimSpace(pe.Name)
			if label == "" {
				label = fmt.Sprintf("index %d", idx)
			}
			result.Errors = append(result.Errors, fmt.Sprintf("project %q: %s", label, err.Error()))
			continue
		}
		result.Imported++
		result.ImportedProjectIDs = append(result.ImportedProjectIDs, newID)
	}

	return result, nil
}

// importSingleProject membuat satu project baru (beserta relasinya) di dalam
// transaksi terpisah agar kegagalan satu project tidak membatalkan yang lain.
func (s *importExportService) importSingleProject(pe entity.ProjectExport, actorID uint64) (uint64, error) {
	name := strings.TrimSpace(pe.Name)
	if name == "" {
		return 0, apperrors.New(apperrors.ErrValidation, "name is required")
	}

	var newProjectID uint64
	err := s.db.Transaction(func(tx *gorm.DB) error {
		projectCode, err := generateProjectCode(tx, time.Now().Year())
		if err != nil {
			return err
		}

		project := &projectentity.Project{
			ProjectCode:    projectCode,
			Name:           name,
			Description:    pe.Description,
			Category:       strings.TrimSpace(pe.Category),
			InitiationType: normalizeInitiation(pe.InitiationType),
			Priority:       normalizePriority(pe.Priority),
			Notes:          strings.TrimSpace(pe.Notes),
			StartDate:      pe.StartDate,
			EndDate:        pe.EndDate,
			Status:         normalizeStatus(pe.Status),
			Health:         pe.Health,
			CreatedBy:      actorID,
		}
		if err := tx.Create(project).Error; err != nil {
			return err
		}
		newProjectID = project.ID

		// Milestones — simpan mapping ref_id lama -> id baru untuk relasi task.
		milestoneIDMap := make(map[uint64]uint64)
		for _, m := range pe.Milestones {
			milestone := &milestoneentity.Milestone{
				ProjectID:   project.ID,
				Name:        strings.TrimSpace(m.Name),
				Description: m.Description,
				OrderIndex:  m.OrderIndex,
				StartDate:   m.StartDate,
				EndDate:     m.EndDate,
				Status:      normalizeMilestoneStatus(m.Status),
			}
			if milestone.Name == "" {
				continue
			}
			if err := tx.Create(milestone).Error; err != nil {
				return err
			}
			if m.RefID != 0 {
				milestoneIDMap[m.RefID] = milestone.ID
			}
		}

		// Tasks.
		for _, t := range pe.Tasks {
			title := strings.TrimSpace(t.Title)
			if title == "" {
				continue
			}
			var milestoneID *uint64
			if t.MilestoneRefID != nil {
				if mapped, ok := milestoneIDMap[*t.MilestoneRefID]; ok {
					milestoneID = &mapped
				}
			}
			task := &taskentity.Task{
				ProjectID:   project.ID,
				MilestoneID: milestoneID,
				Title:       title,
				Description: t.Description,
				Priority:    normalizeTaskPriority(t.Priority),
				Status:      normalizeTaskStatus(t.Status),
				Progress:    clampProgress(t.Progress),
				OrderIndex:  t.OrderIndex,
				StartDate:   t.StartDate,
				DueDate:     t.DueDate,
				CreatedBy:   actorID,
			}
			if err := tx.Create(task).Error; err != nil {
				return err
			}
		}

		// Budget + transactions.
		if pe.Budget != nil {
			budget := &budgetentity.Budget{
				ProjectID:       project.ID,
				BudgetType:      normalizeBudgetType(pe.Budget.BudgetType),
				BudgetName:      strings.TrimSpace(pe.Budget.BudgetName),
				AllocatedBudget: pe.Budget.AllocatedBudget,
				CreatedBy:       actorID,
			}
			if err := tx.Create(budget).Error; err != nil {
				return err
			}
			for _, tr := range pe.Budget.Transactions {
				transaction := &budgetentity.BudgetTransaction{
					BudgetID:        budget.ID,
					Type:            normalizeTransactionType(tr.Type),
					AdjustmentType:  normalizeAdjustmentType(tr.AdjustmentType),
					Amount:          tr.Amount,
					Reason:          tr.Reason,
					Description:     tr.Description,
					TransactionDate: transactionDate(tr.TransactionDate),
					CreatedBy:       actorID,
				}
				if err := tx.Create(transaction).Error; err != nil {
					return err
				}
			}
		}

		// Members — hanya user yang benar-benar ada yang dibuat.
		hasActivePM := false
		for _, m := range pe.Members {
			if m.UserID == 0 || !userExists(tx, m.UserID) {
				continue
			}
			role := normalizeMemberRole(m.ProjectRole)
			status := normalizeMemberStatus(m.Status)
			member := &projectentity.ProjectMember{
				ProjectID:   project.ID,
				UserID:      m.UserID,
				ProjectRole: role,
				Status:      status,
				JoinedAt:    time.Now(),
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
			if role == projectentity.RoleProjectManager && status == projectentity.MemberActive {
				hasActivePM = true
			}
		}

		// Pastikan selalu ada minimal satu PROJECT_MANAGER aktif: fallback ke actor.
		if !hasActivePM {
			var existing int64
			if err := tx.Model(&projectentity.ProjectMember{}).
				Where("project_id = ? AND user_id = ?", project.ID, actorID).
				Count(&existing).Error; err != nil {
				return err
			}
			if existing == 0 {
				fallbackPM := &projectentity.ProjectMember{
					ProjectID:   project.ID,
					UserID:      actorID,
					ProjectRole: projectentity.RoleProjectManager,
					Status:      projectentity.MemberActive,
					JoinedAt:    time.Now(),
				}
				if err := tx.Create(fallbackPM).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newProjectID, nil
}

// =========================================================================
// HELPERS
// =========================================================================

func userExists(tx *gorm.DB, userID uint64) bool {
	var count int64
	if err := tx.Table("users").Where("id = ? AND deleted_at IS NULL", userID).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func transactionDate(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now()
	}
	return t
}

func clampProgress(p int) int {
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

func normalizeInitiation(value *string) *string {
	if value == nil {
		return nil
	}
	v := strings.ToUpper(strings.TrimSpace(*value))
	switch v {
	case "NEW_INITIATIVE", "RENEWAL", "ENHANCEMENT":
		return &v
	}
	return nil
}

func normalizePriority(value string) string {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch v {
	case "LOW", "MEDIUM", "HIGH", "URGENT":
		return v
	}
	return "MEDIUM"
}

func normalizeStatus(value string) projectentity.ProjectStatus {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch projectentity.ProjectStatus(v) {
	case projectentity.ProjectPlanned, projectentity.ProjectActive, projectentity.ProjectOnHold,
		projectentity.ProjectCompleted, projectentity.ProjectCancelled:
		return projectentity.ProjectStatus(v)
	}
	return projectentity.ProjectPlanned
}

func normalizeMilestoneStatus(value string) milestoneentity.MilestoneStatus {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch milestoneentity.MilestoneStatus(v) {
	case milestoneentity.MilestonePlanned, milestoneentity.MilestoneActive,
		milestoneentity.MilestoneCompleted, milestoneentity.MilestoneCancelled:
		return milestoneentity.MilestoneStatus(v)
	}
	return milestoneentity.MilestonePlanned
}

func normalizeTaskPriority(value string) taskentity.TaskPriority {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch taskentity.TaskPriority(v) {
	case taskentity.PriorityLow, taskentity.PriorityMedium, taskentity.PriorityHigh, taskentity.PriorityUrgent:
		return taskentity.TaskPriority(v)
	}
	return taskentity.PriorityMedium
}

func normalizeTaskStatus(value string) taskentity.TaskStatus {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch taskentity.TaskStatus(v) {
	case taskentity.StatusTodo, taskentity.StatusInProgress, taskentity.StatusDone, taskentity.StatusCancelled:
		return taskentity.TaskStatus(v)
	}
	return taskentity.StatusTodo
}

func normalizeMemberRole(value string) projectentity.ProjectRole {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch projectentity.ProjectRole(v) {
	case projectentity.RoleProjectManager, projectentity.RoleMember, projectentity.RoleObserver:
		return projectentity.ProjectRole(v)
	}
	return projectentity.RoleMember
}

func normalizeMemberStatus(value string) projectentity.MemberStatus {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch projectentity.MemberStatus(v) {
	case projectentity.MemberActive, projectentity.MemberSuspended,
		projectentity.MemberLeft, projectentity.MemberRemoved:
		return projectentity.MemberStatus(v)
	}
	return projectentity.MemberActive
}

func normalizeBudgetType(value *string) *string {
	if value == nil {
		return nil
	}
	v := strings.ToUpper(strings.TrimSpace(*value))
	if v == "CAPEX" || v == "OPEX" {
		return &v
	}
	return nil
}

func normalizeTransactionType(value string) budgetentity.TransactionType {
	v := strings.ToUpper(strings.TrimSpace(value))
	switch budgetentity.TransactionType(v) {
	case budgetentity.TypeExpense, budgetentity.TypeRefund, budgetentity.TypeAdjustment:
		return budgetentity.TransactionType(v)
	}
	return budgetentity.TypeExpense
}

func normalizeAdjustmentType(value *string) *budgetentity.AdjustmentType {
	if value == nil {
		return nil
	}
	v := strings.ToUpper(strings.TrimSpace(*value))
	switch budgetentity.AdjustmentType(v) {
	case budgetentity.AdjustmentErrorCorrection, budgetentity.AdjustmentBudgetReallocation,
		budgetentity.AdjustmentAuditCorrection, budgetentity.AdjustmentManualOverride:
		adj := budgetentity.AdjustmentType(v)
		return &adj
	}
	return nil
}

// generateProjectCode meniru logika di provisioning_service (PRJ-YYYY-0001)
// dengan retry loop untuk menghindari race pada nomor urut.
func generateProjectCode(tx *gorm.DB, year int) (string, error) {
	prefix := fmt.Sprintf("PRJ-%d-", year)
	var lastCode string

	if err := tx.Model(&projectentity.Project{}).
		Select("project_code").
		Where("project_code LIKE ?", prefix+"%").
		Order("project_code DESC").
		Limit(1).
		Scan(&lastCode).Error; err != nil {
		return "", err
	}

	next := 1
	if lastCode != "" {
		if parsed, err := strconv.Atoi(strings.TrimPrefix(lastCode, prefix)); err == nil {
			next = parsed + 1
		}
	}

	for attempt := 0; attempt < 100; attempt++ {
		code := fmt.Sprintf("%s%04d", prefix, next+attempt)
		var count int64
		if err := tx.Model(&projectentity.Project{}).Where("project_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}

	return "", apperrors.New(apperrors.ErrConflict, "unable to generate unique project code")
}

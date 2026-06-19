package membership

import (
	"fmt"

	attachmentservice "github.com/Kal-el21/backend/internal/domain/attachment/service"
	"gorm.io/gorm"
)

// checker adalah adapter konkret yang query langsung ke tabel terkait
// untuk resolve project_id dari berbagai entity_type, dan mengecek
// project_members. Ditempatkan di infrastructure layer (bukan domain)
// karena query-nya lintas domain (PROJECT_REQUEST, TASK, MILESTONE, dll).
type checker struct {
	db *gorm.DB
}

func NewMembershipChecker(db *gorm.DB) attachmentservice.MembershipChecker {
	return &checker{db: db}
}

func (c *checker) GetProjectIDByEntity(entityType string, entityID uint64) (uint64, error) {
	var projectID uint64
	var err error

	switch entityType {
	case "PROJECT":
		err = c.db.Table("projects").Select("id").Where("id = ?", entityID).Scan(&projectID).Error
	case "MILESTONE":
		err = c.db.Table("milestones").Select("project_id").Where("id = ?", entityID).Scan(&projectID).Error
	case "TASK":
		err = c.db.Table("tasks").Select("project_id").Where("id = ?", entityID).Scan(&projectID).Error
	case "HANDOVER":
		err = c.db.Table("handovers").Select("project_id").Where("id = ?", entityID).Scan(&projectID).Error
	case "BUDGET_TRANSACTION":
		err = c.db.Table("budget_transactions").
			Select("budgets.project_id").
			Joins("JOIN budgets ON budgets.id = budget_transactions.budget_id").
			Where("budget_transactions.id = ?", entityID).
			Scan(&projectID).Error
	case "PROJECT_REQUEST":
		// Project request belum punya project_id (sebelum approved).
		// Untuk kasus ini, ownership divalidasi via requester_id, bukan project membership.
		// Resolver mengembalikan error khusus yang ditangani terpisah di service layer.
		return 0, fmt.Errorf("PROJECT_REQUEST uses requester-based ownership, not project membership")
	default:
		return 0, fmt.Errorf("unsupported entity_type: %s", entityType)
	}

	if err != nil {
		return 0, err
	}
	if projectID == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return projectID, nil
}

func (c *checker) IsActiveMember(projectID uint64, userID uint64) (bool, error) {
	var count int64
	err := c.db.Table("project_members").
		Where("project_id = ? AND user_id = ? AND status = 'ACTIVE'", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

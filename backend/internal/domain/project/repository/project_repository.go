package repository

import (
	"time"

	"github.com/Kal-el21/backend/internal/domain/project/entity"
	"gorm.io/gorm"
)

// ProjectListFilter berisi semua dimensi filter & sorting untuk halaman
// portfolio project (Phase 5). MemberUserID, bila diisi, membatasi hasil ke
// project yang diikuti user tersebut (konsisten dengan permission USER).
type ProjectListFilter struct {
	Search         string
	Status         string
	InitiationType string
	Priority       string
	BudgetType     string
	Sort           string
	Page           int
	Limit          int
	MemberUserID   *uint64
}

// ProjectListRow adalah proyek ditambah ringkasan budget (allocated/used)
// yang di-join untuk ditampilkan di mini-card tiap baris.
type ProjectListRow struct {
	entity.Project
	BudgetType      *string  `gorm:"column:budget_type"`
	BudgetAllocated float64  `gorm:"column:budget_allocated"`
	BudgetUsed      float64  `gorm:"column:budget_used"`
}

// ProjectDeadlineRow adalah proyek dengan end_date di dalam window deadline.
type ProjectDeadlineRow struct {
	ID          uint64     `gorm:"column:id"`
	ProjectCode string     `gorm:"column:project_code"`
	Name        string     `gorm:"column:name"`
	EndDate     *time.Time `gorm:"column:end_date"`
	Status      string     `gorm:"column:status"`
}

type ProjectRepository interface {
	Create(project *entity.Project) error
	FindByID(id uint64) (*entity.Project, error)
	FindAll(page, limit int, statusFilter string) ([]entity.Project, int64, error)
	FindByMemberUserID(userID uint64, page, limit int) ([]entity.Project, int64, error)
	FindWithFilters(filter ProjectListFilter) ([]ProjectListRow, int64, error)
	FindDeadlineProjects(window string, now time.Time) ([]ProjectDeadlineRow, error)
	UpdateHealth(projectID uint64, health string) error
	UpdateWithVersionCheck(project *entity.Project, expectedVersion int) (int64, error)
	Delete(id uint64, deletedBy uint64) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *entity.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uint64) (*entity.Project, error) {
	var project entity.Project
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindAll(page, limit int, statusFilter string) ([]entity.Project, int64, error) {
	var projects []entity.Project
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.Project{}).Where("deleted_at IS NULL")
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&projects).Error
	return projects, total, err
}

// FindByMemberUserID mengambil semua project di mana user adalah member aktif.
// Digunakan untuk USER/VIEWER yang hanya boleh lihat project yang mereka ikuti.
func (r *projectRepository) FindByMemberUserID(userID uint64, page, limit int) ([]entity.Project, int64, error) {
	var projects []entity.Project
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&entity.Project{}).
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ? AND project_members.status = 'ACTIVE' AND projects.deleted_at IS NULL", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("projects.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&projects).Error

	return projects, total, err
}

// FindWithFilters mengambil project portfolio dengan filter, sorting, dan
// ringkasan budget (allocated/used) per baris.
func (r *projectRepository) FindWithFilters(f ProjectListFilter) ([]ProjectListRow, int64, error) {
	var rows []ProjectListRow
	var total int64

	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	offset := (f.Page - 1) * f.Limit

	usedSub := r.db.Table("budget_transactions").
		Select("budget_id, COALESCE(SUM(amount), 0) as used").
		Group("budget_id")

	query := r.db.Model(&entity.Project{}).
		Select(`projects.*, b.budget_type, COALESCE(b.allocated_budget, 0) as budget_allocated, COALESCE(bt.used, 0) as budget_used`).
		Joins("LEFT JOIN budgets b ON b.project_id = projects.id AND b.deleted_at IS NULL").
		Joins("LEFT JOIN (?) bt ON bt.budget_id = b.id", usedSub).
		Where("projects.deleted_at IS NULL")

	if f.MemberUserID != nil {
		query = query.
			Joins("JOIN project_members pm ON pm.project_id = projects.id").
			Where("pm.user_id = ? AND pm.status = 'ACTIVE'", *f.MemberUserID)
	}
	if f.Search != "" {
		query = query.Where("projects.name ILIKE ?", "%"+f.Search+"%")
	}
	if f.Status != "" {
		query = query.Where("projects.status = ?", f.Status)
	}
	if f.InitiationType != "" {
		query = query.Where("projects.initiation_type = ?", f.InitiationType)
	}
	if f.Priority != "" {
		query = query.Where("projects.priority = ?", f.Priority)
	}
	if f.BudgetType != "" {
		query = query.Where("b.budget_type = ?", f.BudgetType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Progress sorting is handled in the service layer because progress is
	// derived from tasks/milestones, not a stored column.
	if f.Sort == "progress" {
		// Skip SQL-level sort; service layer will sort after progress is calculated.
	} else {
		switch f.Sort {
		case "code":
			query = query.Order("projects.project_code ASC NULLS LAST")
		case "name":
			query = query.Order("projects.name ASC NULLS LAST")
		case "end_date":
			query = query.Order("projects.end_date ASC NULLS LAST")
		case "budget":
			query = query.Order("budget_allocated DESC NULLS LAST")
		default:
			query = query.Order("projects.created_at DESC")
		}
	}

	err := query.Limit(f.Limit).Offset(offset).Scan(&rows).Error
	return rows, total, err
}

// FindDeadlineProjects mengambil project yang end_date-nya masuk dalam window
// (overdue | 30 | 60 | 90 hari ke depan). Project yang sudah COMPLETED/CANCELLED
// dikecualikan.
func (r *projectRepository) FindDeadlineProjects(window string, now time.Time) ([]ProjectDeadlineRow, error) {
	var rows []ProjectDeadlineRow

	base := r.db.Model(&entity.Project{}).
		Select("id, project_code, name, end_date, status").
		Where("projects.deleted_at IS NULL").
		Where("projects.status NOT IN ('COMPLETED', 'CANCELLED')").
		Where("projects.end_date IS NOT NULL")

	switch window {
	case "overdue":
		base = base.Where("projects.end_date < ?", startOfDay(now))
	case "30":
		base = base.Where("projects.end_date >= ? AND projects.end_date <= ?", startOfDay(now), startOfDay(now).AddDate(0, 0, 30))
	case "60":
		base = base.Where("projects.end_date >= ? AND projects.end_date <= ?", startOfDay(now), startOfDay(now).AddDate(0, 0, 60))
	case "90":
		base = base.Where("projects.end_date >= ? AND projects.end_date <= ?", startOfDay(now), startOfDay(now).AddDate(0, 0, 90))
	default:
		base = base.Where("projects.end_date >= ?", startOfDay(now))
	}

	err := base.Order("projects.end_date ASC").Scan(&rows).Error
	return rows, err
}

// startOfDay dipakai agar perbandingan tanggal tidak bergantung pada jam.
func startOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

// UpdateHealth menyimpan hasil perhitungan kesehatan project.
func (r *projectRepository) UpdateHealth(projectID uint64, health string) error {
	return r.db.Model(&entity.Project{}).
		Where("id = ?", projectID).
		Update("health", health).Error
}

func (r *projectRepository) UpdateWithVersionCheck(project *entity.Project, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Project{}).
		Where("id = ? AND version = ?", project.ID, expectedVersion).
		Updates(map[string]interface{}{
			"name":           project.Name,
			"description":    project.Description,
			"category":       project.Category,
			"initiation_type": project.InitiationType,
			"priority":       project.Priority,
			"notes":          project.Notes,
			"start_date":     project.StartDate,
			"end_date":       project.EndDate,
			"status":         project.Status,
			"version":        gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}

func (r *projectRepository) Delete(id uint64, deletedBy uint64) error {
	return r.db.Model(&entity.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("now()"),
			"deleted_by": deletedBy,
		}).Error
}

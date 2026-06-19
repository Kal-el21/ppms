package repository

import (
	"github.com/Kal-el21/backend/internal/domain/search/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SearchRepository memanfaatkan pg_trgm extension + GIN index yang sudah
// didefinisikan di migration Phase 0 (idx_*_title_trgm / idx_*_name_trgm)
// untuk fuzzy search yang performant tanpa Elasticsearch.
type SearchRepository interface {
	SearchProjects(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error)
	SearchTasks(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error)
	SearchRequests(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error)
}

type searchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchProjects(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error) {
	var results []dto.SearchResult

	q := r.db.Table("projects").
		Select("'PROJECT' as entity_type, projects.id as entity_id, projects.name as title, projects.description as snippet").
		Where("projects.name % ? AND projects.deleted_at IS NULL", query). // operator % dari pg_trgm: similarity match
		Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:  "similarity(projects.name, ?) DESC",
				Vars: []interface{}{query},
			},
		}).
		Limit(limit)

	if !isAdmin {
		q = q.Joins("JOIN project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	err := q.Scan(&results).Error
	return results, err
}

func (r *searchRepository) SearchTasks(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error) {
	var results []dto.SearchResult

	q := r.db.Table("tasks").
		Select("'TASK' as entity_type, tasks.id as entity_id, tasks.title as title, tasks.description as snippet").
		Where("tasks.title % ? AND tasks.deleted_at IS NULL", query).
		Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:  "similarity(tasks.title, ?) DESC",
				Vars: []interface{}{query},
			},
		}).
		Limit(limit)

	if !isAdmin {
		q = q.Joins("JOIN project_members ON project_members.project_id = tasks.project_id").
			Where("project_members.user_id = ? AND project_members.status = 'ACTIVE'", userID)
	}

	err := q.Scan(&results).Error
	return results, err
}

func (r *searchRepository) SearchRequests(query string, userID uint64, isAdmin bool, limit int) ([]dto.SearchResult, error) {
	var results []dto.SearchResult

	q := r.db.Table("project_requests").
		Select("'PROJECT_REQUEST' as entity_type, project_requests.id as entity_id, project_requests.title as title, project_requests.description as snippet").
		Where("project_requests.title % ? AND project_requests.deleted_at IS NULL", query).
		Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:  "similarity(project_requests.title, ?) DESC",
				Vars: []interface{}{query},
			},
		}).
		Limit(limit)

	if !isAdmin {
		q = q.Where("project_requests.requester_id = ?", userID)
	}

	err := q.Scan(&results).Error
	return results, err
}

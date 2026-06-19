package service

import (
	"github.com/Kal-el21/backend/internal/domain/search/dto"
	"github.com/Kal-el21/backend/internal/domain/search/repository"
)

type SearchService interface {
	Search(query string, userID uint64, isAdmin bool) (*dto.SearchResponse, error)
}

type searchService struct {
	repo repository.SearchRepository
}

func NewSearchService(repo repository.SearchRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) Search(query string, userID uint64, isAdmin bool) (*dto.SearchResponse, error) {
	const perCategoryLimit = 5

	projects, err := s.repo.SearchProjects(query, userID, isAdmin, perCategoryLimit)
	if err != nil {
		return nil, err
	}

	tasks, err := s.repo.SearchTasks(query, userID, isAdmin, perCategoryLimit)
	if err != nil {
		return nil, err
	}

	requests, err := s.repo.SearchRequests(query, userID, isAdmin, perCategoryLimit)
	if err != nil {
		return nil, err
	}

	return &dto.SearchResponse{
		Projects: projects,
		Tasks:    tasks,
		Requests: requests,
	}, nil
}

package article

import "grabber/internal/repository"

type Service interface {
	GetArticles(q string, page, rowsPerPage int64) ([]*repository.Article, error)
}

type service struct {
	repo repository.Repository
}

// NewService creates new reports service instance with specified repository
func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// GetArticles returns articles from the database by specified search phrase and paging
func (s *service) GetArticles(q string, page, rowsPerPage int64) ([]*repository.Article, error) {
	if page <= 0 {
		page = 1
	}

	if rowsPerPage <= 0 || rowsPerPage > 20 {
		rowsPerPage = 10
	}
	return s.repo.GetArticles(q, page, rowsPerPage)
}

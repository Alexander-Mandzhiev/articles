package service

import (
	"articles/internal/models"
	"articles/internal/repository"
)

type ArticleService struct {
	repo repository.Article
}

func NewArticleService(repo repository.Article) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) Create(article *models.Article) (*models.Article, error) {
	return s.repo.Create(article)
}

func (s *ArticleService) GetAll() ([]models.Article, error) {
	return s.repo.GetAll()
}

func (s *ArticleService) GetOne(id string) (models.Article, error) {
	return s.repo.GetOne(id)
}

func (s *ArticleService) Update(id string, article models.Article) (models.Article, error) {
	return s.repo.Update(id, article)
}
func (s *ArticleService) Delete(id string) error {
	return s.repo.Delete(id)
}

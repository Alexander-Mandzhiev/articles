package service

import (
	"articles/internal/models"
	"articles/internal/repository"
)

type Article interface {
	Create(article *models.Article) (*models.Article, error)
	GetAll() ([]models.Article, error)
	GetOne(id string) (models.Article, error)
	Update(id string, article models.Article) (models.Article, error)
	Delete(id string) error
}

type Service struct {
	Article
}

func NewService(repos repository.Repository) *Service {
	return &Service{
		Article: NewArticleService(repos.Article),
	}
}

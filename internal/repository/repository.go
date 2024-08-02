package repository

import (
	"articles/internal/models"

	"github.com/jmoiron/sqlx"
	//"database/sql"
)

type Article interface {
	Create(article *models.Article) (*models.Article, error)
	GetAll() ([]models.Article, error)
	Update(id string, article models.Article) (models.Article, error)
	GetOne(id string) (models.Article, error)
	Delete(id string) error
}
type Repository struct {
	Article
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Article: NewArticlePostgres(db),
	}
}

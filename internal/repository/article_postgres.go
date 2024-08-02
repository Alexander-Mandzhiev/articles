package repository

import (
	"articles/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ArticlePostgres struct {
	db *sqlx.DB
}

func NewArticlePostgres(db *sqlx.DB) *ArticlePostgres {
	return &ArticlePostgres{db: db}
}

func (r *ArticlePostgres) Create(article *models.Article) (*models.Article, error) {
	id := uuid.New().String()
	query := "INSERT INTO articles (id, title, text, authors) VALUES ($1, $2, $3, $4) RETURNING *"
	if err := r.db.QueryRow(query, id, article.Title, article.Text, article.Authors).
		Scan(&article.ID, &article.Title, &article.Text, &article.Authors, &article.CreatedAt); err != nil {
		return nil, err
	}
	return article, nil
}

func (r *ArticlePostgres) GetAll() ([]*models.Article, error) {
	var articles []*models.Article
	query := "SELECT * FROM articles ORDER BY created_at DESC"

	if err := r.db.Select(&articles, query); err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *ArticlePostgres) GetOne(id string) (*models.Article, error) {
	var a models.Article
	query := `SELECT * FROM articles WHERE id = $1`
	if err := r.db.Get(&a, query, id); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticlePostgres) Update(id string, article *models.Article) (*models.Article, error) {
	query := "UPDATE articles SET title=$2, text=$3, authors=$4 WHERE id = $1 RETURNING *"
	if err := r.db.QueryRow(query, id, article.Title, article.Text, article.Authors).
		Scan(&article.ID, &article.Title, &article.Text, &article.Authors, &article.CreatedAt); err != nil {
		return nil, err
	}
	return article, nil
}

func (r *ArticlePostgres) Delete(id string) error {
	query := "DELETE FROM articles WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

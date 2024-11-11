package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"news/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=Post
type Post interface {
	CreatePost(post models.Post)
	Posts(quantity int) ([]models.Post, error)
}

type Repository struct {
	Post
}

// Создание экземпляра репозитория
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Post: NewPostPostgres(db),
	}
}

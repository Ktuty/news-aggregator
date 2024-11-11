package repository

import (
	"context"
	"log"
	"news/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostPostgres struct {
	db *pgxpool.Pool
}

// конструктор для создания экземпляра БД
func NewPostPostgres(db *pgxpool.Pool) *PostPostgres {
	return &PostPostgres{db: db}
}

// Добавление поста в БД
func (p *PostPostgres) CreatePost(post models.Post) {
	_, err := p.db.Exec(context.Background(), `
	INSERT INTO posts (title, content, published_at, link) VALUES ($1, $2, $3, $4) `,
		post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		// ошибки не логирую, чтобы не захламлять терминал с логами такого типа
		// Error inserting post: ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности "posts_link_key" (SQLSTATE 23505)

		// log.Println("Error inserting post:", err)
		// return err
		return
	}
}

// Получение постов из БД
func (p *PostPostgres) Posts(quantity int) ([]models.Post, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT id, title, content, published_at, link FROM posts ORDER BY published_at DESC LIMIT $1`, quantity)
	if err != nil {
		log.Println("Error getting posts:", err)
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			log.Println("Error getting posts:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

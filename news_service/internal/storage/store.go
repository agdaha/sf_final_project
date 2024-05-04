package storage

import "github.com/agdaha/sf_final_project/news_service/internal/models"

type Store interface {
	UpdatePosts([]models.Post) error
	Posts(string, int) (models.Response, error)
	Post(int) (models.Post, error)
	Close() error
}

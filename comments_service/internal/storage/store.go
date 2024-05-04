package storage

import "github.com/agdaha/sf_final_project/comments_service/internal/models"

type Store interface {
	Get(newsId int) ([]models.Comment, error)
	Post(models.NewComment) (uint64, error)
	Close() error
}

package commentservice

import (
	"context"

	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
)

type CommentService interface {
	GetCommentsForNews(ctx context.Context, newsId int) ([]byte, error)
	CreateComment(ctx context.Context, comment models.NewComment) (int, error)
}

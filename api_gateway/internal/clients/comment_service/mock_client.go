package commentservice

import (
	"context"
	"errors"

	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
)

var _ CommentService = &Mock_client{}

type Mock_client struct {
}

func (c *Mock_client) GetCommentsForNews(ctx context.Context, newsId int) ([]byte, error) {
	if newsId == 13 {
		return []byte(""), errors.New("error")
	}
	return []byte(`[
        {
            "id": 3,
            "author": "Avast Avastovich",
            "comment_text": " New 2500 comment",
            "news_id": 26,
            "parent_id": null,
            "level": 1
        },
        {
            "id": 4,
            "author": "Avast Avastovich",
            "comment_text": " New 2500 comment",
            "news_id": 26,
            "parent_id": null,
            "level": 1
        },
        {
            "id": 5,
            "author": "Avast Avastovich",
            "comment_text": " New 2500 comment",
            "news_id": 26,
            "parent_id": 4,
            "level": 2
        },
        {
            "id": 6,
            "author": "Avast Avastovich",
            "comment_text": " New 2500 comment",
            "news_id": 26,
            "parent_id": 4,
            "level": 2
        }
    ]`), nil
}

func (c *Mock_client) CreateComment(ctx context.Context, comment models.NewComment) (int, error) {
	if comment.Author == "onegin" {
		return 0, errors.New("error")
	}
	return 1, nil
}

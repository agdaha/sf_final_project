package newsservice

import "context"

type NewsService interface {
	GetNews(ctx context.Context, page int, search string) ([]byte, error)
	GetNewsDetailed(ctx context.Context, id int) ([]byte, error)
}

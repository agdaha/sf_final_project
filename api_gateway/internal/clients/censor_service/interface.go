package censorservice

import (
	"context"
)

type CensorService interface {
	CheckComment(ctx context.Context, commentText string) (int, error)
}

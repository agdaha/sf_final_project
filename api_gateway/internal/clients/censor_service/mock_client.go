package censorservice

import (
	"context"
	"strings"
)

var _ CensorService = &Mock_client{}

type Mock_client struct {
}

func (c *Mock_client) CheckComment(ctx context.Context, text string) (int, error) {
	if strings.Contains(text, "qwerty") {
		return 400, nil
	}
	return 200, nil
}

package censorservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
)

var _ CensorService = &client{}

type client struct {
	URL        string
	HTTPClient *http.Client
	Logger     *slog.Logger
}

func New(URL string, logger *slog.Logger) CensorService {
	return &client{
		URL: URL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger: logger.With(slog.String("op", "clents.censor_service")),
	}
}

func (c *client) CheckComment(ctx context.Context, text string) (int, error) {
	var code int

	c.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(map[string]string{"comment_text": text})
	if err != nil {
		return code, fmt.Errorf("failed to marshal text")
	}

	c.Logger.Debug("create new request")
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(dataBytes))
	if err != nil {
		return code, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	c.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Request-Id", middleware.GetReqID(ctx))

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return code, fmt.Errorf("failed to send request. error: %w", err)
	}
	defer response.Body.Close()
	return response.StatusCode, nil
}

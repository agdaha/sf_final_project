package commentservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
)

var _ CommentService = &client{}

type client struct {
	URL        string
	Resource   string
	HTTPClient *http.Client
	Logger     *slog.Logger
}

func New(URL string, logger *slog.Logger) CommentService {
	return &client{
		URL: URL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger: logger,
	}
}

func (c *client) GetCommentsForNews(ctx context.Context, newsId int) ([]byte, error) {
	var comments []byte

	uri := fmt.Sprintf("%s/news/%d", c.URL, newsId)
	c.Logger.Debug("create new request")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %v", err)
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
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusBadRequest {
		c.Logger.Debug("read body")
		comments, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body")
		}
		return comments, nil
	}
	return comments, fmt.Errorf(" ошибка со статусом %v", response.StatusCode)
}

func (c *client) CreateComment(ctx context.Context, comment models.NewComment) (int, error) {
	var code int

	c.Logger.Debug("marshal comment to bytes")
	commentBytes, err := json.Marshal(comment)
	if err != nil {
		return code, fmt.Errorf("failed to marshal comment")
	}

	c.Logger.Debug("create new request")
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(commentBytes))
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
	var id int
	err = json.NewDecoder(response.Body).Decode(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

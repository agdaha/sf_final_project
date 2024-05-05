package newsservice

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
)

var _ NewsService = &client{}

type client struct {
	URL        string
	HTTPClient *http.Client
	Logger     *slog.Logger
}

func New(URL string, logger *slog.Logger) NewsService {
	return &client{
		URL: URL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger: logger.With(slog.String("op", "clents.news_service")),
	}
}

func (c *client) GetNews(ctx context.Context, page int, search string) ([]byte, error) {
	var news []byte

	c.Logger.Debug("create new request")
	uri := fmt.Sprintf("%s?s=%s&page=%d", c.URL, search, page)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Request-Id", middleware.GetReqID(ctx))

	c.Logger.Debug("send request")
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK && response.StatusCode >= http.StatusBadRequest {
		return news, fmt.Errorf(" ошибка со статусом %v", response.StatusCode)
	}

	c.Logger.Debug("read body")
	news, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return news, nil
}

func (c *client) GetNewsDetailed(ctx context.Context, id int) ([]byte, error) {
	var news []byte

	c.Logger.Debug("create new request")
	uri := fmt.Sprintf("%s/%d", c.URL, id)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("X-Request-Id", middleware.GetReqID(ctx))

	c.Logger.Debug("send request")
	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK && response.StatusCode >= http.StatusBadRequest {
		return news, fmt.Errorf(" ошибка со статусом %v", response.StatusCode)
	}

	c.Logger.Debug("read body")
	news, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return news, nil
}

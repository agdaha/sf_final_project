package getnewslist

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	newsservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/news_service"
	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
)

func TestGetNewsHandler(t *testing.T) {

	handlerToTest := New(&newsservice.Mock_client{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	req := httptest.NewRequest("GET", "http://testing", nil)

	w := httptest.NewRecorder()

	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(" Not StatusOK")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")
	if !isFound {
		t.Errorf(" Not application/json type content")
	}

	var n models.Response
	err := json.NewDecoder(resp.Body).Decode(&n)
	if err != nil {
		t.Errorf(" Not decoded")
	}
}

func find(ct []string, s string) bool {
	isFound := false
	for _, v := range ct {
		if v == s {
			isFound = true
			break
		}
	}
	return isFound
}

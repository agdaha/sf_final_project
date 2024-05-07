package getnewsdetailed

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	commentservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/comment_service"
	newsservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/news_service"
	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
	"github.com/julienschmidt/httprouter"
)

func TestGetNewsDetailedHandler(t *testing.T) {
	const id = 1

	handlerToTest := New(&newsservice.Mock_client{}, &commentservice.Mock_client{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
	)

	req := httptest.NewRequest("GET", "http://testing/api/"+strconv.Itoa(id), nil)

	w := httptest.NewRecorder()
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/:newsid", handlerToTest)
	router.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf(" Not StatusOk")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
	}

	var nd models.NewsDetailed
	err := json.NewDecoder(resp.Body).Decode(&nd)
	if err != nil {
		t.Errorf(" Not decoded")
	}
}

func TestGetNewsDetailedWithErrorHandler(t *testing.T) {
	const id = 13

	handlerToTest := New(&newsservice.Mock_client{}, &commentservice.Mock_client{},
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})),
	)

	req := httptest.NewRequest("GET", "http://testing/api/"+strconv.Itoa(id), nil)

	w := httptest.NewRecorder()
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/:newsid", handlerToTest)
	router.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf(" Not Status error")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
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

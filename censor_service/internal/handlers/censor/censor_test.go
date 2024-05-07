package censor

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agdaha/sf_final_project/censor_service/internal/storage/memdb"
)

func TestCensorValidHandle(t *testing.T) {
	db, _ := memdb.New()
	handlerToTest := New(db, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	dataBytes, _ := json.Marshal(map[string]string{"comment_text": "Avast Avastovich"})

	req := httptest.NewRequest("POST", "http://testing", bytes.NewBuffer(dataBytes))

	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("wrong statusCode not 200")
	}
}

func TestCensorNoValidHandle(t *testing.T) {
	db, _ := memdb.New()
	handlerToTest := New(db, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	dataBytes, _ := json.Marshal(map[string]string{"comment_text": "qwerty"})

	req := httptest.NewRequest("POST", "http://testing", bytes.NewBuffer(dataBytes))

	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("wrong statusCode not 400")
	}
}

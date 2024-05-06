package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	buf := new(bytes.Buffer)
	thandler := slog.NewTextHandler(buf, nil)
	log := slog.New(thandler)
	logMiddleware := Logger(log)

	handlerToTest := logMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	logString := buf.String()
	if !strings.Contains(logString, "op=middleware/logger") {
		t.Error("No op string")
	}
	if !strings.Contains(logString, "logger middleware enabled") {
		t.Error("No enabled string")
	}
	if !strings.Contains(logString, "method=") {
		t.Error("No method string")
	}
	if !strings.Contains(logString, "path=") {
		t.Error("No path string")
	}
	if !strings.Contains(logString, "remote_addr=") {
		t.Error("No remote_addr string")
	}
	if !strings.Contains(logString, "user_agent=") {
		t.Error("No remouser_agentte_addr string")
	}
	if !strings.Contains(logString, "request_id=") {
		t.Error("No request_id string")
	}
	if !strings.Contains(logString, "level=INFO") {
		t.Error("No level info")
	}
	if !strings.Contains(logString, "request completed") {
		t.Error("No request completed")
	}
	if !strings.Contains(logString, "status=") {
		t.Error("No status string")
	}
	if !strings.Contains(logString, "bytes=") {
		t.Error("No bytes string")
	}
	if !strings.Contains(logString, "duration=") {
		t.Error("No duration string")
	}

}

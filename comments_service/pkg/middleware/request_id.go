package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0
const RequestIDHeader string = "X-Request-Id"

// middleware получение request id
// - из строки запроса, иначе:
// - из заголовка запроса, иначе:
// - генерация id
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("request_id")
		if id == "" {
			id = r.Header.Get(RequestIDHeader)
			if id == "" {
				id = uuid.New().String()
			}
		}
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, id)))
	})
}

// функция получения request id ранее сохраненного в контексте запроса
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

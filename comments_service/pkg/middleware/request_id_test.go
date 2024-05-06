package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(RequestIDKey)
		if val == nil {
			t.Error("request_id not present")
		}
		valStr, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		if valStr != "1234" {
			t.Error("wrong request_id")
		}
	})

	handlerToTest := RequestID(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing?request_id=1234", nil)
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestGeneratingRequestID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(RequestIDKey)
		if val == nil {
			t.Error("request_id not present")
		}
		_, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		// if valStr != "1234" {
		// 	t.Error("wrong request_id")
		// }
	})

	handlerToTest := RequestID(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

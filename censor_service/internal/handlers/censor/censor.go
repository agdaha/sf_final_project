package censor

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/agdaha/sf_final_project/censor_service/internal/models"
	"github.com/agdaha/sf_final_project/censor_service/internal/storage"
	"github.com/agdaha/sf_final_project/censor_service/pkg/middleware"
)

func New(store storage.Store, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.censor"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Request

		log.Debug("decode comment")
		err := json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) || req.CommentText == "" {
			log.Error("request body is empty")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request body is empty"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to decode request body"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		valid := true
		for _, w := range store.Words() {
			if strings.Contains(strings.ToLower(req.CommentText), w) {
				valid = false
				break
			}
		}

		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

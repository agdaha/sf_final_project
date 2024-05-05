package getcomments

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/pkg/middleware"
	"github.com/julienschmidt/httprouter"
)

func New(db storage.Store, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.getcomments"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		newsIdStr := params.ByName("newsId")
		id, err := strconv.Atoi(newsIdStr)
		if err != nil {
			log.Error(" failed parse newsid", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debug("Получение id новости", slog.Int("id", id))

		log.Debug("Чтение комментариев из БД")
		comments, err := db.Get(id)
		if err != nil {
			log.Error(" error get comments from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug("маршализация комментариев")
		commentsBytes, err := json.Marshal(comments)
		if err != nil {
			log.Error(" error marshalling ", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(commentsBytes)
	}
}

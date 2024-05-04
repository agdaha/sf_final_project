package getposts

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/agdaha/sf_final_project/news_service/internal/storage"
	"github.com/agdaha/sf_final_project/news_service/pkg/middleware"
)

func New(db storage.Store, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.getposts"
		var (
			page   int
			search string
			err    error
		)
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")
		pageString := r.URL.Query().Get("page")
		if pageString == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(pageString)
			if err != nil {
				page = 1
			}
		}
		search = r.URL.Query().Get("s")

		log.Debug("Запрос в БД с параметрами:", slog.Int("page", page), slog.String("search", search))
		res, err := db.Posts(search, page)
		if err != nil {
			log.Error("Error get posts from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug("Маршализация данных")
		postsBytes, err := json.Marshal(res)
		if err != nil {
			log.Error("Error marshalling ", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug("Запись данных в ответ")
		w.WriteHeader(http.StatusOK)
		w.Write(postsBytes)
	}
}

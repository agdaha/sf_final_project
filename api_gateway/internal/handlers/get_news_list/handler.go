package getnewslist

import (
	"log/slog"
	"net/http"
	"strconv"

	newsservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/news_service"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
)

func New(newsService newsservice.NewsService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.get_news_list"
		var page int
		var search string
		var err error

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pageString := r.URL.Query().Get("page")
		if pageString == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(pageString)
			if err != nil {
				page = 1
			}
		}
		if page <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("номера страниц не могут быть отрицательными"))
		}
		search = r.URL.Query().Get("s")

		news, err := newsService.GetNews(r.Context(), page, search)
		if err != nil {
			log.Error("Ups", slog.Any("error", err))
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(news)
	}
}

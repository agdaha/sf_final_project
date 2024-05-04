package getonepost

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/agdaha/sf_final_project/news_service/internal/storage"
	"github.com/agdaha/sf_final_project/news_service/pkg/middleware"
	"github.com/julienschmidt/httprouter"
)

func New(db storage.Store, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.getonepost"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		postIds := params.ByName("newsid")
		id, err := strconv.Atoi(postIds)
		if err != nil {
			log.Error(" failed parse newsid", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Error":"failed parse newsid"}`))
			return
		}

		log.Debug("Запрос в БД с параметром:", slog.Int("news id", id))
		res, err := db.Post(id)
		if err != nil {
			log.Error(" error get posts from db", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debug("Маршализация данных")
		postBytes, err := json.Marshal(res)
		if err != nil {
			log.Error(" error marshalling ", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debug("Запись данных в ответ")
		w.WriteHeader(http.StatusOK)
		w.Write(postBytes)
	}
}

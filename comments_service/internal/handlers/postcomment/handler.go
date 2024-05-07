package postcomment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/agdaha/sf_final_project/comments_service/internal/models"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/pkg/middleware"
	"github.com/go-playground/validator/v10"
)

func New(db storage.Store, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.postcomment"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		w.Header().Set("Content-Type", "application/json")

		log.Debug("декодирование комментария")
		var comment models.NewComment
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			log.Error(" ошибка декодирования комментария", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(" ошибка декодирования комментария"))
			return
		}

		log.Debug("Валидация полученной структуры")
		err = validator.New().Struct(comment)
		if err != nil {
			log.Error(" No valid data ", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Debug("Сохранение комментария из БД")
		id, err := db.Post(comment)
		if err != nil {
			log.Error(" error insert comment to db ", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debug(" id нового комментария", slog.Uint64("id", id))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(int(id))))
	}
}

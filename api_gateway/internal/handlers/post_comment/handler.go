package postcomment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	censorservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/censor_service"
	commentservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/comment_service"
	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
)

func New(commentService commentservice.CommentService, censorService censorservice.CensorService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.post_comment"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		w.Header().Set("Content-Type", "application/json")

		log.Debug("decode create comment")
		var comment models.NewComment
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			log.Error(" ошибка декодирования комментария", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(" ошибка декодирования комментария"))
			return
		}

		log.Debug(" Проверка коммента", slog.Any("comment", comment))
		code, err := censorService.CheckComment(r.Context(), comment.CommentText)
		if err != nil {
			w.WriteHeader(code)
			switch code {
			case http.StatusBadRequest:
				w.Write([]byte("проверьте направляемые данные"))
			case http.StatusInternalServerError:

			case http.StatusServiceUnavailable:
				w.Write([]byte("недоступен сервис цензурирования"))
			}
			return
		}
		if code != http.StatusOK {
			w.WriteHeader(code)
			w.Write([]byte("пост не прошел цензор"))
			return
		}

		id, err := commentService.CreateComment(r.Context(), comment)
		if err != nil {
			log.Debug("ошибка при сохранении", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Debug(" комментарий сохранен")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(int(id))))
	}
}

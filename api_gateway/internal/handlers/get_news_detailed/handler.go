package getnewsdetailed

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	commentservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/comment_service"
	newsservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/news_service"
	"github.com/agdaha/sf_final_project/api_gateway/internal/models"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func New(newsService newsservice.NewsService, commentService commentservice.CommentService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.get_news_detailed"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		w.Header().Set("Content-Type", "application/json")

		id, err := exctractNewsId(r)
		if err != nil {
			log.Error(" failed parse newsid", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Error":"failed parse newsid"}`))
			return
		}

		var wg sync.WaitGroup
		ch := make(chan interface{}, 2)

		wg.Add(1)
		go func(ch chan interface{}, c context.Context) {
			defer wg.Done()
			newsB, err := newsService.GetNewsDetailed(c, id)
			if err != nil {
				ch <- models.RoutineNews{
					News: models.News{},
					Err:  errors.Wrap(err, "get news"),
				}
				return
			}
			var newsOne models.News
			err = json.NewDecoder(bytes.NewBuffer(newsB)).Decode(&newsOne)
			if err != nil {
				ch <- models.RoutineNews{
					News: models.News{},
					Err:  errors.Wrap(err, "decoder news"),
				}
				return
			}
			ch <- models.RoutineNews{
				News: newsOne,
				Err:  nil,
			}
		}(ch, r.Context())

		wg.Add(1)
		go func(ch chan interface{}, c context.Context) {
			defer wg.Done()
			commentsB, err := commentService.GetCommentsForNews(r.Context(), id)
			if err != nil {
				ch <- models.RoutineComments{
					Comments: []models.Comment{},
					Err:      errors.Wrap(err, "get comments for news"),
				}
				return
			}
			var comments []models.Comment
			err = json.NewDecoder(bytes.NewBuffer(commentsB)).Decode(&comments)
			if err != nil {
				ch <- models.RoutineComments{
					Comments: []models.Comment{},
					Err:      errors.Wrap(err, "decoder comments"),
				}
				return
			}
			ch <- models.RoutineComments{
				Comments: comments,
				Err:      nil,
			}
		}(ch, r.Context())

		newsD := models.NewsDetailed{}

		wg.Wait()
		close(ch)

		for result := range ch {
			switch result := result.(type) {
			case models.RoutineNews:
				if result.Err != nil {
					log.Error("Что-то пошло не так news", slog.Any("error", result.Err))
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Что-то пошло не так news"))
					return
				}
				newsD.News = result.News
			case models.RoutineComments:
				if result.Err != nil {
					log.Error("Что-то пошло не так com", slog.Any("error", result.Err))
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Что-то пошло не так com"))
					return
				}
				newsD.Comments = result.Comments
			}
		}

		newsDetailedBytes, err := json.Marshal(newsD)
		if err != nil {
			log.Error("Error marshalling ", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(newsDetailedBytes)
	}
}

func exctractNewsId(r *http.Request) (int, error) {
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	newsIds := params.ByName("newsid")
	id, err := strconv.Atoi(newsIds)
	return id, err
}

package reader

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/agdaha/sf_final_project/news_service/internal/config"
	"github.com/agdaha/sf_final_project/news_service/internal/models"
	"github.com/agdaha/sf_final_project/news_service/internal/storage"
)

type reader struct {
	config *config.Config
	store  storage.Store
	log    *slog.Logger
}

func New(config *config.Config, store storage.Store, log *slog.Logger) *reader {
	return &reader{
		config: config,
		store:  store,
		log:    log,
	}
}

func (r *reader) Start() (chan []models.Post, chan error) {
	// запуск парсинга новостей в отдельном потоке для каждой ссылки
	chPosts := make(chan []models.Post)
	chErrs := make(chan error)

	for _, url := range r.config.Rss {
		go r.getPosts(url, chPosts, chErrs)
	}

	// запись потока новостей в БД
	go func() {
		for posts := range chPosts {
			err := r.store.UpdatePosts(posts)
			if err != nil {
				chErrs <- err
			}
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range chErrs {
			r.log.Error(" ошибка получения RSS", slog.Any("error", err))
		}
	}()
	return chPosts, chErrs
}

// Получение и обработка rss подписок
func (r *reader) getPosts(url string, posts chan<- []models.Post, errs chan<- error) {
	for {
		news, err := ParseRss(url)
		if err != nil {
			errs <- err
			continue
		}
		r.log.Info(fmt.Sprintf(" чтение новостей с %v, получено %v", url, len(news)))
		posts <- news
		time.Sleep(time.Minute * time.Duration(r.config.RequestPeriod))
	}
}

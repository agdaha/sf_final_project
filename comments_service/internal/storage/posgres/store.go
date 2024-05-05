package posgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/agdaha/sf_final_project/comments_service/internal/models"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/pkg/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ storage.Store = &store{}

type store struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

// Функция-конструктор хранилища на postgresql
func New(constr string, log *slog.Logger) (*store, error) {
	if constr == "" {
		return nil, errors.New("не указано подключение к БД")
	}

	pool, err := pgxpool.New(context.Background(), constr)
	if err != nil {
		return nil, err
	}

	log.Debug("проверка доступности БД")
	r := helper.Retry(pool.Ping, 5, 10*time.Second)
	err = r(context.Background())
	if err != nil {
		return nil, err
	}

	return &store{pool: pool, log: log.With(slog.String("op", "store.postgres"))}, nil
}

// Получение комментариев к новости
func (s *store) Get(newsId int) ([]models.Comment, error) {

	s.log.Debug("получение комментариев по id новости", slog.Int("newsId", newsId))
	rows, err := s.pool.Query(context.Background(), `
	SELECT id, author, text, news_id, parent_id, cardinality(path) 
	FROM comments_structure 
	WHERE news_id=$1 
	ORDER BY path
	`,
		newsId,
	)
	if err != nil {
		return nil, err
	}

	s.log.Debug("получение структуры комментариев")
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(
			&comment.Id,
			&comment.Author,
			&comment.CommentText,
			&comment.NewsId,
			&comment.ParentId,
			&comment.Level,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Сохранение нового коментария
func (s *store) Post(c models.NewComment) (uint64, error) {

	s.log.Debug("сохранение комментария", slog.Any("Comment", c))
	row := s.pool.QueryRow(context.Background(), `
	INSERT INTO comments (author, text, news_id, parent_id)
	VALUES ($1,$2,$3,$4)
	RETURNING id;`,
		c.Author,
		c.CommentText,
		c.NewsId,
		c.ParentId,
	)

	var id uint64
	err := row.Scan(&id)
	if err != nil {
		s.log.Error("Unable to INSERT:", slog.Any("error", err))
		return 0, err
	}
	s.log.Debug("получение id нового комментария", slog.Uint64("id", id))
	return id, nil
}

func (s *store) Close() error {
	s.pool.Close()
	return nil
}

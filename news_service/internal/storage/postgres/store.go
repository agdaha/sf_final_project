package postgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
	"github.com/agdaha/sf_final_project/news_service/pkg/helper"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postsPerPage        = 15
	lenShortDescription = 400
)

// Реализация интерфейса хранилища на postgresql
type Store struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

// Функция-конструктор хранилища на postgresql
func New(constr string, log *slog.Logger) (*Store, error) {
	if constr == "" {
		return nil, errors.New("не указано подключение к БД")
	}

	pool, err := pgxpool.New(context.Background(), constr)
	if err != nil {
		return nil, err
	}

	//Ожидание доступности БД на старте. 5 попыток с интервалом в 3 сек.
	r := helper.Retry(pool.Ping, 5, 10*time.Second)
	err = r(context.Background())
	if err != nil {
		return nil, err
	}

	return &Store{pool: pool, log: log}, nil
}

// Запись/Обновление новостей в хранилище
func (s *Store) UpdatePosts(posts []models.Post) error {
	for _, post := range posts {
		_, err := s.pool.Exec(context.Background(), `
		INSERT INTO news (title, description, link, pub_date, author, guid)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (link)
		DO UPDATE SET title=$1, description=$2, link=$3, pub_date=$4, author=$5, guid=$6;`,
			post.Title,
			post.Description,
			post.Link,
			post.PubDate,
			post.Author,
			post.Guid,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Получение свежих новостей
// search - строка поиска в названии новости
// page - номер страницы
func (s *Store) Posts(search string, page int) (models.Response, error) {
	s.log.Debug("параметры запроса", slog.Any("search", search), slog.Any("page", page))
	row := s.pool.QueryRow(context.Background(), `
		SELECT count(id) as cnt
		FROM news
		WHERE title ILIKE '%`+search+`%'`)

	var postsCount int
	err := row.Scan(&postsCount)
	if err != nil {
		return models.Response{}, err
	}

	pagesCount := postsCount / postsPerPage
	if postsCount%postsPerPage > 0 {
		pagesCount++
	}

	if page > pagesCount && pagesCount != 0 {
		page = pagesCount
	}

	rows, err := s.pool.Query(context.Background(), `
	SELECT id, title, description, pub_date 
	FROM news
	WHERE title ILIKE '%`+search+`%'
	ORDER BY pub_date DESC
	LIMIT $1 OFFSET $2
	`,
		postsPerPage,
		(page-1)*postsPerPage,
	)

	if err != nil {
		return models.Response{}, err
	}

	var posts []models.ShortPost
	for rows.Next() {
		var p models.ShortPost
		err = rows.Scan(
			&p.Id,
			&p.Title,
			&p.Description,
			&p.PubDate,
		)
		if err != nil {
			return models.Response{}, err
		}
		p.Description = strip.StripTags(p.Description)
		if len(p.Description) > lenShortDescription {
			p.Description = p.Description[:lenShortDescription+3] + " ..."
		}
		posts = append(posts, p)
	}

	res := models.Response{
		Posts: posts,
		Pages: models.Pages{
			Total:   pagesCount,
			Current: page,
		},
	}
	return res, rows.Err()
}

func (s *Store) Post(id int) (models.Post, error) {
	row := s.pool.QueryRow(context.Background(), `
		SELECT id, title, description, link, pub_date, author
		FROM news
		WHERE id=$1
	`, id)
	p := models.Post{}
	err := row.Scan(
		&p.Id,
		&p.Title,
		&p.Description,
		&p.Link,
		&p.PubDate,
		&p.Author,
	)
	if err != nil {
		return models.Post{}, err
	}
	return p, nil
}

func (s *Store) Close() error {
	s.pool.Close()
	return nil
}

package memdb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
)

type Store struct {
	Db []models.Post
}

// Создание нового хранилища
func New(count int) (*Store, error) {
	var posts []models.Post
	for i := 0; i < count; i++ {
		p := models.Post{
			Id:          i,
			Title:       fmt.Sprintf("%v %v", "Title", i),
			Description: fmt.Sprintf("%v %v", "Description", i),
			Link:        fmt.Sprintf("%v %v", "Link", i),
			Author:      fmt.Sprintf("%v %v", "Author", i),
			Guid:        fmt.Sprintf("%v %v", "Guid", i),
		}
		posts = append(posts, p)
	}
	return &Store{Db: posts}, nil
}

// Запись/Обновление новостей в хранилище
func (s *Store) UpdatePosts(posts []models.Post) error {
	return nil
}

// Получение
func (s *Store) Posts(search string, page int) (models.Response, error) {

	filteredPosts := []models.ShortPost{}

	for _, v := range s.Db {
		if strings.Contains(v.Title, search) {
			shPost := models.ShortPost{
				Id:          v.Id,
				Title:       v.Title,
				Description: v.Description,
				PubDate:     v.PubDate,
			}
			filteredPosts = append(filteredPosts, shPost)
		}
	}

	postsCount := len(filteredPosts)

	postsPerPage := 5

	pagesCount := postsCount / postsPerPage
	if postsCount%postsPerPage > 0 {
		pagesCount++
	}
	if page > pagesCount && pagesCount != 0 {
		page = pagesCount
	}

	posts := []models.ShortPost{}

	start := (page - 1) * postsPerPage
	finish := page * postsPerPage
	if finish > postsCount {
		finish = postsCount
	}

	for _, v := range filteredPosts[start:finish] {
		if strings.Contains(v.Title, search) {
			posts = append(posts, v)
		}
	}

	resp := models.Response{
		News: posts,
		Pages: models.Pages{
			Total:   pagesCount,
			Current: page,
		},
	}

	return resp, nil
}

func (s *Store) Post(id int) (models.Post, error) {
	if id >= len(s.Db) {
		return models.Post{}, errors.New("index error")
	}
	return s.Db[id], nil
}

func (s *Store) Close() error {
	return nil
}

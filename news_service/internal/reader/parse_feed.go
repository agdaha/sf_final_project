package reader

import (
	"strings"
	"time"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
)

// Получение среза постов/нововстей
func ParseRss(url string) ([]models.Post, error) {
	f, err := GetRss(url)
	if err != nil {
		return nil, err
	}
	var posts []models.Post

	for _, item := range f.Chanel.Items {
		var p models.Post
		p.Title = item.Title
		p.Description = item.Description
		//p.Description = strip.StripTags(p.Description)
		p.Link = item.Link
		p.Author = item.Author
		p.Guid = item.Guid

		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubDate = t.Unix()
		}
		posts = append(posts, p)
	}
	return posts, nil
}

package models

// Модель данных поста/новости
type Post struct {
	Id          int
	Title       string
	Description string
	Link        string
	PubDate     int64
	Author      string
	Guid        string
}

type ShortPost struct {
	Id          int
	Title       string
	Description string
	PubDate     int64
}

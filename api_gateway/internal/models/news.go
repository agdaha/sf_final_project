package models

type NewsDetailed struct {
	News     News
	Comments Comments
}

type News struct {
	Id          int
	Title       string
	Description string
	Link        string
	PubDate     int
	Author      string
	Guid        string
}

type RoutineNews struct {
	News News
	Err  error
}

type ShortNews struct {
	Id          int
	Title       string
	Description string
	PubDate     int64
}

type Response struct {
	News  []ShortNews
	Pages Pages
}

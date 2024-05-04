package reader

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Подписка
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  Channel  `xml:"channel"`
}

// Метаданные Канала новостей
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Отдельная новость в подписке
type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
	Guid        string `xml:"guid"`
}

// Получение подписки новостей
func GetRss(url string) (*Feed, error) {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get rss error: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body rss error: %v", err)
	}

	var f Feed

	err = xml.Unmarshal(body, &f)
	if err != nil {
		return nil, fmt.Errorf("unmarshal rss error: %v", err)
	}

	return &f, nil
}

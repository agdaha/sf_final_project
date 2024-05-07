package getposts

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
	"github.com/agdaha/sf_final_project/news_service/internal/storage/memdb"
)

const postsPerPage = 5
const postCounts = 15

func TestGetPosts(t *testing.T) {
	db, _ := memdb.New(postCounts)
	handlerToTest := New(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	req := httptest.NewRequest("GET", "http://testing?page=2", nil)

	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf(" Not StatusOk")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
	}

	var r models.Response

	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Error("error decoder json ")
	}

	if len(r.News) != postsPerPage {
		t.Errorf("count of news error got:%v want:%v", len(r.News), postsPerPage)
	}

	total := postCounts / postsPerPage
	if r.Pages.Total != total {
		t.Errorf("wrong number total pages: got:%v want:%v ", r.Pages.Total, total)
	}

	if r.Pages.Current != 2 {
		t.Errorf("wrong number current: got:%v want:%v ", r.Pages.Current, 2)
	}
}

func TestGetPostsWithSearch(t *testing.T) {
	db, _ := memdb.New(postCounts)
	handlerToTest := New(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	req := httptest.NewRequest("GET", "http://testing?s=1", nil)

	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf(" Not StatusOk")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
	}

	var r models.Response

	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Error("error decoder json ")
	}
	log.Println(r)
	if len(r.News) != 5 {
		t.Errorf("count of news error got:%v want:%v", len(r.News), postsPerPage)
	}

	total := 2
	if r.Pages.Total != total {
		t.Errorf("wrong number total pages: got:%v want:%v ", r.Pages.Total, total)
	}

	if r.Pages.Current != 1 {
		t.Errorf("wrong number current: got:%v want:%v ", r.Pages.Current, 1)
	}
}

func find(ct []string, s string) bool {
	isFound := false
	for _, v := range ct {
		if v == s {
			isFound = true
			break
		}
	}
	return isFound
}

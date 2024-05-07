package getonepost

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
	"github.com/agdaha/sf_final_project/news_service/internal/storage/memdb"
	"github.com/julienschmidt/httprouter"
)

const postsPerPage = 5
const postCounts = 15

func TestGetPostById(t *testing.T) {
	const id = 1
	db, _ := memdb.New(postCounts)
	handlerToTest := New(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	req := httptest.NewRequest("GET", "http://testing/api/"+strconv.Itoa(id), nil)

	w := httptest.NewRecorder()

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/:newsid", handlerToTest)
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf(" Not StatusOk")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
	}

	var post models.Post

	err := json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		t.Error("error decoder json ")
	}

	if post.Id != id {
		t.Errorf(" not get correct id got:%v want:%v", post.Id, id)
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

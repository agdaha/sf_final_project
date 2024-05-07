package postcomment

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage/postgres"
	"github.com/stretchr/testify/suite"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	pgContainer *postgres.PostgresContainer
	store       storage.Store //проверка интерфейса
	ctx         context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := postgres.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	store, err := postgres.New(suite.pgContainer.ConnectionString, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err != nil {
		log.Fatal(err)
	}
	suite.store = store
}

func (suite *CustomerRepoTestSuite) TestPostCommentHandler() {
	t := suite.T()

	handlerToTest := New(suite.store, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))

	rec := strings.NewReader(`{
		"author": "Avast Avastovich",
		"comment_text": " Fisrt comment",
		"news_id":1,
		"parent_id":null
	}`)

	req := httptest.NewRequest("POST", "http://testing", rec)

	w := httptest.NewRecorder()
	handlerToTest.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf(" Not StatusOk")
	}
	ct := resp.Header["Content-Type"]
	isFound := find(ct, "application/json")

	if !isFound {
		t.Errorf(" Not application/json type content")
	}
	var r int

	err := json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Error("error decoder json ")
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

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}

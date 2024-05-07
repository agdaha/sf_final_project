package getcomments

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/agdaha/sf_final_project/comments_service/internal/models"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage/postgres"
	"github.com/julienschmidt/httprouter"
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

func (suite *CustomerRepoTestSuite) TestGetCommentsHandler() {
	t := suite.T()

	const id = 1

	handlerToTest := New(suite.store, slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))

	req := httptest.NewRequest("GET", "http://testing/api/"+strconv.Itoa(id), nil)

	w := httptest.NewRecorder()

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/:newsId", handlerToTest)
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

	var coms []models.Comment
	err := json.NewDecoder(resp.Body).Decode(&coms)
	if err != nil {
		t.Error("error decoder json ")
	}

	if len(coms) != 4 {
		t.Error("wrong count comments")
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

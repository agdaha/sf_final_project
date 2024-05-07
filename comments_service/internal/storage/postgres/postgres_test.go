package postgres

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/agdaha/sf_final_project/comments_service/internal/models"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	pgContainer *PostgresContainer
	store       storage.Store //проверка интерфейса
	ctx         context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	store, err := New(suite.pgContainer.ConnectionString, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err != nil {
		log.Fatal(err)
	}
	suite.store = store
}

func (suite *CustomerRepoTestSuite) TestAddAndGetComments() {
	t := suite.T()

	id, err := suite.store.Post(models.NewComment{
		Author:      "Rothua",
		CommentText: "TextComment",
		NewsId:      1,
		ParentId:    models.NullInt64{NullInt64: sql.NullInt64{Valid: false}},
	})

	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)

	coms, err := suite.store.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, coms[0].NewsId)
	assert.Equal(t, 5, len(coms))
	assert.Equal(t, int(id), coms[len(coms)-1].Id)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}

package postgres

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/agdaha/sf_final_project/news_service/internal/models"
	"github.com/agdaha/sf_final_project/news_service/internal/storage"
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

func (suite *CustomerRepoTestSuite) TestAddPosts() {
	t := suite.T()

	err := suite.store.UpdatePosts([]models.Post{
		{
			Title:       "First title",
			Description: "First Description",
			Link:        "link_123",
			PubDate:     111,
			Author:      "Author 1",
		},
		{
			Title:       "Second title",
			Description: "Second Description",
			Link:        "link_122",
			PubDate:     222,
			Author:      "Author 2",
		},
	})
	assert.NoError(t, err)

	posts, err := suite.store.Posts("", 1)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Equal(t, 2, len(posts.News))

	err = suite.store.UpdatePosts([]models.Post{
		{
			Title:       "Third title",
			Description: "Third Description",
			Link:        "link_125",
			PubDate:     333,
			Author:      "Author 3",
		},
		{
			Title:       "Second title update",
			Description: "Second Description update",
			Link:        "link_122",
			PubDate:     555,
			Author:      "Author 2",
		},
	})
	if err != nil {
		t.Error("Not updated posts")
	}
	posts, err = suite.store.Posts("", 1)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Equal(t, 3, len(posts.News))

	assert.Equal(t, "Second title update", posts.News[0].Title)
	assert.Equal(t, "Second Description update", posts.News[0].Description)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}

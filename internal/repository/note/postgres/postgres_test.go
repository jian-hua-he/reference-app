package postgres_test

import (
	"context"
	"database/sql"
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	dbpostgres "github.com/jian-hua-he/reference-app/internal/adapter/database/postgres"
	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/jian-hua-he/reference-app/internal/repository"
	"github.com/jian-hua-he/reference-app/internal/repository/note/postgres"
)

//go:embed testdata/schema.sql
var schemaSQL string

type PostgresRepoSuite struct {
	suite.Suite
	container *tcpostgres.PostgresContainer
	db        *sql.DB
	repo      *postgres.Repo
}

func TestPostgresRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoSuite))
}

func (s *PostgresRepoSuite) SetupSuite() {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(s.T(), err)
	s.container = container

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	db, err := dbpostgres.NewDBFromConnString(connStr)
	require.NoError(s.T(), err)

	_, err = db.ExecContext(ctx, schemaSQL)
	require.NoError(s.T(), err)

	s.db = db
	s.repo = postgres.NewRepo(db)
}

func (s *PostgresRepoSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.container != nil {
		s.container.Terminate(context.Background()) //nolint:errcheck
	}
}

func (s *PostgresRepoSuite) TearDownTest() {
	_, err := s.db.ExecContext(context.Background(), "DELETE FROM notes")
	require.NoError(s.T(), err)
}

func (s *PostgresRepoSuite) TestCreate() {
	ctx := context.Background()

	got, err := s.repo.Create(ctx, "hello postgres")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), got)

	assert.NotEmpty(s.T(), got.ID)
	assert.Equal(s.T(), "hello postgres", got.Text)
	assert.False(s.T(), got.CreatedAt.IsZero())

	// verify it was persisted
	var count int
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM notes WHERE id = $1", got.ID).Scan(&count)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, count)
}

func (s *PostgresRepoSuite) TestList_Empty() {
	ctx := context.Background()

	got, err := s.repo.List(ctx)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), []entity.Note{}, got)
}

func (s *PostgresRepoSuite) TestList_WithNotes() {
	ctx := context.Background()

	created1, err := s.repo.Create(ctx, "note one")
	require.NoError(s.T(), err)

	created2, err := s.repo.Create(ctx, "note two")
	require.NoError(s.T(), err)

	got, err := s.repo.List(ctx)
	require.NoError(s.T(), err)
	assert.Len(s.T(), got, 2)

	assert.ElementsMatch(s.T(), []entity.Note{*created1, *created2}, got)
}

func (s *PostgresRepoSuite) TestDelete_Existing() {
	ctx := context.Background()

	created, err := s.repo.Create(ctx, "to delete")
	require.NoError(s.T(), err)

	err = s.repo.Delete(ctx, created.ID)
	assert.NoError(s.T(), err)

	// verify it was deleted
	got, err := s.repo.List(ctx)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), []entity.Note{}, got)
}

func (s *PostgresRepoSuite) TestDelete_NonExistent() {
	ctx := context.Background()

	err := s.repo.Delete(ctx, "non-existent-id")
	assert.ErrorIs(s.T(), err, repository.ErrNotFound)
}

package migration_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	dbpostgres "github.com/jian-hua-he/reference-app/internal/adapter/database/postgres"
	"github.com/jian-hua-he/reference-app/internal/adapter/database/postgres/migration"
)

type MigrationSuite struct {
	suite.Suite
	container *tcpostgres.PostgresContainer
	db        *sql.DB
}

func TestMigrationSuite(t *testing.T) {
	suite.Run(t, new(MigrationSuite))
}

func (s *MigrationSuite) SetupSuite() {
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
	s.db = db
}

func (s *MigrationSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.container != nil {
		s.container.Terminate(context.Background()) //nolint:errcheck
	}
}

func (s *MigrationSuite) TearDownTest() {
	// Clean up schema_migrations and notes table between tests
	s.db.ExecContext(context.Background(), "DROP TABLE IF EXISTS schema_migrations") //nolint:errcheck
	s.db.ExecContext(context.Background(), "DROP TABLE IF EXISTS notes")             //nolint:errcheck
}

func (s *MigrationSuite) TestUp() {
	err := migration.Up(s.db)
	require.NoError(s.T(), err)

	// Verify notes table exists by inserting a row
	_, err = s.db.ExecContext(context.Background(),
		"INSERT INTO notes (id, text) VALUES ($1, $2)", "test-id", "test note")
	assert.NoError(s.T(), err)
}

func (s *MigrationSuite) TestUp_Idempotent() {
	err := migration.Up(s.db)
	require.NoError(s.T(), err)

	// Running Up again should not error
	err = migration.Up(s.db)
	assert.NoError(s.T(), err)
}

func (s *MigrationSuite) TestDown() {
	err := migration.Up(s.db)
	require.NoError(s.T(), err)

	err = migration.Down(s.db)
	require.NoError(s.T(), err)

	// Verify notes table no longer exists
	_, err = s.db.ExecContext(context.Background(),
		"INSERT INTO notes (id, text) VALUES ($1, $2)", "test-id", "test note")
	assert.Error(s.T(), err)
}

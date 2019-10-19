package testutils

import (
	"database/sql"
	"testing"

	"github.com/cockroachdb/errors"
	sqlMigrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"

	"personaapp/internal/migrate"
	"personaapp/internal/server/auth/storage"
	"personaapp/pkg/dockertest"
	"personaapp/pkg/postgresql"
)

func Migrate(db *sql.DB, md sqlMigrate.MigrationDirection) error {
	source := sqlMigrate.MemoryMigrationSource{
		Migrations: migrate.GetMigrations(),
	}
	if _, err := sqlMigrate.Exec(db, "postgres", source, md); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func EnsurePostgres(t *testing.T) *postgresql.Storage {
	cfg, err := dockertest.EnsurePostgres()
	require.NoError(t, err)

	pg, err := postgresql.New(&postgresql.Config{
		Host:               "localhost",
		Port:               uint16(cfg.Port),
		User:               cfg.User,
		Password:           cfg.Password,
		Database:           cfg.Database,
		MaxOpenConnections: 16,
		MaxIdleConnections: 16,
	})
	require.NoError(t, err)

	return pg
}

func InitStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
	pg := EnsurePostgres(t)
	require.NoError(t, Migrate(pg.DB, sqlMigrate.Up))

	return storage.New(pg), pg.Close
}

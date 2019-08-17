package register_test

import (
	"database/sql"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/cockroachdb/errors"
	cmdmigrate "personaapp/cmd/migrate"
	"github.com/stretchr/testify/require"
	storage "personaapp/internal/server/storage/register"
	"personaapp/pkg/dockertest"
	"personaapp/pkg/postgresql"
	"testing"
)

func initStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
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

	require.NoError(t, migrateUp(pg.DB))

	return storage.New(pg), pg.Close
}

func migrateUp(db *sql.DB) error {
	source := migrate.MemoryMigrationSource{
		Migrations: cmdmigrate.GetMigrationsForTests(),
	}
	if _, err := migrate.Exec(db, "postgres", source, migrate.Up); err != nil {
		return errors.WithStack(err)
	}
	return nil
}


func TestRegisterCompany(t *testing.T) {
	t.Run("normal flow", func(t *testing.T) {

	})

	t.Run("normal flow with dirty company name", func(t *testing.T) {

	})

	t.Run("normal flow with dirty email", func(t *testing.T) {

	})

	t.Run("normal flow with dirty phone", func(t *testing.T) {

	})

	t.Run("already existing", func(t *testing.T) {

	})

	t.Run("short company name", func(t *testing.T) {

	})

	t.Run("long company name", func(t *testing.T) {

	})

	t.Run("short email", func(t *testing.T) {

	})

	t.Run("long email", func(t *testing.T) {

	})

	t.Run("invalid email format", func(t *testing.T) {

	})

	t.Run("short phone", func(t *testing.T) {

	})

	t.Run("long phone", func(t *testing.T) {

	})

	t.Run("invalid phone format", func(t *testing.T) {

	})

	t.Run("short password", func(t *testing.T) {

	})

	t.Run("long password", func(t *testing.T) {

	})
}
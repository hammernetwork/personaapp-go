package controller_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cockroachdb/errors"
	migrate "github.com/rubenv/sql-migrate"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	cmdmigrate "personaapp/cmd/migrate"
	"personaapp/internal/server/controller"
	"personaapp/internal/server/storage"
	"personaapp/pkg/dockertest"
	"personaapp/pkg/postgresql"
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

// TestPingPong is a bit dirty test, but for initial example should be good
func TestPingPong(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
	}()

	c := controller.New(s)

	t.Run("not found", func(t *testing.T) {
		p, err := c.GetPing(context.Background(), "sample_key")
		require.Equal(t, controller.ErrNotFound, err)
		require.Nil(t, p)
	})

	t.Run("normal flow", func(t *testing.T) {
		key := uuid.NewV4().String()
		value := uuid.NewV4().String()
		err := c.SetPing(context.Background(), &controller.SetPing{
			Key:   key,
			Value: value,
		})
		require.Nil(t, err)

		p, err := c.GetPing(context.Background(), key)
		require.Nil(t, err)
		require.Equal(t, key, p.Key)
		require.Equal(t, value, p.Value)
	})

	t.Run("normal flow with overwrite", func(t *testing.T) {
		key := uuid.NewV4().String()
		value := uuid.NewV4().String()
		err1 := c.SetPing(context.Background(), &controller.SetPing{
			Key:   key,
			Value: value,
		})
		require.Nil(t, err1)

		p1, err2 := c.GetPing(context.Background(), key)
		require.Nil(t, err2)
		require.Equal(t, key, p1.Key)
		require.Equal(t, value, p1.Value)

		// overwrite with new value
		newValue := uuid.NewV4().String()
		err3 := c.SetPing(context.Background(), &controller.SetPing{
			Key:   key,
			Value: newValue,
		})
		require.Nil(t, err3)

		// check that new value
		p2, err4 := c.GetPing(context.Background(), key)
		require.Nil(t, err4)
		require.Equal(t, key, p2.Key)
		require.Equal(t, newValue, p2.Value)
	})
}

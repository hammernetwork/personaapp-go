package migrate_test

import (
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"

	"personaapp/internal/testutils"
)

func TestMigrations(t *testing.T) {
	pg := testutils.EnsurePostgres(t)
	defer func() {
		if err := pg.Close(); err != nil {
			t.Error(err)
		}
	}()

	require.NoError(t, testutils.Migrate(pg.DB, migrate.Up))
	require.NoError(t, testutils.Migrate(pg.DB, migrate.Down))
}

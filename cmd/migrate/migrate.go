package migrate

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	pkgmigrate "personaapp/pkg/migrate"
)

func Command() *cobra.Command {
	var migrations = []*migrate.Migration{
		{
			Id: "01 - Initial",
			Up: []string{
				`CREATE TABLE IF NOT EXISTS pingpong (
					key              VARCHAR(128)            PRIMARY KEY,
					value            VARCHAR(128),
					created_at       TIMESTAMPTZ             NOT NULL
					updated_at       TIMESTAMPTZ             NOT NULL
				);`,
			},

			Down: []string{
				`DROP TABLE IF EXISTS pingpong;`,
			},
		},
	}

	return pkgmigrate.Command("migrate_personaapp", migrations)
}

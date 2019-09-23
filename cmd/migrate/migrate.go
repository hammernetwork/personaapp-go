package migrate

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"

	pkgmigrate "personaapp/pkg/migrate"
)

var migrations = []*migrate.Migration{
	{
		Id: "01 - Initial",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS pingpong (
					key              VARCHAR(128)            PRIMARY KEY,
					value            VARCHAR(128),
					created_at       TIMESTAMPTZ             NOT NULL,
					updated_at       TIMESTAMPTZ             NOT NULL
				);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS pingpong;`,
		},
	},
	{
		Id: "02 - Create uuid extension",
		Up: []string{
			`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		},
		Down: []string{
			`DROP EXTENSION IF EXISTS "uuid-ossp";`,
		},
	},
	{
		Id: "03 - Create companies table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS company (
					id            uuid            PRIMARY KEY     DEFAULT uuid_generate_v4(),
					name          VARCHAR(100)    NOT NULL,
					email         VARCHAR(255)    NOT NULL,
					phone         VARCHAR(30)     NOT NULL,
					password      VARCHAR(100)    NOT NULL,
					created_at    TIMESTAMPTZ     NOT NULL,
					updated_at    TIMESTAMPTZ     NOT NULL
				);`,
			`CREATE UNIQUE INDEX company_email_idx ON company (email);`,
			`CREATE UNIQUE INDEX company_phone_idx ON company (phone);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS company;`,
		},
	},
	{
		Id: "04 - Create personas table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS persona (
					id            uuid            PRIMARY KEY     DEFAULT uuid_generate_v4(),
					first_name    VARCHAR(100)    NOT NULL,
					last_name     VARCHAR(100)    NOT NULL,
					email         VARCHAR(255)    NULL,
					phone         VARCHAR(30)     NOT NULL,
					password      VARCHAR(100)    NOT NULL,
					created_at    TIMESTAMPTZ     NOT NULL,
					updated_at    TIMESTAMPTZ     NOT NULL
				);`,
			`CREATE UNIQUE INDEX persona_email_idx ON persona (email);`,
			`CREATE UNIQUE INDEX persona_phone_idx ON persona (phone);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS persona;`,
		},
	},
}

func Command() *cobra.Command {
	return pkgmigrate.Command("migrate_personaapp", migrations)
}

func GetMigrationsForTests() []*migrate.Migration {
	return migrations
}

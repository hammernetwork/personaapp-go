package migrate

import (
	migrate "github.com/rubenv/sql-migrate"
)

var migrations = []*migrate.Migration{
	{
		Id: "01 - Create uuid extension",
		Up: []string{
			`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		},
		Down: []string{
			`DROP EXTENSION IF EXISTS "uuid-ossp";`,
		},
	},
	{
		Id: "02 - Create auth table",
		Up: []string{
			`CREATE TYPE e_account_type AS ENUM (
					'account_type_company',
					'account_type_persona'
			);`,

			`CREATE TABLE IF NOT EXISTS auth (
					account_id       uuid            PRIMARY KEY,
					account_type     e_account_type  NOT NULL,
					email            VARCHAR(255)    NOT NULL,
					phone            VARCHAR(30)     NOT NULL,
					password_hash    VARCHAR(100)    NOT NULL,
					created_at       TIMESTAMPTZ     NOT NULL,
					updated_at       TIMESTAMPTZ     NOT NULL
			);`,

			`CREATE UNIQUE INDEX auth_email_idx ON auth (email);`,
			`CREATE UNIQUE INDEX auth_phone_idx ON auth (phone);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS auth_phone_idx;`,
			`DROP INDEX IF EXISTS auth_email_idx;`,
			`DROP TABLE IF EXISTS auth;`,
			`DROP TYPE IF EXISTS e_account_type;`,
		},
	},
}

func GetMigrations() []*migrate.Migration {
	return migrations
}

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
	{
		Id: "03 - Create company scope table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS company_scope (
					company_scope_id	uuid			PRIMARY KEY,
					title				VARCHAR(255)	NOT NULL,
					alias				VARCHAR(255)	NOT NULL,
					created_at      	TIMESTAMPTZ     NOT NULL,
					updated_at      	TIMESTAMPTZ     NOT NULL
			);`,
			`CREATE UNIQUE INDEX company_scope_title_idx ON company_scope (title);`,
			`CREATE UNIQUE INDEX company_scope_alias_idx ON company_scope (alias);`,
			//TODO: insert temporary data
		},
		Down: []string{
			`DROP INDEX IF EXISTS company_scope_title_idx;`,
			`DROP INDEX IF EXISTS company_scope_alias_idx;`,
			`DROP TABLE IF EXISTS company_scope;`,
		},
	},
	{
		Id: "04 - Create company table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS company (
					auth_id			uuid			NOT NULL REFERENCES auth(account_id) ON DELETE CASCADE,
					scope_id		uuid			NULL REFERENCES company_scope(company_scope_id) ON DELETE RESTRICT,
					title			VARCHAR(255)	NULL,
					description		VARCHAR(255)	NULL,
					logo_url		VARCHAR(255)	NULL,
					created_at      TIMESTAMPTZ     NOT NULL,
					updated_at      TIMESTAMPTZ     NOT NULL
			);`,
			`CREATE UNIQUE INDEX company_auth_id_idx ON company (auth_id)`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS company_aut_id_idx;`,
			`DROP TABLE IF EXISTS company;`,
		},
	},
}

func GetMigrations() []*migrate.Migration {
	return migrations
}

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

			`CREATE UNIQUE INDEX auth_email_idx ON auth (email) WHERE NULLIF(TRIM(email),'') IS NOT NULL;`,
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
		Id: "03 - Create activity field table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS activity_field (
				id            uuid			   PRIMARY KEY,
				title	      VARCHAR(255)	   NOT NULL,
				icon_url	  VARCHAR(255)     NOT NULL,
				created_at    TIMESTAMPTZ      NOT NULL,
				updated_at    TIMESTAMPTZ      NOT NULL
			);`,
			`CREATE UNIQUE INDEX activity_field_title_idx ON activity_field (title);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS activity_field_title_idx;`,
			`DROP TABLE IF EXISTS activity_field;`,
		},
	},
	{
		Id: "04 - Create company table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS company (
				auth_id			uuid			NOT NULL REFERENCES auth (account_id) ON DELETE CASCADE,
				title			VARCHAR(255)	NULL,
				description		VARCHAR(255)	NULL,
				logo_url		VARCHAR(255)	NULL,
				created_at      TIMESTAMPTZ     NOT NULL,
				updated_at      TIMESTAMPTZ     NOT NULL
			);`,
			`CREATE UNIQUE INDEX company_auth_id_idx ON company (auth_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS company_auth_id_idx;`,
			`DROP TABLE IF EXISTS company;`,
		},
	},
	{
		Id: "05 - Create company activity fields table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS company_activity_fields (
				company_id           uuid            REFERENCES company (auth_id) ON DELETE CASCADE,
				activity_field_id    uuid            REFERENCES activity_field (id) ON DELETE CASCADE,
				CONSTRAINT company_activity_fields_pkey PRIMARY KEY (company_id, activity_field_id)
			);`,
			`CREATE INDEX company_activity_fields_company_id_idx ON company_activity_fields (company_id);`,
			`CREATE INDEX company_activity_fields_activity_field_id_idx ON company_activity_fields (activity_field_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS company_activity_fields_company_id_idx;`,
			`DROP INDEX IF EXISTS company_activity_fields_activity_field_id_idx;`,
			`DROP TABLE IF EXISTS company_activity_fields;`,
		},
	},
	{
		Id: "06 - Create vacancy category table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancy_category (
				id            uuid			   PRIMARY KEY,
				title		  VARCHAR(255)	   NOT NULL,
				icon_url	  VARCHAR(255)	   NOT NULL,
				created_at    TIMESTAMPTZ     NOT NULL,
				updated_at    TIMESTAMPTZ     NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS vacancy_category;`,
		},
	},
	{
		Id: "07 - Create vacancy table",
		Up: []string{
			`CREATE EXTENSION IF NOT EXISTS "postgis";`,
			`CREATE TABLE IF NOT EXISTS vacancy (
				id            			uuid					PRIMARY KEY,
				title		  			VARCHAR(255)			NOT NULL,
				description	  			VARCHAR(255)			NOT NULL,
				phone					VARCHAR(30)     		NOT NULL,
				image_url				VARCHAR(255)     		NOT NULL,
				min_salary	  			INTEGER					NULL,
				max_salary	  			INTEGER					NULL,
				location	  			GEOGRAPHY(POINT,4326)	NULL,
				work_months_experience  INTEGER					NOT NULL,
				work_schedule			VARCHAR(100)			NOT NULL,
				company_id	  			uuid					REFERENCES company (auth_id) ON DELETE CASCADE,
				position				SERIAL					NOT NULL,
				created_at       		TIMESTAMPTZ     		NOT NULL,
				updated_at       		TIMESTAMPTZ     		NOT NULL
			);`,
			`CREATE INDEX created_at_position_vacancy_idx ON vacancy (created_at, position);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS created_at_position_vacancy_idx;`,
			`DROP TABLE IF EXISTS vacancy;`,
		},
	},
	{
		Id: "08 - Create vacancies categories table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancies_categories (
				vacancy_id           uuid            REFERENCES vacancy (id) ON DELETE CASCADE,
				category_id    		 uuid            REFERENCES vacancy_category (id) ON DELETE CASCADE,
				CONSTRAINT vacancies_categories_pkey PRIMARY KEY (vacancy_id, category_id)
			);`,
			`CREATE UNIQUE INDEX vacancies_categories_unique_idx ON vacancies_categories (vacancy_id, category_id);`,
			`CREATE INDEX vacancies_categories_vacancy_id_idx ON vacancies_categories (vacancy_id);`,
			`CREATE INDEX vacancies_categories_category_id_idx ON vacancies_categories (category_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS vacancies_categories_unique_idx;`,
			`DROP INDEX IF EXISTS vacancies_categories_vacancy_id_idx;`,
			`DROP INDEX IF EXISTS vacancies_categories_category_id_idx;`,
			`DROP TABLE IF EXISTS vacancies_categories;`,
		},
	},
	{
		Id: "09 - Create vacancies images table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancies_images (
				vacancy_id         	uuid            	REFERENCES vacancy (id) ON DELETE CASCADE,
				position	  		INTEGER				NOT NULL,
				image_url			VARCHAR(255)     	NOT NULL
			);`,
			`CREATE INDEX vacancies_images_vacancy_id_idx ON vacancies_images (vacancy_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS vacancies_images_vacancy_id_idx;`,
			`DROP TABLE IF EXISTS vacancies_images;`,
		},
	},
}

func GetMigrations() []*migrate.Migration {
	return migrations
}

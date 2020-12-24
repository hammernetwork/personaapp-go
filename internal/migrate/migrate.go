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
				'account_type_persona',
				'account_type_admin'
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
			`CREATE UNIQUE INDEX auth_phone_idx ON auth (phone) WHERE NULLIF(TRIM(phone),'') IS NOT NULL;`,
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
				description		TEXT			NULL,
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
		Id: "05 - Create persona table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS persona (
				auth_id			uuid			NOT NULL REFERENCES auth (account_id) ON DELETE CASCADE,
				name			VARCHAR(255)	NULL,
				avatar_url		VARCHAR(255)	NULL,
				created_at      TIMESTAMPTZ     NOT NULL,
				updated_at      TIMESTAMPTZ     NOT NULL
			);`,
			`CREATE UNIQUE INDEX persona_auth_id_idx ON persona (auth_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS persona_auth_id_idx;`,
			`DROP TABLE IF EXISTS persona;`,
		},
	},
	{
		Id: "06 - Create extra auth email table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS auth_email (
				auth_id			uuid			NOT NULL REFERENCES auth (account_id) ON DELETE CASCADE,
				name			VARCHAR(255)	NOT NULL,
				email         	VARCHAR(255)    NOT NULL,
				CONSTRAINT auth_email_pkey PRIMARY KEY (auth_id, email)
			);`,
			`CREATE INDEX auth_email_auth_id_idx ON auth_email (auth_id);`,
			`CREATE INDEX auth_email_email_idx ON auth_email (email);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS auth_email_auth_id_idx;`,
			`DROP INDEX IF EXISTS auth_email_email_idx;`,
			`DROP TABLE IF EXISTS auth_email;`,
		},
	},
	{
		Id: "07 - Create extra auth phone table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS auth_phone (
				auth_id			uuid			NOT NULL REFERENCES auth (account_id) ON DELETE CASCADE,
				name			VARCHAR(255)	NOT NULL,
				phone         	VARCHAR(255)    NOT NULL,
				CONSTRAINT auth_phone_pkey PRIMARY KEY (auth_id, phone)
			);`,
			`CREATE INDEX auth_phone_auth_id_idx ON auth_phone (auth_id);`,
			`CREATE INDEX auth_phone_phone_idx ON auth_phone (phone);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS auth_phone_auth_id_idx;`,
			`DROP INDEX IF EXISTS auth_phone_phone_idx;`,
			`DROP TABLE IF EXISTS auth_phone;`,
		},
	},
	{
		Id: "08 - Create company activity fields table",
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
		Id: "09 - Create vacancy category table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancy_category (
				id            uuid			   PRIMARY KEY,
				title		  VARCHAR(255)	   NOT NULL,
				icon_url	  VARCHAR(255)	   NOT NULL,
				rating		  INTEGER		   NOT NULL,
				created_at    TIMESTAMPTZ      NOT NULL,
				updated_at    TIMESTAMPTZ      NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS vacancy_category;`,
		},
	},
	{
		Id: "10 - Create vacancy table",
		Up: []string{
			`CREATE EXTENSION IF NOT EXISTS "postgis";`,
			`CREATE TYPE e_vacancy_type AS ENUM (
				'vacancy_type_remote',
				'vacancy_type_normal'
			);`,
			`CREATE TABLE IF NOT EXISTS vacancy (
				id            			uuid					PRIMARY KEY,
				title		  			VARCHAR(255)			NOT NULL,
				description	  			TEXT					NOT NULL,
				phone					VARCHAR(30)     		NOT NULL,
				min_salary	  			INTEGER					NULL,
				max_salary	  			INTEGER					NULL,
				location	  			GEOGRAPHY(POINT,4326)	NULL,
				work_months_experience  INTEGER					NOT NULL,
				work_schedule			TEXT					NOT NULL,
				company_id	  			uuid					REFERENCES company (auth_id) ON DELETE CASCADE,
				position				SERIAL					NOT NULL,
				type					e_vacancy_type			NOT NULL,
				address					TEXT					NULL,
				country_code			INTEGER					NOT NULL,
				created_at       		TIMESTAMPTZ     		NOT NULL,
				updated_at       		TIMESTAMPTZ     		NOT NULL
			);`,
			`CREATE INDEX created_at_position_vacancy_idx ON vacancy (created_at, position);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS created_at_position_vacancy_idx;`,
			`DROP TABLE IF EXISTS vacancy;`,
			`DROP TYPE IF EXISTS e_vacancy_type;`,
		},
	},
	{
		Id: "11 - Create vacancies categories table",
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
		Id: "12 - Create vacancies images table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancies_images (
				vacancy_id         	uuid            	REFERENCES vacancy (id) ON DELETE CASCADE,
				position	  		INTEGER				NOT NULL,
				image_url			VARCHAR(255)     	NOT NULL
			);`,
			`CREATE UNIQUE INDEX vacancies_images_unique_idx ON vacancies_images (vacancy_id, position);`,
			`CREATE INDEX vacancies_images_vacancy_id_idx ON vacancies_images (vacancy_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS vacancies_images_unique_idx;`,
			`DROP INDEX IF EXISTS vacancies_images_vacancy_id_idx;`,
			`DROP TABLE IF EXISTS vacancies_images;`,
		},
	},
	{
		Id: "13 - Create city table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS city (
				id         			uuid            	PRIMARY KEY,
				name				VARCHAR(255)     	NOT NULL,
				country_code	  	INTEGER				NOT NULL,
				rating	  			INTEGER				NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS city;`,
		},
	},
	{
		Id: "14 - Create vacancy cities table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS vacancy_cities (
				vacancy_id           uuid            REFERENCES vacancy (id) ON DELETE CASCADE,
				city_id    		 	 uuid            REFERENCES city (id) ON DELETE CASCADE,
				CONSTRAINT vacancy_cities_pkey PRIMARY KEY (vacancy_id, city_id)
			);`,
			`CREATE UNIQUE INDEX vacancy_cities_unique_idx ON vacancy_cities (vacancy_id, city_id);`,
			`CREATE INDEX vacancy_cities_vacancy_id_idx ON vacancy_cities (vacancy_id);`,
			`CREATE INDEX vacancy_cities_city_id_idx ON vacancy_cities (city_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS vacancy_cities_unique_idx;`,
			`DROP INDEX IF EXISTS vacancy_cities_vacancy_id_idx;`,
			`DROP INDEX IF EXISTS vacancy_cities_city_id_idx;`,
			`DROP TABLE IF EXISTS vacancy_cities;`,
		},
	},
	{
		Id: "15 - Create auth_secret for password recovery table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS auth_secret (
				email					 VARCHAR(255)			PRIMARY KEY,
				secret           		 uuid            		NOT NULL,
				attempts    		 	 INTEGER            	NOT NULL,
				expiresAt    		 	 TIMESTAMPTZ            NOT NULL
			);`,
			`CREATE UNIQUE INDEX auth_secret_secret_idx ON auth_secret (secret);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS auth_secret_secret_idx;`,
			`DROP TABLE IF EXISTS auth_secret;`,
		},
	},
	{
		Id: "16 - Create cv table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS cv (
				id            			uuid					PRIMARY KEY,
				persona_id	  			uuid					REFERENCES auth (account_id) ON DELETE CASCADE,
				position				VARCHAR(255)			NULL,
				work_months_experience  INTEGER					NULL,
				min_salary	  			INTEGER					NULL,
				max_salary	  			INTEGER					NULL,
				created_at       		TIMESTAMPTZ     		NOT NULL,
				updated_at       		TIMESTAMPTZ     		NOT NULL
			);`,
			`CREATE INDEX created_at_position_id_cv_idx ON cv (created_at, persona_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS created_at_position_id_cv_idx;`,
			`DROP TABLE IF EXISTS cv;`,
		},
	},
	{
		Id: "17 - Create job type table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS job_type (
				id            			uuid					PRIMARY KEY,
				name					TEXT					NOT NULL,
				created_at    			TIMESTAMPTZ      		NOT NULL,
				updated_at    			TIMESTAMPTZ      		NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS job_type;`,
		},
	},
	{
		Id: "18 - Create cv job types table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS  cv_job_types  (
				cv_id           		uuid            		REFERENCES cv (id) ON DELETE CASCADE,
				job_type_id    			uuid            		REFERENCES job_type (id) ON DELETE CASCADE,
				CONSTRAINT cv_job_types_pkey PRIMARY KEY (cv_id, job_type_id)
			);`,
			`CREATE UNIQUE INDEX cv_job_types_unique_idx ON cv_job_types (cv_id, job_type_id);`,
			`CREATE INDEX cv_job_types_cv_id_idx ON cv_job_types (cv_id);`,
			`CREATE INDEX cv_job_types_job_type_id_idx ON cv_job_types (job_type_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS cv_job_types_unique_idx;`,
			`DROP INDEX IF EXISTS cv_job_types_cv_id_idx;`,
			`DROP INDEX IF EXISTS cv_job_types_job_type_id_idx;`,
			`DROP TABLE IF EXISTS cv_job_types;`,
		},
	},
	{
		Id: "19 - Create job_kind table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS job_kind (
				id            			uuid					PRIMARY KEY,
				name					TEXT					NOT NULL,
				created_at    			TIMESTAMPTZ      		NOT NULL,
				updated_at    			TIMESTAMPTZ      		NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS job_kind;`,
		},
	},
	{
		Id: "20 - Create cv kinds of job table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS cv_job_kinds  (
				cv_id           		uuid            		REFERENCES cv (id) ON DELETE CASCADE,
				job_kind_id    			uuid            		REFERENCES job_kind (id) ON DELETE CASCADE,
				CONSTRAINT cv_job_kinds_pkey PRIMARY KEY (cv_id, job_kind_id)
			);`,
			`CREATE UNIQUE INDEX cv_job_kinds_unique_idx ON cv_job_kinds (cv_id, job_kind_id);`,
			`CREATE INDEX cv_job_kinds_cv_id_idx ON cv_job_kinds (cv_id);`,
			`CREATE INDEX cv_job_kinds_job_kind_id_idx ON cv_job_kinds (job_kind_id);`,
		},
		Down: []string{
			`DROP INDEX IF EXISTS cv_job_kinds_unique_idx;`,
			`DROP INDEX IF EXISTS cv_job_kinds_cv_id_idx;`,
			`DROP INDEX IF EXISTS cv_job_kinds_job_kind_id_idx;`,
			`DROP TABLE IF EXISTS cv_job_kinds;`,
		},
	},
	{
		Id: "21 - Create experience table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS experience (
				id            			uuid					PRIMARY KEY,
				cv_id	  				uuid					REFERENCES cv (id) ON DELETE CASCADE,
				company_name			TEXT					NOT NULL,
				date_from       		TIMESTAMPTZ     		NULL,
				date_till       		TIMESTAMPTZ     		NULL,
				position				TEXT					NOT NULL,
				description				TEXT					NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS experience;`,
		},
	},
	{
		Id: "22 - Create education table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS education (
				id            			uuid					PRIMARY KEY,
				cv_id	  				uuid					REFERENCES cv (id) ON DELETE CASCADE,
				institution				TEXT					NOT NULL,
				date_from       		TIMESTAMPTZ     		NULL,
				date_till       		TIMESTAMPTZ     		NULL,
				speciality				TEXT					NOT NULL,
				description				TEXT					NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS education;`,
		},
	},
	{
		Id: "23 - Create custom_section table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS custom_section (
				id            			uuid					PRIMARY KEY,
				cv_id	  				uuid					REFERENCES cv (id) ON DELETE CASCADE,
				description				TEXT					NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS custom_section;`,
		},
	},
	{
		Id: "24 - Create story table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS story (
				id            			uuid					PRIMARY KEY,
				cv_id	  				uuid					REFERENCES cv (id) ON DELETE CASCADE,
				chapter_name			TEXT					NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS story cascade;`,
		},
	},
	{
		Id: "25 - Create story_episodes table",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS story_episodes (
				id            			uuid					PRIMARY KEY,
				stories_id	  			uuid					REFERENCES story (id) ON DELETE CASCADE,
				media_url				TEXT					NOT NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS stories_episode;`,
		},
	},
}

func GetMigrations() []*migrate.Migration {
	return migrations
}

package db

import (
	"fmt"
	"os"
	"persona/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx"
)

type DB struct {
	ConnectionPool *pgx.ConnPool
}

var db *DB

func Init() {

	m, err := migrate.New(
		"file://migrations",
		"postgres://localhost:5432/my_database?sslmode=disable",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set migration for database: %v\n", err)
		os.Exit(1)
	}

	m.Up()

	config := config.GetConfig()
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host: config.DatabaseHost,
			// Port              uint16 // default: 5432
			Database: config.DatabaseName,
			User:     "",
			Password: "",
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	db = &DB{
		ConnectionPool: connPool,
	}
}

func GetDB() *DB {
	return db
}

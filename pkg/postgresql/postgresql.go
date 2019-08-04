package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // register pg driver
	"github.com/pkg/errors"
)

type Storage struct {
	*sql.DB
}

func New(cfg *Config) (*Storage, error) {
	const format = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	dataSourceName := fmt.Sprintf(format, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open PostgreSQL, name: %s", dataSourceName)
	}
	if cfg.MaxOpenConnections != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}

	return &Storage{DB: db}, nil
}

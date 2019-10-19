package dockertest

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ory/dockertest"
)

const postgresComponent componentIdentifier = "postgres"

const (
	postgresDatabase = "personaapp"
	postgresUser     = "personaapp"
	postgresPassword = "personaapp"
)

type PostgresConfig struct {
	Port     int
	Database string
	User     string
	Password string
}

func init() {
	components.registerComponent(postgresComponent, initPostgresComponent)
}

func EnsurePostgres(modifiers ...OptionModifier) (*PostgresConfig, error) {
	port, err := components.ensureComponent(postgresComponent, modifiers...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &PostgresConfig{
		Port:     port,
		Database: postgresDatabase,
		User:     postgresUser,
		Password: postgresPassword,
	}, nil
}

func pgRunOptions(internalPort string) *dockertest.RunOptions {
	return &dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "9.6",
		ExposedPorts: []string{internalPort},
		Env: []string{
			"POSTGRES_USER=" + postgresUser,
			"POSTGRES_PASSWORD=" + postgresPassword,
			"POSTGRES_DB=" + postgresDatabase,
		},
	}
}

func initPostgresComponent(
	pool *dockertest.Pool,
	modifiers ...OptionModifier,
) (rport int, _ *dockertest.Resource, rerr error) {
	const (
		expireResource = 10 * time.Minute
		internalPort   = "5432"
		portID         = internalPort + "/tcp"
	)

	options := pgRunOptions(internalPort)
	applyOptionsModifiers(options, modifiers...)

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	defer func() {
		if rerr == nil {
			return
		}
		if err := pool.Purge(resource); err != nil {
			log.Println("failed to purge a pg resource", errors.WithStack(err))
		}
	}()

	var portStr string

	if err := pool.Retry(func() error {
		portStr = resource.GetPort(portID)
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			postgresUser,
			postgresPassword,
			"localhost",
			portStr,
			postgresDatabase,
		)
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(db.Ping())
	}); err != nil {
		return 0, nil, errors.WithStack(err)
	}

	if err := resource.Expire(uint(expireResource / time.Second)); err != nil {
		return 0, nil, errors.WithStack(err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	return port, resource, nil
}

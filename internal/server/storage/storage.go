package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"

	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type Ping struct {
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) TxPutPing(ctx context.Context, tx pkgtx.Tx, p *Ping) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
			UPDATE pingpong SET
				key = $1,
				value = $2,
				created_at = $3,
				updated_at = $4
			WHERE key = $1
			RETURNING *
		)
		INSERT INTO pingpong (key, value, created_at, updated_at)
		SELECT $1, $2, $3, $4
		WHERE NOT EXISTS (SELECT * FROM upsert)`,
		p.Key,
		p.Value,
		p.CreatedAt,
		p.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Storage) TxGetPing(ctx context.Context, tx pkgtx.Tx, key string) (*Ping, error) {
	var p Ping
	err := s.QueryRowContext(
		ctx,
		`SELECT * FROM pingpong WHERE
			key = $1`,
		key,
	).Scan(&p.Key, &p.Value, &p.CreatedAt, &p.UpdatedAt)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &p, nil
}

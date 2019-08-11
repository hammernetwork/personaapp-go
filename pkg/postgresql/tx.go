package postgresql

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"

	pkgtx "personaapp/pkg/tx"
)

// Client describes common pg operations which can be executed either in transaction or without it.
type Client interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func FromTx(tx pkgtx.Tx) Client {
	return tx.(Client)
}

type Tx struct {
	Client
	tx *sql.Tx
}

func (tx *Tx) Rollback() error {
	return errors.WithStack(tx.tx.Rollback())
}

func (tx *Tx) Commit() error {
	switch err := tx.tx.Commit(); err {
	case nil:
	case sql.ErrTxDone:
	default:
		return errors.WithStack(err)
	}
	return nil
}

type NoTx struct {
	Client
	db *sql.DB
}

func (tx *NoTx) Commit() error {
	return nil
}

func (tx *NoTx) Rollback() error {
	return nil
}

func (s *Storage) NoTx() pkgtx.Tx {
	return &NoTx{Client: s.DB, db: s.DB}
}

func (s *Storage) BeginTx(ctx context.Context) (pkgtx.Tx, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // TODO: weaken requirements
		ReadOnly:  false,                 // TODO: some Tx could be read only
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Tx{tx: tx, Client: tx}, nil
}

package controller

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"personaapp/internal/server/storage"
	pkgtx "personaapp/pkg/tx"
)

var ErrNotFound = errors.New("not found")

type SetPing struct {
	Key   string
	Value string
}

type Ping struct {
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface {
	TxPutPing(ctx context.Context, tx pkgtx.Tx, p *storage.Ping) error
	TxGetPing(ctx context.Context, tx pkgtx.Tx, key string) (*storage.Ping, error)

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

func (c *Controller) SetPing(ctx context.Context, sp *SetPing) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(i context.Context, tx pkgtx.Tx) error {
		var ping *storage.Ping
		switch p, err := c.s.TxGetPing(ctx, tx, sp.Key); err {
		case nil:
			ping = &storage.Ping{
				Key:       p.Key,
				Value:     sp.Value,
				CreatedAt: p.CreatedAt,
				UpdatedAt: time.Now(),
			}
		case storage.ErrNotFound:
			now := time.Now()
			ping = &storage.Ping{
				Key:       sp.Key,
				Value:     sp.Value,
				CreatedAt: now,
				UpdatedAt: now,
			}
		default:
			return errors.WithStack(err)
		}

		return errors.WithStack(c.s.TxPutPing(ctx, tx, ping))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) GetPing(ctx context.Context, key string) (*Ping, error) {
	p, err := c.s.TxGetPing(ctx, c.s.NoTx(), key)
	switch err {
	case nil:
	case storage.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}

	return &Ping{
		Key:       p.Key,
		Value:     p.Value,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

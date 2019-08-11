package tx

import (
	"context"
	"database/sql/driver"

	"github.com/pkg/errors"
)

type Tx interface {
	driver.Tx
}

var ErrConcurrentTx = errors.New("concurrent transaction")

// TxBeginner returns a new transaction.
type TxBeginner interface {
	BeginTx(ctx context.Context) (Tx, error)
}

// TxOption configures the way a transaction is executed.
type TxOption interface {
	apply(*txSettings)
}

type txSettings struct {
	attempts int
}

// RunInTx runs f in a transaction.
// Since f may be called multiple times, f should usually be idempotent.
func RunInTx(ctx context.Context, txer TxBeginner, fn func(context.Context, Tx) error, opts ...TxOption) error {
	settings := newTxSettings(opts)
	for n := 0; n < settings.attempts; n++ {
		tx, err := txer.BeginTx(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to begin tx")
		}

		if err = fn(ctx, tx); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != ErrConcurrentTx {
			return err
		}
	}

	return ErrConcurrentTx
}

func newTxSettings(opts []TxOption) *txSettings {
	s := &txSettings{attempts: 3}
	for _, o := range opts {
		o.apply(s)
	}
	return s
}

// MaxAttempts returns a TxOption that overrides the default 3 attempt times.
func MaxAttempts(attempts int) TxOption {
	return maxAttempts(attempts)
}

type maxAttempts int

func (w maxAttempts) apply(s *txSettings) {
	if w > 0 {
		s.attempts = int(w)
	}
}

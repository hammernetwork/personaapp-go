package register

import (
	"personaapp/pkg/postgresql"
	"context"
	pkgtx "personaapp/pkg/tx"
	"github.com/cockroachdb/errors"
	"time"
)

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type Company struct {
	Name string
	Email string
	Phone string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) TxCheckCompanyIsUnique(ctx context.Context, tx pkgtx.Tx, name string, email string) (bool, error) {
	c := postgresql.FromTx(tx)

	var count int
	err := c.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM company WHERE
			name = $1 OR email = $2`,
		name,
		email,
	).Scan(&count)

	if err != nil {
		return false, errors.WithStack(err)
	}
	return count == 0, nil
}

func (s *Storage) TxCreateCompany(ctx context.Context, tx pkgtx.Tx, cp *Company) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`INSERT INTO company (name, email, phone, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
		cp.Name,
		cp.Email,
		cp.Phone,
		cp.Password,
		cp.CreatedAt,
		cp.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
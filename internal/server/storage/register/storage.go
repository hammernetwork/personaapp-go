package register

import (
	"database/sql"
	"personaapp/pkg/postgresql"
	"context"
	pkgtx "personaapp/pkg/tx"
	"github.com/cockroachdb/errors"
	"time"
)

var ErrNotFound = errors.New("not found")

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

type Persona struct {
	FirstName string
	LastName string
	Email *string
	Phone string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) TxGetCompanyByEmailOrPhone(ctx context.Context, tx pkgtx.Tx, phone string, email string) (*Company, error) {
	c := postgresql.FromTx(tx)

	var cp Company
	err := c.QueryRowContext(
		ctx,
		`SELECT name, email, phone, password, created_at, updated_at 
			FROM company 
			WHERE phone = $1 OR email = $2`,
		phone,
		email,
	).Scan(&cp.Name, &cp.Email, &cp.Phone, &cp.Password, &cp.CreatedAt, &cp.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &cp, nil
}

func (s *Storage) TxPutCompany(ctx context.Context, tx pkgtx.Tx, cp *Company) error {
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

func (s *Storage) TxGetPersonaByPhone(ctx context.Context, tx pkgtx.Tx, phone string) (*Persona, error) {
	c := postgresql.FromTx(tx)

	var p Persona
	err := c.QueryRowContext(
		ctx,
		`SELECT first_name, last_name, email, phone, password, created_at, updated_at 
			FROM persona 
			WHERE phone = $1`,
		phone,
	).Scan(&p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Password, &p.CreatedAt, &p.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &p, nil
}

func (s *Storage) TxGetPersonaByEmailOrPhone(ctx context.Context, tx pkgtx.Tx, phone string, email string) (*Persona, error) {
	c := postgresql.FromTx(tx)

	var p Persona
	err := c.QueryRowContext(
		ctx,
		`SELECT first_name, last_name, email, phone, password, created_at, updated_at 
			FROM persona 
			WHERE phone = $1 OR email = $2`,
		phone,
		email,
	).Scan(&p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Password, &p.CreatedAt, &p.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &p, nil
}

func (s *Storage) TxPutPersona(ctx context.Context, tx pkgtx.Tx, p *Persona) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`INSERT INTO persona (first_name, last_name, email, phone, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		p.FirstName,
		p.LastName,
		p.Email,
		p.Phone,
		p.Password,
		p.CreatedAt,
		p.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
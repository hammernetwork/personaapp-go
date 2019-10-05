package storage

import (
	"context"
	"database/sql"
	"github.com/cockroachdb/errors"
	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
	"time"
)

type AccountType string

const (
	AccountTypeCompany AccountType = "account_type_company"
	AccountTypePersona AccountType = "account_type_persona"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type AuthData struct {
	AccountID    string
	Account      AccountType
	Email        string
	Phone        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *Storage) TxPutAuth(ctx context.Context, tx pkgtx.Tx, ad *AuthData) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
			UPDATE auth SET
				email = $2,
				phone = $3,
				password_hash = $4,
				created_at = $5,
				updated_at = $6
			WHERE account_id = $1
			RETURNING *
		)
		INSERT INTO auth (account_id, account_type, email, phone, password_hash, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $6
		WHERE NOT EXISTS (SELECT * FROM upsert)`,
		ad.AccountID,
		ad.Email,
		ad.Phone,
		ad.PasswordHash,
		ad.CreatedAt,
		ad.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Storage) TxGetAuthDataByAccountID(ctx context.Context, tx pkgtx.Tx, accountID string) (*AuthData, error) {
	c := postgresql.FromTx(tx)

	var ad AuthData
	err := c.QueryRowContext(
		ctx,
		`SELECT account_id, account_type, email, phone, password_hash, created_at, updated_at 
			FROM auth 
			WHERE account_id = $1`,
		accountID,
	).Scan(&ad.AccountID, &ad.Account, &ad.Email, &ad.Phone, &ad.PasswordHash, &ad.CreatedAt, &ad.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &ad, nil
}

func (s *Storage) TxGetAuthDataByPhoneOrEmail(
	ctx context.Context,
	tx pkgtx.Tx,
	phone string,
	email string,
) (*AuthData, error) {
	c := postgresql.FromTx(tx)

	var ad AuthData
	err := c.QueryRowContext(
		ctx,
		`SELECT account_id, account_type, email, phone, password_hash, created_at, updated_at 
			FROM auth 
			WHERE phone = $1 OR email = $2`,
		phone,
		email,
	).Scan(&ad.AccountID, &ad.Account, &ad.Email, &ad.Phone, &ad.PasswordHash, &ad.CreatedAt, &ad.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &ad, nil
}

func (s *Storage) TxGetAuthDataByPhone(
	ctx context.Context,
	tx pkgtx.Tx,
	phone string,
) (*AuthData, error) {
	c := postgresql.FromTx(tx)

	var ad AuthData
	err := c.QueryRowContext(
		ctx,
		`SELECT account_id, account_type, email, phone, password_hash, created_at, updated_at 
			FROM auth
			WHERE phone = $1`,
		phone,
	).Scan(&ad.AccountID, &ad.Account, &ad.Email, &ad.Phone, &ad.PasswordHash, &ad.CreatedAt, &ad.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &ad, nil
}

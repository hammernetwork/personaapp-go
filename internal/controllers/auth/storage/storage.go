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
	AccountTypeAdmin   AccountType = "account_type_admin"
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

type AuthSecret struct {
	Email     string
	Secret    string
	Attempts  int
	ExpiresAt time.Time
}

func (s *Storage) TxPutAuth(ctx context.Context, tx pkgtx.Tx, ad *AuthData) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
			UPDATE auth SET
				account_type = $2,
				email = $3,
				phone = $4,
				password_hash = $5,
				created_at = $6,
				updated_at = $7
			WHERE account_id = $1
			RETURNING *
		)
		INSERT INTO auth (account_id, account_type, email, phone, password_hash, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT * FROM upsert)`,
		ad.AccountID,
		ad.Account,
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

func (s *Storage) TxGetAuthDataByEmail(
	ctx context.Context,
	tx pkgtx.Tx,
	email string,
) (*AuthData, error) {
	c := postgresql.FromTx(tx)

	var ad AuthData
	err := c.QueryRowContext(
		ctx,
		`SELECT account_id, account_type, email, phone, password_hash, created_at, updated_at 
			FROM auth
			WHERE email = $1`,
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

func (s *Storage) TxGetAuthDataByID(ctx context.Context, tx pkgtx.Tx, accountID string) (*AuthData, error) {
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

func (s *Storage) TxGetAuthSecretByEmail(ctx context.Context, tx pkgtx.Tx, email string) (*AuthSecret, error) {
	c := postgresql.FromTx(tx)

	var as AuthSecret
	err := c.QueryRowContext(
		ctx,
		`SELECT email, secret, attempts, expiresAt
			FROM auth_secret
			WHERE email = $1`,
		email,
	).Scan(&as.Email, &as.Secret, &as.Attempts, &as.ExpiresAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}

	return &as, nil
}

func (s *Storage) TxGetAuthSecretBySecret(ctx context.Context, tx pkgtx.Tx, secret string) (*AuthSecret, error) {
	c := postgresql.FromTx(tx)

	var as AuthSecret
	err := c.QueryRowContext(
		ctx,
		`SELECT email, secret, attempts, expiresAt
			FROM auth_secret
			WHERE secret = $1`,
		secret,
	).Scan(&as.Email, &as.Secret, &as.Attempts, &as.ExpiresAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}

	return &as, nil
}

func (s *Storage) TxPutAuthSecretByEmail(ctx context.Context, tx pkgtx.Tx, authSecret *AuthSecret) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE auth_secret SET
					secret = $2,
					attempts = $3,
					expiresAt = $4
				WHERE email = $1
				RETURNING email, secret, attempts, expiresAt
			)
			INSERT INTO auth_secret (email, secret, attempts, expiresAt)
			SELECT $1, $2, $3, $4
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		authSecret.Email,
		authSecret.Secret,
		authSecret.Attempts,
		authSecret.ExpiresAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxDeleteAuthSecret(ctx context.Context, tx pkgtx.Tx, email string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM auth_secret 
			WHERE email = $1`,
		email,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

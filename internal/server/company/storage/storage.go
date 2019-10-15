package storage

import (
	"context"
	"database/sql"
	"github.com/cockroachdb/errors"
	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
	"time"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type CompanyData struct {
	AuthID      string
	Title       string
	Description string
	LogoURL     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ActivityField struct {
	ID    string
	Title string
	Alias string
}

func (s *Storage) TxGetCompanyByID(ctx context.Context, tx pkgtx.Tx, authID string) (*CompanyData, error) {
	c := postgresql.FromTx(tx)

	var cd CompanyData
	err := c.QueryRowContext(
		ctx,
		`SELECT auth_id, title, description, logo_url, created_at, updated_at 
			FROM company 
			WHERE auth_id = $1;`,
		authID,
	).Scan(&cd.AuthID, &cd.Title, &cd.Description, &cd.LogoURL, &cd.CreatedAt, &cd.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}
	return &cd, nil
}

func (s *Storage) TxPutCompany(ctx context.Context, tx pkgtx.Tx, cd *CompanyData) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
			UPDATE company SET
				title = $2,
				description = $3,
				logo_url = $4,
				created_at = $5,
				updated_at = $6
			WHERE auth_id = $1
			RETURNING *
		)
		INSERT INTO auth (auth_id, title, description, logo_url, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $6
		WHERE NOT EXISTS (SELECT * FROM upsert)`,
		cd.AuthID,
		cd.Title,
		cd.Description,
		cd.LogoURL,
		cd.CreatedAt,
		cd.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Storage) TxGetCompanyActivityFieldsByID(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
) ([]*ActivityField, error) {
	//TODO: implement
	return nil, nil
}
func (s *Storage) TxPutCompanyActivityFields(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
	activityFields []*ActivityField,
) error {
	//TODO: implement
	return nil
}

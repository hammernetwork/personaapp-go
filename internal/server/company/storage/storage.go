package storage

import (
	"context"
	"github.com/cockroachdb/errors"
	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
	"time"
)

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type Fields uint8

const (
	FieldScopeID Fields = 1 << iota
	FieldTitle
	FieldDescription
	FieldLogoURL
)

type CompanyData struct {
	Fields      Fields
	AuthID      string
	ScopeID     string
	Title       string
	Description string
	LogoURL     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Storage) TxPutCompany(ctx context.Context, tx pkgtx.Tx, cd *CompanyData) error {
	// TODO: add bit mask check

	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
			UPDATE company SET
				scope_id = $2,
				title = $3,
				description = $4,
				logo_url = $5,
				created_at = $6,
				updated_at = $7
			WHERE auth_id = $1
			RETURNING *
		)
		INSERT INTO auth (auth_id, scope_id, title, description, logo_url, created_at, updated_at)
		SELECT $1, $2, $3, $4, $5, $6, $7
		WHERE NOT EXISTS (SELECT * FROM upsert)`,
		cd.AuthID,
		cd.ScopeID,
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

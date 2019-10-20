package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cockroachdb/errors"
	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
	"strings"
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
		WHERE NOT EXISTS (SELECT * FROM upsert);`,
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

func (s *Storage) TxGetActivityFieldsByCompanyID(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
) ([]*ActivityField, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT caf.activity_field_id, af.title, af.alias
			FROM company_activity_fields AS caf
			INNER JOIN activity_field AS af
			ON caf.activity_field_id = af.id
			WHERE caf.company_id = $1;`,
		authID,
	)

	switch err {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	afs := make([]*ActivityField, 0)

	for rows.Next() {
		var af ActivityField
		if err := rows.Scan(&af.ID, &af.Title, &af.Alias); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		afs = append(afs, &af)
	}

	return afs, nil
}

func (s *Storage) TxPutCompanyActivityFields(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
	activityFieldsIDs []string,
) error {
	c := postgresql.FromTx(tx)

	now := time.Now()

	queryFormat := `INSERT 
		INTO company_activity_fields (company_id, activity_field_id, created_at, updated_at)
		VALUES %s;`

	columns := 4
	valueStrings := make([]string, 0, len(activityFieldsIDs))
	valueArgs := make([]interface{}, 0, len(activityFieldsIDs)*columns)

	for i := 0; i < len(activityFieldsIDs); i++ {
		offset := i * columns
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4))
		valueArgs = append(valueArgs, authID, activityFieldsIDs[i], now, now)
	}

	if _, err := c.ExecContext(
		ctx,
		fmt.Sprintf(queryFormat, strings.Join(valueStrings, ",")),
		valueArgs...,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxDeleteCompanyActivityFieldsByCompanyID(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE 
			FROM company_activity_fields
			WHERE company_id = $1`,
		authID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

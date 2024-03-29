package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
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
	ID          string
	Title       string
	Description string
	LogoURL     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ActivityField struct {
	ID        string
	Title     string
	IconURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Storage) TxGetCompanyByID(ctx context.Context, tx pkgtx.Tx, authID string) (*CompanyData, error) {
	c := postgresql.FromTx(tx)

	var cd CompanyData
	err := c.QueryRowContext(
		ctx,
		`SELECT auth_id, title, description, logo_url, created_at, updated_at 
			FROM company 
			WHERE auth_id = $1`,
		authID,
	).Scan(&cd.ID, &cd.Title, &cd.Description, &cd.LogoURL, &cd.CreatedAt, &cd.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}

	return &cd, nil
}

func (s *Storage) TxGetCompaniesByID(
	ctx context.Context,
	tx pkgtx.Tx,
	companyIDs []string,
) (_ []*CompanyData, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT auth_id, title, description, logo_url, created_at, updated_at
			FROM company
			WHERE auth_id = ANY($1::uuid[])`,
		pq.Array(companyIDs),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			if rerr != nil {
				rerr = errors.WithSecondaryError(rerr, err)
				return
			}

			rerr = errors.WithStack(err)
		}
	}()

	var cs []*CompanyData

	for rows.Next() {
		var cd CompanyData

		err := rows.Scan(&cd.ID, &cd.Title, &cd.Description, &cd.LogoURL, &cd.CreatedAt, &cd.UpdatedAt)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		cs = append(cs, &cd)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return cs, nil
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
			INSERT INTO company (auth_id, title, description, logo_url, created_at, updated_at)
			SELECT $1, $2, $3, $4, $5, $6
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		cd.ID,
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
) (_ []*ActivityField, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT caf.activity_field_id, af.title, af.icon_url
			FROM company_activity_fields AS caf
			INNER JOIN activity_field AS af
			ON caf.activity_field_id = af.id
			WHERE caf.company_id = $1`,
		authID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			if rerr != nil {
				rerr = errors.WithSecondaryError(rerr, err)
				return
			}

			rerr = errors.WithStack(err)
		}
	}()

	afs := make([]*ActivityField, 0)

	for rows.Next() {
		var af ActivityField
		if err := rows.Scan(&af.ID, &af.Title, &af.IconURL); err != nil {
			return nil, errors.WithStack(err)
		}

		afs = append(afs, &af)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return afs, nil
}

func (s *Storage) TxGetActivityFieldsByID(
	ctx context.Context,
	tx pkgtx.Tx,
	activityFieldID string,
) (_ *ActivityField, rerr error) {
	c := postgresql.FromTx(tx)

	var af ActivityField
	err := c.QueryRowContext(
		ctx,
		`SELECT af.id, af.title, af.icon_url, af.created_at, af.updated_at
			FROM activity_field AS af
			WHERE id = $1`,
		activityFieldID,
	).Scan(&af.ID, &af.Title, &af.IconURL, &af.CreatedAt, &af.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, errors.WithStack(err)
	}

	return &af, nil
}

func (s *Storage) TxPutCompanyActivityFields(
	ctx context.Context,
	tx pkgtx.Tx,
	authID string,
	activityFieldsIDs []string,
) error {
	c := postgresql.FromTx(tx)

	queryFormat := `INSERT 
		INTO company_activity_fields (company_id, activity_field_id) 
		VALUES %s`

	columns := 2
	length := len(activityFieldsIDs)
	valueStrings := make([]string, length)
	valueArgs := make([]interface{}, length*columns)

	for i := 0; i < length; i++ {
		offset := i * columns
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", offset+1, offset+2)
		valueArgs[offset] = authID
		valueArgs[offset+1] = activityFieldsIDs[i]
	}

	query := fmt.Sprintf(queryFormat, strings.TrimSuffix(strings.Join(valueStrings, ","), ","))
	if _, err := c.ExecContext(
		ctx,
		query,
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

func (s *Storage) TxPutActivityField(ctx context.Context, tx pkgtx.Tx, af *ActivityField) error {
	c := postgresql.FromTx(tx)
	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE activity_field SET
					title = $2,
					icon_url = $3,
					created_at = $4,
					updated_at = $5
				WHERE id = $1
				RETURNING *
			)
			INSERT INTO activity_field (id, title, icon_url, created_at, updated_at)
			SELECT $1, $2, $3, $4, $5
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		af.ID,
		af.Title,
		af.IconURL,
		af.CreatedAt,
		af.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetActivityFields(
	ctx context.Context,
	tx pkgtx.Tx,
) (_ []*ActivityField, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT af.id, af.title, af.icon_url
			FROM activity_field AS af`,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			if rerr != nil {
				rerr = errors.WithSecondaryError(rerr, err)
				return
			}

			rerr = errors.WithStack(err)
		}
	}()

	afs := make([]*ActivityField, 0)

	for rows.Next() {
		var af ActivityField
		if err := rows.Scan(&af.ID, &af.Title, &af.IconURL); err != nil {
			return nil, errors.WithStack(err)
		}

		afs = append(afs, &af)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return afs, nil
}

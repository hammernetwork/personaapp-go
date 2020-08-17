package storage

import (
	"context"
	"database/sql"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
	"personaapp/pkg/postgresql"
	pkgtx "personaapp/pkg/tx"
)

var ErrNotFound = errors.New("not found")

type Storage struct {
	*postgresql.Storage
}

func New(db *postgresql.Storage) *Storage {
	return &Storage{db}
}

type City struct {
	ID          string
	Name        string
	CountryCode int32
	Rating      int32
}

/**
Cities part start
*/
func (s *Storage) TxGetCitiesList(
	ctx context.Context,
	tx pkgtx.Tx,
	countryCodes []int32,
	rating int32,
	filter string,
) (_ []*City, rerr error) {
	c := postgresql.FromTx(tx)

	like := "%" + filter + "%"
	rows, err := c.QueryContext(
		ctx,
		`SELECT id, name, country_code, rating
				FROM city
				WHERE ($1 = '{}' OR country_code = ANY($1::INTEGER[]))
					AND rating >= $2
					AND ($3 = '' OR name LIKE $3)
				ORDER BY name ASC`,
		pq.Array(countryCodes),
		rating,
		like,
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

	cities := make([]*City, 0)

	for rows.Next() {
		var city City
		if err := rows.Scan(&city.ID, &city.Name, &city.CountryCode, &city.Rating); err != nil {
			return nil, errors.WithStack(err)
		}

		cities = append(cities, &city)
	}
	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return cities, nil
}

func (s *Storage) TxPutCity(ctx context.Context, tx pkgtx.Tx, city *City) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE city SET
					name = $2,
					country_code = $3,
					rating = $4
				WHERE id = $1
				RETURNING id, name, country_code, rating
			)
			INSERT INTO city (id, name, country_code, rating)
			SELECT $1, $2, $3, $4
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		city.ID,
		city.Name,
		city.CountryCode,
		city.Rating,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetCity(ctx context.Context, tx pkgtx.Tx, cityID string) (*City, error) {
	c := postgresql.FromTx(tx)

	var city City
	err := c.QueryRowContext(
		ctx,
		`SELECT id, name, country_code, rating
				FROM city
				WHERE id = $1`,
		cityID,
	).Scan(&city.ID, &city.Name, &city.CountryCode, &city.Rating)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, errors.WithStack(ErrNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &city, nil
}

func (s *Storage) TxDeleteCity(ctx context.Context, tx pkgtx.Tx, cityID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM city 
			WHERE id = $1`,
		cityID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Cities part end
*/

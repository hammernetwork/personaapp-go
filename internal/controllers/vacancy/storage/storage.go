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

type VacancyCategory struct {
	ID        string
	Title     string
	IconURL   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VacancyCategoryShort struct {
	VacancyID string
	ID        string
	Title     string
}

type Vacancy struct {
	ID        string
	Title     string
	Phone     string
	MinSalary int32
	MaxSalary int32
	CompanyID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VacancyDetails struct {
	Vacancy
	Description          string
	WorkMonthsExperience int32
	WorkSchedule         string
	LocationLatitude     float32
	LocationLongitude    float32
}

type Cursor struct {
	PrevCreatedAt time.Time
	PrevPosition  int
}

func (s *Storage) TxGetVacanciesCategoriesList(ctx context.Context, tx pkgtx.Tx) (_ []*VacancyCategory, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(ctx, `SELECT id, title, icon_url FROM vacancy_category`)
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

	vcs := make([]*VacancyCategory, 0)

	for rows.Next() {
		var vc VacancyCategory
		if err := rows.Scan(&vc.ID, &vc.Title, &vc.IconURL); err != nil {
			return nil, errors.WithStack(err)
		}

		vcs = append(vcs, &vc)
	}
	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return vcs, nil
}

func (s *Storage) TxPutVacancyCategory(ctx context.Context, tx pkgtx.Tx, category *VacancyCategory) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE vacancy_category SET
					title = $2,
					icon_url = $3,
					updated_at = $5
				WHERE id = $1
				RETURNING id, title, icon_url
			)
			INSERT INTO vacancy_category (id, title, icon_url, created_at, updated_at)
			SELECT $1, $2, $3, $4, $5
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		category.ID,
		category.Title,
		category.IconURL,
		category.CreatedAt,
		category.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetVacancyCategory(ctx context.Context, tx pkgtx.Tx, categoryID string) (*VacancyCategory, error) {
	c := postgresql.FromTx(tx)

	var vc VacancyCategory
	err := c.QueryRowContext(
		ctx,
		`SELECT id, title, icon_url
				FROM vacancy_category
				WHERE id = $1`,
		categoryID,
	).Scan(&vc.ID, &vc.Title, &vc.IconURL)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, errors.WithStack(ErrNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &vc, nil
}

func (s *Storage) TxGetVacancyDetails(ctx context.Context, tx pkgtx.Tx, vacancyID string) (*VacancyDetails, error) {
	c := postgresql.FromTx(tx)

	var vd VacancyDetails
	err := c.QueryRowContext(
		ctx,
		`SELECT id, title, description, phone, min_salary, max_salary, company_id, work_months_experience, 
					work_schedule, ST_X(location::geometry), ST_Y(location::geometry), created_at, updated_at
				FROM vacancy
				WHERE id = $1`,
		vacancyID,
	).Scan(&vd.ID, &vd.Title, &vd.Description, &vd.Phone, &vd.MinSalary, &vd.MaxSalary, &vd.CompanyID,
		&vd.WorkMonthsExperience, &vd.WorkSchedule, &vd.LocationLongitude, &vd.LocationLatitude,
		&vd.CreatedAt, &vd.UpdatedAt)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, errors.WithStack(ErrNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &vd, nil
}

func (s *Storage) TxPutVacancy(ctx context.Context, tx pkgtx.Tx, vacancy *VacancyDetails) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE vacancy SET
					title = $2,
					description = $3,
					phone = $4,
					min_salary = $5,
					max_salary = $6,
					company_id = $7,
					work_months_experience = $8,
					work_schedule = $9,
					location = ST_SetSRID(ST_MakePoint($10, $11), 4326),
					updated_at = $13
				WHERE id = $1
				RETURNING id, title, description, phone, min_salary, max_salary, company_id, work_months_experience, 
					work_schedule, location, created_at, updated_at
			)
			INSERT INTO vacancy (id, title, description, phone, min_salary, max_salary, company_id, work_months_experience, 
					work_schedule, location, created_at, updated_at)
			SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, ST_SetSRID(ST_MakePoint($10, $11), 4326), $12, $13
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		vacancy.ID,
		vacancy.Title,
		vacancy.Description,
		vacancy.Phone,
		vacancy.MinSalary,
		vacancy.MaxSalary,
		vacancy.CompanyID,
		vacancy.WorkMonthsExperience,
		vacancy.WorkSchedule,
		vacancy.LocationLongitude,
		vacancy.LocationLatitude,
		vacancy.CreatedAt,
		vacancy.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// nolint:funlen // will rework
func (s *Storage) TxGetVacanciesList(
	ctx context.Context,
	tx pkgtx.Tx,
	categoriesIDs []string,
	limit int,
	cursor *Cursor,
) (_ []*Vacancy, _ *Cursor, rerr error) {
	cursorCreatedAt := time.Now()
	cursorPosition := -1

	if cursor != nil {
		cursorCreatedAt = cursor.PrevCreatedAt
		cursorPosition = cursor.PrevPosition
	}

	c := postgresql.FromTx(tx)
	rows, err := c.QueryContext(
		ctx,
		`WITH filtered_categories AS (
			SELECT vacancy_id
			FROM vacancies_categories
			WHERE category_id = ANY($1::uuid[])
		)
		SELECT v.id, v.title, v.phone, v.min_salary, v.max_salary, v.company_id, v.position, v.created_at
		FROM vacancy AS v
		WHERE ($1 = '{}' OR v.id IN (SELECT vacancy_id FROM filtered_categories))
		AND ($3 < 0 OR (v.created_at, v.position) < ($4, $3))
		ORDER BY v.created_at DESC, v.position DESC
		LIMIT $2`,
		pq.Array(categoriesIDs),
		limit,
		cursorPosition,
		cursorCreatedAt,
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
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

	vs := make([]*Vacancy, 0)

	var lastCreatedAt time.Time
	var lastPosition int

	for rows.Next() {
		var v Vacancy

		err := rows.Scan(
			&v.ID,
			&v.Title,
			&v.Phone,
			&v.MinSalary,
			&v.MaxSalary,
			&v.CompanyID,
			&lastPosition,
			&lastCreatedAt,
		)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		vs = append(vs, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	if len(vs) == 0 || len(vs) < limit {
		return vs, nil, nil
	}

	return vs, &Cursor{
		PrevCreatedAt: lastCreatedAt,
		PrevPosition:  lastPosition,
	}, nil
}

func (s *Storage) TxGetVacanciesCategories(
	ctx context.Context,
	tx pkgtx.Tx,
	vacancyIDs []string,
) ([]*VacancyCategoryShort, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT vscs.vacancy_id, vc.id, vc.title
			FROM vacancies_categories AS vscs
			INNER JOIN vacancy_category AS vc
			ON vscs.category_id = vc.id
			WHERE vacancy_id = ANY($1::uuid[])`,
		pq.Array(vacancyIDs),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vcs := make([]*VacancyCategoryShort, 0)

	for rows.Next() {
		var vc VacancyCategoryShort
		if err := rows.Scan(&vc.VacancyID, &vc.ID, &vc.Title); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		vcs = append(vcs, &vc)
	}

	return vcs, nil
}

func (s *Storage) TxPutVacancyCategories(
	ctx context.Context,
	tx pkgtx.Tx,
	vacancyID string,
	categoriesIDs []string,
) error {
	c := postgresql.FromTx(tx)

	queryFormat := `INSERT 
		INTO vacancies_categories (vacancy_id, category_id)
		VALUES %s`

	columns := 2
	valueStrings := make([]string, len(categoriesIDs))
	valueArgs := make([]interface{}, len(categoriesIDs)*columns)

	for i := 0; i < len(categoriesIDs); i++ {
		offset := i * columns
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", offset+1, offset+2)
		valueArgs[offset] = vacancyID
		valueArgs[offset+1] = categoriesIDs[i]
	}

	if _, err := c.ExecContext(
		ctx,
		fmt.Sprintf(queryFormat, strings.TrimSuffix(strings.Join(valueStrings, ","), ",")),
		valueArgs...,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxDeleteVacancyCategories(ctx context.Context, tx pkgtx.Tx, vacancyID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM vacancies_categories 
			WHERE vacancy_id = $1`,
		vacancyID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Vacancy images part start
*/

func (s *Storage) TxGetVacanciesImages(
	ctx context.Context,
	tx pkgtx.Tx,
	vacancyIDs []string,
) (map[string][]string, error) {
	c := postgresql.FromTx(tx)

	vacancyImageMap := make(map[string][]string, len(vacancyIDs))

	if len(vacancyIDs) == 0 {
		return vacancyImageMap, nil
	}

	rows, err := c.QueryContext(
		ctx,
		// nolint:gosec // will rework
		`SELECT vacancy_id, image_url
			FROM vacancies_images
			WHERE vacancy_id = ANY($1::uuid[])
			ORDER BY vacancy_id ASC, position ASC`,
		pq.Array(vacancyIDs),
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	for rows.Next() {
		var (
			vacancyID string
			imageURL  string
		)

		err := rows.Scan(&vacancyID, &imageURL)
		if err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		vacancyImageMap[vacancyID] = append(vacancyImageMap[vacancyID], imageURL)
	}

	return vacancyImageMap, nil
}

func (s *Storage) TxPutVacancyImages(
	ctx context.Context,
	tx pkgtx.Tx,
	vacancyID string,
	imageUrls []string,
) error {
	c := postgresql.FromTx(tx)

	queryFormat := `INSERT 
		INTO vacancies_images (vacancy_id, position, image_url)
		VALUES %s`

	columns := 3
	valueStrings := make([]string, len(imageUrls))
	valueArgs := make([]interface{}, len(imageUrls)*columns)

	for i := 0; i < len(imageUrls); i++ {
		offset := i * columns
		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d)", offset+1, offset+2, offset+3)
		valueArgs[offset] = vacancyID
		valueArgs[offset+1] = i
		valueArgs[offset+2] = imageUrls[i]
	}

	if _, err := c.ExecContext(
		ctx,
		fmt.Sprintf(queryFormat, strings.TrimSuffix(strings.Join(valueStrings, ","), ",")),
		valueArgs...,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxDeleteVacancyImages(ctx context.Context, tx pkgtx.Tx, vacancyID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM vacancies_images 
			WHERE vacancy_id = $1`,
		vacancyID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Vacancy images part end
*/
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

type JobType struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CVJobType struct {
	ID   string
	Name string
}

type JobKind struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CVJobKind struct {
	ID   string
	Name string
}

type CVExperience struct {
	ID          string
	CompanyName string
	DateFrom    time.Time
	DateTill    time.Time
	Position    string
	Description string
}

type CVEducation struct {
	ID          string
	Institution string
	DateFrom    time.Time
	DateTill    time.Time
	Speciality  string
	Description string
}

type CVCustomSection struct {
	ID          string
	Description string
}

type CVCustomStory struct {
	ID          string
	ChapterName string
	MediaURL    string
}

type StoryEpisode struct {
	ID       string
	StoryID  string
	MediaURL string
}

type CV struct {
	ID                   string
	PersonaID            string
	Position             string
	WorkMonthsExperience int32
	MinSalary            int32
	MaxSalary            int32
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CVShort struct {
	ID                   string
	Position             string
	WorkMonthsExperience int32
	MinSalary            int32
	MaxSalary            int32
}

/**
Job types part start
*/

func (s *Storage) TxPutJobType(ctx context.Context, tx pkgtx.Tx, jobType *JobType) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE job_type SET
					name = $2,
					updated_at = $4
				WHERE id = $1
				RETURNING id, name, created_at, updated_at
			)
			INSERT INTO job_type (id, name, created_at, updated_at)
			SELECT $1, $2, $3, $4
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		jobType.ID,
		jobType.Name,
		jobType.CreatedAt,
		jobType.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetJobTypes(ctx context.Context, tx pkgtx.Tx) (_ []*JobType, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, name, created_at, updated_at
				FROM job_type
				ORDER BY name ASC`,
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

	jts := make([]*JobType, 0)

	for rows.Next() {
		var jt JobType
		if err := rows.Scan(&jt.ID, &jt.Name, &jt.CreatedAt, &jt.UpdatedAt); err != nil {
			return nil, errors.WithStack(err)
		}

		jts = append(jts, &jt)
	}
	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return jts, nil
}

func (s *Storage) TxDeleteJobType(ctx context.Context, tx pkgtx.Tx, jobTypeID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM job_type 
			WHERE id = $1`,
		jobTypeID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxPutCVJobTypes(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
	jobTypesIDs []string,
) error {
	c := postgresql.FromTx(tx)

	queryFormat := `INSERT 
		INTO cv_job_types (cv_id, job_type_id)
		VALUES %s`

	columns := 2
	valueStrings := make([]string, len(jobTypesIDs))
	valueArgs := make([]interface{}, len(jobTypesIDs)*columns)

	for i := 0; i < len(jobTypesIDs); i++ {
		offset := i * columns
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", offset+1, offset+2)
		valueArgs[offset] = cvID
		valueArgs[offset+1] = jobTypesIDs[i]
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

func (s *Storage) TxGetCVJobTypes(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVJobType, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT jt.id, jt.name
			FROM cv_job_types AS cjt
			INNER JOIN job_type AS jt
			ON cjt.job_type_id = jt.id
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvjts := make([]*CVJobType, 0)

	for rows.Next() {
		var cvjt CVJobType
		if err := rows.Scan(&cvjt.ID, &cvjt.Name); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cvjts = append(cvjts, &cvjt)
	}

	return cvjts, nil
}

func (s *Storage) TxDeleteCVJobTypes(ctx context.Context, tx pkgtx.Tx, cvID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM cv_job_types 
			WHERE cv_id = $1`,
		cvID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Job types part end
*/

/**
Job kinds part start
*/

func (s *Storage) TxPutJobKind(ctx context.Context, tx pkgtx.Tx, jobKind *JobKind) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE job_kind SET
					name = $2,
					updated_at = $4
				WHERE id = $1
				RETURNING id, name, created_at, updated_at
			)
			INSERT INTO job_kind (id, name, created_at, updated_at)
			SELECT $1, $2, $3, $4
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		jobKind.ID,
		jobKind.Name,
		jobKind.CreatedAt,
		jobKind.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetJobKinds(
	ctx context.Context,
	tx pkgtx.Tx,
) (_ []*JobKind, rerr error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, name, created_at, updated_at
				FROM job_kind
				ORDER BY name ASC`,
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

	jks := make([]*JobKind, 0)

	for rows.Next() {
		var jk JobKind
		if err := rows.Scan(&jk.ID, &jk.Name, &jk.CreatedAt, &jk.UpdatedAt); err != nil {
			return nil, errors.WithStack(err)
		}

		jks = append(jks, &jk)
	}
	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return jks, nil
}

func (s *Storage) TxDeleteJobKind(ctx context.Context, tx pkgtx.Tx, jobKindID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM job_kind 
			WHERE id = $1`,
		jobKindID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxPutCVJobKinds(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
	jobKindsIDs []string,
) error {
	c := postgresql.FromTx(tx)

	queryFormat := `INSERT 
		INTO cv_job_kinds (cv_id, job_kind_id)
		VALUES %s`

	columns := 2
	valueStrings := make([]string, len(jobKindsIDs))
	valueArgs := make([]interface{}, len(jobKindsIDs)*columns)

	for i := 0; i < len(jobKindsIDs); i++ {
		offset := i * columns
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", offset+1, offset+2)
		valueArgs[offset] = cvID
		valueArgs[offset+1] = jobKindsIDs[i]
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

func (s *Storage) TxGetCVJobKinds(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVJobKind, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT jk.id, jk.name
			FROM cv_job_kinds AS cjk
			INNER JOIN job_kind AS jk
			ON cjk.job_kind_id = jk.id
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvjks := make([]*CVJobKind, 0)

	for rows.Next() {
		var cvjk CVJobKind
		if err := rows.Scan(&cvjk.ID, &cvjk.Name); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cvjks = append(cvjks, &cvjk)
	}

	return cvjks, nil
}

func (s *Storage) TxDeleteCVJobKinds(ctx context.Context, tx pkgtx.Tx, cvID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM cv_job_kinds 
			WHERE cv_id = $1`,
		cvID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Job kinds part end
*/

/**
Experience part start
*/

func (s *Storage) TxPutExperience(ctx context.Context, tx pkgtx.Tx, cvID string, experience *CVExperience) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE experience SET
					cv_id = $2,
					company_name = $3,
					date_from = $4,
					date_till = $5,
					position = $6,
					description = $7
				WHERE id = $1
				RETURNING id, cv_id, company_name, date_from, date_till, position, description
			)
			INSERT INTO experience (id, cv_id, company_name, date_from, date_till, position, description)
			SELECT $1, $2, $3, $4, $5, $6, $7
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		experience.ID,
		cvID,
		experience.CompanyName,
		experience.DateFrom,
		experience.DateTill,
		experience.Position,
		experience.Description,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetExperiences(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVExperience, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, company_name, date_from, date_till, position, description
			FROM experience
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cves := make([]*CVExperience, 0)

	for rows.Next() {
		var cve CVExperience
		if err := rows.Scan(
			&cve.ID,
			&cve.CompanyName,
			&cve.DateFrom,
			&cve.DateTill,
			&cve.Position,
			&cve.Description,
		); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cves = append(cves, &cve)
	}

	return cves, nil
}

func (s *Storage) TxDeleteExperience(ctx context.Context, tx pkgtx.Tx, experienceID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM experience 
			WHERE id = $1`,
		experienceID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Experience part end
*/

/**
Education part start
*/
func (s *Storage) TxPutEducation(ctx context.Context, tx pkgtx.Tx, cvID string, education *CVEducation) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE education SET
					cv_id = $2,
					institution = $3,
					date_from = $4,
					date_till = $5,
					speciality = $6,
					description = $7
				WHERE id = $1
				RETURNING id, cv_id, institution, date_from, date_till, speciality, description
			)
			INSERT INTO education (id, cv_id, institution, date_from, date_till, speciality, description)
			SELECT $1, $2, $3, $4, $5, $6, $7
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		education.ID,
		cvID,
		education.Institution,
		education.DateFrom,
		education.DateTill,
		education.Speciality,
		education.Description,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetEducations(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVEducation, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, institution, date_from, date_till, speciality, description
			FROM education
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cves := make([]*CVEducation, 0)

	for rows.Next() {
		var cve CVEducation
		if err := rows.Scan(
			&cve.ID,
			&cve.Institution,
			&cve.DateFrom,
			&cve.DateTill,
			&cve.Speciality,
			&cve.Description,
		); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cves = append(cves, &cve)
	}

	return cves, nil
}

func (s *Storage) TxDeleteEducation(ctx context.Context, tx pkgtx.Tx, educationID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM education 
			WHERE id = $1`,
		educationID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Education part end
*/

/**
Custom sections part start
*/
func (s *Storage) TxPutCustomSection(ctx context.Context, tx pkgtx.Tx, cvID string, education *CVCustomSection) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE custom_section SET
					cv_id = $2,
					description = $3
				WHERE id = $1
				RETURNING id, cv_id, description
			)
			INSERT INTO custom_section (id, cv_id, description)
			SELECT $1, $2, $3
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		education.ID,
		cvID,
		education.Description,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetCustomSections(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVCustomSection, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, description
			FROM custom_section
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvcss := make([]*CVCustomSection, 0)

	for rows.Next() {
		var cvcs CVCustomSection
		if err := rows.Scan(
			&cvcs.ID,
			&cvcs.Description,
		); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cvcss = append(cvcss, &cvcs)
	}

	return cvcss, nil
}

func (s *Storage) TxDeleteCustomSection(ctx context.Context, tx pkgtx.Tx, sectionID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM custom_section 
			WHERE id = $1`,
		sectionID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Custom sections end
*/

/**
Story part start
*/
func (s *Storage) TxPutStory(ctx context.Context, tx pkgtx.Tx, cvID string, story *CVCustomStory) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE story SET
					cv_id = $2,
					chapter_name = $3
				WHERE id = $1
				RETURNING id, cv_id, chapter_name
			)
			INSERT INTO story (id, cv_id, chapter_name)
			SELECT $1, $2, $3
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		story.ID,
		cvID,
		story.ChapterName,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetStories(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*CVCustomStory, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, chapter_name
			FROM story
			WHERE cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvcss := make([]*CVCustomStory, 0)

	for rows.Next() {
		var cvcs CVCustomStory
		if err := rows.Scan(
			&cvcs.ID,
			&cvcs.ChapterName,
		); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		cvcss = append(cvcss, &cvcs)
	}

	return cvcss, nil
}

func (s *Storage) TxDeleteStory(ctx context.Context, tx pkgtx.Tx, storyID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM story 
			WHERE id = $1`,
		storyID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Story end
*/

/**
Stories episodes start
*/
func (s *Storage) TxPutStoriesEpisode(ctx context.Context, tx pkgtx.Tx, story *StoryEpisode) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE story_episodes SET
					stories_id = $2,
					media_url = $3
				WHERE id = $1
				RETURNING id, stories_id, media_url
			)
			INSERT INTO story_episodes (id, stories_id, media_url)
			SELECT $1, $2, $3
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		story.ID,
		story.StoryID,
		story.MediaURL,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetStoriesEpisodesByID(
	ctx context.Context,
	tx pkgtx.Tx,
	storyEpisodeID string,
) ([]*StoryEpisode, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT se.stories_id, se.id, se.media_url
			FROM story_episodes AS se
			WHERE se.id = $1`,
		storyEpisodeID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ses := make([]*StoryEpisode, 0)

	for rows.Next() {
		var se StoryEpisode
		if err := rows.Scan(&se.StoryID, &se.ID, &se.MediaURL); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		ses = append(ses, &se)
	}

	return ses, nil
}

func (s *Storage) TxGetStoriesEpisodes(
	ctx context.Context,
	tx pkgtx.Tx,
	cvID string,
) ([]*StoryEpisode, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT se.stories_id, se.id, se.media_url
			FROM story_episodes AS se
			INNER JOIN story AS s
			ON se.stories_id = s.id
			WHERE s.cv_id = $1`,
		cvID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ses := make([]*StoryEpisode, 0)

	for rows.Next() {
		var se StoryEpisode
		if err := rows.Scan(&se.StoryID, &se.ID, &se.MediaURL); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		ses = append(ses, &se)
	}

	return ses, nil
}

func (s *Storage) TxDeleteStoriesEpisode(ctx context.Context, tx pkgtx.Tx, episodeID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM story_episodes 
			WHERE id = $1`,
		episodeID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/**
Stories episodes end
*/

/*
CV start
*/
func (s *Storage) TxPutCV(ctx context.Context, tx pkgtx.Tx, cv *CV) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`WITH upsert AS (
				UPDATE cv SET
					persona_id = $2,
					position = $3,
					work_months_experience = $4,
					min_salary = $5,
					max_salary = $6,
					updated_at = $8
				WHERE id = $1
				RETURNING id, persona_id, position, work_months_experience, min_salary, max_salary, created_at, updated_at
			)
			INSERT INTO cv (id, persona_id, position, work_months_experience, min_salary, max_salary, created_at, updated_at)
			SELECT $1, $2, $3, $4, $5, $6, $7, $8
			WHERE NOT EXISTS (SELECT * FROM upsert)`,
		cv.ID,
		cv.PersonaID,
		cv.Position,
		cv.WorkMonthsExperience,
		cv.MinSalary,
		cv.MaxSalary,
		cv.CreatedAt,
		cv.UpdatedAt,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Storage) TxGetCV(ctx context.Context, tx pkgtx.Tx, cvID string) (*CV, error) {
	c := postgresql.FromTx(tx)

	var cv CV
	err := c.QueryRowContext(
		ctx,
		`SELECT id, persona_id, position, work_months_experience, min_salary, max_salary, created_at, updated_at
				FROM cv
				WHERE id = $1`,
		cvID,
	).Scan(
		&cv.ID,
		&cv.PersonaID,
		&cv.Position,
		&cv.WorkMonthsExperience,
		&cv.MinSalary,
		&cv.MaxSalary,
		&cv.CreatedAt,
		&cv.UpdatedAt,
	)

	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, errors.WithStack(ErrNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &cv, nil
}

func (s *Storage) TxGetCVs(
	ctx context.Context,
	tx pkgtx.Tx,
	personaID string,
) ([]*CVShort, error) {
	c := postgresql.FromTx(tx)

	rows, err := c.QueryContext(
		ctx,
		`SELECT id, position, work_months_experience, min_salary, max_salary, 
			FROM cv
			WHERE cv.persona_id = $1`,
		personaID,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	scvss := make([]*CVShort, 0)

	for rows.Next() {
		var cv CVShort
		if err := rows.Scan(&cv.ID, &cv.Position, &cv.WorkMonthsExperience, &cv.MinSalary, &cv.MaxSalary); err != nil {
			_ = rows.Close()
			return nil, errors.WithStack(err)
		}

		scvss = append(scvss, &cv)
	}

	return scvss, nil
}

func (s *Storage) TxDeleteCV(ctx context.Context, tx pkgtx.Tx, cvID string) error {
	c := postgresql.FromTx(tx)

	if _, err := c.ExecContext(
		ctx,
		`DELETE FROM cv 
			WHERE id = $1`,
		cvID,
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

/*
CV end
*/

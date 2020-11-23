package controller

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
	"personaapp/internal/controllers/company/storage"
	pkgtx "personaapp/pkg/tx"
	"time"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint:godox // TODO: Implement CDN link check
		return true
	})
}

var (
	ErrCompanyNotFound          = errors.New("company not found")
	ErrCategoryNotFound         = errors.New("category not found")
	ErrActivityFieldNotFound    = errors.New("activity field not found")
	ErrCityNotFound             = errors.New("city not found")
	ErrInvalidTitle             = errors.New("invalid title")
	ErrInvalidTitleLength       = errors.New("invalid title length")
	ErrInvalidDescription       = errors.New("invalid description")
	ErrInvalidDescriptionLength = errors.New("invalid description length")
	ErrInvalidLogoURL           = errors.New("invalid logo_url")
	ErrInvalidLogoURLLength     = errors.New("invalid logo_url length")
	ErrInvalidLogoURLFormat     = errors.New("invalid logo_url format")
)

type Storage interface {
	TxGetCompanyByID(ctx context.Context, tx pkgtx.Tx, authID string) (*storage.CompanyData, error)
	TxGetCompaniesByID(ctx context.Context, tx pkgtx.Tx, companyIDs []string) ([]*storage.CompanyData, error)
	TxPutCompany(ctx context.Context, tx pkgtx.Tx, cs *storage.CompanyData) error
	TxPutActivityField(ctx context.Context, tx pkgtx.Tx, af *storage.ActivityField) error
	TxGetActivityFields(
		ctx context.Context,
		tx pkgtx.Tx,
	) (_ []*storage.ActivityField, rerr error)
	TxGetActivityFieldsByID(
		ctx context.Context,
		tx pkgtx.Tx,
		activityFieldID string,
	) (_ *storage.ActivityField, rerr error)
	TxGetActivityFieldsByCompanyID(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
	) ([]*storage.ActivityField, error)
	TxPutCompanyActivityFields(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
		activityFieldsIDs []string,
	) error
	TxDeleteCompanyActivityFieldsByCompanyID(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
	) error

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

type CompanyData struct {
	ID          string
	Title       *string `valid:"stringlength(0|100)"`
	Description *string `valid:"stringlength(0|255)"`
	LogoURL     *string `valid:"stringlength(0|255),media_link"`
}

type Company struct {
	ID             string
	ActivityFields []string
	Title          string
	Description    string
	LogoURL        string
}

type ActivityField struct {
	ID      string
	Title   string
	IconURL string
}

func (cd *CompanyData) validate() error {
	var fieldErrors = []struct {
		Field        string
		Errors       map[string]error
		DefaultError error
	}{
		{
			Field: "Title",
			Errors: map[string]error{
				"stringlength": ErrInvalidTitleLength,
			},
			DefaultError: ErrInvalidTitle,
		},
		{
			Field: "Description",
			Errors: map[string]error{
				"stringlength": ErrInvalidDescriptionLength,
			},
			DefaultError: ErrInvalidDescription,
		},
		{
			Field: "LogoURL",
			Errors: map[string]error{
				"stringlength": ErrInvalidLogoURLLength,
				"media_link":   ErrInvalidLogoURLFormat,
			},
			DefaultError: ErrInvalidLogoURL,
		},
	}

	if valid, err := govalidator.ValidateStruct(cd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				validatorError, ok := err.(govalidator.Error)
				if !ok {
					return errors.Wrap(fe.DefaultError, msg)
				}

				if err, ok := fe.Errors[validatorError.Validator]; ok {
					return errors.Wrap(err, msg)
				}

				return errors.Wrap(fe.DefaultError, msg)
			}
		}

		return errors.New("company data struct is filled with some invalid data")
	}

	return nil
}

func (c *Controller) Update(ctx context.Context, cd *CompanyData) error {
	if err := cd.validate(); err != nil {
		return errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		var scd *storage.CompanyData
		now := time.Now()

		switch company, err := c.s.TxGetCompanyByID(ctx, tx, cd.ID); err {
		case nil:
			scd = &storage.CompanyData{
				ID:          company.ID,
				Title:       company.Title,
				Description: company.Description,
				LogoURL:     company.LogoURL,
				CreatedAt:   company.CreatedAt,
				UpdatedAt:   now,
			}
		case storage.ErrNotFound:
			scd = &storage.CompanyData{
				ID:          cd.ID,
				Title:       "",
				Description: "",
				LogoURL:     "",
				CreatedAt:   now,
				UpdatedAt:   now,
			}
		default:
			return errors.WithStack(err)
		}

		if cd.Title != nil {
			scd.Title = *cd.Title
		}
		if cd.Description != nil {
			scd.Description = *cd.Description
		}
		if cd.LogoURL != nil {
			scd.LogoURL = *cd.LogoURL
		}

		if err := c.s.TxPutCompany(ctx, tx, scd); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) UpdateActivityFields(ctx context.Context, companyID string, activityFields []string) error {
	if activityFields == nil {
		return nil
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if err := c.s.TxDeleteCompanyActivityFieldsByCompanyID(ctx, tx, companyID); err != nil {
			return errors.WithStack(err)
		}

		if len(activityFields) > 0 {
			if err := c.s.TxPutCompanyActivityFields(ctx, tx, companyID, activityFields); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) Get(ctx context.Context, companyID string) (*Company, error) {
	var company *Company

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		cd, err := c.s.TxGetCompanyByID(ctx, tx, companyID)

		switch errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCompanyNotFound)
		default:
			return errors.WithStack(err)
		}

		afs, err := c.s.TxGetActivityFieldsByCompanyID(ctx, tx, companyID)

		if err != nil {
			return errors.WithStack(err)
		}

		activityFieldsIDs := make([]string, len(afs))
		for i := 0; i < len(afs); i++ {
			activityFieldsIDs[i] = afs[i].Title
		}

		company = &Company{
			ID:             cd.ID,
			ActivityFields: activityFieldsIDs,
			Title:          cd.Title,
			Description:    cd.Description,
			LogoURL:        cd.LogoURL,
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return company, nil
}

func (c *Controller) GetCompaniesList(ctx context.Context, companyIDs []string) ([]*Company, error) {
	cds, err := c.s.TxGetCompaniesByID(ctx, c.s.NoTx(), companyIDs)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	companies := make([]*Company, 0, len(cds))
	for _, company := range cds {
		companies = append(companies, &Company{
			ID:             company.ID,
			ActivityFields: nil,
			Title:          company.Title,
			Description:    company.Description,
			LogoURL:        company.LogoURL,
		})
	}

	return companies, nil
}

func (c *Controller) UpdateActivityField(ctx context.Context, activityFieldID *string, cd *ActivityField) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		var ID string

		if activityFieldID != nil {
			ID = *activityFieldID
		} else {
			ID = uuid.NewV4().String()
		}

		now := time.Now()

		svc := &storage.ActivityField{
			ID:        ID,
			Title:     cd.Title,
			IconURL:   cd.IconURL,
			CreatedAt: now,
			UpdatedAt: now,
		}

		return errors.WithStack(c.s.TxPutActivityField(ctx, tx, svc))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) GetActivityField(ctx context.Context, activityFieldID string) (*ActivityField, error) {
	var activityField *ActivityField

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		af, err := c.s.TxGetActivityFieldsByID(ctx, tx, activityFieldID)

		switch errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrActivityFieldNotFound)
		default:
			return errors.WithStack(err)
		}

		activityField = &ActivityField{
			ID:      af.ID,
			Title:   af.Title,
			IconURL: af.IconURL,
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return activityField, nil
}

func (c *Controller) GetActivityFields(ctx context.Context) ([]*ActivityField, error) {
	afs, err := c.s.TxGetActivityFields(ctx, c.s.NoTx())

	if err != nil {
		return nil, errors.WithStack(err)
	}

	activityFields := make([]*ActivityField, 0, len(afs))
	for _, activityField := range afs {
		activityFields = append(activityFields, &ActivityField{
			ID:      activityField.ID,
			Title:   activityField.Title,
			IconURL: activityField.IconURL,
		})
	}

	return activityFields, nil
}

func (c *Controller) DeleteCompanyActivityFieldsByCompanyID(ctx context.Context, authID string) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		return errors.WithStack(c.s.TxDeleteCompanyActivityFieldsByCompanyID(ctx, tx, authID))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

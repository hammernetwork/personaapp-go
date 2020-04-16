package controller

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
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

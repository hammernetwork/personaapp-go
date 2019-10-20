package controller

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	companyStorage "personaapp/internal/server/company/storage"
	pkgtx "personaapp/pkg/tx"
	"time"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint TODO: Implement CDN link check
		return true
	})
}

var (
	ErrCompanyNotFound    = errors.New("company not found")
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidLogoURL     = errors.New("invalid logo_url")
)

type Storage interface {
	TxGetCompanyByID(ctx context.Context, tx pkgtx.Tx, authID string) (*companyStorage.CompanyData, error)
	TxPutCompany(ctx context.Context, tx pkgtx.Tx, cs *companyStorage.CompanyData) error
	TxGetActivityFieldsByCompanyID(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
	) ([]*companyStorage.ActivityField, error)
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
	AuthID         string
	ActivityFields []string
	Title          *string `valid:"stringlength(0|100)"`
	Description    *string `valid:"stringlength(0|255)"`
	LogoURL        *string `valid:"stringlength(0|255),media_link"`
}

type Company struct {
	AuthID         string
	ActivityFields []string
	Title          string
	Description    string
	LogoURL        string
}

func (cd *CompanyData) validate() error {
	var fieldErrors = []struct {
		Field string
		Error error
	}{
		{Field: "Title", Error: ErrInvalidTitle},
		{Field: "Description", Error: ErrInvalidDescription},
		{Field: "LogoURL", Error: ErrInvalidLogoURL},
	}

	if valid, err := govalidator.ValidateStruct(cd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.Wrap(fe.Error, msg)
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
		var scd *companyStorage.CompanyData
		now := time.Now()

		switch company, err := c.s.TxGetCompanyByID(ctx, tx, cd.AuthID); err {
		case nil:
			scd = &companyStorage.CompanyData{
				AuthID:      company.AuthID,
				Title:       company.Title,
				Description: company.Description,
				LogoURL:     company.LogoURL,
				CreatedAt:   company.CreatedAt,
				UpdatedAt:   now,
			}
		case companyStorage.ErrNotFound:
			scd = &companyStorage.CompanyData{
				AuthID:      cd.AuthID,
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

		if err := c.s.TxDeleteCompanyActivityFieldsByCompanyID(ctx, tx, cd.AuthID); err != nil {
			return errors.WithStack(err)
		}

		if err := c.s.TxPutCompanyActivityFields(ctx, tx, cd.AuthID, cd.ActivityFields); err != nil {
			return errors.WithStack(err)
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

		switch err {
		case nil:
		case companyStorage.ErrNotFound:
			return ErrCompanyNotFound
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
			AuthID:         cd.AuthID,
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

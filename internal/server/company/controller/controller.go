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
		//TODO: Implement CDN link check
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
	TxGetCompanyActivityFieldsByID(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
	) ([]*companyStorage.ActivityField, error)
	TxPutCompanyActivityFields(
		ctx context.Context,
		tx pkgtx.Tx,
		authID string,
		activityFields []*companyStorage.ActivityField,
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
		company, err := c.s.TxGetCompanyByID(ctx, tx, cd.AuthID)
		switch err {
		case nil:
		case companyStorage.ErrNotFound:
			return ErrCompanyNotFound
		default:
			return errors.WithStack(err)
		}

		now := time.Now()

		scd := &companyStorage.CompanyData{
			AuthID:      company.AuthID,
			Title:       company.Title,
			Description: company.Description,
			LogoURL:     company.LogoURL,
			CreatedAt:   now,
			UpdatedAt:   now,
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

		//activityFields, err := c.s.TxGetCompanyActivityFieldsByID(ctx, tx, cd.AuthID)
		//if err != nil {
		//	return errors.WithStack(err)
		//}
		//TODO: get company activity fields
		//TODO: update company activity fields
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

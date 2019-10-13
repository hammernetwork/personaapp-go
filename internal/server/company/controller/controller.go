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
	ErrInvalidTitle       = errors.New("invalid title")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidLogoURL     = errors.New("invalid logo_url")
)

type Storage interface {
	TxPutCompany(ctx context.Context, tx pkgtx.Tx, cs *companyStorage.CompanyData) error

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
	AuthID      string
	ScopeID     *string
	Title       *string `valid:"stringlength(0|100)"`
	Description *string `valid:"stringlength(0|255)"`
	LogoURL     *string `valid:"stringlength(0|255),media_link"`
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

	now := time.Now()
	scd := companyStorage.CompanyData{
		Fields:      0,
		AuthID:      cd.AuthID,
		ScopeID:     "",
		Title:       "",
		Description: "",
		LogoURL:     "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if cd.ScopeID != nil {
		scd.ScopeID = *cd.ScopeID
		scd.Fields |= companyStorage.FieldScopeID
	}

	if cd.Title != nil {
		scd.Title = *cd.Title
		scd.Fields |= companyStorage.FieldTitle
	}

	if cd.Description != nil {
		scd.Description = *cd.Description
		scd.Fields |= companyStorage.FieldDescription
	}

	if cd.LogoURL != nil {
		scd.LogoURL = *cd.LogoURL
		scd.Fields |= companyStorage.FieldLogoURL
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		return errors.WithStack(c.s.TxPutCompany(ctx, tx, &scd))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

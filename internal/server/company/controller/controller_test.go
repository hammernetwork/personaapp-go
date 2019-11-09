package controller_test

import (
	"context"
	"testing"
	"time"

	authController "personaapp/internal/server/auth/controller"
	companyController "personaapp/internal/server/company/controller"
	"personaapp/internal/server/company/storage"
	"personaapp/internal/testutils"

	sqlMigrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
)

var authCfg = &authController.Config{
	TokenExpiration:   5 * time.Minute,
	PrivateSigningKey: "signkey",
	TokenValidityGap:  15 * time.Second,
}

func InitStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Up))

	return storage.New(pg), pg.Close
}

func TestGetNonExistingCompany(t *testing.T) {
	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	t.Run("normal flow", func(t *testing.T) {
		company, err := cc.Get(context.Background(), "nonexistingcompany")
		require.Nil(t, company)
		require.Error(t, companyController.ErrCompanyNotFound, err)
	})
}

func TestGetExistingButNotCompletedCompany(t *testing.T) {
	as, authCloser := testutils.InitAuthStorage(t)
	defer func() {
		if err := authCloser(); err != nil {
			t.Error(err)
		}
	}()

	ac := authController.New(authCfg, as)

	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	rd := authController.RegisterData{
		Email:    "companytest@gmail.com",
		Phone:    "+380011234567",
		Account:  authController.AccountTypeCompany,
		Password: "Password",
	}

	token, err := ac.Register(context.Background(), &rd)
	if err != nil {
		t.Error(err)
	}

	t.Run("normal flow", func(t *testing.T) {
		company, err := cc.Get(context.Background(), token.AccountID)
		require.Error(t, companyController.ErrCompanyNotFound, err)
		require.Nil(t, company)
	})
}

func TestUpdateExistingCompany(t *testing.T) {
	as, authCloser := testutils.InitAuthStorage(t)
	defer func() {
		if err := authCloser(); err != nil {
			t.Error(err)
		}
	}()

	ac := authController.New(authCfg, as)

	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	rd := authController.RegisterData{
		Email:    "companytest2@gmail.com",
		Phone:    "+380019988776",
		Account:  authController.AccountTypeCompany,
		Password: "Password2",
	}

	token, err := ac.Register(context.Background(), &rd)
	if err != nil {
		t.Error(err)
	}

	title := "Title"
	description := "Description"
	logoURL := "https://logourl.com"

	cd := companyController.CompanyData{
		AuthID:      token.AccountID,
		Title:       &title,
		Description: &description,
		LogoURL:     &logoURL,
	}

	t.Run("update all fields", func(t *testing.T) {
		require.NoError(t, cc.Update(context.Background(), &cd))

		company, err := cc.Get(context.Background(), token.AccountID)
		require.NoError(t, err)
		require.Equal(t, title, company.Title)
		require.Equal(t, description, company.Description)
		require.Equal(t, logoURL, company.LogoURL)
	})
}

func TestUpdateNonExistingCompany(t *testing.T) {
	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	cd := companyController.CompanyData{
		AuthID:      "nonexistingid",
		Title:       nil,
		Description: nil,
		LogoURL:     nil,
	}

	t.Run("normal flow", func(t *testing.T) {
		require.Error(t, companyController.ErrCompanyNotFound, cc.Update(context.Background(), &cd))
	})
}

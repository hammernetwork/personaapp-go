package controller_test

import (
	"testing"
	"context"

	"personaapp/internal/testutils"
	companyController "personaapp/internal/server/company/controller"
	authController "personaapp/internal/server/auth/controller"
	"personaapp/internal/server/company/storage"

	"github.com/stretchr/testify/require"
	sqlMigrate "github.com/rubenv/sql-migrate"
)

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

	ac := authController.New(nil, as)

	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	rd := authController.RegisterData{
		Email:    "company_test@gmail.com",
		Phone:    "+0112345678",
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

func TestUpdateNonExistingCompany(t *testing.T) {
	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	cd := companyController.CompanyData{
		AuthID:         "nonexistingid",
		ActivityFields: nil,
		Title:          nil,
		Description:    nil,
		LogoURL:        nil,
	}

	t.Run("normal flow", func(t *testing.T) {
		require.Error(t, companyController.ErrCompanyNotFound, cc.Update(context.Background(), &cd))
	})
}

func TestUpdateExistingCompany(t *testing.T) {
	as, authCloser := testutils.InitAuthStorage(t)
	defer func() {
		if err := authCloser(); err != nil {
			t.Error(err)
		}
	}()

	ac := authController.New(nil, as)

	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	rd := authController.RegisterData{
		Email:    "company_test@gmail.com",
		Phone:    "+0112345678",
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

	cd := companyController.CompanyData{
		AuthID:         token.AccountID,
		ActivityFields: nil,
		Title:          nil,
		Description:    nil,
		LogoURL:        nil,
	}

	//t.Run("after update", func(t *testing.T) {
	//	err :cc.Update(context.Background(), &cd)
	//	company, err := cc.Get(context.Background(), token.AccountID)
	//	require.NoError(t, err)
	//	require.Equal(t, company.Title)
	//})
}
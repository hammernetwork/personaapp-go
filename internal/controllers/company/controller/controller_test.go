package controller_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"personaapp/internal/controllers/company/storage"
	"testing"
	"time"

	authController "personaapp/internal/controllers/auth/controller"
	companyController "personaapp/internal/controllers/company/controller"
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
		Phone:    "+380500000001",
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
		Email:    "companytest1@gmail.com",
		Phone:    "+380500000011",
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
		ID:          token.AccountID,
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
		ID:          "nonexistingid",
		Title:       nil,
		Description: nil,
		LogoURL:     nil,
	}

	t.Run("normal flow", func(t *testing.T) {
		require.Error(t, companyController.ErrCompanyNotFound, cc.Update(context.Background(), &cd))
	})
}

func TestUpdateCompanyActivityFields(t *testing.T) {
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
		Phone:    "+380500000002",
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
		ID:          token.AccountID,
		Title:       &title,
		Description: &description,
		LogoURL:     &logoURL,
	}

	activityID := uuid.NewV4().String()
	af := companyController.ActivityField{
		Title:   "Title2",
		IconURL: "https://logourl2.com",
	}

	activityID2 := uuid.NewV4().String()
	af2 := companyController.ActivityField{
		Title:   "Title3",
		IconURL: "https://logourl3.com",
	}

	t.Run("update activity fields", func(t *testing.T) {
		require.NoError(t, cc.UpdateActivityField(context.Background(), &activityID, &af))
		require.NoError(t, cc.UpdateActivityField(context.Background(), &activityID2, &af2))

		fields, err := cc.GetActivityFields(context.Background())
		require.NoError(t, err)
		require.Equal(t, af.Title, fields[0].Title)
		require.Equal(t, af.IconURL, fields[0].IconURL)
		require.Equal(t, af2.Title, fields[1].Title)
		require.Equal(t, af2.IconURL, fields[1].IconURL)
	})

	activityFields := []string{activityID, activityID2}

	t.Run("relate company to activity fields", func(t *testing.T) {
		require.NoError(t, cc.Update(context.Background(), &cd))
		require.NoError(t, cc.UpdateActivityFields(context.Background(), token.AccountID, activityFields))

		company, err := cc.Get(context.Background(), token.AccountID)
		require.NoError(t, err)
		require.Equal(t, af.Title, company.ActivityFields[0])
		require.Equal(t, af2.Title, company.ActivityFields[1])

		require.NoError(t, cc.DeleteCompanyActivityFieldsByCompanyID(context.Background(), token.AccountID))
	})
}

func TestGetAndUpdateActivityFields(t *testing.T) {
	_, authCloser := testutils.InitAuthStorage(t)
	defer func() {
		if err := authCloser(); err != nil {
			t.Error(err)
		}
	}()

	cs, companyCloser := InitStorage(t)
	defer func() {
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
	}()

	cc := companyController.New(cs)

	title := "Title"
	iconURL := "https://logourl.com"

	activityID := uuid.NewV4().String()
	af := companyController.ActivityField{
		Title:   title,
		IconURL: iconURL,
	}

	t.Run("insert activity field", func(t *testing.T) {
		require.NoError(t, cc.UpdateActivityField(context.Background(), &activityID, &af))

		activityField, err := cc.GetActivityField(context.Background(), activityID)
		require.NoError(t, err)
		require.Equal(t, title, activityField.Title)
		require.Equal(t, iconURL, activityField.IconURL)
	})

	title = "Title1"
	iconURL = "https://logourl1.com"

	af = companyController.ActivityField{
		Title:   title,
		IconURL: iconURL,
	}

	t.Run("update activity field", func(t *testing.T) {
		require.NoError(t, cc.UpdateActivityField(context.Background(), &activityID, &af))

		activityField, err := cc.GetActivityField(context.Background(), activityID)
		require.NoError(t, err)
		require.Equal(t, title, activityField.Title)
		require.Equal(t, iconURL, activityField.IconURL)
	})

}

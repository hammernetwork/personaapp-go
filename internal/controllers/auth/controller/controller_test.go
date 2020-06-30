package controller_test

import (
	"context"
	sqlMigrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
	"personaapp/internal/controllers/auth/controller"
	"personaapp/internal/controllers/auth/storage"
	"personaapp/internal/testutils"
	"testing"
	"time"
)

var authCfg = &controller.Config{
	TokenExpiration:   5 * time.Minute,
	PrivateSigningKey: "signkey",
	TokenValidityGap:  15 * time.Second,
}

func InitStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Up))

	return storage.New(pg), pg.Close
}

func TestRegister(t *testing.T) {
	s, closer := InitStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
	}()

	c := controller.New(authCfg, s)

	t.Run("two accounts with empty email", func(t *testing.T) {
		_, err := c.Register(context.TODO(), &controller.RegisterData{
			Email:    "",
			Phone:    "+380500000101",
			Account:  controller.AccountTypeCompany,
			Password: "random_password",
		})
		require.Nil(t, err)

		_, err = c.Register(context.TODO(), &controller.RegisterData{
			Email:    "",
			Phone:    "+380500000102",
			Account:  controller.AccountTypeCompany,
			Password: "random_password",
		})
		require.Nil(t, err)
	})
}

func TestUpdateAuthEmailAndPhone(t *testing.T) {
	as, authCloser := testutils.InitAuthStorage(t)
	defer func() {
		if err := authCloser(); err != nil {
			t.Error(err)
		}
	}()

	ac := controller.New(authCfg, as)

	rd := controller.RegisterData{
		Email:    "companytest1@gmail.com",
		Phone:    "+380500000002",
		Account:  controller.AccountTypePersona,
		Password: "Password2",
	}

	token, err := ac.Register(context.Background(), &rd)
	if err != nil {
		t.Error(err)
	}

	newEmail := "companytest2@gmail.com"

	t.Run("update email", func(t *testing.T) {
		_, err := ac.UpdateEmail(context.Background(), token.AccountID, newEmail, rd.Password, rd.Account)
		require.NoError(t, err)

		self, err := ac.GetSelf(context.Background(), token.AccountID)
		require.NoError(t, err)
		require.Equal(t, newEmail, self.Email)
	})

	newPhone := "+380500000004"

	t.Run("update phone", func(t *testing.T) {
		_, err := ac.UpdatePhone(context.Background(), token.AccountID, newPhone, rd.Password)
		require.NoError(t, err)

		self, err := ac.GetSelf(context.Background(), token.AccountID)
		require.NoError(t, err)
		require.Equal(t, newPhone, self.Phone)
	})

	pd := controller.UpdatePasswordData{
		OldPassword: rd.Password,
		NewPassword: "Password3",
	}

	t.Run("update password", func(t *testing.T) {
		_, err := ac.UpdatePassword(context.Background(), token.AccountID, &pd)
		require.NoError(t, err)
	})
}

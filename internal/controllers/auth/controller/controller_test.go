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

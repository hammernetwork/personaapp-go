package register_test

import (
	"context"
	"database/sql"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/cockroachdb/errors"
	cmdmigrate "personaapp/cmd/migrate"
	"github.com/stretchr/testify/require"
	controller "personaapp/internal/server/controller/register"
	storage "personaapp/internal/server/storage/register"
	"personaapp/pkg/dockertest"
	"personaapp/pkg/postgresql"
	"testing"
)

func initStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
	cfg, err := dockertest.EnsurePostgres()
	require.NoError(t, err)

	pg, err := postgresql.New(&postgresql.Config{
		Host:               "localhost",
		Port:               uint16(cfg.Port),
		User:               cfg.User,
		Password:           cfg.Password,
		Database:           cfg.Database,
		MaxOpenConnections: 16,
		MaxIdleConnections: 16,
	})
	require.NoError(t, err)

	require.NoError(t, migrateUp(pg.DB))

	return storage.New(pg), pg.Close
}

func migrateUp(db *sql.DB) error {
	source := migrate.MemoryMigrationSource{
		Migrations: cmdmigrate.GetMigrationsForTests(),
	}
	if _, err := migrate.Exec(db, "postgres", source, migrate.Up); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func TestRegisterCompanyNormalFlow(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
	}()

	c := controller.New(s)

	t.Run("normal flow", func(t *testing.T) {
		err := c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "normal company",
			Email:    "normal_company@gmail.com",
			Phone:    "+380991234567",
			Password: "SuperPassword",
		})

		require.Nil(t, err)
	})

	t.Run("normal flow with dirty company name", func(t *testing.T) {

	})

	t.Run("normal flow with dirty email", func(t *testing.T) {

	})

	t.Run("normal flow with dirty phone", func(t *testing.T) {

	})
}


func TestRegisterCompanyAlreadyExisting(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
	}()

	c := controller.New(s)

	t.Run("already existing", func(t *testing.T) {
		existingEmail := "some_company@gmail.com"
		nonExistingEmail := "some_other_company@gmail.com"
		existingPhone := "+380990123456"
		nonExistingPhone := "+380996543210"

		require.Nil(t, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "some company",
			Email:    existingEmail,
			Phone:    existingPhone,
			Password: "SuperPassword",
		}))

		require.Error(t, controller.ErrAlreadyExists, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "other company",
			Email:    existingEmail,
			Phone:   nonExistingPhone,
			Password: "OtherSuperPassword",
		}))

		require.Error(t, controller.ErrAlreadyExists, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "some other company",
			Email:    nonExistingEmail,
			Phone:    existingPhone,
			Password: "SomeOtherPassword",
		}))
	})
}

func TestRegisterCompanyInvalidArgument(t *testing.T){
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
	}()

	c := controller.New(s)

	t.Run("short company name", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyNameInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "c",
			Email:    "short_company@gmail.com",
			Phone:    "+380992345678",
			Password: "ShortPassword",
		}))
	})

	t.Run("long company name", func(t *testing.T) {
		name := "Very long company name that should not pass validations and fail this test so that Max will be happy."
		require.Error(t, controller.ErrCompanyNameInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     name,
			Email:    "long_company@gmail.com",
			Phone:    "+380992345678",
			Password: "LongPassword",
		}))
	})

	t.Run("invalid email format", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyEmailInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid email company",
			Email:    "plainemail",
			Phone:    "+38000112233",
			Password: "InvalidEmail",
		}))

		require.Error(t, controller.ErrCompanyEmailInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid email company",
			Email:    "#@%^%#$@#$@#.com",
			Phone:    "+38000112233",
			Password: "InvalidEmail",
		}))

		require.Error(t, controller.ErrCompanyEmailInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid email company",
			Email:    "@domain.com",
			Phone:    "+38000112233",
			Password: "InvalidEmail",
		}))

		require.Error(t, controller.ErrCompanyEmailInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid email company",
			Email:    "email.domain.com",
			Phone:    "+38000112233",
			Password: "InvalidEmail",
		}))
	})

	t.Run("short phone", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "short phone company",
			Email:    "short_phone@gmail.com",
			Phone:    "1234",
			Password: "InvalidPhone",
		}))

		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "short phone company",
			Email:    "short_phone@gmail.com",
			Phone:    "+3805",
			Password: "InvalidPhone",
		}))
	})

	t.Run("long phone", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "long phone company",
			Email:    "long_phone@gmail.com",
			Phone:    "+380501234567890123451",
			Password: "InvalidPhone",
		}))
	})

	t.Run("invalid phone format", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid phone company",
			Email:    "invalid_phone@gmail.com",
			Phone:    "phone",
			Password: "InvalidPhone",
		}))

		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid phone company",
			Email:    "invalid_phone@gmail.com",
			Phone:    "+38(099)4888377",
			Password: "InvalidPhone",
		}))

		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "invalid phone company",
			Email:    "invalid_phone@gmail.com",
			Phone:    "099-488-83-77",
			Password: "InvalidPhone",
		}))
	})

	t.Run("short password", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "short password company",
			Email:    "invalid_password@gmail.com",
			Phone:    "099-488-83-77",
			Password: "abcde",
		}))
	})

	t.Run("long password", func(t *testing.T) {
		require.Error(t, controller.ErrCompanyPhoneInvalid, c.RegisterCompany(context.Background(), &controller.Company{
			Name:     "short password company",
			Email:    "invalid_password@gmail.com",
			Phone:    "099-488-83-77",
			Password: "abcdeABCDEabcdeABCDEabcdeABCDEz",
		}))
	})
}
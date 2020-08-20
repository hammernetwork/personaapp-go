package controller_test

import (
	"context"
	"fmt"
	sqlMigrate "github.com/rubenv/sql-migrate"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	authController "personaapp/internal/controllers/auth/controller"
	"personaapp/internal/controllers/city/controller"
	"personaapp/internal/controllers/city/storage"
	companyStorage "personaapp/internal/controllers/company/storage"
	"personaapp/internal/testutils"
	"testing"
	"time"
)

var authCfg = &authController.Config{
	TokenExpiration:   5 * time.Minute,
	PrivateSigningKey: "signkey",
	TokenValidityGap:  15 * time.Second,
}

func initStorage(t *testing.T) (_ *storage.Storage, closer func() error) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Up))

	return storage.New(pg), pg.Close
}

func initCompanyStorage(t *testing.T) (_ *companyStorage.Storage, closer func() error) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Up))

	return companyStorage.New(pg), pg.Close
}

func cleanup(t *testing.T) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Down))
}

func TestController_PutCity(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)

	t.Run("create new city", func(t *testing.T) {
		cityToCreate := controller.City{
			Name:        "Hamburg",
			CountryCode: 0,
			Rating:      0,
		}
		ID, err := c.PutCity(context.TODO(), nil, &cityToCreate)
		require.NoError(t, err)
		require.NotNil(t, ID)

		createdCities, err := c.GetCities(context.TODO(), []int32{0}, 0, cityToCreate.Name)
		require.NoError(t, err)
		require.NotNil(t, createdCities)
		require.Equal(t, cityToCreate.Name, createdCities[0].Name)
		require.Equal(t, cityToCreate.CountryCode, createdCities[0].CountryCode)
		require.Equal(t, cityToCreate.Rating, createdCities[0].Rating)
	})

	t.Run("update city", func(t *testing.T) {
		cityToCreate := controller.City{
			Name:        "Dnipro",
			CountryCode: 1,
			Rating:      1,
		}
		ID, err := c.PutCity(context.TODO(), nil, &cityToCreate)

		require.NoError(t, err)
		require.NotNil(t, ID)

		cityToUpdate := controller.City{
			Name:        "Lviv",
			CountryCode: 0,
			Rating:      0,
		}

		stringID := string(ID)
		updatedID, err := c.PutCity(context.TODO(), &stringID, &cityToUpdate)
		require.NoError(t, err)
		require.Equal(t, ID, updatedID)

		updatedCities, err := c.GetCities(context.TODO(), []int32{0}, 0, cityToUpdate.Name)
		require.NoError(t, err)
		require.NotNil(t, updatedCities)
		require.Equal(t, string(ID), updatedCities[0].ID)
		require.Equal(t, cityToUpdate.Name, updatedCities[0].Name)
		require.Equal(t, cityToUpdate.CountryCode, updatedCities[0].CountryCode)
		require.Equal(t, cityToUpdate.Rating, updatedCities[0].Rating)
	})

	t.Run("update city with invalid ID", func(t *testing.T) {
		ID := uuid.NewV4().String()
		_, err := c.PutCity(context.TODO(), &ID, &controller.City{
			Name:        "Sidney",
			CountryCode: 0,
			Rating:      0,
		})
		require.EqualError(t, err, controller.ErrCityNotFound.Error())
	})

	t.Run("update city with empty city struct", func(t *testing.T) {
		_, err := c.PutCity(context.TODO(), nil, nil)
		require.EqualError(t, err, controller.ErrInvalidCity.Error())
	})

	t.Run("update city with invalid Name", func(t *testing.T) {
		_, err := c.PutCity(context.TODO(), nil, &controller.City{
			Name:        "",
			CountryCode: 0,
			Rating:      0,
		})
		require.EqualError(t, err, controller.ErrInvalidCityName.Error())

		_, err = c.PutCity(context.TODO(), nil, &controller.City{
			Name: "Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef" +
				"Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef" +
				"Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef" +
				"Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef" +
				"Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef" +
				"Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef",
			CountryCode: 0,
			Rating:      0,
		})
		require.EqualError(t, err, controller.ErrInvalidCityName.Error())
	})
}

func TestController_GetCities(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)

	t.Run("get city list", func(t *testing.T) {
		count := 5
		citiesMap := make(map[string]*controller.City)
		for i := 0; i < count; i++ {
			city := &controller.City{
				Name:        fmt.Sprintf("City %d", i+1),
				CountryCode: int32(i % 2),
				Rating:      0,
			}
			ID, err := c.PutCity(context.TODO(), nil, city)
			require.NoError(t, err)
			require.NotNil(t, ID)

			citiesMap[string(ID)] = city
		}

		cities, err := c.GetCities(context.TODO(), []int32{}, 0, "")
		require.NoError(t, err)
		require.Equal(t, count, len(cities))

		for _, c := range cities {
			require.NotNil(t, citiesMap[c.ID])
			require.Equal(t, citiesMap[c.ID].Name, c.Name)
			require.Equal(t, citiesMap[c.ID].CountryCode, c.CountryCode)
			require.Equal(t, citiesMap[c.ID].Rating, c.Rating)
		}
	})
}

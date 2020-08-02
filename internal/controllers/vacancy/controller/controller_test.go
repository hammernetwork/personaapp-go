package controller_test

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	sqlMigrate "github.com/rubenv/sql-migrate"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	authController "personaapp/internal/controllers/auth/controller"
	companyController "personaapp/internal/controllers/company/controller"
	companyStorage "personaapp/internal/controllers/company/storage"
	"personaapp/internal/controllers/vacancy/controller"
	"personaapp/internal/controllers/vacancy/storage"
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

func TestController_PutVacancyCategory(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)

	t.Run("create new vacancy category", func(t *testing.T) {
		categoryToCreate := controller.VacancyCategory{
			Title:   "Put category",
			IconURL: "https://s3.bucket.org/1.jpg",
		}
		ID, err := c.PutVacancyCategory(context.TODO(), nil, &categoryToCreate)
		require.NoError(t, err)
		require.NotNil(t, ID)

		createdCategory, err := c.GetVacancyCategory(context.TODO(), string(ID))
		require.NoError(t, err)
		require.NotNil(t, createdCategory)
		require.Equal(t, categoryToCreate.Title, createdCategory.Title)
		require.Equal(t, categoryToCreate.IconURL, createdCategory.IconURL)
	})

	t.Run("update vacancy category", func(t *testing.T) {
		ID, err := c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
			Title:   "Put category to update",
			IconURL: "https://s3.bucket.org/create.jpg",
		})

		require.NoError(t, err)
		require.NotNil(t, ID)

		categoryToUpdate := controller.VacancyCategory{
			Title:   "Update category",
			IconURL: "https://s3.bucket.org/update.jpg",
		}

		stringID := string(ID)
		updatedID, err := c.PutVacancyCategory(context.TODO(), &stringID, &categoryToUpdate)
		require.NoError(t, err)
		require.Equal(t, ID, updatedID)

		updatedCategory, err := c.GetVacancyCategory(context.TODO(), stringID)
		require.NoError(t, err)
		require.NotNil(t, updatedCategory)
		require.Equal(t, string(ID), updatedCategory.ID)
		require.Equal(t, categoryToUpdate.Title, updatedCategory.Title)
		require.Equal(t, categoryToUpdate.IconURL, updatedCategory.IconURL)
	})

	t.Run("update vacancy category with invalid ID", func(t *testing.T) {
		ID := uuid.NewV4().String()
		_, err := c.PutVacancyCategory(context.TODO(), &ID, &controller.VacancyCategory{
			Title:   "Valid title",
			IconURL: "https://s3.bucket.org/valid_url.jpg",
		})
		require.EqualError(t, errors.Cause(err), controller.ErrVacancyCategoryNotFound.Error())
	})

	t.Run("update vacancy category with empty vacancy struct", func(t *testing.T) {
		_, err := c.PutVacancyCategory(context.TODO(), nil, nil)
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategory.Error())
	})

	t.Run("update vacancy category with invalid Title", func(t *testing.T) {
		_, err := c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
			Title:   "a",
			IconURL: "https://s3.bucket.org/valid_url.jpg",
		})
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategoryTitle.Error())

		_, err = c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
			Title:   "Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef",
			IconURL: "https://s3.bucket.org/valid_url.jpg",
		})
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategoryTitle.Error())
	})
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
		require.EqualError(t, errors.Cause(err), controller.ErrCityNotFound.Error())
	})

	t.Run("update city with empty city struct", func(t *testing.T) {
		_, err := c.PutCity(context.TODO(), nil, nil)
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidCity.Error())
	})

	t.Run("update city with invalid Name", func(t *testing.T) {
		_, err := c.PutCity(context.TODO(), nil, &controller.City{
			Name:        "",
			CountryCode: 0,
			Rating:      0,
		})
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidCityTitle.Error())

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
		require.EqualError(t, errors.Cause(err), controller.ErrInvalidCityTitle.Error())
	})
}

func TestController_PutVacancy(t *testing.T) {
	s, closer := initStorage(t)
	as, authCloser := testutils.InitAuthStorage(t)
	cs, companyCloser := initCompanyStorage(t)

	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		if err := authCloser(); err != nil {
			t.Error(err)
		}
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)
	ac := authController.New(authCfg, as)
	cc := companyController.New(cs)

	t.Run("create new vacancy", func(t *testing.T) {
		token, err := ac.Register(context.TODO(), &authController.RegisterData{
			Email:    "vacancy@gmail.com",
			Phone:    "+380503000001",
			Account:  authController.AccountTypeCompany,
			Password: "Password1488",
		})
		require.NoError(t, err)
		require.NotNil(t, token)

		claims, err := ac.GetAuthClaims(context.TODO(), token.Token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.NotEmpty(t, claims.AccountID)

		companyTitle := "Title"
		require.NoError(t, cc.Update(context.TODO(), &companyController.CompanyData{
			ID:    claims.AccountID,
			Title: &companyTitle,
		}))

		category := controller.VacancyCategory{
			Title:   "New vacancy category",
			IconURL: "https://s3.bucket.org/new_vacancy_category.jpg",
		}
		categoryID, err := c.PutVacancyCategory(context.TODO(), nil, &category)
		require.NoError(t, err)
		require.NotNil(t, categoryID)

		city := controller.City{
			Name:        "Obukhiv",
			CountryCode: 0,
			Rating:      0,
		}
		cityID, err := c.PutCity(context.TODO(), nil, &city)
		require.NoError(t, err)
		require.NotNil(t, cityID)

		vacancy := controller.VacancyDetails{
			Vacancy: controller.Vacancy{
				Title:     "Put vacancy",
				Phone:     "+380503000002",
				MinSalary: 10000,
				MaxSalary: 20000,
				ImageURLs: []string{
					"https://s3.bucket.org/new_vacancy.jpg",
					"https://s3.bucket.org/new_vacancy.jpg",
				},
				CompanyID: claims.AccountID,
			},
			Description:          "Description",
			WorkMonthsExperience: 10,
			WorkSchedule:         "24 hours, 7 days a week",
			LocationLatitude:     1.027,
			LocationLongitude:    2.055,
			Type:                 controller.VacancyTypeRemote,
			Address:              "Trafalgar sq",
			CountryCode:          0,
		}
		vacancyID, err := c.PutVacancy(context.TODO(), nil, &vacancy, []string{string(categoryID)}, []string{string(cityID)})
		require.NoError(t, err)
		require.NotNil(t, vacancyID)

		vd, err := c.GetVacancyDetails(context.TODO(), string(vacancyID))
		require.NoError(t, err)
		require.NotNil(t, vd)

		require.Equal(t, vacancy.Title, vd.Title)
		require.Equal(t, vacancy.Phone, vd.Phone)
		require.Equal(t, vacancy.MinSalary, vd.MinSalary)
		require.Equal(t, vacancy.MaxSalary, vd.MaxSalary)
		require.Equal(t, vacancy.ImageURLs, vd.ImageURLs)
		require.Equal(t, vacancy.CompanyID, vd.CompanyID)
		require.Equal(t, vacancy.Description, vd.Description)
		require.Equal(t, vacancy.WorkMonthsExperience, vd.WorkMonthsExperience)
		require.Equal(t, vacancy.WorkSchedule, vd.WorkSchedule)
		require.Equal(t, vacancy.LocationLatitude, vd.LocationLatitude)
		require.Equal(t, vacancy.LocationLongitude, vd.LocationLongitude)
	})
}

func TestController_GetVacanciesCategoriesList(t *testing.T) {
	s, closer := initStorage(t)
	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)

	t.Run("get vacancies categories list", func(t *testing.T) {
		count := 5
		categoriesMap := make(map[string]*controller.VacancyCategory)
		for i := 0; i < count; i++ {
			category := &controller.VacancyCategory{
				Title:   fmt.Sprintf("Category %d", i+1),
				IconURL: fmt.Sprintf("https://s3.bucket.org/category_%d.jpg", i+1),
			}
			ID, err := c.PutVacancyCategory(context.TODO(), nil, category)
			require.NoError(t, err)
			require.NotNil(t, ID)

			categoriesMap[string(ID)] = category
		}

		categories, err := c.GetVacanciesCategoriesList(context.TODO())
		require.NoError(t, err)
		require.Equal(t, count, len(categories))

		for _, cat := range categories {
			require.NotNil(t, categoriesMap[cat.ID])
			require.Equal(t, categoriesMap[cat.ID].Title, cat.Title)
			require.Equal(t, categoriesMap[cat.ID].IconURL, cat.IconURL)
		}
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

func TestController_GetVacanciesList(t *testing.T) {
	s, closer := initStorage(t)
	as, authCloser := testutils.InitAuthStorage(t)
	cs, companyCloser := initCompanyStorage(t)

	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		if err := authCloser(); err != nil {
			t.Error(err)
		}
		if err := companyCloser(); err != nil {
			t.Error(err)
		}
		cleanup(t)
	}()

	c := controller.New(s)
	ac := authController.New(authCfg, as)
	cc := companyController.New(cs)

	t.Run("get vacancies list", func(t *testing.T) {
		token, err := ac.Register(context.TODO(), &authController.RegisterData{
			Email:    "company@gmail.com",
			Phone:    "+380503000001",
			Account:  authController.AccountTypeCompany,
			Password: "Password1488",
		})
		require.NoError(t, err)
		require.NotNil(t, token)

		claims, err := ac.GetAuthClaims(context.TODO(), token.Token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.NotEmpty(t, claims.AccountID)

		companyTitle := "Title"
		require.NoError(t, cc.Update(context.TODO(), &companyController.CompanyData{
			ID:    claims.AccountID,
			Title: &companyTitle,
		}))

		categoriesCount := 4
		categoriesIDs := make([]string, categoriesCount)
		categoriesMap := make(map[string]*controller.VacancyCategory)
		for i := 0; i < categoriesCount; i++ {
			category := &controller.VacancyCategory{
				Title:   fmt.Sprintf("Category %d", i+1),
				IconURL: fmt.Sprintf("https://s3.bucket.org/category_%d.jpg", i+1),
			}
			ID, err := c.PutVacancyCategory(context.TODO(), nil, category)
			require.NoError(t, err)
			require.NotNil(t, ID)

			categoriesMap[string(ID)] = category
			categoriesIDs[i] = string(ID)
		}

		citiesCount := 3
		citiesIDs := make([]string, citiesCount)
		citiesMap := make(map[string]*controller.City)
		for i := 0; i < citiesCount; i++ {
			city := &controller.City{
				Name:        fmt.Sprintf("Town %d", i+1),
				CountryCode: int32(i),
				Rating:      0,
			}
			ID, err := c.PutCity(context.TODO(), nil, city)
			require.NoError(t, err)
			require.NotNil(t, ID)

			citiesMap[string(ID)] = city
			citiesIDs[i] = string(ID)
		}

		vacanciesCount := 3
		imagePlaceholder := "https://s3.bucket.org/vacancy_%d.jpg"
		vacanciesIDs := make([]string, vacanciesCount)
		vacanciesMap := make(map[string]*controller.VacancyDetails)
		for i := 0; i < vacanciesCount; i++ {
			vacancy := &controller.VacancyDetails{
				Vacancy: controller.Vacancy{
					Title:     fmt.Sprintf("Vacancy %d", i+1),
					Phone:     fmt.Sprintf("+38050100000%d", i+1),
					MinSalary: int32(i+1) * 1000,
					MaxSalary: int32(i+1) * 2000,
					ImageURLs: []string{
						fmt.Sprintf(imagePlaceholder, i+1),
						fmt.Sprintf(imagePlaceholder, i+2),
					},
					CompanyID: claims.AccountID,
				},
				Description:          fmt.Sprintf("Description %d", i+1),
				WorkMonthsExperience: int32(i + 1),
				WorkSchedule:         fmt.Sprintf("24 hours, %d days a week", i+1),
				LocationLatitude:     float32(i) * 1.027,
				LocationLongitude:    float32(i) * 2.055,
				Type:                 controller.VacancyTypeNormal,
				Address:              fmt.Sprintf("Address %d", i+1),
				CountryCode:          0,
			}

			ID, err := c.PutVacancy(context.TODO(), nil, vacancy, []string{categoriesIDs[i], categoriesIDs[i+1]}, []string{citiesIDs[i]})
			require.NoError(t, err)
			require.NotNil(t, ID)

			vacanciesMap[string(ID)] = vacancy
			vacanciesIDs[i] = string(ID)
		}

		t.Run("get all without filters", func(t *testing.T) {
			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 100)
			require.NoError(t, err)
			// Check vacancies
			require.Equal(t, 3, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
			require.Equal(t, vacanciesIDs[1], vacancies[1].ID)
			require.Equal(t, vacanciesIDs[0], vacancies[2].ID)
			// Check images
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[2].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[2].ImageURLs[1])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[1].ImageURLs[1])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
		})

		t.Run("get categories by vacancy ids", func(t *testing.T) {
			categories, err := c.GetVacanciesCategories(context.TODO(), vacanciesIDs)
			require.NoError(t, err)
			require.Equal(t, categoriesMap[categoriesIDs[3]].Title, categories[5].Title)
			require.Equal(t, categoriesMap[categoriesIDs[2]].Title, categories[4].Title)
			require.Equal(t, categoriesMap[categoriesIDs[2]].Title, categories[3].Title)
			require.Equal(t, categoriesMap[categoriesIDs[1]].Title, categories[2].Title)
			require.Equal(t, categoriesMap[categoriesIDs[1]].Title, categories[1].Title)
			require.Equal(t, categoriesMap[categoriesIDs[0]].Title, categories[0].Title)
		})

		t.Run("get cities by vacancy ids", func(t *testing.T) {
			cities, err := c.GetVacancyCities(context.TODO(), vacanciesIDs)
			require.NoError(t, err)
			require.Equal(t, citiesMap[citiesIDs[2]].Name, cities[2].Name)
			require.Equal(t, citiesMap[citiesIDs[1]].Name, cities[1].Name)
			require.Equal(t, citiesMap[citiesIDs[0]].Name, cities[0].Name)
		})

		t.Run("get all with invalid filter", func(t *testing.T) {
			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{uuid.NewV4().String()}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 0, len(vacancies))
		})

		t.Run("get all with invalid cursor", func(t *testing.T) {
			cursor := controller.Cursor("invalid cursor")
			_, _, err := c.GetVacanciesList(context.TODO(), []string{}, &cursor, 100)
			require.Error(t, err)
			require.EqualError(t, controller.ErrInvalidCursor, err.Error())
		})

		t.Run("get one with changed categories for cursor", func(t *testing.T) {
			_, cursor, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)

			_, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[1]}, cursor, 1)
			require.Error(t, err)
			require.EqualError(t, controller.ErrInvalidCursor, err.Error())
		})

		t.Run("get all by one category filter", func(t *testing.T) {
			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{categoriesIDs[0]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[0], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[1])

			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[1]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 2, len(vacancies))
			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[1])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[1].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[1])

			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[2]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 2, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[1].ImageURLs[1])

			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[3]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
		})

		t.Run("get 1 vacancy without filter and test pagination", func(t *testing.T) {
			vacancies, cursor, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[1])

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[0], vacancies[0].ID)
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[0].ImageURLs[0])
			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[1])

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.Nil(t, cursor)
			require.Equal(t, 0, len(vacancies))
		})
	})
}

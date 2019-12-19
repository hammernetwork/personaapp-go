package controller_test

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	sqlMigrate "github.com/rubenv/sql-migrate"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	authController "personaapp/internal/server/controllers/auth/controller"
	companyController "personaapp/internal/server/controllers/company/controller"
	companyStorage "personaapp/internal/server/controllers/company/storage"
	"personaapp/internal/server/controllers/vacancy/controller"
	"personaapp/internal/server/controllers/vacancy/storage"
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

		vacancy := controller.VacancyDetails{
			Vacancy: controller.Vacancy{
				Title:     "Put vacancy",
				Phone:     "+380503000002",
				MinSalary: 10000,
				MaxSalary: 20000,
				ImageURL:  "https://s3.bucket.org/new_vacancy.jpg",
				CompanyID: claims.AccountID,
			},
			Description:          "Description",
			WorkMonthsExperience: 10,
			WorkSchedule:         "24 hours, 7 days a week",
			LocationLatitude:     1.027,
			LocationLongitude:    2.055,
		}
		vacancyID, err := c.PutVacancy(context.TODO(), nil, &vacancy, []string{string(categoryID)})
		require.NoError(t, err)
		require.NotNil(t, vacancyID)

		vd, err := c.GetVacancyDetails(context.TODO(), string(vacancyID))
		require.NoError(t, err)
		require.NotNil(t, vd)

		require.Equal(t, vacancy.Title, vd.Title)
		require.Equal(t, vacancy.Phone, vd.Phone)
		require.Equal(t, vacancy.MinSalary, vd.MinSalary)
		require.Equal(t, vacancy.MaxSalary, vd.MaxSalary)
		require.Equal(t, vacancy.ImageURL, vd.ImageURL)
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

		categoriesCount := 3
		categoriesIDs := make([]string, 3)
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

		vacanciesCount := 3
		vacanciesIDs := make([]string, 3)
		vacanciesMap := make(map[string]*controller.VacancyDetails)
		for i := 0; i < vacanciesCount; i++ {
			vacancy := &controller.VacancyDetails{
				Vacancy: controller.Vacancy{
					Title:     fmt.Sprintf("Vacancy %d", i+1),
					Phone:     fmt.Sprintf("+38050100000%d", i+1),
					MinSalary: int32(i+1) * 1000,
					MaxSalary: int32(i+1) * 2000,
					ImageURL:  fmt.Sprintf("https://s3.bucket.org/vacancy_%d.jpg", i+1),
					CompanyID: claims.AccountID,
				},
				Description:          fmt.Sprintf("Description %d", i+1),
				WorkMonthsExperience: int32(i + 1),
				WorkSchedule:         fmt.Sprintf("24 hours, %d days a week", i+1),
				LocationLatitude:     float32(i) * 1.027,
				LocationLongitude:    float32(i) * 2.055,
			}

			ID, err := c.PutVacancy(context.TODO(), nil, vacancy, []string{categoriesIDs[i]})
			require.NoError(t, err)
			require.NotNil(t, ID)

			vacanciesMap[string(ID)] = vacancy
			vacanciesIDs[i] = string(ID)
		}

		t.Run("get all without filters", func(t *testing.T) {
			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 3, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
			require.Equal(t, vacanciesIDs[1], vacancies[1].ID)
			require.Equal(t, vacanciesIDs[0], vacancies[2].ID)
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

			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[1]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)

			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[2]}, nil, 100)
			require.NoError(t, err)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
		})

		t.Run("get 1 vacancy without filter and test pagination", func(t *testing.T) {
			vacancies, cursor, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, 1, len(vacancies))
			require.Equal(t, vacanciesIDs[0], vacancies[0].ID)

			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
			require.NoError(t, err)
			require.Nil(t, cursor)
			require.Equal(t, 0, len(vacancies))
		})
	})
}

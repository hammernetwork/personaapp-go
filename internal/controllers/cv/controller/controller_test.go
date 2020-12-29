package controller_test

import (
	"context"
	"fmt"
	sqlMigrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
	authController "personaapp/internal/controllers/auth/controller"
	"personaapp/internal/controllers/cv/controller"
	"personaapp/internal/controllers/cv/storage"
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

//func initCVStorage(t *testing.T) (_ *cvStorage.Storage, closer func() error) {
//	pg := testutils.EnsurePostgres(t)
//	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Up))
//
//	return cvStorage.New(pg), pg.Close
//}

func cleanup(t *testing.T) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Down))
}

//func TestController_PutVacancyCategory(t *testing.T) {
//	s, closer := initStorage(t)
//	defer func() {
//		if err := closer(); err != nil {
//			t.Error(err)
//		}
//		cleanup(t)
//	}()
//
//	c := controller.New(s)
//
//	t.Run("create new vacancy category", func(t *testing.T) {
//		categoryToCreate := controller.VacancyCategory{
//			Title:   "Put category",
//			IconURL: "https://s3.bucket.org/1.jpg",
//			Rating:  0,
//		}
//		ID, err := c.PutVacancyCategory(context.TODO(), nil, &categoryToCreate)
//		require.NoError(t, err)
//		require.NotNil(t, ID)
//
//		createdCategory, err := c.GetVacancyCategory(context.TODO(), string(ID))
//		require.NoError(t, err)
//		require.NotNil(t, createdCategory)
//		require.Equal(t, categoryToCreate.Title, createdCategory.Title)
//		require.Equal(t, categoryToCreate.IconURL, createdCategory.IconURL)
//		require.Equal(t, categoryToCreate.Rating, createdCategory.Rating)
//	})
//
//	t.Run("update vacancy category", func(t *testing.T) {
//		ID, err := c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
//			Title:   "Put category to update",
//			IconURL: "https://s3.bucket.org/create.jpg",
//			Rating:  0,
//		})
//
//		require.NoError(t, err)
//		require.NotNil(t, ID)
//
//		categoryToUpdate := controller.VacancyCategory{
//			Title:   "Update category",
//			IconURL: "https://s3.bucket.org/update.jpg",
//			Rating:  1,
//		}
//
//		stringID := string(ID)
//		updatedID, err := c.PutVacancyCategory(context.TODO(), &stringID, &categoryToUpdate)
//		require.NoError(t, err)
//		require.Equal(t, ID, updatedID)
//
//		updatedCategory, err := c.GetVacancyCategory(context.TODO(), stringID)
//		require.NoError(t, err)
//		require.NotNil(t, updatedCategory)
//		require.Equal(t, string(ID), updatedCategory.ID)
//		require.Equal(t, categoryToUpdate.Title, updatedCategory.Title)
//		require.Equal(t, categoryToUpdate.IconURL, updatedCategory.IconURL)
//		require.Equal(t, categoryToUpdate.Rating, updatedCategory.Rating)
//	})
//
//	t.Run("update vacancy category with invalid ID", func(t *testing.T) {
//		ID := uuid.NewV4().String()
//		_, err := c.PutVacancyCategory(context.TODO(), &ID, &controller.VacancyCategory{
//			Title:   "Valid title",
//			IconURL: "https://s3.bucket.org/valid_url.jpg",
//			Rating:  0,
//		})
//		require.EqualError(t, errors.Cause(err), controller.ErrVacancyCategoryNotFound.Error())
//	})
//
//	t.Run("update vacancy category with empty vacancy struct", func(t *testing.T) {
//		_, err := c.PutVacancyCategory(context.TODO(), nil, nil)
//		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategory.Error())
//	})
//
//	t.Run("update vacancy category with invalid Title", func(t *testing.T) {
//		_, err := c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
//			Title:   "a",
//			IconURL: "https://s3.bucket.org/valid_url.jpg",
//			Rating:  0,
//		})
//		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategoryTitle.Error())
//
//		_, err = c.PutVacancyCategory(context.TODO(), nil, &controller.VacancyCategory{
//			Title:   "Abcd abcd abcd abcd abcd abcd abcd abcd abcd abcdef",
//			IconURL: "https://s3.bucket.org/valid_url.jpg",
//			Rating:  0,
//		})
//		require.EqualError(t, errors.Cause(err), controller.ErrInvalidVacancyCategoryTitle.Error())
//	})
//}

func TestController_PutCV(t *testing.T) {
	s, closer := initStorage(t)
	as, authCloser := testutils.InitAuthStorage(t)
	//cvc, cvCloser := initCVStorage(t)

	defer func() {
		if err := closer(); err != nil {
			t.Error(err)
		}
		if err := authCloser(); err != nil {
			t.Error(err)
		}
		//if err := cvCloser(); err != nil {
		//	t.Error(err)
		//}
		cleanup(t)
	}()

	c := controller.New(s)
	ac := authController.New(authCfg, as)
	//cv := cvController.New(cvc)

	t.Run("create new cv", func(t *testing.T) {
		token, err := ac.Register(context.TODO(), &authController.RegisterData{
			Email:    "cv@gmail.com",
			Phone:    "+380503030001",
			Account:  authController.AccountTypePersona,
			Password: "Password1488",
		})
		require.NoError(t, err)
		require.NotNil(t, token)

		claims, err := ac.GetAuthClaims(context.TODO(), token.Token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.NotEmpty(t, claims.AccountID)

		newCV := controller.CV{
			PersonaID:            claims.AccountID,
			Position:             "Position 1",
			WorkMonthsExperience: 12,
			MinSalary:            1200,
			MaxSalary:            12000,
		}

		cvID, err := c.PutCV(context.TODO(), nil, &newCV)
		require.NoError(t, err)
		require.NotNil(t, cvID)

		cv, err := c.GetCV(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.NotNil(t, cv)

		require.Equal(t, newCV.Position, cv.Position)
		require.Equal(t, newCV.WorkMonthsExperience, cv.WorkMonthsExperience)
		require.Equal(t, newCV.MinSalary, cv.MinSalary)
		require.Equal(t, newCV.MaxSalary, cv.MaxSalary)

		// Custom section
		count := 5
		sectionsMap := make(map[string]*controller.CVCustomSection)
		for i := 0; i < count; i++ {
			section := &controller.CVCustomSection{
				CvID:        string(cvID),
				Description: fmt.Sprintf("Section %d", i+1),
			}
			ID, err := c.PutCustomSection(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			sectionsMap[string(ID)] = section
		}

		sections, err := c.GetCustomSections(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.Equal(t, count, len(sections))

		for _, cat := range sections {
			require.NotNil(t, sectionsMap[cat.ID])
			require.Equal(t, sectionsMap[cat.ID].Description, cat.Description)
			require.Equal(t, sectionsMap[cat.ID].CvID, cat.CvID)
		}

		// Education
		educationsMap := make(map[string]*controller.CVEducation)
		for i := 0; i < count; i++ {
			section := &controller.CVEducation{
				CvID:        string(cvID),
				Institution: fmt.Sprintf("Institution %d", i+1),
				DateFrom:    time.Now(),
				DateTill:    time.Now().Add(1000),
				Speciality:  fmt.Sprintf("Speciality %d", i+1),
				Description: fmt.Sprintf("Description %d", i+1),
			}
			ID, err := c.PutEducation(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			educationsMap[string(ID)] = section
		}

		educations, err := c.GetEducations(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.Equal(t, count, len(educations))

		for _, cat := range educations {
			require.NotNil(t, educationsMap[cat.ID])
			require.Equal(t, educationsMap[cat.ID].Description, cat.Description)
			require.Equal(t, educationsMap[cat.ID].Institution, cat.Institution)
			require.Equal(t, educationsMap[cat.ID].Speciality, cat.Speciality)
		}

		// Experience
		experienceMap := make(map[string]*controller.CVExperience)
		for i := 0; i < count; i++ {
			section := &controller.CVExperience{
				CvID:        string(cvID),
				CompanyName: fmt.Sprintf("CompanyName %d", i+1),
				DateFrom:    time.Now(),
				DateTill:    time.Now().Add(1000),
				Position:    fmt.Sprintf("Position %d", i+1),
				Description: fmt.Sprintf("Description %d", i+1),
			}
			ID, err := c.PutExperience(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			experienceMap[string(ID)] = section
		}

		experiences, err := c.GetExperiences(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.Equal(t, count, len(experiences))

		for _, cat := range experiences {
			require.NotNil(t, experienceMap[cat.ID])
			require.Equal(t, experienceMap[cat.ID].Description, cat.Description)
			require.Equal(t, experienceMap[cat.ID].Position, cat.Position)
			require.Equal(t, experienceMap[cat.ID].CompanyName, cat.CompanyName)
		}

		// Story
		storyMap := make(map[string]*controller.CVCustomStory)
		for i := 0; i < count; i++ {
			section := &controller.CVCustomStory{
				CvID:        string(cvID),
				ChapterName: fmt.Sprintf("ChapterName %d", i+1),
			}
			ID, err := c.PutStory(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			storyMap[string(ID)] = section
		}

		stories, err := c.GetStories(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.Equal(t, count, len(stories))

		for _, cat := range stories {
			require.NotNil(t, storyMap[cat.ID])
			require.Equal(t, storyMap[cat.ID].ChapterName, cat.ChapterName)
		}
		storyID := stories[0].ID

		// Episodes
		episodesMap := make(map[string]*controller.StoryEpisode)
		for i := 0; i < count; i++ {
			episode := &controller.StoryEpisode{
				StoryID:  storyID,
				MediaURL: fmt.Sprintf("MediaURL %d", i+1),
			}
			ID, err := c.PutStoriesEpisode(context.TODO(), nil, episode)
			require.NoError(t, err)
			require.NotNil(t, ID)

			episodesMap[string(ID)] = episode
		}

		episodes, err := c.GetStoriesEpisodes(context.TODO(), string(cvID))
		require.NoError(t, err)
		require.Equal(t, count, len(stories))

		for _, cat := range episodes {
			require.NotNil(t, episodesMap[cat.ID])
			require.Equal(t, episodesMap[cat.ID].StoryID, cat.StoryID)
			require.Equal(t, episodesMap[cat.ID].MediaURL, cat.MediaURL)
		}

		// Job kind
		kindMap := make(map[string]*controller.JobKind)
		for i := 0; i < count; i++ {
			section := &controller.JobKind{
				Name: fmt.Sprintf("Name %d", i+1),
			}
			ID, err := c.PutJobKind(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			kindMap[string(ID)] = section
		}

		kinds, err := c.GetJobKinds(context.TODO())
		require.NoError(t, err)
		require.Equal(t, count, len(kinds))

		for _, cat := range kinds {
			require.NotNil(t, kindMap[cat.ID])
			require.Equal(t, kindMap[cat.ID].Name, cat.Name)
		}

		// CV Job kind
		cvKinds := make([]string, count)
		for i := 0; i < count; i++ {
			cvKinds[i] = kinds[i].ID
		}

		err = c.PutCVJobKinds(context.TODO(), string(cvID), cvKinds)
		require.NoError(t, err)

		cvJobKinds, err := c.GetCVJobKinds(context.TODO(), string(cvID))
		require.NoError(t, err)

		for idx, cat := range cvJobKinds {
			require.NotNil(t, cvJobKinds[idx])
			require.Equal(t, cvJobKinds[idx].ID, cat.ID)
			require.Equal(t, cvJobKinds[idx].Name, cat.Name)
		}

		// Job types
		typeMap := make(map[string]*controller.JobType)
		for i := 0; i < count; i++ {
			section := &controller.JobType{
				Name: fmt.Sprintf("Name %d", i+1),
			}
			ID, err := c.PutJobType(context.TODO(), nil, section)
			require.NoError(t, err)
			require.NotNil(t, ID)

			typeMap[string(ID)] = section
		}

		tpes, err := c.GetJobTypes(context.TODO())
		require.NoError(t, err)
		require.Equal(t, count, len(tpes))

		for _, cat := range kinds {
			require.NotNil(t, kindMap[cat.ID])
			require.Equal(t, kindMap[cat.ID].Name, cat.Name)
		}

		// CV Job types
		cvTypes := make([]string, count)
		for i := 0; i < count; i++ {
			cvTypes[i] = tpes[i].ID
		}

		err = c.PutCVJobTypes(context.TODO(), string(cvID), cvTypes)
		require.NoError(t, err)

		cvJobTypes, err := c.GetCVJobTypes(context.TODO(), string(cvID))
		require.NoError(t, err)

		for idx, cat := range cvJobTypes {
			require.NotNil(t, cvJobTypes[idx])
			require.Equal(t, cvJobTypes[idx].ID, cat.ID)
			require.Equal(t, cvJobTypes[idx].Name, cat.Name)
		}

	})

}

//func TestController_GetVacanciesCategoriesList(t *testing.T) {
//	s, closer := initStorage(t)
//	defer func() {
//		if err := closer(); err != nil {
//			t.Error(err)
//		}
//		cleanup(t)
//	}()

//c := controller.New(s)

//t.Run("get vacancies categories list", func(t *testing.T) {
//	count := 5
//	categoriesMap := make(map[string]*controller.VacancyCategory)
//	for i := 0; i < count; i++ {
//		category := &controller.VacancyCategory{
//			Title:   fmt.Sprintf("Category %d", i+1),
//			IconURL: fmt.Sprintf("https://s3.bucket.org/category_%d.jpg", i+1),
//			Rating:  int32(i % 2),
//		}
//		ID, err := c.PutVacancyCategory(context.TODO(), nil, category)
//		require.NoError(t, err)
//		require.NotNil(t, ID)
//
//		categoriesMap[string(ID)] = category
//	}
//
//	categories, err := c.GetVacanciesCategoriesList(context.TODO(), nil)
//	require.NoError(t, err)
//	require.Equal(t, count, len(categories))
//
//	for _, cat := range categories {
//		require.NotNil(t, categoriesMap[cat.ID])
//		require.Equal(t, categoriesMap[cat.ID].Title, cat.Title)
//		require.Equal(t, categoriesMap[cat.ID].IconURL, cat.IconURL)
//		require.Equal(t, categoriesMap[cat.ID].Rating, cat.Rating)
//	}
//
//	popularCategoriesCount := 2
//	rating := new(int32)
//	*rating = 1
//	categories, err = c.GetVacanciesCategoriesList(context.TODO(), rating)
//	require.NoError(t, err)
//	require.Equal(t, popularCategoriesCount, len(categories))
//})
//}

//func TestController_GetVacanciesList(t *testing.T) {
//	s, closer := initStorage(t)
//	as, authCloser := testutils.InitAuthStorage(t)
//	cs, companyCloser := initCompanyStorage(t)
//	cys, cityCloser := initCityStorage(t)
//
//	defer func() {
//		if err := closer(); err != nil {
//			t.Error(err)
//		}
//		if err := authCloser(); err != nil {
//			t.Error(err)
//		}
//		if err := companyCloser(); err != nil {
//			t.Error(err)
//		}
//		if err := cityCloser(); err != nil {
//			t.Error(err)
//		}
//		cleanup(t)
//	}()
//
//	c := controller.New(s)
//	ac := authController.New(authCfg, as)
//	cc := companyController.New(cs)
//	cy := cityController.New(cys)
//
//	t.Run("get vacancies list", func(t *testing.T) {
//		token, err := ac.Register(context.TODO(), &authController.RegisterData{
//			Email:    "company@gmail.com",
//			Phone:    "+380503000001",
//			Account:  authController.AccountTypeCompany,
//			Password: "Password1488",
//		})
//		require.NoError(t, err)
//		require.NotNil(t, token)
//
//		claims, err := ac.GetAuthClaims(context.TODO(), token.Token)
//		require.NoError(t, err)
//		require.NotNil(t, claims)
//		require.NotEmpty(t, claims.AccountID)
//
//		companyTitle := "Title"
//		require.NoError(t, cc.Update(context.TODO(), &companyController.CompanyData{
//			ID:    claims.AccountID,
//			Title: &companyTitle,
//		}))
//
//		categoriesCount := 4
//		categoriesIDs := make([]string, categoriesCount)
//		categoriesMap := make(map[string]*controller.VacancyCategory)
//		for i := 0; i < categoriesCount; i++ {
//			category := &controller.VacancyCategory{
//				Title:   fmt.Sprintf("Category %d", i+1),
//				IconURL: fmt.Sprintf("https://s3.bucket.org/category_%d.jpg", i+1),
//				Rating:  int32(i % 2),
//			}
//			ID, err := c.PutVacancyCategory(context.TODO(), nil, category)
//			require.NoError(t, err)
//			require.NotNil(t, ID)
//
//			categoriesMap[string(ID)] = category
//			categoriesIDs[i] = string(ID)
//		}
//
//		citiesCount := 3
//		citiesIDs := make([]string, citiesCount)
//		citiesMap := make(map[string]*cityController.City)
//		for i := 0; i < citiesCount; i++ {
//			city := &cityController.City{
//				Name:        fmt.Sprintf("Town %d", i+1),
//				CountryCode: int32(i),
//				Rating:      0,
//			}
//			ID, err := cy.PutCity(context.TODO(), nil, city)
//			require.NoError(t, err)
//			require.NotNil(t, ID)
//
//			citiesMap[string(ID)] = city
//			citiesIDs[i] = string(ID)
//		}
//
//		vacanciesCount := 3
//		imagePlaceholder := "https://s3.bucket.org/vacancy_%d.jpg"
//		vacanciesIDs := make([]string, vacanciesCount)
//		vacanciesMap := make(map[string]*controller.VacancyDetails)
//		for i := 0; i < vacanciesCount; i++ {
//			vacancy := &controller.VacancyDetails{
//				Vacancy: controller.Vacancy{
//					Title:     fmt.Sprintf("Vacancy %d", i+1),
//					Phone:     fmt.Sprintf("+38050100000%d", i+1),
//					MinSalary: int32(i+1) * 1000,
//					MaxSalary: int32(i+1) * 2000,
//					ImageURLs: []string{
//						fmt.Sprintf(imagePlaceholder, i+1),
//						fmt.Sprintf(imagePlaceholder, i+2),
//					},
//					CompanyID: claims.AccountID,
//				},
//				Description:          fmt.Sprintf("Description %d", i+1),
//				WorkMonthsExperience: int32(i + 1),
//				WorkSchedule:         fmt.Sprintf("24 hours, %d days a week", i+1),
//				LocationLatitude:     float32(i) * 1.027,
//				LocationLongitude:    float32(i) * 2.055,
//				Type:                 controller.VacancyTypeNormal,
//				Address:              fmt.Sprintf("Address %d", i+1),
//				CountryCode:          0,
//			}
//
//			ID, err := c.PutVacancy(context.TODO(), nil, vacancy, []string{categoriesIDs[i], categoriesIDs[i+1]}, []string{citiesIDs[i]})
//			require.NoError(t, err)
//			require.NotNil(t, ID)
//
//			vacanciesMap[string(ID)] = vacancy
//			vacanciesIDs[i] = string(ID)
//		}
//
//		t.Run("get all without filters", func(t *testing.T) {
//			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 100)
//			require.NoError(t, err)
//			// Check vacancies
//			require.Equal(t, 3, len(vacancies))
//			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
//			require.Equal(t, vacanciesIDs[1], vacancies[1].ID)
//			require.Equal(t, vacanciesIDs[0], vacancies[2].ID)
//			// Check images
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[2].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[2].ImageURLs[1])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[1].ImageURLs[1])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
//		})
//
//		t.Run("get categories by vacancy ids", func(t *testing.T) {
//			categories, err := c.GetVacanciesCategories(context.TODO(), vacanciesIDs)
//			require.NoError(t, err)
//			require.Equal(t, categoriesMap[categoriesIDs[3]].Title, categories[5].Title)
//			require.Equal(t, categoriesMap[categoriesIDs[2]].Title, categories[4].Title)
//			require.Equal(t, categoriesMap[categoriesIDs[2]].Title, categories[3].Title)
//			require.Equal(t, categoriesMap[categoriesIDs[1]].Title, categories[2].Title)
//			require.Equal(t, categoriesMap[categoriesIDs[1]].Title, categories[1].Title)
//			require.Equal(t, categoriesMap[categoriesIDs[0]].Title, categories[0].Title)
//		})
//
//		t.Run("get cities by vacancy ids", func(t *testing.T) {
//			cities, err := c.GetVacancyCities(context.TODO(), vacanciesIDs)
//			require.NoError(t, err)
//			require.Equal(t, citiesMap[citiesIDs[2]].Name, cities[2].Name)
//			require.Equal(t, citiesMap[citiesIDs[1]].Name, cities[1].Name)
//			require.Equal(t, citiesMap[citiesIDs[0]].Name, cities[0].Name)
//		})
//
//		t.Run("get all with invalid filter", func(t *testing.T) {
//			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{uuid.NewV4().String()}, nil, 100)
//			require.NoError(t, err)
//			require.Equal(t, 0, len(vacancies))
//		})
//
//		t.Run("get all with invalid cursor", func(t *testing.T) {
//			cursor := controller.Cursor("invalid cursor")
//			_, _, err := c.GetVacanciesList(context.TODO(), []string{}, &cursor, 100)
//			require.Error(t, err)
//			require.EqualError(t, controller.ErrInvalidCursor, err.Error())
//		})
//
//		t.Run("get one with changed categories for cursor", func(t *testing.T) {
//			_, cursor, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 1)
//			require.NoError(t, err)
//			require.NotNil(t, cursor)
//
//			_, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[1]}, cursor, 1)
//			require.Error(t, err)
//			require.EqualError(t, controller.ErrInvalidCursor, err.Error())
//		})
//
//		t.Run("get all by one category filter", func(t *testing.T) {
//			vacancies, _, err := c.GetVacanciesList(context.TODO(), []string{categoriesIDs[0]}, nil, 100)
//			require.NoError(t, err)
//			require.Equal(t, 1, len(vacancies))
//			require.Equal(t, vacanciesIDs[0], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[1])
//
//			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[1]}, nil, 100)
//			require.NoError(t, err)
//			require.Equal(t, 2, len(vacancies))
//			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[1])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[1].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[1])
//
//			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[2]}, nil, 100)
//			require.NoError(t, err)
//			require.Equal(t, 2, len(vacancies))
//			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[1].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[1].ImageURLs[1])
//
//			vacancies, _, err = c.GetVacanciesList(context.TODO(), []string{categoriesIDs[3]}, nil, 100)
//			require.NoError(t, err)
//			require.Equal(t, 1, len(vacancies))
//			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
//		})
//
//		t.Run("get 1 vacancy without filter and test pagination", func(t *testing.T) {
//			vacancies, cursor, err := c.GetVacanciesList(context.TODO(), []string{}, nil, 1)
//			require.NoError(t, err)
//			require.NotNil(t, cursor)
//			require.Equal(t, 1, len(vacancies))
//			require.Equal(t, vacanciesIDs[2], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 4), vacancies[0].ImageURLs[1])
//
//			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
//			require.NoError(t, err)
//			require.NotNil(t, cursor)
//			require.Equal(t, 1, len(vacancies))
//			require.Equal(t, vacanciesIDs[1], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 3), vacancies[0].ImageURLs[1])
//
//			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
//			require.NoError(t, err)
//			require.NotNil(t, cursor)
//			require.Equal(t, 1, len(vacancies))
//			require.Equal(t, vacanciesIDs[0], vacancies[0].ID)
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 1), vacancies[0].ImageURLs[0])
//			require.Equal(t, fmt.Sprintf(imagePlaceholder, 2), vacancies[0].ImageURLs[1])
//
//			vacancies, cursor, err = c.GetVacanciesList(context.TODO(), []string{}, cursor, 1)
//			require.NoError(t, err)
//			require.Nil(t, cursor)
//			require.Equal(t, 0, len(vacancies))
//		})
//	})
//}

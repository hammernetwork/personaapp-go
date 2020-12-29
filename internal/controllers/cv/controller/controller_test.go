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

func cleanup(t *testing.T) {
	pg := testutils.EnsurePostgres(t)
	require.NoError(t, testutils.Migrate(pg.DB, sqlMigrate.Down))
}

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
		cleanup(t)
	}()

	c := controller.New(s)
	ac := authController.New(authCfg, as)

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

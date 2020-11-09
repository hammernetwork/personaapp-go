package controller

import (
	"context"
	"personaapp/internal/controllers/cv/storage"
	pkgtx "personaapp/pkg/tx"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint:godox // TODO: Implement CDN link check
		return true
	})
}

var (
	ErrCVNotFound              = errors.New("cv not found")
	ErrStoriesEpisodesNotFound = errors.New("stories episodes not found")
	ErrStoriesNotFound         = errors.New("stories not found")
	ErrCustomSectionsNotFound  = errors.New("custom sections not found")
	ErrCustomSectionNotFound   = errors.New("custom section not found")

	ErrInvalidStoriesEpisodesMediaURL = errors.New("invalid stories episodes media url")
	ErrInvalidStoryMediaURL           = errors.New("invalid story media url")

	ErrVacancyCategoryNotFound            = errors.New("vacancy category not found")
	ErrVacancyImagesNotFound              = errors.New("vacancy image not found")
	ErrInvalidCursor                      = errors.New("invalid cursor")
	ErrInvalidStoryEpisode                = errors.New("invalid story episode struct")
	ErrInvalidStory                       = errors.New("invalid story struct")
	ErrInvalidVacancyCategoryTitle        = errors.New("invalid vacancy category title")
	ErrInvalidVacancy                     = errors.New("invalid vacancy struct")
	ErrInvalidVacancyTitle                = errors.New("invalid vacancy title")
	ErrInvalidVacancyPhone                = errors.New("invalid vacancy phone")
	ErrInvalidVacancyMinSalary            = errors.New("invalid vacancy min salary")
	ErrInvalidVacancyMaxSalary            = errors.New("invalid vacancy max salary")
	ErrInvalidVacancyImageURL             = errors.New("invalid vacancy image url")
	ErrInvalidVacancyCompanyID            = errors.New("invalid vacancy company id")
	ErrInvalidVacancyDescription          = errors.New("invalid vacancy description")
	ErrInvalidVacancyWorkMonthsExperience = errors.New("invalid vacancy work months experience")
	ErrInvalidVacancyWorkSchedule         = errors.New("invalid vacancy work schedule")
	ErrInvalidVacancyLocationLatitude     = errors.New("invalid vacancy location latitude")
	ErrInvalidVacancyLocationLongitude    = errors.New("invalid vacancy location longitude")
)

type Storage interface {
	TxPutJobType(ctx context.Context, tx pkgtx.Tx, jobType *storage.JobType) error
	TxGetJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
	) (_ []*storage.JobType, rerr error)
	TxPutCVJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
		jobTypesIDs []string,
	) error
	TxGetCVJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVJobType, error)
	TxDeleteCVJobTypes(ctx context.Context, tx pkgtx.Tx, CVID string)

	TxPutJobKind(ctx context.Context, tx pkgtx.Tx, jobKind *storage.JobKind) error
	TxGetJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
	) (_ []*storage.JobKind, rerr error)
	TxPutCVJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
		jobKindsIDs []string,
	) error
	TxGetCVJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVJobKind, error)
	TxDeleteCVJobKinds(ctx context.Context, tx pkgtx.Tx, CVID string) error

	TxPutCVExperience(ctx context.Context, tx pkgtx.Tx, CVID string, experience *storage.CVExperience) error
	TxGetCVExperience(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVExperience, error)
	TxDeleteCVExperience(ctx context.Context, tx pkgtx.Tx, experienceID string) error

	TxPutCVEducations(ctx context.Context, tx pkgtx.Tx, CVID string, education *storage.CVEducation) error
	TxGetCVEducations(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVEducation, error)
	TxDeleteCVEducations(ctx context.Context, tx pkgtx.Tx, educationID string) error

	TxPutCustomSections(ctx context.Context, tx pkgtx.Tx, CVID string, education *storage.CVCustomSection) error
	TxGetCustomSections(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVCustomSection, error)
	TxDeleteCustomSection(ctx context.Context, tx pkgtx.Tx, sectionID string) error

	TxPutStory(ctx context.Context, tx pkgtx.Tx, CVID string, story *storage.CVCustomStory) error
	TxGetStories(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVCustomStory, error)
	TxDeleteStory(ctx context.Context, tx pkgtx.Tx, storyID string) error

	TxPutStoriesEpisodes(ctx context.Context, tx pkgtx.Tx, story *storage.StoryEpisode) error
	TxGetStoriesEpisodesById(
		ctx context.Context,
		tx pkgtx.Tx,
		storyEpisodeID string,
	) ([]*StoryEpisode, error)
	TxGetStoriesEpisodes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.StoryEpisode, error)
	TxDeleteStoriesEpisodes(ctx context.Context, tx pkgtx.Tx, episodeID string) error

	TxPutCV(ctx context.Context, tx pkgtx.Tx, cv *storage.CV) error
	TxGetCV(ctx context.Context, tx pkgtx.Tx, CVID string) (*storage.CV, error)
	TxGetCVs(
		ctx context.Context,
		tx pkgtx.Tx,
		personaID string,
	) ([]*storage.CVShort, error)
	TxDeleteCV(ctx context.Context, tx pkgtx.Tx, cvID string) error

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

type VacancyCategoryID string

type VacancyCategory struct {
	ID      string
	Title   string `valid:"stringlength(2|50)"`
	IconURL string `valid:"stringlength(10|255),media_link"`
	Rating  int32
}

type VacancyID string

// Vacancy models for put
type Vacancy struct {
	ID        string
	Title     string   `valid:"stringlength(5|80),required"`
	Phone     string   `valid:"phone,required"`
	MinSalary int32    `valid:"range(0|1000000000),required"`
	MaxSalary int32    `valid:"range(0|1000000000),required"`
	ImageURLs []string `valid:"stringlength(0|255),media_link"`
	CompanyID string   `valid:"required"`
}

type VacancyDetails struct {
	Vacancy
	Description          string
	WorkMonthsExperience int32   `valid:"range(0|1200),required"`
	WorkSchedule         string  `valid:"stringlength(0|100)"`
	LocationLatitude     float32 `valid:"latitude"`
	LocationLongitude    float32 `valid:"longitude"`
	Type                 VacancyType
	Address              string
	CountryCode          int32
}

type VacancyCategoryShort struct {
	VacancyID string
	ID        string
	Title     string
}

type VacancyCity struct {
	VacancyID   string
	ID          string
	Name        string
	CountryCode int32
	Rating      int32
}

type Cursor string

type StoryEpisode struct {
	StoryID  string
	ID       string
	MediaUrl string `valid:"stringlength(10|2000),media_link"`
}

type CVCustomStory struct {
	ID          string
	ChapterName string
	MediaUrl    string `valid:"stringlength(10|2000),media_link"`
}

type CVCustomSection struct {
	ID          string
	Description string
}

type StoriesEpisodeID string

type StoryID string

type CustomSectionID string

func (c Cursor) String() string {
	return string(c)
}

func (vc *StoryEpisode) validate() error {
	if vc == nil {
		return ErrInvalidStoryEpisode
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "MediaUrl", DefaultError: ErrInvalidStoriesEpisodesMediaURL},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("story episode struct is filled with some invalid data")
	}

	return nil
}

func (vc *CVCustomStory) validate() error {
	if vc == nil {
		return ErrInvalidStory
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "MediaUrl", DefaultError: ErrInvalidStoryMediaURL},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("story struct is filled with some invalid data")
	}

	return nil
}

func (vc *VacancyCategory) validate() error {
	if vc == nil {
		return ErrInvalidVacancyCategory
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "Title", DefaultError: ErrInvalidVacancyCategoryTitle},
		{Field: "IconURL", DefaultError: ErrInvalidVacancyCategoryIconURL},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("vacancy category struct is filled with some invalid data")
	}

	return nil
}

func (vd *VacancyDetails) validate() error {
	if vd == nil {
		return ErrInvalidVacancy
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "Title", DefaultError: ErrInvalidVacancyTitle},
		{Field: "Phone", DefaultError: ErrInvalidVacancyPhone},
		{Field: "MinSalary", DefaultError: ErrInvalidVacancyMinSalary},
		{Field: "MaxSalary", DefaultError: ErrInvalidVacancyMaxSalary},
		{Field: "ImageURL", DefaultError: ErrInvalidVacancyImageURL},
		{Field: "CompanyID", DefaultError: ErrInvalidVacancyCompanyID},
		{Field: "Description", DefaultError: ErrInvalidVacancyDescription},
		{Field: "WorkMonthsExperience", DefaultError: ErrInvalidVacancyWorkMonthsExperience},
		{Field: "WorkSchedule", DefaultError: ErrInvalidVacancyWorkSchedule},
		{Field: "LocationLatitude", DefaultError: ErrInvalidVacancyLocationLatitude},
		{Field: "LocationLongitude", DefaultError: ErrInvalidVacancyLocationLongitude},
	}

	if valid, err := govalidator.ValidateStruct(vd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.Wrap(fe.DefaultError, msg)
			}
		}

		return errors.New("vacancy details struct is filled with some invalid data")
	}

	return nil
}

func (c *Controller) PutVacancyCategory(
	ctx context.Context,
	categoryID *string,
	category *VacancyCategory,
) (VacancyCategoryID, error) {
	var ID VacancyCategoryID

	if err := category.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if categoryID != nil {
			switch _, err := c.s.TxGetVacancyCategory(ctx, tx, *categoryID); errors.Cause(err) {
			case nil:
				ID = VacancyCategoryID(*categoryID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrVacancyCategoryNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = VacancyCategoryID(uuid.NewV4().String())
		}

		now := time.Now()

		svc := storage.VacancyCategory{
			ID:        string(ID),
			Title:     category.Title,
			IconURL:   category.IconURL,
			Rating:    category.Rating,
			CreatedAt: now,
			UpdatedAt: now,
		}

		return errors.WithStack(c.s.TxPutVacancyCategory(ctx, tx, &svc))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetVacancyCategory(ctx context.Context, categoryID string) (*VacancyCategory, error) {
	vc, err := c.s.TxGetVacancyCategory(ctx, c.s.NoTx(), categoryID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	return &VacancyCategory{
		ID:      vc.ID,
		Title:   vc.Title,
		IconURL: vc.IconURL,
		Rating:  vc.Rating,
	}, nil
}

func (c *Controller) GetVacanciesCategoriesList(ctx context.Context, rating *int32) ([]*VacancyCategory, error) {
	var r int32 = 0
	if rating != nil {
		r = *rating
	}

	vcs, err := c.s.TxGetVacanciesCategoriesList(ctx, c.s.NoTx(), r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvcs := make([]*VacancyCategory, len(vcs))
	for idx, vc := range vcs {
		cvcs[idx] = &VacancyCategory{
			ID:      vc.ID,
			Title:   vc.Title,
			IconURL: vc.IconURL,
			Rating:  vc.Rating,
		}
	}

	return cvcs, nil
}

func (c *Controller) DeleteVacancyCategory(
	ctx context.Context,
	categoryID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteVacancyCategory(ctx, tx, categoryID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrVacancyCategoryNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) PutVacancy(
	ctx context.Context,
	vacancyID *string,
	vacancy *VacancyDetails,
	categoryIDs []string,
	cityIDs []string,
) (VacancyID, error) {
	var vid VacancyID

	if err := vacancy.validate(); err != nil {
		return vid, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		// Look for vacancy id
		if vacancyID != nil {
			switch _, err := c.s.TxGetVacancyDetails(ctx, tx, *vacancyID); errors.Cause(err) {
			case nil:
				vid = VacancyID(*vacancyID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrVacancyNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			vid = VacancyID(uuid.NewV4().String())
		}

		vacancyType, err := toStorageVacancyType(vacancy.Type)
		if err != nil {
			return errors.WithStack(err)
		}

		v := toStorageVacancyDetails(string(vid), vacancyType, vacancy)

		// Update vacancy
		if err := c.s.TxPutVacancy(ctx, tx, v); err != nil {
			return errors.WithStack(err)
		}

		// Update vacancy categories
		if err := updateVacancyCategories(ctx, tx, c, vid, categoryIDs); err != nil {
			return errors.WithStack(err)
		}

		// Update vacancy cities
		if err := updateVacancyCities(ctx, tx, c, vid, cityIDs); err != nil {
			return errors.WithStack(err)
		}

		// Update vacancy images
		if err := updateVacancyImages(ctx, c, tx, vid, vacancy); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return vid, errors.WithStack(err)
	}

	return vid, nil
}

func updateVacancyCategories(
	ctx context.Context,
	tx pkgtx.Tx,
	c *Controller,
	vid VacancyID,
	categoryIDs []string,
) error {
	if err := c.s.TxDeleteVacancyCategories(ctx, tx, string(vid)); err != nil {
		return errors.WithStack(err)
	}

	if len(categoryIDs) > 0 {
		if err := c.s.TxPutVacancyCategories(ctx, tx, string(vid), categoryIDs); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func updateVacancyCities(
	ctx context.Context,
	tx pkgtx.Tx,
	c *Controller,
	vid VacancyID,
	cityIDs []string,
) error {
	if err := c.s.TxDeleteVacancyCities(ctx, tx, string(vid)); err != nil {
		return errors.WithStack(err)
	}

	if len(cityIDs) > 0 {
		if err := c.s.TxPutVacancyCities(ctx, tx, string(vid), cityIDs); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func updateVacancyImages(
	ctx context.Context,
	c *Controller,
	tx pkgtx.Tx,
	vid VacancyID,
	vacancy *VacancyDetails,
) error {
	gvi, err := c.s.TxGetVacanciesImages(ctx, tx, []string{string(vid)})
	if err != nil {
		return errors.WithStack(err)
	}

	if equal(gvi[string(vid)], vacancy.ImageURLs) {
		return nil
	}

	if err := c.s.TxDeleteVacancyImages(ctx, tx, string(vid)); err != nil {
		return errors.WithStack(err)
	}

	if len(vacancy.ImageURLs) > 0 {
		return errors.WithStack(c.s.TxPutVacancyImages(ctx, tx, string(vid), vacancy.ImageURLs))
	}

	return nil
}

func (c *Controller) GetVacancyDetails(ctx context.Context, vacancyID string) (*VacancyDetails, error) {
	// Get vacancy details from DB
	vd, err := c.s.TxGetVacancyDetails(ctx, c.s.NoTx(), vacancyID)

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	// Get vacancy images from DB
	vi, err := c.s.TxGetVacanciesImages(ctx, c.s.NoTx(), []string{vacancyID})

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyImagesNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	vacancyType, err := fromStorageVacancyType(vd.Type)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &VacancyDetails{
		Vacancy: Vacancy{
			ID:        vd.ID,
			Title:     vd.Title,
			Phone:     vd.Phone,
			MinSalary: vd.MinSalary,
			MaxSalary: vd.MaxSalary,
			ImageURLs: vi[vacancyID],
			CompanyID: vd.CompanyID,
		},
		Description:          vd.Description,
		WorkMonthsExperience: vd.WorkMonthsExperience,
		WorkSchedule:         vd.WorkSchedule,
		LocationLatitude:     vd.LocationLatitude,
		LocationLongitude:    vd.LocationLongitude,
		Type:                 vacancyType,
		Address:              vd.Address,
		CountryCode:          vd.CountryCode,
	}, nil
}

func (c *Controller) GetVacanciesList(
	ctx context.Context,
	categoriesIDs []string,
	cursor *Cursor,
	limit int,
) ([]*Vacancy, *Cursor, error) {
	cursorData, err := toCursorData(cursor)
	if err != nil || (cursorData != nil && !equal(cursorData.CategoriesIDs, categoriesIDs)) {
		return nil, nil, errors.WithStack(ErrInvalidCursor)
	}

	maxLimit := 100
	if limit > maxLimit || limit <= 0 {
		limit = maxLimit
	}

	vcs, storageCursor, err := c.s.TxGetVacanciesList(
		ctx,
		c.s.NoTx(),
		categoriesIDs,
		limit,
		toStorageCursor(cursorData),
	)

	switch err {
	case nil:
	default:
		return nil, nil, errors.WithStack(err)
	}

	controllerCursor, err := toCursor(storageCursor, categoriesIDs)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	vacancyIDs := extractVacancyIDs(vcs)
	vacanciesImagesMap, err := c.s.TxGetVacanciesImages(ctx, c.s.NoTx(), vacancyIDs)

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, nil, errors.WithStack(ErrVacancyImagesNotFound)
	default:
		return nil, nil, errors.WithStack(err)
	}

	controllerVacancies := make([]*Vacancy, len(vcs))

	for idx, v := range vcs {
		controllerVacancies[idx] = &Vacancy{
			ID:        v.ID,
			Title:     v.Title,
			Phone:     v.Phone,
			MinSalary: v.MinSalary,
			MaxSalary: v.MaxSalary,
			ImageURLs: vacanciesImagesMap[v.ID],
			CompanyID: v.CompanyID,
		}
	}

	return controllerVacancies, controllerCursor, nil
}

func extractVacancyIDs(vcs []*storage.Vacancy) []string {
	vacancyIDs := make([]string, len(vcs))
	for idx := range vcs {
		vacancyIDs[idx] = vcs[idx].ID
	}

	return vacancyIDs
}

func (c *Controller) GetVacanciesCategories(
	ctx context.Context,
	vacancyIDs []string,
) ([]*VacancyCategoryShort, error) {
	// Get categories
	vscs, err := c.s.TxGetVacanciesCategories(ctx, c.s.NoTx(), vacancyIDs)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	categories := make([]*VacancyCategoryShort, len(vscs))
	for idx, vacancy := range vscs {
		categories[idx] = &VacancyCategoryShort{
			VacancyID: vacancy.VacancyID,
			ID:        vacancy.ID,
			Title:     vacancy.Title,
		}
	}

	return categories, nil
}

func (c *Controller) GetVacancyCities(
	ctx context.Context,
	vacancyIDs []string,
) ([]*VacancyCity, error) {
	// Get categories
	vcs, err := c.s.TxGetVacancyCities(ctx, c.s.NoTx(), vacancyIDs)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	cities := make([]*VacancyCity, len(vcs))
	for idx, city := range vcs {
		cities[idx] = &VacancyCity{
			VacancyID:   city.VacancyID,
			ID:          city.ID,
			Name:        city.Name,
			CountryCode: city.CountryCode,
			Rating:      city.Rating,
		}
	}

	return cities, nil
}

// Custom section
func (c *Controller) PutCustomSections(
	ctx context.Context,
	CVID string,
	education *CVCustomSection,
) (CustomSectionID, error) {
	var ID CustomSectionID

	//if err := education.validate(); err != nil {
	//	return ID, errors.WithStack(err)
	//}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &education.ID != nil {
			switch _, err := c.s.TxGetStories(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = CustomSectionID(education.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrCustomSectionsNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = CustomSectionID(uuid.NewV4().String())
		}

		customSection := storage.CVCustomSection{
			ID:          string(ID),
			Description: education.Description,
		}

		return errors.WithStack(c.s.TxPutCustomSections(ctx, tx, CVID, &customSection))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetCustomSections(
	ctx context.Context,
	CVID string,
) ([]*CVCustomSection, error) {
	// Get categories
	se, err := c.s.TxGetCustomSections(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	sections := make([]*CVCustomSection, len(se))
	for idx, section := range se {
		sections[idx] = &CVCustomSection{
			ID:          section.ID,
			Description: section.Description,
		}
	}

	return sections, nil
}

func (c *Controller) DeleteCustomSection(
	ctx context.Context,
	sectionID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCustomSection(ctx, tx, sectionID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCustomSectionNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Stories
func (c *Controller) PutStory(
	ctx context.Context,
	CVID string,
	story *CVCustomStory,
) (StoryID, error) {
	var ID StoryID

	if err := story.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &story.ID != nil {
			switch _, err := c.s.TxGetStories(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = StoryID(story.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrStoriesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = StoryID(uuid.NewV4().String())
		}

		customStory := storage.CVCustomStory{
			ID:          string(ID),
			ChapterName: story.ChapterName,
			MediaUrl:    story.MediaUrl,
		}

		return errors.WithStack(c.s.TxPutStory(ctx, tx, CVID, &customStory))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetStories(
	ctx context.Context,
	CVID string,
) ([]*CVCustomStory, error) {
	// Get categories
	se, err := c.s.TxGetStories(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	stories := make([]*CVCustomStory, len(se))
	for idx, story := range se {
		stories[idx] = &CVCustomStory{
			ID:          story.ID,
			ChapterName: story.ChapterName,
			MediaUrl:    story.MediaUrl,
		}
	}

	return stories, nil
}

func (c *Controller) DeleteStory(
	ctx context.Context,
	storyID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteStory(ctx, tx, storyID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrStoriesNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Stories episodes
func (c *Controller) PutStoriesEpisodes(
	ctx context.Context,
	storyEpisode *StoryEpisode,
) (StoriesEpisodeID, error) {
	var ID StoriesEpisodeID

	if err := storyEpisode.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &storyEpisode.ID != nil {
			switch _, err := c.s.TxGetStoriesEpisodesById(ctx, tx, storyEpisode.ID); errors.Cause(err) {
			case nil:
				ID = StoriesEpisodeID(storyEpisode.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrStoriesEpisodesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = StoriesEpisodeID(uuid.NewV4().String())
		}

		svc := storage.StoryEpisode{
			ID:       string(ID),
			StoryID:  storyEpisode.StoryID,
			MediaUrl: storyEpisode.MediaUrl,
		}

		return errors.WithStack(c.s.TxPutStoriesEpisodes(ctx, tx, &svc))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetStoriesEpisodes(
	ctx context.Context,
	CVID string,
) ([]*StoryEpisode, error) {
	// Get categories
	se, err := c.s.TxGetStoriesEpisodes(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	episodes := make([]*StoryEpisode, len(se))
	for idx, episode := range se {
		episodes[idx] = &StoryEpisode{
			StoryID:  episode.StoryID,
			ID:       episode.ID,
			MediaUrl: episode.MediaUrl,
		}
	}

	return episodes, nil
}

func (c *Controller) DeleteStoriesEpisodes(
	ctx context.Context,
	episodeID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteStoriesEpisodes(ctx, tx, episodeID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrStoriesEpisodesNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CV
func (c *Controller) DeleteCV(
	ctx context.Context,
	CVID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCV(ctx, tx, CVID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCVNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Mappings

func toStorageVacancyDetails(vid string, vacancyType storage.VacancyType, vd *VacancyDetails) *storage.VacancyDetails {
	now := time.Now()

	return &storage.VacancyDetails{
		Vacancy: storage.Vacancy{
			ID:        vid,
			Title:     vd.Title,
			Phone:     vd.Phone,
			MinSalary: vd.MinSalary,
			MaxSalary: vd.MaxSalary,
			CompanyID: vd.CompanyID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Description:          vd.Description,
		WorkMonthsExperience: vd.WorkMonthsExperience,
		WorkSchedule:         vd.WorkSchedule,
		LocationLatitude:     vd.LocationLatitude,
		LocationLongitude:    vd.LocationLongitude,
		Type:                 vacancyType,
		Address:              vd.Address,
		CountryCode:          vd.CountryCode,
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

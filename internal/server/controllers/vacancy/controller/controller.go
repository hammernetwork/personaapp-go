package controller

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
	"personaapp/internal/server/controllers/vacancy/storage"
	pkgtx "personaapp/pkg/tx"
	"time"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint TODO: Implement CDN link check
		return true
	})
}

var (
	ErrVacancyNotFound                    = errors.New("vacancy not found")
	ErrVacancyCategoryNotFound            = errors.New("vacancy category not found")
	ErrInvalidCursor                      = errors.New("invalid cursor")
	ErrInvalidVacancyCategory             = errors.New("invalid vacancy category struct")
	ErrInvalidVacancyCategoryTitle        = errors.New("invalid vacancy category title")
	ErrInvalidVacancyCategoryIconURL      = errors.New("invalid vacancy category icon url")
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
	TxGetVacancyCategory(ctx context.Context, tx pkgtx.Tx, categoryID string) (*storage.VacancyCategory, error)
	TxPutVacancyCategory(ctx context.Context, tx pkgtx.Tx, category *storage.VacancyCategory) error
	TxGetVacanciesCategoriesList(ctx context.Context, tx pkgtx.Tx) ([]*storage.VacancyCategory, error)

	TxGetVacanciesCategories(ctx context.Context, tx pkgtx.Tx, vacancyIDs []string) ([]*storage.VacancyCategoryExt, error)
	TxPutVacancyCategories(ctx context.Context, tx pkgtx.Tx, vacancyID string, categoriesIDs []string) error
	TxDeleteVacancyCategories(ctx context.Context, tx pkgtx.Tx, vacancyID string) error

	TxGetVacancyDetails(ctx context.Context, tx pkgtx.Tx, vacancyID string) (*storage.VacancyDetails, error)
	TxPutVacancy(ctx context.Context, tx pkgtx.Tx, vacancy *storage.VacancyDetails) error
	TxGetVacanciesList(
		ctx context.Context,
		tx pkgtx.Tx,
		categoriesIDs []string,
		limit int,
		cursor *storage.Cursor,
	) ([]*storage.Vacancy, *storage.Cursor, error)

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
}

type VacancyID string

// Vacancy models for put
type Vacancy struct {
	ID         string
	Title      string `valid:"stringlength(5|80),required"`
	Phone      string `valid:"phone,required"`
	MinSalary  int32  `valid:"range(0|1000000000),required"`
	MaxSalary  int32  `valid:"range(0|1000000000),required"`
	ImageURL   string `valid:"stringlength(0|255),media_link"`
	CompanyID  string `valid:"required"`
	Categories map[string]string
}

type VacancyDetails struct {
	Vacancy
	Description          string
	WorkMonthsExperience int32   `valid:"range(0|1200),required"`
	WorkSchedule         string  `valid:"stringlength(0|100)"`
	LocationLatitude     float32 `valid:"latitude"`
	LocationLongitude    float32 `valid:"longitude"`
}

// Vacancy models for get
type VacancyExt struct {
	ID         string
	Title      string `valid:"stringlength(5|80),required"`
	Phone      string `valid:"phone,required"`
	MinSalary  int32  `valid:"range(0|1000000000),required"`
	MaxSalary  int32  `valid:"range(0|1000000000),required"`
	ImageURL   string `valid:"stringlength(0|255),media_link"`
	CompanyID  string `valid:"required"`
	Categories map[string]string
}

type VacancyDetailsExt struct {
	VacancyExt
	Description          string
	WorkMonthsExperience int32   `valid:"range(0|1200),required"`
	WorkSchedule         string  `valid:"stringlength(0|100)"`
	LocationLatitude     float32 `valid:"latitude"`
	LocationLongitude    float32 `valid:"longitude"`
}

type Cursor string

func (c Cursor) String() string {
	return string(c)
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
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyCategoryNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &VacancyCategory{
		ID:      vc.ID,
		Title:   vc.Title,
		IconURL: vc.IconURL,
	}, nil
}

func (c *Controller) GetVacanciesCategoriesList(ctx context.Context) ([]*VacancyCategory, error) {
	vcs, err := c.s.TxGetVacanciesCategoriesList(ctx, c.s.NoTx())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cvcs := make([]*VacancyCategory, len(vcs))
	for idx, vc := range vcs {
		cvcs[idx] = &VacancyCategory{
			ID:      vc.ID,
			Title:   vc.Title,
			IconURL: vc.IconURL,
		}
	}

	return cvcs, nil
}

func (c *Controller) PutVacancy(
	ctx context.Context,
	vacancyID *string,
	vacancy *VacancyDetails,
	categories []string,
) (VacancyID, error) {
	var vid VacancyID

	if err := vacancy.validate(); err != nil {
		return vid, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
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

		v := toStorageVacancyDetails(string(vid), vacancy)

		if err := c.s.TxPutVacancy(ctx, tx, v); err != nil {
			return errors.WithStack(err)
		}

		if err := c.s.TxDeleteVacancyCategories(ctx, tx, string(vid)); err != nil {
			return errors.WithStack(err)
		}

		if len(categories) > 0 {
			return errors.WithStack(c.s.TxPutVacancyCategories(ctx, tx, string(vid), categories))
		}

		return nil
	}); err != nil {
		return vid, errors.WithStack(err)
	}

	return vid, nil
}

func (c *Controller) GetVacancyDetails(ctx context.Context, vacancyID string) (*VacancyDetailsExt, error) {
	// Get vacancy details
	vd, err := c.s.TxGetVacancyDetails(ctx, c.s.NoTx(), vacancyID)

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	// Get categories
	vacancyIDs := make([]string, 1)
	vacancyIDs[0] = vacancyID

	vscs, err := c.s.TxGetVacanciesCategories(ctx, c.s.NoTx(), vacancyIDs)

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyCategoryNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	categoryMap := toVacancyCategoriesMap(vscs, vacancyID)

	return &VacancyDetailsExt{
		VacancyExt: VacancyExt{
			ID:         vd.ID,
			Title:      vd.Title,
			Phone:      vd.Phone,
			MinSalary:  vd.MinSalary,
			MaxSalary:  vd.MaxSalary,
			ImageURL:   vd.ImageURL,
			CompanyID:  vd.CompanyID,
			Categories: categoryMap,
		},
		Description:          vd.Description,
		WorkMonthsExperience: vd.WorkMonthsExperience,
		WorkSchedule:         vd.WorkSchedule,
		LocationLatitude:     vd.LocationLatitude,
		LocationLongitude:    vd.LocationLongitude,
	}, nil
}

func toVacancyCategoriesMap(cs []*storage.VacancyCategoryExt, vacancyID string) map[string]string {
	categoryMap := make(map[string]string)

	for _, category := range cs {
		if category.VacancyID == vacancyID {
			categoryMap[category.ID] = category.Title
		}
	}

	return categoryMap
}

func (c *Controller) GetVacanciesList(
	ctx context.Context,
	categoriesIDs []string,
	cursor *Cursor,
	limit int,
) ([]*VacancyExt, *Cursor, error) {
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

	// Get categories by vacancy ids
	vacancyIDs := toVacancyIDs(vcs)
	vscs, err := c.s.TxGetVacanciesCategories(ctx, c.s.NoTx(), vacancyIDs)

	switch err {
	case nil:
	default:
		return nil, nil, errors.WithStack(err)
	}

	controllerVacancies := make([]*VacancyExt, len(vcs))

	for idx, v := range vcs {
		controllerVacancies[idx] = &VacancyExt{
			ID:         v.ID,
			Title:      v.Title,
			Phone:      v.Phone,
			MinSalary:  v.MinSalary,
			MaxSalary:  v.MaxSalary,
			ImageURL:   v.ImageURL,
			CompanyID:  v.CompanyID,
			Categories: toVacancyCategoriesMap(vscs, v.ID),
		}
	}

	return controllerVacancies, controllerCursor, nil
}

func toVacancyIDs(vcs []*storage.Vacancy) []string {
	vacancyIDs := make([]string, len(vcs))
	for idx, v := range vcs {
		vacancyIDs[idx] = v.ID
	}

	return vacancyIDs
}

// Mappings

func toStorageVacancyDetails(vid string, vd *VacancyDetails) *storage.VacancyDetails {
	now := time.Now()

	return &storage.VacancyDetails{
		Vacancy: storage.Vacancy{
			ID:        vid,
			Title:     vd.Title,
			Phone:     vd.Phone,
			MinSalary: vd.MinSalary,
			MaxSalary: vd.MaxSalary,
			ImageURL:  vd.ImageURL,
			CompanyID: vd.CompanyID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Description:          vd.Description,
		WorkMonthsExperience: vd.WorkMonthsExperience,
		WorkSchedule:         vd.WorkSchedule,
		LocationLatitude:     vd.LocationLatitude,
		LocationLongitude:    vd.LocationLongitude,
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

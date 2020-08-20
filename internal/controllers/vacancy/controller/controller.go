package controller

import (
	"context"
	"personaapp/internal/controllers/vacancy/storage"
	pkgtx "personaapp/pkg/tx"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
)

type VacancyType string

const (
	VacancyTypeRemote VacancyType = "remote"
	VacancyTypeNormal VacancyType = "normal"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint:godox // TODO: Implement CDN link check
		return true
	})
}

var (
	ErrVacancyNotFound                    = errors.New("vacancy not found")
	ErrVacancyCategoryNotFound            = errors.New("vacancy category not found")
	ErrVacancyImagesNotFound              = errors.New("vacancy image not found")
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
	TxGetVacanciesCategoriesList(ctx context.Context, tx pkgtx.Tx, rating int32) ([]*storage.VacancyCategory, error)
	TxDeleteVacancyCategory(ctx context.Context, tx pkgtx.Tx, categoryID string) error

	TxGetVacanciesCategories(
		ctx context.Context,
		tx pkgtx.Tx,
		vacancyIDs []string,
	) ([]*storage.VacancyCategoryShort, error)
	TxPutVacancyCategories(ctx context.Context, tx pkgtx.Tx, vacancyID string, categoriesIDs []string) error
	TxDeleteVacancyCategories(ctx context.Context, tx pkgtx.Tx, vacancyID string) error

	TxGetVacanciesImages(ctx context.Context, tx pkgtx.Tx, vacancyIDs []string) (map[string][]string, error)
	TxPutVacancyImages(ctx context.Context, tx pkgtx.Tx, vacancyID string, imageUrls []string) error
	TxDeleteVacancyImages(ctx context.Context, tx pkgtx.Tx, vacancyID string) error

	TxGetVacancyDetails(ctx context.Context, tx pkgtx.Tx, vacancyID string) (*storage.VacancyDetails, error)
	TxPutVacancy(ctx context.Context, tx pkgtx.Tx, vacancy *storage.VacancyDetails) error
	TxGetVacanciesList(
		ctx context.Context,
		tx pkgtx.Tx,
		categoriesIDs []string,
		limit int,
		cursor *storage.Cursor,
	) ([]*storage.Vacancy, *storage.Cursor, error)
	TxDeleteVacancy(ctx context.Context, tx pkgtx.Tx, vacancyID string) error

	TxGetVacancyCities(
		ctx context.Context,
		tx pkgtx.Tx,
		vacancyIDs []string,
	) ([]*storage.VacancyCity, error)
	TxPutVacancyCities(
		ctx context.Context,
		tx pkgtx.Tx,
		vacancyID string,
		cityIDs []string,
	) error
	TxDeleteVacancyCities(ctx context.Context, tx pkgtx.Tx, vacancyID string) error

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

func toStorageVacancyType(vt VacancyType) (storage.VacancyType, error) {
	switch vt {
	case VacancyTypeNormal:
		return storage.VacancyTypeNormal, nil
	case VacancyTypeRemote:
		return storage.VacancyTypeRemote, nil
	default:
		return "", errors.New("wrong vacancy type")
	}
}

func fromStorageVacancyType(vt storage.VacancyType) (VacancyType, error) {
	switch vt {
	case storage.VacancyTypeNormal:
		return VacancyTypeNormal, nil
	case storage.VacancyTypeRemote:
		return VacancyTypeRemote, nil
	default:
		return "", errors.New("wrong vacancy type")
	}
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
			switch _, err := c.s.TxGetVacancyCategory(ctx, tx, *categoryID); errors.Unwrap(err) {
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

	switch err {
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
		switch err := c.s.TxDeleteVacancyCategory(ctx, tx, categoryID); err {
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
			switch _, err := c.s.TxGetVacancyDetails(ctx, tx, *vacancyID); err {
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

	switch err {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrVacancyNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	// Get vacancy images from DB
	vi, err := c.s.TxGetVacanciesImages(ctx, c.s.NoTx(), []string{vacancyID})

	switch err {
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

	switch err {
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

	switch err {
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

	switch err {
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

func (c *Controller) DeleteVacancy(
	ctx context.Context,
	vacancyID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteVacancy(ctx, tx, vacancyID); err {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrVacancyNotFound)
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

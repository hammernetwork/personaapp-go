package controller

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
	"personaapp/internal/controllers/city/storage"
	pkgtx "personaapp/pkg/tx"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint:godox // TODO: Implement CDN link check
		return true
	})
}

var (
	ErrCityNotFound    = errors.New("city not found")
	ErrCitiesNotFound  = errors.New("cities not found")
	ErrInvalidCity     = errors.New("invalid city struct")
	ErrInvalidCityName = errors.New("invalid city name")
)

type Storage interface {
	TxGetCitiesList(
		ctx context.Context,
		tx pkgtx.Tx,
		countryCodes []int32,
		rating int32,
		filter string,
	) (_ []*storage.City, rerr error)
	TxPutCity(ctx context.Context, tx pkgtx.Tx, city *storage.City) error
	TxGetCity(ctx context.Context, tx pkgtx.Tx, cityID string) (*storage.City, error)
	TxDeleteCity(ctx context.Context, tx pkgtx.Tx, cityID string) error

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

type VacancyID string

type VacancyCity struct {
	VacancyID   string
	ID          string
	Name        string
	CountryCode int32
	Rating      int32
}

type CityID string

type City struct {
	ID          string
	Name        string `valid:"stringlength(1|255),required"`
	CountryCode int32
	Rating      int32
}

func (vc *City) validate() error {
	if vc == nil {
		return ErrInvalidCity
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "Name", DefaultError: ErrInvalidCityName},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("city struct is filled with some invalid data")
	}

	return nil
}

func (c *Controller) GetCities(
	ctx context.Context,
	countryCodes []int32,
	rating int32,
	filter string,
) ([]*City, error) {
	// Get cities
	cs, err := c.s.TxGetCitiesList(ctx, c.s.NoTx(), countryCodes, rating, filter)

	switch err {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	cities := make([]*City, len(cs))
	for idx, city := range cs {
		cities[idx] = &City{
			ID:          city.ID,
			Name:        city.Name,
			CountryCode: city.CountryCode,
			Rating:      city.Rating,
		}
	}

	return cities, nil
}

func (c *Controller) PutCity(
	ctx context.Context,
	cityID *string,
	city *City,
) (CityID, error) {
	var ID CityID

	if err := city.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if cityID != nil {
			switch _, err := c.s.TxGetCity(ctx, tx, *cityID); errors.Unwrap(err) {
			case nil:
				ID = CityID(*cityID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrCityNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = CityID(uuid.NewV4().String())
		}

		svc := storage.City{
			ID:          string(ID),
			Name:        city.Name,
			CountryCode: city.CountryCode,
			Rating:      city.Rating,
		}

		return errors.WithStack(c.s.TxPutCity(ctx, tx, &svc))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) DeleteCity(
	ctx context.Context,
	id string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCity(ctx, tx, id); errors.Unwrap(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCityNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

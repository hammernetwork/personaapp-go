package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cityController "personaapp/internal/controllers/city/controller"
	cityapi "personaapp/pkg/grpcapi/city"
)

type CityController interface {
	GetCities(ctx context.Context, countryCodes []int32, rating int32, filter string) ([]*cityController.City, error)
	PutCity(
		ctx context.Context,
		cityID *string,
		category *cityController.City,
	) (cityController.CityID, error)
	DeleteCity(ctx context.Context, cityID string) error
}

func (s *Server) GetCities(
	ctx context.Context,
	req *cityapi.GetCitiesRequest,
) (*cityapi.GetCitiesResponse, error) {
	_, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cs, err := s.cy.GetCities(ctx, []int32{}, req.Rating.GetValue(), req.Filter.GetValue())
	switch errors.Cause(err) {
	case nil:
	case cityController.ErrCitiesNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	cities := make([]*cityapi.City, len(cs))
	for idx, c := range cs {
		cities[idx] = &cityapi.City{
			Id:          c.ID,
			Name:        c.Name,
			CountryCode: c.CountryCode,
			Rating:      c.Rating,
		}
	}

	return &cityapi.GetCitiesResponse{
		Cities: cities,
	}, nil
}

func (s *Server) UpsertCity(
	ctx context.Context,
	req *cityapi.UpsertCityRequest,
) (*cityapi.UpsertCityResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cityID, err := s.cy.PutCity(ctx, getOptionalString(req.Id), &cityController.City{
		Name:        req.Name,
		CountryCode: req.CountryCode,
		Rating:      req.Rating,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cityapi.UpsertCityResponse{
		Id: string(cityID),
	}, nil
}

func (s *Server) DeleteCity(
	ctx context.Context,
	req *cityapi.DeleteCityRequest,
) (*cityapi.DeleteCityResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cy.DeleteCity(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cityController.ErrCitiesNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cityapi.DeleteCityResponse{}, nil
}

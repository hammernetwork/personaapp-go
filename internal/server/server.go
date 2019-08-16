package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"

	"personaapp/internal/server/controller"
	registerController "personaapp/internal/server/controller/register"
	"personaapp/pkg/grpcapi/personaappapi"
)

type Controller interface {
	SetPing(ctx context.Context, sp *controller.SetPing) error
	GetPing(ctx context.Context, key string) (*controller.Ping, error)
}

type RegisterController interface {
	RegisterCompany(ctx context.Context, cp *registerController.Company) error
}

type Server struct {
	c Controller
	rc RegisterController
}

func New(c Controller, rc RegisterController) *Server {
	return &Server{c: c, rc: rc}
}

func (s *Server) SetPing(ctx context.Context, req *personaappapi.SetPingRequest) (*personaappapi.SetPingResponse, error) {
	if err := s.c.SetPing(ctx, &controller.SetPing{
		Key:   req.Key,
		Value: req.Value,
	}); err != nil {
		return nil, errors.WithStack(err)
	}
	return &personaappapi.SetPingResponse{
		Ping: nil, // TODO fill response
	}, nil
}

func (s *Server) GetPing(ctx context.Context, req *personaappapi.GetPingRequest) (*personaappapi.GetPingResponse, error) {
	ping, err := s.c.GetPing(ctx, req.GetKey())
	switch err {
	case nil:
	case controller.ErrNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, errors.WithStack(err)
	}
	createdAt, err := ptypes.TimestampProto(ping.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	updatedAt, err := ptypes.TimestampProto(ping.UpdatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &personaappapi.GetPingResponse{
		Ping: &personaappapi.Ping{
			Key:       ping.Key,
			Value:     ping.Value,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}, nil
}

func (s *Server) RegisterCompany(ctx context.Context, req *personaappapi.RegisterCompanyRequest) (*personaappapi.RegisterCompanyResponse, error) {
	err := s.rc.RegisterCompany(ctx, &registerController.Company{
		Name: 		 req.GetCompanyName(),
		Email:       req.GetEmail(),
		Phone:       req.GetPhone(),
		Password:    req.GetPassword(),
	})

	switch err {
	case nil:
	case registerController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	default:
		return nil, errors.WithStack(err)
	}

	return &personaappapi.RegisterCompanyResponse{}, nil
}

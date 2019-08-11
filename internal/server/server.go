package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"

	"personaapp/internal/server/controller"
	"personaapp/pkg/grpcapi/personaappapi"
)

type Controller interface {
	SetPing(ctx context.Context, sp *controller.SetPing) error
	GetPing(ctx context.Context, key string) (*controller.Ping, error)
}

type Server struct {
	c Controller
}

func New(c Controller) *Server {
	return &Server{c: c}
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

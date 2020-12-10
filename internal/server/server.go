package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	authController "personaapp/internal/controllers/auth/controller"
	apiauth "personaapp/pkg/grpcapi/auth"
)

type Server struct {
	ac AuthController
	cc CompanyController
	vc VacancyController
	cy CityController
	cv CVController
}

func New(ac AuthController, cc CompanyController, vc VacancyController, cy CityController, cv CVController) *Server {
	return &Server{ac: ac, cc: cc, vc: vc, cy: cy, cv: cv}
}

func (s *Server) getAuthClaims(ctx context.Context) (*authController.AuthClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no bearer provided")
	}
	token := md.Get("Bearer")[0] // nolint TODO
	claims, err := s.ac.GetAuthClaims(ctx, token)

	switch errors.Cause(err) {
	case nil:
		return claims, nil
	case authController.ErrUnauthorized:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}
}

func (s *Server) isCompanyAccountType(c *authController.AuthClaims) bool {
	return toServerAccount(c.AccountType) == apiauth.AccountType_ACCOUNT_TYPE_COMPANY
}

func (s *Server) isAdminAccountType(c *authController.AuthClaims) bool {
	return toServerAccount(c.AccountType) == apiauth.AccountType_ACCOUNT_TYPE_ADMIN
}

func getOptionalString(sw *wrappers.StringValue) *string {
	if sw == nil {
		return nil
	}

	return &sw.Value
}

func getOptionalInt32(sw *wrappers.Int32Value) *int32 {
	if sw == nil {
		return nil
	}

	return &sw.Value
}

func fieldViolationStatus(fieldViolation *errdetails.BadRequest_FieldViolation) *status.Status {
	st, err := status.New(codes.InvalidArgument, fieldViolation.Description).WithDetails(fieldViolation)
	if err != nil {
		return status.New(codes.Internal, err.Error())
	}

	return st
}

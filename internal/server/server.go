package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authController "personaapp/internal/server/auth/controller"
	companyController "personaapp/internal/server/company/controller"
	"personaapp/pkg/grpcapi/personaappapi"
)

type AuthController interface {
	Register(ctx context.Context, rd *authController.RegisterData) (*authController.AuthToken, error)
	Login(ctx context.Context, ld *authController.LoginData) (*authController.AuthToken, error)
	Refresh(ctx context.Context, tokenStr string) (*authController.AuthToken, error)
	GetAuthClaims(ctx context.Context, tokenStr string) (*authController.AuthClaims, error)
}

type CompanyController interface {
	Update(ctx context.Context, cd *companyController.CompanyData) error
}

type Server struct {
	ac AuthController
	cc CompanyController
}

func New(ac AuthController, cc CompanyController) *Server {
	return &Server{ac: ac, cc: cc}
}

func toControllerAccount(at personaappapi.AccountType) (authController.AccountType, error) {
	switch at {
	case personaappapi.AccountType_ACCOUNT_TYPE_UNKNOWN:
		return "", errors.New("default unknown account type")
	case personaappapi.AccountType_ACCOUNT_TYPE_COMPANY:
		return authController.AccountTypeCompany, nil
	case personaappapi.AccountType_ACCOUNT_TYPE_PERSONA:
		return authController.AccountTypePersona, nil
	default:
		return "", errors.New("unknown account type")
	}
}

func toServerAccount(at authController.AccountType) personaappapi.AccountType {
	switch at {
	case authController.AccountTypeCompany:
		return personaappapi.AccountType_ACCOUNT_TYPE_COMPANY
	case authController.AccountTypePersona:
		return personaappapi.AccountType_ACCOUNT_TYPE_PERSONA
	default:
		return personaappapi.AccountType_ACCOUNT_TYPE_UNKNOWN
	}
}

func toServerToken(at *authController.AuthToken) (*personaappapi.Token, error) {
	expiresAt, err := ptypes.TimestampProto(at.ExpiresAt)
	if err != nil {
		return nil, errors.New("invalid expires at")
	}

	return &personaappapi.Token{
		Token:       at.Token,
		ExpiresAt:   expiresAt,
		AccountType: toServerAccount(at.AccountType),
	}, nil
}

func (s *Server) Register(
	ctx context.Context,
	req *personaappapi.RegisterRequest,
) (*personaappapi.RegisterResponse, error) {
	cat, err := toControllerAccount(req.GetAccountType())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authToken, err := s.ac.Register(ctx, &authController.RegisterData{
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Account:  cat,
		Password: req.GetPassword(),
	})

	switch errors.Cause(err) {
	case nil:
	case authController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case authController.ErrInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidLogin:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidEmail:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPhone:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidAccount:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPassword:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &personaappapi.RegisterResponse{Token: sat}, nil
}

func (s *Server) Login(ctx context.Context, req *personaappapi.LoginRequest) (*personaappapi.LoginResponse, error) {
	authToken, err := s.ac.Login(ctx, &authController.LoginData{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})

	switch errors.Cause(err) {
	case nil:
	case authController.ErrInvalidLogin:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPassword:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrUnauthorized:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &personaappapi.LoginResponse{Token: sat}, nil
}

func (s *Server) Logout(context.Context, *personaappapi.LogoutRequest) (*personaappapi.LogoutResponse, error) {
	return &personaappapi.LogoutResponse{}, nil
}

func (s *Server) Refresh(
	ctx context.Context,
	req *personaappapi.RefreshRequest,
) (*personaappapi.RefreshResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no bearer provided")
	}

	token := md.Get("Bearer")[0] // TODO

	authToken, err := s.ac.Refresh(ctx, token)

	switch errors.Cause(err) {
	case nil:
	case authController.ErrUnauthorized:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &personaappapi.RefreshResponse{Token: sat}, nil
}

func getOptionalString(sw *wrappers.StringValue) *string {
	if sw == nil {
		return nil
	}
	return &sw.Value
}

func (s *Server) getAuthClaims(ctx context.Context) (*authController.AuthClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no bearer provided")
	}
	token := md.Get("Bearer")[0] // TODO
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

func (s *Server) UpdateCompany(
	ctx context.Context,
	req *personaappapi.UpdateCompanyRequest,
) (*personaappapi.UpdateCompanyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, err
	}
	if toServerAccount(claims.AccountType) != personaappapi.AccountType_ACCOUNT_TYPE_COMPANY {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	activityFields := make([]string, len(req.ActivityFields))
	i := 0
	for k := range req.ActivityFields {
		activityFields[i] = k
		i++
	}

	err = s.cc.Update(ctx, &companyController.CompanyData{
		AuthID:         claims.AccountID,
		ActivityFields: activityFields,
		Title:          getOptionalString(req.Title),
		Description:    getOptionalString(req.Description),
		LogoURL:        getOptionalString(req.LogoUrl),
	})

	switch errors.Cause(err) {
	case nil:
	case companyController.ErrInvalidTitle:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidDescription:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURL:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &personaappapi.UpdateCompanyResponse{}, nil
}

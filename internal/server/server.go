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
	apiauth "personaapp/pkg/grpcapi/auth"
	apicompany "personaapp/pkg/grpcapi/company"
	apientities "personaapp/pkg/grpcapi/entities"
)

type AuthController interface {
	Register(ctx context.Context, rd *authController.RegisterData) (*authController.AuthToken, error)
	Login(ctx context.Context, ld *authController.LoginData) (*authController.AuthToken, error)
	Refresh(ctx context.Context, tokenStr string) (*authController.AuthToken, error)
	GetAuthClaims(ctx context.Context, tokenStr string) (*authController.AuthClaims, error)
}

type CompanyController interface {
	Get(ctx context.Context, companyID string) (*companyController.Company, error)
	Update(ctx context.Context, cd *companyController.CompanyData) error
}

type Server struct {
	ac AuthController
	cc CompanyController
}

func New(ac AuthController, cc CompanyController) *Server {
	return &Server{ac: ac, cc: cc}
}

func toControllerAccount(at apientities.AccountType) (authController.AccountType, error) {
	switch at {
	case apientities.AccountType_ACCOUNT_TYPE_UNKNOWN:
		return "", errors.New("default unknown account type")
	case apientities.AccountType_ACCOUNT_TYPE_COMPANY:
		return authController.AccountTypeCompany, nil
	case apientities.AccountType_ACCOUNT_TYPE_PERSONA:
		return authController.AccountTypePersona, nil
	default:
		return "", errors.New("unknown account type")
	}
}

func toServerAccount(at authController.AccountType) apientities.AccountType {
	switch at {
	case authController.AccountTypeCompany:
		return apientities.AccountType_ACCOUNT_TYPE_COMPANY
	case authController.AccountTypePersona:
		return apientities.AccountType_ACCOUNT_TYPE_PERSONA
	default:
		return apientities.AccountType_ACCOUNT_TYPE_UNKNOWN
	}
}

func toServerToken(at *authController.AuthToken) (*apiauth.Token, error) {
	expiresAt, err := ptypes.TimestampProto(at.ExpiresAt)
	if err != nil {
		return nil, errors.New("invalid expires at")
	}

	return &apiauth.Token{
		Token:       at.Token,
		ExpiresAt:   expiresAt,
		AccountType: toServerAccount(at.AccountType),
	}, nil
}

func (s *Server) Register(
	ctx context.Context,
	req *apiauth.RegisterRequest,
) (*apiauth.RegisterResponse, error) {
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

	return &apiauth.RegisterResponse{Token: sat}, nil
}

func (s *Server) Login(ctx context.Context, req *apiauth.LoginRequest) (*apiauth.LoginResponse, error) {
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

	return &apiauth.LoginResponse{Token: sat}, nil
}

func (s *Server) Logout(context.Context, *apiauth.LogoutRequest) (*apiauth.LogoutResponse, error) {
	return &apiauth.LogoutResponse{}, nil
}

func (s *Server) Refresh(
	ctx context.Context,
	req *apiauth.RefreshRequest,
) (*apiauth.RefreshResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "no bearer provided")
	}

	token := md.Get("Bearer")[0] //nolint TODO

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

	return &apiauth.RefreshResponse{Token: sat}, nil
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

func companyControllerErrorToServerErrors(err error) (apicompany.ErrorCode, error) {
	errorCode := apicompany.ErrorCode_UNKNOWN_ERROR_CODE

	var statusErr error

	switch errors.Cause(err) {
	case nil:
	case companyController.ErrInvalidTitle:
		errorCode = apicompany.ErrorCode_INVALID_TITLE_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidTitleLength:
		errorCode = apicompany.ErrorCode_INVALID_TITLE_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidDescription:
		errorCode = apicompany.ErrorCode_INVALID_DESCRIPTION_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidDescriptionLength:
		errorCode = apicompany.ErrorCode_INVALID_DESCRIPTION_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURL:
		errorCode = apicompany.ErrorCode_INVALID_LOGO_URL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURLFormat:
		errorCode = apicompany.ErrorCode_INVALID_LOGO_URL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURLLength:
		errorCode = apicompany.ErrorCode_INVALID_LOGO_URL_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrCompanyNotFound:
		return errorCode, status.Error(codes.NotFound, err.Error())
	default:
		return errorCode, status.Error(codes.Internal, err.Error())
	}

	return errorCode, statusErr
}

func (s *Server) UpdateCompany(
	ctx context.Context,
	req *apicompany.UpdateCompanyRequest,
) (*apicompany.UpdateCompanyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, err
	}

	if toServerAccount(claims.AccountType) != apientities.AccountType_ACCOUNT_TYPE_COMPANY {
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

	if errorCode, statusErr := companyControllerErrorToServerErrors(err); statusErr != nil {
		return &apicompany.UpdateCompanyResponse{
			Response: &apicompany.UpdateCompanyResponse_ErrorCode{ErrorCode: errorCode},
		}, statusErr
	}

	return &apicompany.UpdateCompanyResponse{}, nil
}

func (s *Server) GetCompany(
	ctx context.Context,
	req *apicompany.GetCompanyRequest,
) (*apicompany.GetCompanyResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, err
	}

	company, err := s.cc.Get(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCompanyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	activityFields := make(map[string]*apicompany.CompanyActivityField, len(company.ActivityFields))
	for _, af := range company.ActivityFields {
		activityFields[af] = &apicompany.CompanyActivityField{}
	}

	return &apicompany.GetCompanyResponse{
		Company: &apicompany.Company{
			Id:             company.AuthID,
			Title:          company.Title,
			Description:    company.Description,
			LogoUrl:        company.LogoURL,
			ActivityFields: activityFields,
		},
	}, nil
}

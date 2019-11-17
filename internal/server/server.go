package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"personaapp/pkg/grpcapi/vacancy"

	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authController "personaapp/internal/server/auth/controller"
	companyController "personaapp/internal/server/company/controller"
	apiauth "personaapp/pkg/grpcapi/auth"
	apicompany "personaapp/pkg/grpcapi/company"
)

type AuthController interface {
	Register(ctx context.Context, rd *authController.RegisterData) (*authController.AuthToken, error)
	Login(ctx context.Context, ld *authController.LoginData) (*authController.AuthToken, error)
	Refresh(ctx context.Context, tokenStr string) (*authController.AuthToken, error)
	GetAuthClaims(ctx context.Context, tokenStr string) (*authController.AuthClaims, error)
	UpdateEmail(ctx context.Context, accountID string, email string) (*authController.AuthToken, error)
}

type CompanyController interface {
	Get(ctx context.Context, companyID string) (*companyController.Company, error)
	Update(ctx context.Context, cd *companyController.CompanyData) error
	UpdateActivityFields(ctx context.Context, companyID string, activityFields []string) error
}

type Server struct {
	ac AuthController
	cc CompanyController
}

func New(ac AuthController, cc CompanyController) *Server {
	return &Server{ac: ac, cc: cc}
}

func toControllerAccount(at apiauth.AccountType) (authController.AccountType, error) {
	switch at {
	case apiauth.AccountType_ACCOUNT_TYPE_UNKNOWN:
		return "", errors.New("default unknown account type")
	case apiauth.AccountType_ACCOUNT_TYPE_COMPANY:
		return authController.AccountTypeCompany, nil
	case apiauth.AccountType_ACCOUNT_TYPE_PERSONA:
		return authController.AccountTypePersona, nil
	default:
		return "", errors.New("unknown account type")
	}
}

func toServerAccount(at authController.AccountType) apiauth.AccountType {
	switch at {
	case authController.AccountTypeCompany:
		return apiauth.AccountType_ACCOUNT_TYPE_COMPANY
	case authController.AccountTypePersona:
		return apiauth.AccountType_ACCOUNT_TYPE_PERSONA
	default:
		return apiauth.AccountType_ACCOUNT_TYPE_UNKNOWN
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

// nolint: funlen
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

	errorCode := apiauth.RegisterResponse_UNKNOWN_ERROR_CODE

	var statusErr error

	switch errors.Cause(err) {
	case nil:
	case authController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case authController.ErrUnauthorized:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case authController.ErrInvalidEmail:
		errorCode = apiauth.RegisterResponse_INVALID_EMAIL
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidEmailFormat:
		errorCode = apiauth.RegisterResponse_INVALID_EMAIL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidEmailLength:
		errorCode = apiauth.RegisterResponse_INVALID_EMAIL_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPhone:
		errorCode = apiauth.RegisterResponse_INVALID_PHONE
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPhoneFormat:
		errorCode = apiauth.RegisterResponse_INVALID_PHONE_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidAccount:
		errorCode = apiauth.RegisterResponse_INVALID_ACCOUNT_TYPE
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPassword:
		errorCode = apiauth.RegisterResponse_INVALID_PASSWORD
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPasswordLength:
		errorCode = apiauth.RegisterResponse_INVALID_PASSWORD_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	if statusErr != nil {
		return &apiauth.RegisterResponse{
			Response: &apiauth.RegisterResponse_ErrorCode_{ErrorCode: errorCode},
		}, statusErr
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.RegisterResponse{
		Response: &apiauth.RegisterResponse_Body_{
			Body: &apiauth.RegisterResponse_Body{Token: sat},
		},
	}, nil
}

func (s *Server) Login(ctx context.Context, req *apiauth.LoginRequest) (*apiauth.LoginResponse, error) {
	authToken, err := s.ac.Login(ctx, &authController.LoginData{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})

	errorCode := apiauth.LoginResponse_UNKNOWN_ERROR_CODE

	var statusErr error

	switch errors.Cause(err) {
	case nil:
	case authController.ErrInvalidLogin:
		errorCode = apiauth.LoginResponse_INVALID_LOGIN
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidLoginLength:
		errorCode = apiauth.LoginResponse_INVALID_LOGIN_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPassword:
		errorCode = apiauth.LoginResponse_INVALID_PASSWORD
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case authController.ErrInvalidPasswordLength:
		errorCode = apiauth.LoginResponse_INVALID_PASSWORD_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	if statusErr != nil {
		return &apiauth.LoginResponse{
			Response: &apiauth.LoginResponse_ErrorCode_{ErrorCode: errorCode},
		}, statusErr
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.LoginResponse{
		Response: &apiauth.LoginResponse_Body_{
			Body: &apiauth.LoginResponse_Body{Token: sat},
		},
	}, nil
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
	case authController.ErrInvalidToken:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	sat, err := toServerToken(authToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.RefreshResponse{
		Response: &apiauth.RefreshResponse_Body_{
			Body: &apiauth.RefreshResponse_Body{Token: sat},
		},
	}, nil
}

func (s *Server) UpdateEmail(
	ctx context.Context,
	req *apiauth.UpdateEmailRequest,
) (*apiauth.UpdateEmailResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, err
	}

	errorCode := apiauth.UpdateEmailResponse_UNKNOWN_ERROR_CODE

	var statusErr error

	token, updateErr := s.ac.UpdateEmail(ctx, claims.AccountID, req.Email)

	switch errors.Cause(updateErr) {
	case nil:
	case authController.ErrInvalidEmailFormat:
		errorCode = apiauth.UpdateEmailResponse_INVALID_EMAIL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, updateErr.Error())
	case authController.ErrInvalidEmailLength:
		errorCode = apiauth.UpdateEmailResponse_INVALID_EMAIL_LENGTH
		statusErr = status.Error(codes.InvalidArgument, updateErr.Error())
	case authController.ErrInvalidEmail:
		errorCode = apiauth.UpdateEmailResponse_INVALID_EMAIL
		statusErr = status.Error(codes.InvalidArgument, updateErr.Error())
	case authController.ErrAuthEntityNotFound:
		statusErr = status.Error(codes.NotFound, updateErr.Error())
	default:
		statusErr = status.Error(codes.Internal, updateErr.Error())
	}

	if statusErr != nil {
		return &apiauth.UpdateEmailResponse{
			Response: &apiauth.UpdateEmailResponse_ErrorCode_{ErrorCode: errorCode},
		}, statusErr
	}

	sat, err := toServerToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.UpdateEmailResponse{
		Response: &apiauth.UpdateEmailResponse_Body_{
			Body: &apiauth.UpdateEmailResponse_Body{Token: sat},
		},
	}, nil
}

func (s *Server) UpdatePhone(context.Context, *apiauth.UpdatePhoneRequest) (*apiauth.UpdatePhoneResponse, error) {
	// nolint: TODO: implement
	return nil, nil
}

func (s *Server) UpdatePassword(
	context.Context,
	*apiauth.UpdatePasswordRequest,
) (*apiauth.UpdatePasswordResponse, error) {
	// nolint: TODO: implement
	return nil, nil
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

// Company
func (s *Server) GetCompaniesActivityFieldsList(
	context.Context,
	*apicompany.GetCompaniesActivityFieldsListRequest,
) (*apicompany.GetCompaniesActivityFieldsListResponse, error) {
	// nolint: TODO: implement
	return nil, nil
}

// nolint: funlen
func (s *Server) UpdateCompany(
	ctx context.Context,
	req *apicompany.UpdateCompanyRequest,
) (*apicompany.UpdateCompanyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, err
	}

	if toServerAccount(claims.AccountType) != apiauth.AccountType_ACCOUNT_TYPE_COMPANY {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cc.Update(ctx, &companyController.CompanyData{
		AuthID:      claims.AccountID,
		Title:       getOptionalString(req.Title),
		Description: getOptionalString(req.Description),
		LogoURL:     getOptionalString(req.LogoUrl),
	})

	var statusErr error

	errorCode := apicompany.UpdateCompanyResponse_UNKNOWN_ERROR_CODE

	switch errors.Cause(err) {
	case nil:
	case companyController.ErrInvalidTitle:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_TITLE_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidTitleLength:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_TITLE_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidDescription:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_DESCRIPTION_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidDescriptionLength:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_DESCRIPTION_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURL:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_LOGO_URL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURLFormat:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_LOGO_URL_FORMAT
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrInvalidLogoURLLength:
		errorCode = apicompany.UpdateCompanyResponse_INVALID_LOGO_URL_LENGTH
		statusErr = status.Error(codes.InvalidArgument, err.Error())
	case companyController.ErrCompanyNotFound:
		errorCode = apicompany.UpdateCompanyResponse_UNKNOWN_ERROR_CODE
		statusErr = status.Error(codes.NotFound, err.Error())
	default:
		errorCode = apicompany.UpdateCompanyResponse_UNKNOWN_ERROR_CODE
		statusErr = status.Error(codes.Internal, err.Error())
	}

	if statusErr != nil {
		return &apicompany.UpdateCompanyResponse{
			Response: &apicompany.UpdateCompanyResponse_ErrorCode_{
				ErrorCode: errorCode,
			},
		}, statusErr
	}

	return &apicompany.UpdateCompanyResponse{}, nil
}

func (s *Server) UpdateCompanyActivityFields(
	ctx context.Context,
	req *apicompany.UpdateCompanyActivityFieldsRequest,
) (*apicompany.UpdateCompanyActivityFieldsResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, err
	}

	if toServerAccount(claims.AccountType) != apiauth.AccountType_ACCOUNT_TYPE_COMPANY {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	activityFields := make([]string, len(req.ActivityFields))
	i := 0

	for k := range req.ActivityFields {
		activityFields[i] = k
		i++
	}

	err = s.cc.UpdateActivityFields(ctx, claims.AccountID, activityFields)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apicompany.UpdateCompanyActivityFieldsResponse{}, nil
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

	activityFields := make(map[string]*apicompany.GetCompanyResponse_CompanyActivityField, len(company.ActivityFields))
	for _, af := range company.ActivityFields {
		activityFields[af] = &apicompany.GetCompanyResponse_CompanyActivityField{}
	}

	return &apicompany.GetCompanyResponse{
		Response: &apicompany.GetCompanyResponse_Body_{
			Body: &apicompany.GetCompanyResponse_Body{
				Company: &apicompany.GetCompanyResponse_Company{
					Id:             company.AuthID,
					Title:          company.Title,
					Description:    company.Description,
					LogoUrl:        company.LogoURL,
					ActivityFields: activityFields,
				},
			},
		},
	}, nil
}

// Vacancy

func (s *Server) GetVacanciesFiltersList(
	context.Context,
	*vacancy.GetVacanciesFiltersListRequest,
) (*vacancy.GetVacanciesFiltersListResponse, error) {
	// nolint: TODO: implement
	return nil, nil
}

func (s *Server) GetVacanciesList(
	context.Context,
	*vacancy.GetVacanciesListRequest,
) (*vacancy.GetVacanciesListResponse, error) {
	// nolint: TODO: implement
	return nil, nil
}

func (s *Server) GetVacancyDetails(
	context.Context,
	*vacancy.GetVacancyDetailsRequest,
) (*vacancy.GetVacancyDetailsResponse, error) {
	// nolint: TODO: implement
	return nil, nil
}

package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	authController "personaapp/internal/controllers/auth/controller"
	companyController "personaapp/internal/controllers/company/controller"
	apiauth "personaapp/pkg/grpcapi/auth"
)

type AuthController interface {
	Register(ctx context.Context, rd *authController.RegisterData) (*authController.AuthToken, error)
	Login(ctx context.Context, ld *authController.LoginData) (*authController.AuthToken, error)
	Refresh(ctx context.Context, tokenStr string) (*authController.AuthToken, error)
	GetAuthClaims(ctx context.Context, tokenStr string) (*authController.AuthClaims, error)
	GetAuth(ctx context.Context, accountID string) (*authController.AuthData, error)
	UpdateEmail(
		ctx context.Context,
		accountID string,
		email string,
		password string,
		ac authController.AccountType,
	) (*authController.AuthToken, error)
	UpdatePhone(ctx context.Context, accountID string, phone string, password string) (*authController.AuthToken, error)
	UpdatePassword(
		ctx context.Context,
		accountID string,
		upd *authController.UpdatePasswordData,
	) (*authController.AuthToken, error)
}

func toControllerAccount(at apiauth.AccountType) (authController.AccountType, error) {
	switch at {
	case apiauth.AccountType_ACCOUNT_TYPE_UNKNOWN:
		return "", errors.New("default unknown account type")
	case apiauth.AccountType_ACCOUNT_TYPE_COMPANY:
		return authController.AccountTypeCompany, nil
	case apiauth.AccountType_ACCOUNT_TYPE_PERSONA:
		return authController.AccountTypePersona, nil
	case apiauth.AccountType_ACCOUNT_TYPE_ADMIN:
		return authController.AccountTypeAdmin, nil
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
	case authController.AccountTypeAdmin:
		return apiauth.AccountType_ACCOUNT_TYPE_ADMIN
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

// nolint:funlen // will rework
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

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(err); causeErr {
	case nil:
	case authController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case authController.ErrUnauthorized:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case authController.ErrInvalidEmail:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidEmailFormat:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidEmailLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidPhone:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Phone", Description: causeErr.Error()}
	case authController.ErrInvalidPhoneFormat:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Phone", Description: causeErr.Error()}
	case authController.ErrInvalidAccount:
		fv = &errdetails.BadRequest_FieldViolation{Field: "AccountType", Description: causeErr.Error()}
	case authController.ErrInvalidPassword:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	case authController.ErrInvalidPasswordLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
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

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(err); causeErr {
	case nil:
	case authController.ErrInvalidLogin:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Login", Description: causeErr.Error()}
	case authController.ErrInvalidLoginLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Login", Description: causeErr.Error()}
	case authController.ErrInvalidPassword:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	case authController.ErrInvalidPasswordLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
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
	case authController.ErrInvalidToken:
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

func (s *Server) GetSelf(
	ctx context.Context,
	req *apiauth.GetSelfRequest,
) (*apiauth.GetSelfResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	self, err := s.ac.GetAuth(ctx, claims.AccountID)
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCompanyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.GetSelfResponse{
		Email:       self.Email,
		Phone:       self.Phone,
		AccountType: toServerAccount(self.Account),
	}, nil
}

// nolint:dupl // will rework
func (s *Server) UpdateEmail(
	ctx context.Context,
	req *apiauth.UpdateEmailRequest,
) (*apiauth.UpdateEmailResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	token, updateErr := s.ac.UpdateEmail(ctx, claims.AccountID, req.Email, req.Password, claims.AccountType)

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(updateErr); causeErr {
	case nil:
	case authController.ErrInvalidEmailFormat:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidEmailLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidEmail:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Email", Description: causeErr.Error()}
	case authController.ErrInvalidPassword:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	case authController.ErrAuthEntityNotFound:
		return nil, status.Error(codes.NotFound, updateErr.Error())
	case authController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, updateErr.Error())
	default:
		return nil, status.Error(codes.Internal, updateErr.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
	}

	sat, err := toServerToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.UpdateEmailResponse{Token: sat}, nil
}

// nolint:dupl // will rework
func (s *Server) UpdatePhone(
	ctx context.Context,
	req *apiauth.UpdatePhoneRequest,
) (*apiauth.UpdatePhoneResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	token, updateErr := s.ac.UpdatePhone(ctx, claims.AccountID, req.Phone, req.Password)

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(updateErr); causeErr {
	case nil:
	case authController.ErrInvalidPhoneFormat:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Phone", Description: causeErr.Error()}
	case authController.ErrInvalidPhoneRequired:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Phone", Description: causeErr.Error()}
	case authController.ErrInvalidPhone:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Phone", Description: causeErr.Error()}
	case authController.ErrAuthEntityNotFound:
		return nil, status.Error(codes.NotFound, updateErr.Error())
	case authController.ErrAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, updateErr.Error())
	default:
		return nil, status.Error(codes.Internal, updateErr.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
	}

	sat, err := toServerToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.UpdatePhoneResponse{Token: sat}, nil
}

// nolint:dupl // will rework
func (s *Server) UpdatePassword(
	ctx context.Context,
	req *apiauth.UpdatePasswordRequest,
) (*apiauth.UpdatePasswordResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	upd := &authController.UpdatePasswordData{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	token, updateErr := s.ac.UpdatePassword(ctx, claims.AccountID, upd)

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(updateErr); causeErr {
	case nil:
	case authController.ErrInvalidOldPassword:
		fv = &errdetails.BadRequest_FieldViolation{Field: "OldPassword", Description: causeErr.Error()}
	case authController.ErrInvalidOldPasswordLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "OldPassword", Description: causeErr.Error()}
	case authController.ErrInvalidOldPasswordNotMatch:
		fv = &errdetails.BadRequest_FieldViolation{Field: "OldPassword", Description: causeErr.Error()}
	case authController.ErrInvalidPassword:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	case authController.ErrInvalidPasswordLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Password", Description: causeErr.Error()}
	case authController.ErrAuthEntityNotFound:
		return nil, status.Error(codes.NotFound, updateErr.Error())
	default:
		return nil, status.Error(codes.Internal, updateErr.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
	}

	sat, err := toServerToken(token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &apiauth.UpdatePasswordResponse{Token: sat}, nil
}

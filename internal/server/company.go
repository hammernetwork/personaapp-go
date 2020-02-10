package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	companyController "personaapp/internal/controllers/company/controller"
	apicompany "personaapp/pkg/grpcapi/company"
)

type CompanyController interface {
	Get(ctx context.Context, companyID string) (*companyController.Company, error)
	GetCompaniesList(ctx context.Context, companyIDs []string) ([]*companyController.Company, error)
	Update(ctx context.Context, cd *companyController.CompanyData) error
	UpdateActivityFields(ctx context.Context, companyID string, activityFields []string) error
}

// Company
func (s *Server) GetCompaniesActivityFieldsList(
	context.Context,
	*apicompany.GetCompaniesActivityFieldsListRequest,
) (*apicompany.GetCompaniesActivityFieldsListResponse, error) {
	// nolint:godox // TODO: implement
	return nil, nil
}

// nolint:funlen // will rework
func (s *Server) UpdateCompany(
	ctx context.Context,
	req *apicompany.UpdateCompanyRequest,
) (*apicompany.UpdateCompanyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isCompanyAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cc.Update(ctx, &companyController.CompanyData{
		ID:          claims.AccountID,
		Title:       getOptionalString(req.Title),
		Description: getOptionalString(req.Description),
		LogoURL:     getOptionalString(req.LogoUrl),
	})

	var fv *errdetails.BadRequest_FieldViolation

	switch causeErr := errors.Cause(err); causeErr {
	case nil:
	case companyController.ErrInvalidTitle:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Title", Description: causeErr.Error()}
	case companyController.ErrInvalidTitleLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Title", Description: causeErr.Error()}
	case companyController.ErrInvalidDescription:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Description", Description: causeErr.Error()}
	case companyController.ErrInvalidDescriptionLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "Description", Description: causeErr.Error()}
	case companyController.ErrInvalidLogoURL:
		fv = &errdetails.BadRequest_FieldViolation{Field: "LogoUrl", Description: causeErr.Error()}
	case companyController.ErrInvalidLogoURLFormat:
		fv = &errdetails.BadRequest_FieldViolation{Field: "LogoUrl", Description: causeErr.Error()}
	case companyController.ErrInvalidLogoURLLength:
		fv = &errdetails.BadRequest_FieldViolation{Field: "LogoUrl", Description: causeErr.Error()}
	case companyController.ErrCompanyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	if fv != nil {
		return nil, fieldViolationStatus(fv).Err()
	}

	return &apicompany.UpdateCompanyResponse{}, nil
}

func (s *Server) UpdateCompanyActivityFields(
	ctx context.Context,
	req *apicompany.UpdateCompanyActivityFieldsRequest,
) (*apicompany.UpdateCompanyActivityFieldsResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isCompanyAccountType(claims) {
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
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
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
		Company: &apicompany.GetCompanyResponse_Company{
			Id:             company.ID,
			Title:          company.Title,
			Description:    company.Description,
			LogoUrl:        company.LogoURL,
			ActivityFields: activityFields,
		},
	}, nil
}

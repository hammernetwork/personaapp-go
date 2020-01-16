package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	companyController "personaapp/internal/server/controllers/company/controller"
	vacancyController "personaapp/internal/server/controllers/vacancy/controller"
	vacancyapi "personaapp/pkg/grpcapi/vacancy"
)

type VacancyController interface {
	GetVacancyCategory(ctx context.Context, categoryID string) (*vacancyController.VacancyCategory, error)
	PutVacancyCategory(
		ctx context.Context,
		categoryID *string,
		category *vacancyController.VacancyCategory,
	) (vacancyController.VacancyCategoryID, error)
	GetVacanciesCategoriesList(ctx context.Context) ([]*vacancyController.VacancyCategory, error)
	PutVacancy(
		ctx context.Context,
		vacancyID *string,
		vacancy *vacancyController.VacancyDetails,
		categories []string,
	) (vacancyController.VacancyID, error)
	GetVacancyDetails(ctx context.Context, vacancyID string) (*vacancyController.VacancyDetailsExt, error)
	GetVacanciesList(
		ctx context.Context,
		categoriesIDs []string,
		cursor *vacancyController.Cursor,
		limit int,
	) ([]*vacancyController.VacancyExt, *vacancyController.Cursor, error)
}

// Vacancy

func (s *Server) GetVacancyCategoriesList(
	ctx context.Context,
	req *vacancyapi.GetVacancyCategoriesListRequest,
) (*vacancyapi.GetVacancyCategoriesListResponse, error) {
	_, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vcs, err := s.vc.GetVacanciesCategoriesList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	svcs := map[string]*vacancyapi.GetVacancyCategoriesListResponse_VacancyCategory{}
	for _, vc := range vcs {
		svcs[vc.ID] = &vacancyapi.GetVacancyCategoriesListResponse_VacancyCategory{
			Id:      vc.ID,
			Title:   vc.Title,
			IconUrl: vc.IconURL,
		}
	}

	return &vacancyapi.GetVacancyCategoriesListResponse{VacancyCategories: svcs}, nil
}

func (s *Server) GetVacancyDetails(
	ctx context.Context,
	req *vacancyapi.GetVacancyDetailsRequest,
) (*vacancyapi.GetVacancyDetailsResponse, error) {
	_, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vd, err := s.vc.GetVacancyDetails(ctx, req.VacancyId)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	cd, err := s.cc.Get(ctx, vd.CompanyID)
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCompanyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &vacancyapi.GetVacancyDetailsResponse{
		Vacancy: toServerVacancy(vd),
		Image:   &vacancyapi.GetVacancyDetailsResponse_VacancyImage{ImageUrl: vd.ImageURL},
		Location: &vacancyapi.GetVacancyDetailsResponse_VacancyLocation{
			Latitude:  vd.LocationLatitude,
			Longitude: vd.LocationLongitude,
		},
		Description: &vacancyapi.GetVacancyDetailsResponse_VacancyDescription{
			Description:          vd.Description,
			WorkMonthsExperience: uint32(vd.WorkMonthsExperience),
			WorkSchedule:         vd.WorkSchedule,
		},
		Company: toServerCompany(cd),
	}, nil
}

// nolint: funlen
func (s *Server) GetVacanciesList(
	ctx context.Context,
	req *vacancyapi.GetVacanciesListRequest,
) (*vacancyapi.GetVacanciesListResponse, error) {
	categoriesIds := make([]string, len(req.CategoriesIds))
	for id := range req.CategoriesIds {
		categoriesIds = append(categoriesIds, id)
	}

	vcs, cursor, err := s.vc.GetVacanciesList(
		ctx,
		categoriesIds,
		toControllerCursor(req.Cursor),
		int(req.Count.GetValue()),
	)

	switch causeErr := errors.Cause(err); causeErr {
	case nil:
	case vacancyController.ErrInvalidCursor:
		fv := &errdetails.BadRequest_FieldViolation{Field: "Cursor", Description: causeErr.Error()}
		return nil, fieldViolationStatus(fv).Err()
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	vacanciesIDs := make([]string, len(vcs))
	vacancies := map[string]*vacancyapi.GetVacanciesListResponse_VacancyDetails{}

	companyIdsMap := make(map[string]bool)

	for idx, v := range vcs {
		vacanciesIDs[idx] = v.ID
		vacancies[v.ID] = &vacancyapi.GetVacanciesListResponse_VacancyDetails{
			Vacancy: &vacancyapi.Vacancy{
				Id:         v.ID,
				Title:      v.Title,
				Phone:      v.Phone,
				MinSalary:  v.MinSalary,
				MaxSalary:  v.MaxSalary,
				CompanyId:  v.CompanyID,
				Currency:   vacancyapi.Currency_CURRENCY_UAH,
				Categories: v.Categories,
			},
			ImageUrl: "",
		}
		companyIdsMap[v.CompanyID] = true
	}

	companyIds := make([]string, 0)
	for companyID := range companyIdsMap {
		companyIds = append(companyIds, companyID)
	}

	companies, err := s.cc.GetCompaniesList(ctx, companyIds)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	companiesMap := make(map[string]*vacancyapi.Company)
	for _, company := range companies {
		companiesMap[company.ID] = &vacancyapi.Company{
			Id:      company.ID,
			Title:   company.Title,
			LogoUrl: company.Description,
		}
	}

	return &vacancyapi.GetVacanciesListResponse{
		VacanciesIds: vacanciesIDs,
		Vacancies:    vacancies,
		Companies:    companiesMap,
		Cursor:       toServerCursor(cursor),
	}, nil
}

// Mappings
func toServerVacancy(vd *vacancyController.VacancyDetailsExt) *vacancyapi.Vacancy {
	return &vacancyapi.Vacancy{
		Id:         vd.ID,
		Title:      vd.Title,
		Phone:      vd.Phone,
		MinSalary:  vd.MinSalary,
		MaxSalary:  vd.MaxSalary,
		CompanyId:  vd.CompanyID,
		Currency:   vacancyapi.Currency_CURRENCY_UAH,
		Categories: vd.Categories,
	}
}

func toServerCompany(cd *companyController.Company) *vacancyapi.GetVacancyDetailsResponse_VacancyCompany {
	afs := make(map[string]*vacancyapi.Empty)

	for _, af := range cd.ActivityFields {
		afs[af] = &vacancyapi.Empty{}
	}

	return &vacancyapi.GetVacancyDetailsResponse_VacancyCompany{
		Company: &vacancyapi.Company{
			Id:      cd.ID,
			Title:   cd.Title,
			LogoUrl: cd.LogoURL,
		},
		Description: &vacancyapi.GetVacancyDetailsResponse_CompanyDescription{
			Description: cd.Description,
		},
	}
}

func toControllerCursor(cursor *wrappers.StringValue) *vacancyController.Cursor {
	if cursor == nil {
		return nil
	}

	vc := vacancyController.Cursor(cursor.Value)

	return &vc
}

func toServerCursor(cursor *vacancyController.Cursor) *wrappers.StringValue {
	if cursor == nil {
		return nil
	}

	return &wrappers.StringValue{Value: cursor.String()}
}
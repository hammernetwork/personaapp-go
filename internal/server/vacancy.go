package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	companyController "personaapp/internal/controllers/company/controller"
	vacancyController "personaapp/internal/controllers/vacancy/controller"
	vacancyapi "personaapp/pkg/grpcapi/vacancy"
)

type VacancyController interface {
	GetVacancyCategory(ctx context.Context, categoryID string) (*vacancyController.VacancyCategory, error)
	GetVacanciesCategoriesList(ctx context.Context, rating *int32) ([]*vacancyController.VacancyCategory, error)
	PutVacancyCategory(
		ctx context.Context,
		categoryID *string,
		category *vacancyController.VacancyCategory,
	) (vacancyController.VacancyCategoryID, error)
	DeleteVacancyCategory(ctx context.Context, categoryID string) error

	PutVacancy(
		ctx context.Context,
		vacancyID *string,
		vacancy *vacancyController.VacancyDetails,
		categories []string,
		cityIDs []string,
	) (vacancyController.VacancyID, error)
	GetVacancyDetails(ctx context.Context, vacancyID string) (*vacancyController.VacancyDetails, error)
	GetVacanciesList(
		ctx context.Context,
		categoriesIDs []string,
		cursor *vacancyController.Cursor,
		limit int,
	) ([]*vacancyController.Vacancy, *vacancyController.Cursor, error)
	GetVacanciesCategories(ctx context.Context, vacancyIDs []string) ([]*vacancyController.VacancyCategoryShort, error)
	GetVacancyCities(ctx context.Context, vacancyIDs []string) ([]*vacancyController.VacancyCity, error)
	DeleteVacancy(ctx context.Context, vacancyID string) error
}

// Vacancy

func (s *Server) GetVacancyCategory(
	ctx context.Context,
	req *vacancyapi.GetVacancyCategoryRequest,
) (*vacancyapi.GetVacancyCategoryResponse, error) {
	_, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vc, err := s.vc.GetVacancyCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &vacancyapi.GetVacancyCategoryResponse{
		Category: &vacancyapi.VacancyCategory{
			Id:      vc.ID,
			Title:   vc.Title,
			IconUrl: vc.IconURL,
		},
	}, nil
}

func (s *Server) GetVacancyCategoriesList(
	ctx context.Context,
	req *vacancyapi.GetVacancyCategoriesListRequest,
) (*vacancyapi.GetVacancyCategoriesListResponse, error) {
	_, err := s.getAuthClaims(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vcs, err := s.vc.GetVacanciesCategoriesList(ctx, getOptionalInt32(req.GetRating()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	svcs := map[string]*vacancyapi.VacancyCategory{}
	for _, vc := range vcs {
		svcs[vc.ID] = &vacancyapi.VacancyCategory{
			Id:      vc.ID,
			Title:   vc.Title,
			IconUrl: vc.IconURL,
			Rating:  vc.Rating,
		}
	}

	return &vacancyapi.GetVacancyCategoriesListResponse{VacancyCategories: svcs}, nil
}

func (s *Server) UpdateVacancyCategory(
	ctx context.Context,
	req *vacancyapi.UpdateVacancyCategoryRequest,
) (*vacancyapi.UpdateVacancyCategoryResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	categoryID, err := s.vc.PutVacancyCategory(ctx, getOptionalString(req.Id), &vacancyController.VacancyCategory{
		Title:   req.Title,
		IconURL: req.IconUrl,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &vacancyapi.UpdateVacancyCategoryResponse{
		Id: string(categoryID),
	}, nil
}

func (s *Server) DeleteVacancyCategory(
	ctx context.Context,
	req *vacancyapi.DeleteVacancyCategoryRequest,
) (*vacancyapi.DeleteVacancyCategoryResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.vc.DeleteVacancyCategory(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyCategoryNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &vacancyapi.DeleteVacancyCategoryResponse{}, nil
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

	// Get companies
	cd, err := s.cc.Get(ctx, vd.CompanyID)
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCompanyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Get vacancy categories
	categoriesMap, vc, err := getVacancyCategoriesFromStorage(ctx, req, s)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Get vacancy categories
	c, err := getVacancyCityFromStorage(ctx, req, s)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &vacancyapi.GetVacancyDetailsResponse{
		Vacancy: toServerVacancy(vd, vc),
		Image:   &vacancyapi.GetVacancyDetailsResponse_VacancyImage{ImageUrls: vd.ImageURLs},
		Location: &vacancyapi.GetVacancyDetailsResponse_VacancyLocation{
			Latitude:  vd.LocationLatitude,
			Longitude: vd.LocationLongitude,
		},
		Description: &vacancyapi.GetVacancyDetailsResponse_VacancyDescription{
			Description:          vd.Description,
			WorkMonthsExperience: uint32(vd.WorkMonthsExperience),
			WorkSchedule:         vd.WorkSchedule,
			Type:                 toServerVacancyType(vd.Type),
			Address:              vd.Address,
			CountryCode:          vd.CountryCode,
		},
		Company:    toServerCompany(cd),
		Categories: categoriesMap,
		City:       c,
	}, nil
}

func getVacancyCategoriesFromStorage(
	ctx context.Context,
	req *vacancyapi.GetVacancyDetailsRequest,
	s *Server,
) (categoriesMap map[string]*vacancyapi.VacancyCategoryShort, vc []string, err error) {
	categories, err := s.vc.GetVacanciesCategories(ctx, []string{req.VacancyId})
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCategoryNotFound:
		return nil, nil, status.Error(codes.NotFound, err.Error())
	}

	categoriesMap = map[string]*vacancyapi.VacancyCategoryShort{}
	vc = make([]string, len(categories))

	for idx, c := range categories {
		vc[idx] = c.Title
		categoriesMap[c.ID] = &vacancyapi.VacancyCategoryShort{
			Title: c.Title,
		}
	}
	return categoriesMap, vc, nil
}

func getVacancyCityFromStorage(
	ctx context.Context,
	req *vacancyapi.GetVacancyDetailsRequest,
	s *Server,
) (*vacancyapi.City, error) {
	cities, err := s.vc.GetVacancyCities(ctx, []string{req.VacancyId})
	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCityNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}
	vacancyCity := cities[0]
	c := &vacancyapi.City{
		Id:          vacancyCity.ID,
		Name:        vacancyCity.Name,
		CountryCode: vacancyCity.CountryCode,
		Rating:      vacancyCity.Rating,
	}
	return c, nil
}

// nolint:funlen // will rework
func (s *Server) GetVacanciesList(
	ctx context.Context,
	req *vacancyapi.GetVacanciesListRequest,
) (*vacancyapi.GetVacanciesListResponse, error) {
	categoriesIds := make([]string, 0, len(req.CategoriesIds))
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
				Id:            v.ID,
				Title:         v.Title,
				Phone:         v.Phone,
				MinSalary:     v.MinSalary,
				MaxSalary:     v.MaxSalary,
				CompanyId:     v.CompanyID,
				Currency:      vacancyapi.Currency_CURRENCY_UAH,
				CategoriesIds: []string{},
			},
			ImageUrls: v.ImageURLs,
		}
		companyIdsMap[v.CompanyID] = true
	}

	// Get vacancy categories
	categories, err := s.vc.GetVacanciesCategories(
		ctx,
		vacanciesIDs,
	)

	switch errors.Cause(err) {
	case nil:
	case companyController.ErrCategoryNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	categoriesMap := map[string]*vacancyapi.VacancyCategoryShort{}

	for _, c := range categories {
		vc := vacancies[c.VacancyID].Vacancy.CategoriesIds
		vca := append(vc, c.Title)
		vacancies[c.VacancyID].Vacancy.CategoriesIds = vca
		categoriesMap[c.ID] = &vacancyapi.VacancyCategoryShort{
			Title: c.Title,
		}
	}

	// Get companies
	companyIds := make([]string, 0, len(companyIdsMap))
	for companyID := range companyIdsMap {
		companyIds = append(companyIds, companyID)
	}

	companies, err := s.cc.GetCompaniesList(ctx, companyIds)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	companiesMap := make(map[string]*vacancyapi.Company, len(companies))
	for _, company := range companies {
		companiesMap[company.ID] = &vacancyapi.Company{
			Id:      company.ID,
			Title:   company.Title,
			LogoUrl: company.LogoURL,
		}
	}

	return &vacancyapi.GetVacanciesListResponse{
		VacanciesIds: vacanciesIDs,
		Vacancies:    vacancies,
		Companies:    companiesMap,
		Cursor:       toServerCursor(cursor),
	}, nil
}

func (s *Server) UpdateVacancy(
	ctx context.Context,
	req *vacancyapi.UpdateVacancyRequest,
) (*vacancyapi.UpdateVacancyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) || !s.isCompanyAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if req.Vacancy.CompanyId != claims.AccountID {
		return nil, status.Error(codes.PermissionDenied, "wrong account")
	}

	vacancyType, err := toControllerVacancyType(req.Description.Type)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vacancyID, err := s.vc.PutVacancy(
		ctx,
		getOptionalString(req.Vacancy.Id),
		&vacancyController.VacancyDetails{
			Vacancy: vacancyController.Vacancy{
				Title:     req.Vacancy.Title,
				Phone:     req.Vacancy.Phone,
				MinSalary: req.Vacancy.MinSalary,
				MaxSalary: req.Vacancy.MaxSalary,
				ImageURLs: req.ImageURLs,
				CompanyID: req.Vacancy.CompanyId,
			},
			Description:          req.Description.Description,
			WorkMonthsExperience: int32(req.Description.WorkMonthsExperience),
			WorkSchedule:         req.Description.WorkSchedule,
			LocationLatitude:     req.Location.Latitude,
			LocationLongitude:    req.Location.Longitude,
			Type:                 vacancyType,
			Address:              req.Description.Address,
			CountryCode:          req.Description.CountryCode,
		},
		req.CategoryIDs,
		req.CityIDs,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &vacancyapi.UpdateVacancyResponse{
		Id: string(vacancyID),
	}, nil
}

func (s *Server) DeleteVacancy(
	ctx context.Context,
	req *vacancyapi.DeleteVacancyRequest,
) (*vacancyapi.DeleteVacancyResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vd, err := s.vc.GetVacancyDetails(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if vd.CompanyID != claims.AccountID {
		return nil, status.Error(codes.PermissionDenied, "wrong account")
	}

	err = s.vc.DeleteVacancy(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &vacancyapi.DeleteVacancyResponse{}, nil
}

// Mappings
func toServerVacancy(vd *vacancyController.VacancyDetails, vc []string) *vacancyapi.Vacancy {
	return &vacancyapi.Vacancy{
		Id:            vd.ID,
		Title:         vd.Title,
		Phone:         vd.Phone,
		MinSalary:     vd.MinSalary,
		MaxSalary:     vd.MaxSalary,
		CompanyId:     vd.CompanyID,
		Currency:      vacancyapi.Currency_CURRENCY_UAH,
		CategoriesIds: vc,
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

func toServerVacancyType(at vacancyController.VacancyType) vacancyapi.VacancyType {
	switch at {
	case vacancyController.VacancyTypeNormal:
		return vacancyapi.VacancyType_VACANCY_TYPE_NORMAL
	case vacancyController.VacancyTypeRemote:
		return vacancyapi.VacancyType_VACANCY_TYPE_REMOTE
	default:
		return vacancyapi.VacancyType_VACANCY_TYPE_UNKNOWN
	}
}

func toControllerVacancyType(at vacancyapi.VacancyType) (vacancyController.VacancyType, error) {
	switch at {
	case vacancyapi.VacancyType_VACANCY_TYPE_UNKNOWN:
		return "", errors.New("default unknown account type")
	case vacancyapi.VacancyType_VACANCY_TYPE_NORMAL:
		return vacancyController.VacancyTypeNormal, nil
	case vacancyapi.VacancyType_VACANCY_TYPE_REMOTE:
		return vacancyController.VacancyTypeRemote, nil
	default:
		return "", errors.New("unknown account type")
	}
}

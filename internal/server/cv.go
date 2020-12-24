package server

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cvController "personaapp/internal/controllers/cv/controller"
	vacancyController "personaapp/internal/controllers/vacancy/controller"
	cvapi "personaapp/pkg/grpcapi/cv"
)

type CVController interface {
	PutJobType(
		ctx context.Context,
		jobType *cvController.JobType,
	) (cvController.JobTypeID, error)
	GetJobTypes(ctx context.Context) ([]*cvController.JobType, error)
	DeleteJobType(
		ctx context.Context,
		jobTypeID string,
	) error

	PutCVJobTypes(
		ctx context.Context,
		cvID string,
		jobTypesIDs []string,
	) error
	GetCVJobTypes(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVJobType, error)
	DeleteCVJobTypes(
		ctx context.Context,
		cvID string,
	) error

	PutJobKind(
		ctx context.Context,
		jobKind *cvController.JobKind,
	) (cvController.JobKindID, error)
	GetJobKinds(
		ctx context.Context,
	) ([]*cvController.JobKind, error)
	DeleteJobKind(
		ctx context.Context,
		jobKindID string,
	) error

	PutCVJobKinds(
		ctx context.Context,
		cvID string,
		jobKindsIDs []string,
	) error
	GetCVJobKinds(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVJobKind, error)
	DeleteCVJobKinds(
		ctx context.Context,
		cvID string,
	) error

	PutExperience(
		ctx context.Context,
		experienceID *string,
		experience *cvController.CVExperience,
	) (cvController.ExperienceID, error)
	GetExperiences(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVExperience, error)
	DeleteExperience(
		ctx context.Context,
		experienceID string,
	) error

	PutEducation(
		ctx context.Context,
		educationID *string,
		education *cvController.CVEducation,
	) (cvController.EducationID, error)
	GetEducations(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVEducation, error)
	DeleteEducation(
		ctx context.Context,
		educationID string,
	) error

	PutCustomSection(
		ctx context.Context,
		sectionID *string,
		customSection *cvController.CVCustomSection,
	) (cvController.CustomSectionID, error)
	GetCustomSections(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVCustomSection, error)
	DeleteCustomSection(
		ctx context.Context,
		sectionID string,
	) error

	PutStory(
		ctx context.Context,
		storyID *string,
		story *cvController.CVCustomStory,
	) (cvController.StoryID, error)
	GetStories(
		ctx context.Context,
		cvID string,
	) ([]*cvController.CVCustomStory, error)
	DeleteStory(
		ctx context.Context,
		storyID string,
	) error

	PutStoriesEpisode(
		ctx context.Context,
		episodeID *string,
		storyEpisode *cvController.StoryEpisode,
	) (cvController.StoriesEpisodeID, error)
	GetStoriesEpisodes(
		ctx context.Context,
		cvID string,
	) ([]*cvController.StoryEpisode, error)
	DeleteStoriesEpisode(
		ctx context.Context,
		episodeID string,
	) error

	PutCV(
		ctx context.Context,
		cvID *string,
		cv *cvController.CV,
	) (cvController.CVID, error)
	GetCV(ctx context.Context, cvID string) (*cvController.CV, error)
	GetCVs(
		ctx context.Context,
		personaID string,
	) ([]*cvController.CVShort, error)
	DeleteCV(
		ctx context.Context,
		cvID string,
	) error
}

// Job Type
func (s *Server) UpsertJobType(
	ctx context.Context,
	req *cvapi.UpsertJobTypeRequest,
) (*cvapi.UpsertJobTypeResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutJobType(ctx, &cvController.JobType{
		Name: req.Name,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertJobTypeResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetJobTypes(
	ctx context.Context,
	_ *cvapi.GetJobTypesRequest,
) (*cvapi.GetJobTypesResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jts, err := s.cv.GetJobTypes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cafs := make([]*cvapi.JobType, 0, len(jts))
	for i, vc := range jts {
		cafs[i] = &cvapi.JobType{
			Id:   vc.ID,
			Name: vc.Name,
		}
	}

	return &cvapi.GetJobTypesResponse{JobType: cafs}, nil
}

func (s *Server) DeleteJobType(
	ctx context.Context,
	req *cvapi.DeleteJobTypeRequest,
) (*cvapi.DeleteJobTypeResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteJobType(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyCategoryNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteJobTypeResponse{}, nil
}

// CV job types
func (s *Server) UpsertCVJobTypes(
	ctx context.Context,
	req *cvapi.UpsertCVJobTypesRequest,
) (*cvapi.UpsertCVJobTypesResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.PutCVJobTypes(ctx, req.CvId, req.JobTypesIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertCVJobTypesResponse{}, nil
}

func (s *Server) GetCVJobTypes(
	ctx context.Context,
	req *cvapi.GetCVJobTypesRequest,
) (*cvapi.GetCVJobTypesResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cvjts, err := s.cv.GetCVJobTypes(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cafs := make([]*cvapi.CVJobType, 0, len(cvjts))
	for i, vc := range cvjts {
		cafs[i] = &cvapi.CVJobType{
			Id:   vc.ID,
			Name: vc.Name,
		}
	}

	return &cvapi.GetCVJobTypesResponse{CvJobType: cafs}, nil
}

func (s *Server) DeleteCVJobTypes(
	ctx context.Context,
	req *cvapi.DeleteCVJobTypesRequest,
) (*cvapi.DeleteCVJobTypesResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteCVJobTypes(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrCVJobTypesNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteCVJobTypesResponse{}, nil
}

// Job Type
func (s *Server) UpsertJobKind(
	ctx context.Context,
	req *cvapi.UpsertJobKindRequest,
) (*cvapi.UpsertJobKindResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutJobKind(ctx, &cvController.JobKind{
		Name: req.Name,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertJobKindResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetJobKinds(
	ctx context.Context,
	_ *cvapi.GetJobKindsRequest,
) (*cvapi.GetJobKindsResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetJobKinds(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cafs := make([]*cvapi.JobKind, 0, len(jks))
	for i, vc := range jks {
		cafs[i] = &cvapi.JobKind{
			Id:   vc.ID,
			Name: vc.Name,
		}
	}

	return &cvapi.GetJobKindsResponse{JobKind: cafs}, nil
}

func (s *Server) DeleteJobKind(
	ctx context.Context,
	req *cvapi.DeleteJobKindRequest,
) (*cvapi.DeleteJobKindResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteJobKind(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case vacancyController.ErrVacancyCategoryNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteJobKindResponse{}, nil
}

// CV job types
func (s *Server) UpsertCVJobKinds(
	ctx context.Context,
	req *cvapi.UpsertCVJobKindsRequest,
) (*cvapi.UpsertCVJobKindsResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.PutCVJobKinds(ctx, req.CvId, req.JobKindsIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertCVJobKindsResponse{}, nil
}

func (s *Server) GetCVJobKinds(
	ctx context.Context,
	req *cvapi.GetCVJobKindsRequest,
) (*cvapi.GetCVJobKindsResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cvjts, err := s.cv.GetCVJobKinds(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cafs := make([]*cvapi.CVJobKind, 0, len(cvjts))
	for i, vc := range cvjts {
		cafs[i] = &cvapi.CVJobKind{
			Id:   vc.ID,
			Name: vc.Name,
		}
	}

	return &cvapi.GetCVJobKindsResponse{CvJobKind: cafs}, nil
}

func (s *Server) DeleteCVJobKinds(
	ctx context.Context,
	req *cvapi.DeleteCVJobKindsRequest,
) (*cvapi.DeleteCVJobKindsResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteCVJobKinds(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrCVJobKindNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteCVJobKindsResponse{}, nil
}

// Experience
func (s *Server) UpsertExperience(
	ctx context.Context,
	req *cvapi.UpsertExperienceRequest,
) (*cvapi.UpsertExperienceResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	dateFrom, err := ptypes.Timestamp(req.DateFrom)
	if err != nil {
		return nil, errors.New("invalid date from")
	}

	dateTill, err := ptypes.Timestamp(req.DateTill)
	if err != nil {
		return nil, errors.New("invalid date till")
	}

	jobTypeID, err := s.cv.PutExperience(ctx, req.CvId, &cvController.CVExperience{
		CompanyName: req.CompanyName,
		DateFrom:    dateFrom,
		DateTill:    dateTill,
		Position:    req.Position,
		Description: req.Description,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertExperienceResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetExperiences(
	ctx context.Context,
	req *cvapi.GetExperiencesRequest,
) (*cvapi.GetExperiencesResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetExperiences(ctx, req.CvId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.CVExperience, 0, len(jks))

	for i, vc := range jks {
		dateFrom, err := ptypes.TimestampProto(vc.DateFrom)
		if err != nil {
			return nil, errors.New("invalid date from")
		}

		dateTill, err := ptypes.TimestampProto(vc.DateTill)
		if err != nil {
			return nil, errors.New("invalid date till")
		}

		cve[i] = &cvapi.CVExperience{
			Id:          vc.ID,
			CompanyName: vc.CompanyName,
			DateFrom:    dateFrom,
			DateTill:    dateTill,
			Position:    vc.Position,
			Description: vc.Description,
		}
	}

	return &cvapi.GetExperiencesResponse{CvExperience: cve}, nil
}

func (s *Server) DeleteExperience(
	ctx context.Context,
	req *cvapi.DeleteExperienceRequest,
) (*cvapi.DeleteExperienceResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteExperience(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrExperienceNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteExperienceResponse{}, nil
}

// Education
func (s *Server) UpsertEducation(
	ctx context.Context,
	req *cvapi.UpsertEducationRequest,
) (*cvapi.UpsertEducationResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	dateFrom, err := ptypes.Timestamp(req.DateFrom)
	if err != nil {
		return nil, errors.New("invalid date from")
	}

	dateTill, err := ptypes.Timestamp(req.DateTill)
	if err != nil {
		return nil, errors.New("invalid date till")
	}

	jobTypeID, err := s.cv.PutEducation(ctx, getOptionalString(req.Id), &cvController.CVEducation{
		ID:          *getOptionalString(req.Id),
		CvID:        req.CvId,
		Institution: req.Institution,
		DateFrom:    dateFrom,
		DateTill:    dateTill,
		Speciality:  req.Speciality,
		Description: req.Description,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertEducationResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetEducations(
	ctx context.Context,
	req *cvapi.GetEducationsRequest,
) (*cvapi.GetEducationsResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	es, err := s.cv.GetEducations(ctx, req.CvId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.CVEducation, 0, len(es))
	for i, vc := range es {
		dateFrom, err := ptypes.TimestampProto(vc.DateFrom)
		if err != nil {
			return nil, errors.New("invalid date from")
		}

		dateTill, err := ptypes.TimestampProto(vc.DateTill)
		if err != nil {
			return nil, errors.New("invalid date till")
		}

		cve[i] = &cvapi.CVEducation{
			Id:          vc.ID,
			Institution: vc.Institution,
			DateFrom:    dateFrom,
			DateTill:    dateTill,
			Speciality:  vc.Speciality,
			Description: vc.Description,
		}
	}

	return &cvapi.GetEducationsResponse{CvEducation: cve}, nil
}

func (s *Server) DeleteEducation(
	ctx context.Context,
	req *cvapi.DeleteEducationRequest,
) (*cvapi.DeleteEducationResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteEducation(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrEducationNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteEducationResponse{}, nil
}

// Custom section
func (s *Server) UpsertCustomSection(
	ctx context.Context,
	req *cvapi.UpsertCustomSectionRequest,
) (*cvapi.UpsertCustomSectionResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutCustomSection(ctx, getOptionalString(req.Id), &cvController.CVCustomSection{
		ID:          *getOptionalString(req.Id),
		CvID:        req.CvId,
		Description: req.Description,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertCustomSectionResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetCustomSections(
	ctx context.Context,
	req *cvapi.GetCustomSectionsRequest,
) (*cvapi.GetCustomSectionsResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetCustomSections(ctx, req.CvId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.CVCustomSection, 0, len(jks))
	for i, vc := range jks {
		cve[i] = &cvapi.CVCustomSection{
			Id:          vc.ID,
			Description: vc.Description,
		}
	}

	return &cvapi.GetCustomSectionsResponse{CvCustomSection: cve}, nil
}

func (s *Server) DeleteCustomSection(
	ctx context.Context,
	req *cvapi.DeleteCustomSectionRequest,
) (*cvapi.DeleteCustomSectionResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteCustomSection(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrCustomSectionNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteCustomSectionResponse{}, nil
}

// Story
func (s *Server) UpsertStory(
	ctx context.Context,
	req *cvapi.UpsertStoryRequest,
) (*cvapi.UpsertStoryResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutStory(ctx, getOptionalString(req.Id), &cvController.CVCustomStory{
		ID:          *getOptionalString(req.Id),
		ChapterName: req.ChapterName,
		MediaURL:    req.MediaUrl,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertStoryResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetStories(
	ctx context.Context,
	req *cvapi.GetStoriesRequest,
) (*cvapi.GetStoriesResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetStories(ctx, req.CvId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.CVCustomStory, 0, len(jks))
	for i, vc := range jks {
		cve[i] = &cvapi.CVCustomStory{
			Id:          vc.ID,
			ChapterName: vc.ChapterName,
			MediaUrl:    vc.MediaURL,
		}
	}

	return &cvapi.GetStoriesResponse{CvCustomStory: cve}, nil
}

func (s *Server) DeleteStory(
	ctx context.Context,
	req *cvapi.DeleteStoryRequest,
) (*cvapi.DeleteStoryResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteStory(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrStoriesNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteStoryResponse{}, nil
}

// Stories episode
func (s *Server) UpsertStoriesEpisode(
	ctx context.Context,
	req *cvapi.UpsertStoriesEpisodeRequest,
) (*cvapi.UpsertStoriesEpisodeResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutStoriesEpisode(ctx, &cvController.StoryEpisode{
		ID:       *getOptionalString(req.Id),
		StoryID:  req.StoryId,
		MediaURL: req.MediaUrl,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertStoriesEpisodeResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetStoriesEpisodes(
	ctx context.Context,
	req *cvapi.GetStoriesEpisodesRequest,
) (*cvapi.GetStoriesEpisodesResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetStoriesEpisodes(ctx, req.CvId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.StoryEpisode, 0, len(jks))
	for i, vc := range jks {
		cve[i] = &cvapi.StoryEpisode{
			Id:       vc.ID,
			StoryId:  vc.StoryID,
			MediaUrl: vc.MediaURL,
		}
	}

	return &cvapi.GetStoriesEpisodesResponse{StoryEpisode: cve}, nil
}

func (s *Server) DeleteStoriesEpisode(
	ctx context.Context,
	req *cvapi.DeleteStoriesEpisodeRequest,
) (*cvapi.DeleteStoriesEpisodeResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteStory(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrStoriesEpisodesNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteStoriesEpisodeResponse{}, nil
}

// CV
func (s *Server) UpsertCV(
	ctx context.Context,
	req *cvapi.UpsertCVRequest,
) (*cvapi.UpsertCVResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jobTypeID, err := s.cv.PutCV(ctx, getOptionalString(req.Id), &cvController.CV{
		//ID:                   *getOptionalString(req.Id),
		PersonaID:            req.PersonaId,
		Position:             req.Position,
		WorkMonthsExperience: req.WorkMonthsExperience,
		MinSalary:            req.MinSalary,
		MaxSalary:            req.MaxSalary,
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.UpsertCVResponse{
		Id: string(jobTypeID),
	}, nil
}

func (s *Server) GetCV(
	ctx context.Context,
	req *cvapi.GetCVRequest,
) (*cvapi.GetCVResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	vc, err := s.cv.GetCV(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cvapi.GetCVResponse{
		Cv: &cvapi.CV{
			Id:                   vc.ID,
			PersonaId:            vc.PersonaID,
			Position:             vc.Position,
			WorkMonthsExperience: vc.WorkMonthsExperience,
			MinSalary:            vc.MinSalary,
			MaxSalary:            vc.MaxSalary,
		},
	}, nil
}

func (s *Server) GetCVs(
	ctx context.Context,
	req *cvapi.GetCVsRequest,
) (*cvapi.GetCVsResponse, error) {
	if _, err := s.getAuthClaims(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	jks, err := s.cv.GetCVs(ctx, req.PersonaId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cve := make([]*cvapi.CVShort, 0, len(jks))
	for i, vc := range jks {
		cve[i] = &cvapi.CVShort{
			Id:                   vc.ID,
			Position:             vc.Position,
			WorkMonthsExperience: vc.WorkMonthsExperience,
			MinSalary:            vc.MinSalary,
			MaxSalary:            vc.MaxSalary,
		}
	}
	return &cvapi.GetCVsResponse{CvShort: cve}, nil
}

func (s *Server) DeleteCV(
	ctx context.Context,
	req *cvapi.DeleteCVRequest,
) (*cvapi.DeleteCVResponse, error) {
	claims, err := s.getAuthClaims(ctx)
	if err != nil || !s.isAdminAccountType(claims) {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = s.cv.DeleteCV(ctx, req.Id)
	switch errors.Cause(err) {
	case nil:
	case cvController.ErrCVNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &cvapi.DeleteCVResponse{}, nil
}

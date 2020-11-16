package controller

import (
	"context"
	"personaapp/internal/controllers/cv/storage"
	pkgtx "personaapp/pkg/tx"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.CustomTypeTagMap.Set("media_link", func(i interface{}, o interface{}) bool {
		// nolint:godox // TODO: Implement CDN link check
		return true
	})
}

var (
	ErrCVNotFound              = errors.New("cv not found")
	ErrStoriesEpisodesNotFound = errors.New("stories episodes not found")
	ErrStoriesNotFound         = errors.New("stories not found")
	ErrExperiencesNotFound     = errors.New("experiences not found")
	ErrJobKindsNotFound        = errors.New("job kinds not found")
	ErrJobTypesNotFound        = errors.New("job types not found")
	ErrCustomSectionsNotFound  = errors.New("custom sections not found")
	ErrCustomSectionNotFound   = errors.New("custom section not found")
	ErrEducationNotFound       = errors.New("education not found")
	ErrExperienceNotFound      = errors.New("experience not found")
	ErrJobKindNotFound         = errors.New("job kind not found")
	ErrJobTypeNotFound         = errors.New("job type not found")
	ErrCVJobKindNotFound       = errors.New("cv job kind not found")
	ErrCVJobTypesNotFound      = errors.New("cv job type not found")

	ErrInvalidStoriesEpisodesMediaURL = errors.New("invalid stories episodes media url")
	ErrInvalidStoryMediaURL           = errors.New("invalid story media url")

	ErrInvalidStoryEpisode = errors.New("invalid story episode struct")
	ErrInvalidStory        = errors.New("invalid story struct")
)

type Storage interface {
	TxPutJobType(ctx context.Context, tx pkgtx.Tx, jobType *storage.JobType) error
	TxGetJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
	) (_ []*storage.JobType, rerr error)
	TxDeleteJobType(ctx context.Context, tx pkgtx.Tx, jobTypeID string) error

	TxPutCVJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
		jobTypesIDs []string,
	) error
	TxGetCVJobTypes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVJobType, error)
	TxDeleteCVJobTypes(ctx context.Context, tx pkgtx.Tx, CVID string) error

	TxPutJobKind(ctx context.Context, tx pkgtx.Tx, jobKind *storage.JobKind) error
	TxGetJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
	) (_ []*storage.JobKind, rerr error)
	TxDeleteJobKind(ctx context.Context, tx pkgtx.Tx, jobKindID string) error

	TxPutCVJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
		jobKindsIDs []string,
	) error
	TxGetCVJobKinds(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVJobKind, error)
	TxDeleteCVJobKinds(ctx context.Context, tx pkgtx.Tx, CVID string) error

	TxPutExperience(ctx context.Context, tx pkgtx.Tx, CVID string, experience *storage.CVExperience) error
	TxGetExperiences(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVExperience, error)
	TxDeleteExperience(ctx context.Context, tx pkgtx.Tx, experienceID string) error

	TxPutEducation(ctx context.Context, tx pkgtx.Tx, CVID string, education *storage.CVEducation) error
	TxGetEducations(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVEducation, error)
	TxDeleteEducation(ctx context.Context, tx pkgtx.Tx, educationID string) error

	TxPutCustomSections(ctx context.Context, tx pkgtx.Tx, CVID string, education *storage.CVCustomSection) error
	TxGetCustomSections(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVCustomSection, error)
	TxDeleteCustomSection(ctx context.Context, tx pkgtx.Tx, sectionID string) error

	TxPutStory(ctx context.Context, tx pkgtx.Tx, CVID string, story *storage.CVCustomStory) error
	TxGetStories(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.CVCustomStory, error)
	TxDeleteStory(ctx context.Context, tx pkgtx.Tx, storyID string) error

	TxPutStoriesEpisodes(ctx context.Context, tx pkgtx.Tx, story *storage.StoryEpisode) error
	TxGetStoriesEpisodesById(
		ctx context.Context,
		tx pkgtx.Tx,
		storyEpisodeID string,
	) ([]*storage.StoryEpisode, error)
	TxGetStoriesEpisodes(
		ctx context.Context,
		tx pkgtx.Tx,
		CVID string,
	) ([]*storage.StoryEpisode, error)
	TxDeleteStoriesEpisodes(ctx context.Context, tx pkgtx.Tx, episodeID string) error

	TxPutCV(ctx context.Context, tx pkgtx.Tx, cv *storage.CV) error
	TxGetCV(ctx context.Context, tx pkgtx.Tx, CVID string) (*storage.CV, error)
	TxGetCVs(
		ctx context.Context,
		tx pkgtx.Tx,
		personaID string,
	) ([]*storage.CVShort, error)
	TxDeleteCV(ctx context.Context, tx pkgtx.Tx, cvID string) error

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

type StoryEpisode struct {
	StoryID  string
	ID       string
	MediaUrl string `valid:"stringlength(10|2000),media_link"`
}

type CVCustomStory struct {
	ID          string
	ChapterName string
	MediaUrl    string `valid:"stringlength(10|2000),media_link"`
}

type CVCustomSection struct {
	ID          string
	Description string
}

type CVEducation struct {
	ID          string
	Institution string
	DateFrom    time.Time
	DateTill    time.Time
	Speciality  string
	Description string
}

type CVExperience struct {
	ID          string
	CompanyName string
	DateFrom    time.Time
	DateTill    time.Time
	Position    string
	Description string
}

type JobKind struct {
	ID        string
	Name      string
}

type CVJobKind struct {
	ID   string
	Name string
}

type JobType struct {
	ID        string
	Name      string
}

type CVJobType struct {
	ID   string
	Name string
}

type CV struct {
	ID                   string
	PersonaID            string
	Position             string
	WorkMonthsExperience int32
	MinSalary            int32
	MaxSalary            int32
}

type CVShort struct {
	ID                   string
	Position             string
	WorkMonthsExperience int32
	MinSalary            int32
	MaxSalary            int32
}

type StoriesEpisodeID string

type StoryID string

type CustomSectionID string

type EducationID string

type ExperienceID string

type JobKindID string

type JobTypeID string

type CVID string

func (vc *StoryEpisode) validate() error {
	if vc == nil {
		return ErrInvalidStoryEpisode
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "MediaUrl", DefaultError: ErrInvalidStoriesEpisodesMediaURL},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("story episode struct is filled with some invalid data")
	}

	return nil
}

func (vc *CVCustomStory) validate() error {
	if vc == nil {
		return ErrInvalidStory
	}

	var fieldErrors = []struct {
		Field        string
		DefaultError error
	}{
		{Field: "MediaUrl", DefaultError: ErrInvalidStoryMediaURL},
	}

	if valid, err := govalidator.ValidateStruct(vc); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.WithStack(fe.DefaultError)
			}
		}

		return errors.New("story struct is filled with some invalid data")
	}

	return nil
}

// Job Type
func (c *Controller) PutJobType(
	ctx context.Context,
	jobType *JobType,
) (JobTypeID, error) {
	var ID JobTypeID

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &jobType.ID != nil {
			switch _, err := c.s.TxGetJobTypes(ctx, tx); errors.Cause(err) {
			case nil:
				ID = JobTypeID(jobType.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrJobTypesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = JobTypeID(uuid.NewV4().String())
		}

		jobType := storage.JobType{
			ID:        string(ID),
			Name:      jobType.Name,
		}

		return errors.WithStack(c.s.TxPutJobType(ctx, tx, &jobType))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetJobTypes(
	ctx context.Context,
) ([]*JobType, error) {
	jt, err := c.s.TxGetJobTypes(ctx, c.s.NoTx())

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	jobTypes := make([]*JobType, len(jt))
	for idx, jobType := range jt {
		jobTypes[idx] = &JobType{
			ID:        jobType.ID,
			Name:      jobType.Name,
		}
	}

	return jobTypes, nil
}

func (c *Controller) DeleteJobType(
	ctx context.Context,
	jobTypeID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteJobType(ctx, tx, jobTypeID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrJobTypeNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CV Job Types
func (c *Controller) PutCVJobTypes(
	ctx context.Context,
	CVID string,
	jobTypesIDs []string,
) error {

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if err := c.s.TxDeleteCVJobTypes(ctx, tx, CVID); err != nil {
			return errors.WithStack(err)
		}

		if len(jobTypesIDs) > 0 {
			if err := c.s.TxPutCVJobTypes(ctx, tx, CVID, jobTypesIDs); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) GetCVJobTypes(
	ctx context.Context,
	CVID string,
) ([]*CVJobType, error) {
	jk, err := c.s.TxGetCVJobTypes(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	jobTypes := make([]*CVJobType, len(jk))
	for idx, jobType := range jk {
		jobTypes[idx] = &CVJobType{
			ID:   jobType.ID,
			Name: jobType.Name,
		}
	}

	return jobTypes, nil
}

func (c *Controller) DeleteCVJobTypes(
	ctx context.Context,
	CVID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCVJobTypes(ctx, tx, CVID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCVJobTypesNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Job Kinds
func (c *Controller) PutJobKind(
	ctx context.Context,
	jobKind *JobKind,
) (JobKindID, error) {
	var ID JobKindID

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &jobKind.ID != nil {
			switch _, err := c.s.TxGetJobKinds(ctx, tx); errors.Cause(err) {
			case nil:
				ID = JobKindID(jobKind.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrJobKindsNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = JobKindID(uuid.NewV4().String())
		}

		jobKind := storage.JobKind{
			ID:        string(ID),
			Name:      jobKind.Name,
		}

		return errors.WithStack(c.s.TxPutJobKind(ctx, tx, &jobKind))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetJobKinds(
	ctx context.Context,
) ([]*JobKind, error) {
	jk, err := c.s.TxGetJobKinds(ctx, c.s.NoTx())

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	jobKinds := make([]*JobKind, len(jk))
	for idx, jobKind := range jk {
		jobKinds[idx] = &JobKind{
			ID:        jobKind.ID,
			Name:      jobKind.Name,
		}
	}

	return jobKinds, nil
}

func (c *Controller) DeleteJobKind(
	ctx context.Context,
	jobKindID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteJobKind(ctx, tx, jobKindID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrJobKindNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CV Job Kinds
func (c *Controller) PutCVJobKinds(
	ctx context.Context,
	CVID string,
	jobKindsIDs []string,
) error {

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if err := c.s.TxDeleteCVJobKinds(ctx, tx, CVID); err != nil {
			return errors.WithStack(err)
		}

		if len(jobKindsIDs) > 0 {
			if err := c.s.TxPutCVJobKinds(ctx, tx, CVID, jobKindsIDs); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) GetCVJobKinds(
	ctx context.Context,
	CVID string,
) ([]*CVJobKind, error) {
	jk, err := c.s.TxGetCVJobKinds(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	cvJobKinds := make([]*CVJobKind, len(jk))
	for idx, cvJobKind := range jk {
		cvJobKinds[idx] = &CVJobKind{
			ID:   cvJobKind.ID,
			Name: cvJobKind.Name,
		}
	}

	return cvJobKinds, nil
}

func (c *Controller) DeleteCVJobKinds(
	ctx context.Context,
	CVID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCVJobKinds(ctx, tx, CVID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCVJobKindNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Experience
func (c *Controller) PutExperience(
	ctx context.Context,
	CVID string,
	experience *CVExperience,
) (ExperienceID, error) {
	var ID ExperienceID

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &experience.ID != nil {
			switch _, err := c.s.TxGetExperiences(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = ExperienceID(experience.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrExperiencesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = ExperienceID(uuid.NewV4().String())
		}

		experience := storage.CVExperience{
			ID:          string(ID),
			CompanyName: experience.CompanyName,
			DateFrom:    experience.DateFrom,
			DateTill:    experience.DateTill,
			Position:    experience.Position,
			Description: experience.Description,
		}

		return errors.WithStack(c.s.TxPutExperience(ctx, tx, CVID, &experience))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetExperiences(
	ctx context.Context,
	CVID string,
) ([]*CVExperience, error) {
	se, err := c.s.TxGetExperiences(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	experiences := make([]*CVExperience, len(se))
	for idx, experience := range se {
		experiences[idx] = &CVExperience{
			ID:          experience.ID,
			CompanyName: experience.CompanyName,
			DateFrom:    experience.DateFrom,
			DateTill:    experience.DateTill,
			Position:    experience.Position,
			Description: experience.Description,
		}
	}

	return experiences, nil
}

func (c *Controller) DeleteExperience(
	ctx context.Context,
	experienceID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteExperience(ctx, tx, experienceID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrExperienceNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Education
func (c *Controller) PutEducation(
	ctx context.Context,
	CVID string,
	education *CVEducation,
) (EducationID, error) {
	var ID EducationID

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &education.ID != nil {
			switch _, err := c.s.TxGetStories(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = EducationID(education.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrCustomSectionsNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = EducationID(uuid.NewV4().String())
		}

		education := storage.CVEducation{
			ID:          string(ID),
			Institution: education.Institution,
			DateFrom:    education.DateFrom,
			DateTill:    education.DateTill,
			Speciality:  education.Speciality,
			Description: education.Description,
		}

		return errors.WithStack(c.s.TxPutEducation(ctx, tx, CVID, &education))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetEducations(
	ctx context.Context,
	CVID string,
) ([]*CVEducation, error) {
	se, err := c.s.TxGetEducations(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	educations := make([]*CVEducation, len(se))
	for idx, education := range se {
		educations[idx] = &CVEducation{
			ID:          education.ID,
			Institution: education.Institution,
			DateFrom:    education.DateFrom,
			DateTill:    education.DateTill,
			Speciality:  education.Speciality,
			Description: education.Description,
		}
	}

	return educations, nil
}

func (c *Controller) DeleteEducation(
	ctx context.Context,
	educationID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteEducation(ctx, tx, educationID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrEducationNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Custom section
func (c *Controller) PutCustomSections(
	ctx context.Context,
	CVID string,
	customSection *CVCustomSection,
) (CustomSectionID, error) {
	var ID CustomSectionID

	//if err := education.validate(); err != nil {
	//	return ID, errors.WithStack(err)
	//}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &customSection.ID != nil {
			switch _, err := c.s.TxGetStories(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = CustomSectionID(customSection.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrCustomSectionsNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = CustomSectionID(uuid.NewV4().String())
		}

		customSection := storage.CVCustomSection{
			ID:          string(ID),
			Description: customSection.Description,
		}

		return errors.WithStack(c.s.TxPutCustomSections(ctx, tx, CVID, &customSection))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetCustomSections(
	ctx context.Context,
	CVID string,
) ([]*CVCustomSection, error) {
	// Get categories
	se, err := c.s.TxGetCustomSections(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	sections := make([]*CVCustomSection, len(se))
	for idx, section := range se {
		sections[idx] = &CVCustomSection{
			ID:          section.ID,
			Description: section.Description,
		}
	}

	return sections, nil
}

func (c *Controller) DeleteCustomSection(
	ctx context.Context,
	sectionID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCustomSection(ctx, tx, sectionID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCustomSectionNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Stories
func (c *Controller) PutStory(
	ctx context.Context,
	CVID string,
	story *CVCustomStory,
) (StoryID, error) {
	var ID StoryID

	if err := story.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &story.ID != nil {
			switch _, err := c.s.TxGetStories(ctx, tx, CVID); errors.Cause(err) {
			case nil:
				ID = StoryID(story.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrStoriesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = StoryID(uuid.NewV4().String())
		}

		customStory := storage.CVCustomStory{
			ID:          string(ID),
			ChapterName: story.ChapterName,
			MediaUrl:    story.MediaUrl,
		}

		return errors.WithStack(c.s.TxPutStory(ctx, tx, CVID, &customStory))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetStories(
	ctx context.Context,
	CVID string,
) ([]*CVCustomStory, error) {
	// Get categories
	se, err := c.s.TxGetStories(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	stories := make([]*CVCustomStory, len(se))
	for idx, story := range se {
		stories[idx] = &CVCustomStory{
			ID:          story.ID,
			ChapterName: story.ChapterName,
			MediaUrl:    story.MediaUrl,
		}
	}

	return stories, nil
}

func (c *Controller) DeleteStory(
	ctx context.Context,
	storyID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteStory(ctx, tx, storyID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrStoriesNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Stories episodes
func (c *Controller) PutStoriesEpisodes(
	ctx context.Context,
	storyEpisode *StoryEpisode,
) (StoriesEpisodeID, error) {
	var ID StoriesEpisodeID

	if err := storyEpisode.validate(); err != nil {
		return ID, errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &storyEpisode.ID != nil {
			switch _, err := c.s.TxGetStoriesEpisodesById(ctx, tx, storyEpisode.ID); errors.Cause(err) {
			case nil:
				ID = StoriesEpisodeID(storyEpisode.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrStoriesEpisodesNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = StoriesEpisodeID(uuid.NewV4().String())
		}

		svc := storage.StoryEpisode{
			ID:       string(ID),
			StoryID:  storyEpisode.StoryID,
			MediaUrl: storyEpisode.MediaUrl,
		}

		return errors.WithStack(c.s.TxPutStoriesEpisodes(ctx, tx, &svc))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetStoriesEpisodes(
	ctx context.Context,
	CVID string,
) ([]*StoryEpisode, error) {
	// Get categories
	se, err := c.s.TxGetStoriesEpisodes(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	episodes := make([]*StoryEpisode, len(se))
	for idx, episode := range se {
		episodes[idx] = &StoryEpisode{
			StoryID:  episode.StoryID,
			ID:       episode.ID,
			MediaUrl: episode.MediaUrl,
		}
	}

	return episodes, nil
}

func (c *Controller) DeleteStoriesEpisodes(
	ctx context.Context,
	episodeID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteStoriesEpisodes(ctx, tx, episodeID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrStoriesEpisodesNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CV
func (c *Controller) PutCV(
	ctx context.Context,
	cv *CV,
) (CVID, error) {
	var ID CVID

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		if &cv.ID != nil {
			switch _, err := c.s.TxGetCV(ctx, tx, cv.ID); errors.Cause(err) {
			case nil:
				ID = CVID(cv.ID)
			case storage.ErrNotFound:
				return errors.WithStack(ErrCVNotFound)
			default:
				return errors.WithStack(err)
			}
		} else {
			ID = CVID(uuid.NewV4().String())
		}

		cv := storage.CV{
			ID:                   string(ID),
			PersonaID:            cv.PersonaID,
			Position:             cv.Position,
			WorkMonthsExperience: cv.WorkMonthsExperience,
			MinSalary:            cv.MinSalary,
			MaxSalary:            cv.MaxSalary,
		}

		return errors.WithStack(c.s.TxPutCV(ctx, tx, &cv))
	}); err != nil {
		return ID, errors.WithStack(err)
	}

	return ID, nil
}

func (c *Controller) GetCV(ctx context.Context, CVID string) (*CV, error) {
	// Get vacancy details from DB
	cv, err := c.s.TxGetCV(ctx, c.s.NoTx(), CVID)

	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrCVNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	return &CV{
		ID:                   cv.ID,
		PersonaID:            cv.PersonaID,
		Position:             cv.Position,
		WorkMonthsExperience: cv.WorkMonthsExperience,
		MinSalary:            cv.MinSalary,
		MaxSalary:            cv.MaxSalary,
	}, nil
}

func (c *Controller) GetCVs(
	ctx context.Context,
	personaID string,
) ([]*CVShort, error) {
	// Get cvs
	cvs, err := c.s.TxGetCVs(ctx, c.s.NoTx(), personaID)

	switch errors.Cause(err) {
	case nil:
	default:
		return nil, errors.WithStack(err)
	}

	cvShorts := make([]*CVShort, len(cvs))
	for idx, cvShort := range cvs {
		cvShorts[idx] = &CVShort{
			ID:                   cvShort.ID,
			Position:             cvShort.Position,
			WorkMonthsExperience: cvShort.WorkMonthsExperience,
			MinSalary:            cvShort.MinSalary,
			MaxSalary:            cvShort.MaxSalary,
		}
	}

	return cvShorts, nil
}

func (c *Controller) DeleteCV(
	ctx context.Context,
	CVID string,
) error {
	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		switch err := c.s.TxDeleteCV(ctx, tx, CVID); errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrCVNotFound)
		default:
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

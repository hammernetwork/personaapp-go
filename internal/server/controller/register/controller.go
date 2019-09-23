package register

import (
	"context"
	"github.com/asaskevich/govalidator"
	storage "personaapp/internal/server/storage/register"
	pkgtx "personaapp/pkg/tx"
	"github.com/cockroachdb/errors"
	_ "personaapp/internal/validator"
	"regexp"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

var ErrCompanyAlreadyExists = errors.New("company with this email or phone already exists")
var ErrCompanyNameInvalid = errors.New("company name is invalid")
var ErrCompanyEmailInvalid = errors.New("company email is invalid")
var ErrCompanyPhoneInvalid = errors.New("company phone is invalid")
var ErrCompanyPasswordInvalid = errors.New("company password is invalid")

var ErrPersonaAlreadyExists = errors.New("persona with this email or phone already exists")
var ErrPersonaFirstNameInvalid = errors.New("persona first name is invalid")
var ErrPersonaLastNameInvalid = errors.New("persona last name is invalid")
var ErrPersonaEmailInvalid = errors.New("persona email is invalid")
var ErrPersonaPhoneInvalid = errors.New("persona phone is invalid")
var ErrPersonaPasswordInvalid = errors.New("persona password is invalid")

type Company struct {
	Name string `valid:"stringlength(2|100),required"`
	Email string `valid:"stringlength(5|255),email,required"`
	Phone string `valid:"phone,required"`
	Password string `valid:"stringlength(6|30),required"`
}

type Persona struct {
	FirstName string `valid:"stringlength(2|50),alpha,required"`
	LastName string `valid:"stringlength(2|50),alpha,required"`
	Email string `valid:"stringlength(5|255),email"`
	Phone string `valid:"phone,required"`
	Password string `valid:"stringlength(6|30),required"`
}

type Storage interface {
	TxGetCompanyByEmailOrPhone(ctx context.Context, tx pkgtx.Tx, phone string, email string) (*storage.Company, error)
	TxPutCompany(ctx context.Context, tx pkgtx.Tx, cp *storage.Company) error

	TxGetPersonaByPhone(ctx context.Context, tx pkgtx.Tx, phone string) (*storage.Persona, error)
	TxGetPersonaByEmailOrPhone(ctx context.Context, tx pkgtx.Tx, phone string, email string) (*storage.Persona, error)
	TxPutPersona(ctx context.Context, tx pkgtx.Tx, cp *storage.Persona) error

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	s Storage
}

func New(s Storage) *Controller {
	return &Controller{s: s}
}

func (c *Controller) RegisterCompany(ctx context.Context, cp *Company) error {
	r := regexp.MustCompile(`\s+`)
	company := Company{
		Name: r.ReplaceAllString(strings.TrimSpace(cp.Name), " "),
		Email: r.ReplaceAllString(strings.TrimSpace(cp.Email), " "),
		Phone: r.ReplaceAllString(strings.TrimSpace(cp.Phone), " "),
		Password: cp.Password,
	}

	if valid, err := govalidator.ValidateStruct(company); !valid {
		if govalidator.ErrorByField(err, "Name") != "" {
			return errors.WithStack(ErrCompanyNameInvalid)
		} else if govalidator.ErrorByField(err, "Email") != "" {
			return errors.WithStack(ErrCompanyEmailInvalid)
		} else if govalidator.ErrorByField(err, "Phone") != "" {
			return errors.WithStack(ErrCompanyPhoneInvalid)
		} else if govalidator.ErrorByField(err, "Password") != "" {
			return errors.WithStack(ErrCompanyPasswordInvalid)
		}
	}

	password, passwordErr := bcrypt.GenerateFromPassword([]byte(company.Password), bcrypt.DefaultCost)
	if passwordErr != nil {
		return errors.WithStack(passwordErr)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(i context.Context, tx pkgtx.Tx) error {
		switch existingCompany, err := c.s.TxGetCompanyByEmailOrPhone(ctx, tx, company.Phone, company.Email); err {
		case nil:
			if existingCompany != nil {
				return errors.WithStack(ErrCompanyAlreadyExists)
			}
		case storage.ErrNotFound:
		default:
			return errors.WithStack(err)
		}

		return errors.WithStack(c.s.TxPutCompany(ctx, tx, &storage.Company{
			Name:     company.Name,
			Email:    company.Email,
			Phone:    company.Phone,
			Password: string(password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Controller) RegisterPersona(ctx context.Context, p *Persona) error {
	r := regexp.MustCompile(`\s+`)
	persona := Persona{
		FirstName: r.ReplaceAllString(strings.TrimSpace(p.FirstName), " "),
		LastName: r.ReplaceAllString(strings.TrimSpace(p.LastName), " "),
		Email: r.ReplaceAllString(strings.TrimSpace(p.Email), " "),
		Phone: r.ReplaceAllString(strings.TrimSpace(p.Phone), " "),
		Password: p.Password,
	}

	if valid, err := govalidator.ValidateStruct(persona); !valid {
		if govalidator.ErrorByField(err, "FirstName") != "" {
			return errors.WithStack(ErrPersonaFirstNameInvalid)
		} else if govalidator.ErrorByField(err, "LastName") != "" {
			return errors.WithStack(ErrPersonaLastNameInvalid)
		} else if govalidator.ErrorByField(err, "Email") != "" {
			return errors.WithStack(ErrPersonaEmailInvalid)
		} else if govalidator.ErrorByField(err, "Phone") != "" {
			return errors.WithStack(ErrPersonaPhoneInvalid)
		} else if govalidator.ErrorByField(err, "Password") != "" {
			return errors.WithStack(ErrPersonaPasswordInvalid)
		}
	}

	password, passwordErr := bcrypt.GenerateFromPassword([]byte(persona.Password), bcrypt.DefaultCost)
	if passwordErr != nil {
		return errors.WithStack(passwordErr)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(i context.Context, tx pkgtx.Tx) error {
		var existingPersona *storage.Persona
		var err error

		emailIsEmpty := len(persona.Email) == 0
		if emailIsEmpty {
			existingPersona, err = c.s.TxGetPersonaByPhone(ctx, tx, persona.Phone)
		} else {
			existingPersona, err = c.s.TxGetPersonaByEmailOrPhone(ctx, tx, persona.Phone, persona.Email)
		}

		switch err {
		case nil:
			if existingPersona != nil {
				return errors.WithStack(ErrPersonaAlreadyExists)
			}
		case storage.ErrNotFound:
		default:
			return errors.WithStack(err)
		}

		var email *string
		if !emailIsEmpty {
			email = &persona.Email
		}

		return errors.WithStack(c.s.TxPutPersona(ctx, tx, &storage.Persona{
			FirstName:  persona.FirstName,
			LastName:   persona.LastName,
			Email:      email,
			Phone:      persona.Phone,
			Password:   string(password),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
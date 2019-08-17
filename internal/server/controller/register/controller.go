package register

import (
	"context"
	storage "personaapp/internal/server/storage/register"
	"personaapp/internal/validator"
	pkgtx "personaapp/pkg/tx"
	"github.com/cockroachdb/errors"
	"regexp"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

var ErrAlreadyExists = errors.New("company with this email or phone already exists")
var ErrCompanyNameInvalid = errors.New("company name is invalid")
var ErrCompanyEmailInvalid = errors.New("company email is invalid")
var ErrCompanyPhoneInvalid = errors.New("company phone is invalid")
var ErrCompanyPasswordInvalid = errors.New("company password is invalid")

type Company struct {
	Name string
	Email string
	Phone string
	Password string
}

type Storage interface {
	TxCheckCompanyIsUnique(ctx context.Context, tx pkgtx.Tx, name string, email string) (bool, error)
	TxCreateCompany(ctx context.Context, tx pkgtx.Tx, cp *storage.Company) error

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

	if err := validator.ValidateCompany(validator.Company{
		Name:     company.Name,
		Email:    company.Email,
		Phone:    company.Phone,
		Password: company.Password,
	}); err != nil {
		switch err {
		case validator.ErrCompanyNameInvalid:
			return errors.WithStack(ErrCompanyNameInvalid)
		case validator.ErrCompanyEmailInvalid:
			return errors.WithStack(ErrCompanyEmailInvalid)
		case validator.ErrCompanyPhoneInvalid:
			return errors.WithStack(ErrCompanyPhoneInvalid)
		case validator.ErrCompanyPasswordInvalid:
			return errors.WithStack(ErrCompanyPasswordInvalid)
		}
	}

	password, passwordErr := bcrypt.GenerateFromPassword([]byte(company.Password), bcrypt.DefaultCost)
	if passwordErr != nil {
		return errors.WithStack(passwordErr)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(i context.Context, tx pkgtx.Tx) error {
		switch exists, err := c.s.TxCheckCompanyIsUnique(ctx, tx, company.Name, company.Email); err {
		case nil:
			if exists {
				return errors.WithStack(ErrAlreadyExists)
			}
		default:
			return errors.WithStack(err)
		}

		return errors.WithStack(c.s.TxCreateCompany(ctx, tx, &storage.Company{
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
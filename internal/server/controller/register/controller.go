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

var ErrAlreadyExists = errors.New("company with this email or phone already exists")
var ErrCompanyNameInvalid = errors.New("company name is invalid")
var ErrCompanyEmailInvalid = errors.New("company email is invalid")
var ErrCompanyPhoneInvalid = errors.New("company phone is invalid")
var ErrCompanyPasswordInvalid = errors.New("company password is invalid")

type Company struct {
	Name string `valid:"stringlength(2|100),required"`
	Email string `valid:"stringlength(5|255),email,required"`
	Phone string `valid:"phone,required"`
	Password string `valid:"stringlength(6|30),required"`
}

type Storage interface {
	TxGetCompanyByEmailOrPhone(ctx context.Context, tx pkgtx.Tx, phone string, email string) (*storage.Company, error)
	TxPutCompany(ctx context.Context, tx pkgtx.Tx, cp *storage.Company) error

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
		switch existingCompany, err := c.s.TxGetCompanyByEmailOrPhone(ctx, tx, company.Phone, company.Name); err {
		case nil:
			if existingCompany != nil {
				return errors.WithStack(ErrAlreadyExists)
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
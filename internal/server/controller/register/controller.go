package register

import (
	"context"
	storage "personaapp/internal/server/storage/register"
	"personaapp/internal/validator"
	pkgtx "personaapp/pkg/tx"
	"github.com/cockroachdb/errors"
	"time"
)

var ErrAlreadyExists = errors.New("company with this email or phone already exists")

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
	if err := validator.ValidateCompany(validator.Company{
		Name:     cp.Name,
		Email:    cp.Email,
		Phone:    cp.Phone,
		Password: cp.Password,
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := pkgtx.RunInTx(ctx, c.s, func(i context.Context, tx pkgtx.Tx) error {
		switch exists, err := c.s.TxCheckCompanyIsUnique(ctx, tx, cp.Name, cp.Email); err {
		case nil:
			if exists {
				return errors.WithStack(ErrAlreadyExists)
			}
		default:
			return errors.WithStack(err)
		}

		return errors.WithStack(c.s.TxCreateCompany(ctx, tx, &storage.Company{
			Name:     cp.Name,
			Email:    cp.Email,
			Phone:    cp.Phone,
			//TODO: wrap in bcrypt hash function
			Password: cp.Password,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}))
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
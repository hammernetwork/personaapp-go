package controller

import (
	"context"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/bcrypt"

	authStorage "personaapp/internal/server/auth/storage"
	pkgtx "personaapp/pkg/tx"
)

type AccountType string

const (
	AccountTypeCompany AccountType = "company"
	AccountTypePersona AccountType = "persona"
)

func init() {
	govalidator.CustomTypeTagMap.Set("accountType", func(i interface{}, o interface{}) bool {
		if at, ok := i.(AccountType); ok {
			switch at {
			case AccountTypeCompany, AccountTypePersona:
				return true
			}
		}

		return false
	})

	govalidator.CustomTypeTagMap.Set("phone", func(i interface{}, o interface{}) bool {
		phone, ok := i.(string)
		if !ok {
			return false
		}

		r := regexp.MustCompile(`^\+380\d{9}$`)
		return r.MatchString(phone)
	})

	govalidator.CustomTypeTagMap.Set("custom_email", func(i interface{}, o interface{}) bool {
		email, ok := i.(string)
		if !ok {
			return false
		}

		rd, ok := o.(*RegisterData)
		if !ok {
			return false
		}

		if rd.Account == AccountTypePersona && rd.Email == "" {
			return true
		}

		return govalidator.IsEmail(email)
	})
}

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPhone    = errors.New("invalid phone")
	ErrInvalidAccount  = errors.New("invalid account")
	ErrInvalidPassword = errors.New("invalid password")
)

type Config struct {
	TokenExpiration   time.Duration
	PrivateSigningKey string
	TokenValidityGap  time.Duration
}

func (c *Config) Flags(name string) *pflag.FlagSet {
	f := pflag.NewFlagSet(name, pflag.PanicOnError)

	f.DurationVar(&c.TokenExpiration, "token_expiration", 5*time.Minute, "Auth token expiration duration")
	f.StringVar(&c.PrivateSigningKey, "private_signing_key", "", "A private key used for token issuing")
	f.DurationVar(&c.TokenValidityGap, "token_validity_gap", 15*time.Second, "Validity gap duration")

	return f
}

type Storage interface {
	TxPutAuth(ctx context.Context, tx pkgtx.Tx, ad *authStorage.AuthData) error
	TxGetAuthDataByPhoneOrEmail(ctx context.Context, tx pkgtx.Tx, phone, email string) (*authStorage.AuthData, error)
	TxGetAuthDataByPhone(ctx context.Context, tx pkgtx.Tx, phone string) (*authStorage.AuthData, error)

	BeginTx(ctx context.Context) (pkgtx.Tx, error)
	NoTx() pkgtx.Tx
}

type Controller struct {
	cfg *Config
	s   Storage
}

func New(cfg *Config, s Storage) *Controller {
	return &Controller{cfg: cfg, s: s}
}

type RegisterData struct {
	Email    string      `valid:"stringlength(5|255),custom_email"`
	Phone    string      `valid:"phone,required"`
	Account  AccountType `valid:"accountType,required"`
	Password string      `valid:"stringlength(6|30),required"`
}

func (rd *RegisterData) Validate() error {
	var fieldErrors = []struct {
		Field string
		Error error
	}{
		{
			Field: "Email",
			Error: ErrInvalidEmail,
		},
		{
			Field: "Phone",
			Error: ErrInvalidPhone,
		},
		{
			Field: "Account",
			Error: ErrInvalidAccount,
		},
		{
			Field: "Password",
			Error: ErrInvalidPassword,
		},
	}

	if valid, err := govalidator.ValidateStruct(rd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.Wrap(fe.Error, msg)
			}
		}

		return errors.New("auth struct is filled with some invalid data")
	}

	return nil
}

func toStorageAccount(at AccountType) (authStorage.AccountType, error) {
	switch at {
	case AccountTypeCompany:
		return authStorage.AccountTypeCompany, nil
	case AccountTypePersona:
		return authStorage.AccountTypePersona, nil
	default:
		return "", errors.New("wrong account type")
	}
}

func fromStorageAccount(at authStorage.AccountType) (AccountType, error) {
	switch at {
	case authStorage.AccountTypeCompany:
		return AccountTypeCompany, nil
	case authStorage.AccountTypePersona:
		return AccountTypePersona, nil
	default:
		return "", errors.New("wrong account type")
	}
}

type AuthToken struct {
	Token       string
	AccountID   string
	AccountType AccountType
	ExpiresAt   time.Time
}

type AuthClaims struct {
	jwt.StandardClaims
	AccountID   string
	AccountType AccountType
}

func (c *Controller) generateToken(accountID string, accountType AccountType) (*AuthToken, error) {
	expiresAt := time.Now().Add(c.cfg.TokenExpiration)
	claims := &AuthClaims{
		AccountID:   accountID,
		AccountType: accountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.cfg.PrivateSigningKey))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &AuthToken{
		Token:       tokenString,
		AccountID:   accountID,
		AccountType: accountType,
		ExpiresAt:   expiresAt,
	}, nil
}

func (c *Controller) isAuthorized(token string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.cfg.PrivateSigningKey), nil
	})

	switch err {
	case nil:
	case jwt.ErrSignatureInvalid:
		return nil, errors.Wrap(ErrUnauthorized, "signature mismatch")
	default:
		return nil, errors.Wrap(ErrInvalidArgument, "wrong token")
	}

	if !parsedToken.Valid {
		return nil, errors.Wrap(ErrUnauthorized, "invalid token")
	}

	return claims, nil
}

func (c *Controller) refreshToken(tokenStr string) (*AuthToken, error) {
	acs, err := c.isAuthorized(tokenStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// a new token will only be issued if the old token is within TokenValidityGap of expiry.
	if time.Until(time.Unix(acs.ExpiresAt, 0)) > c.cfg.TokenValidityGap {
		return nil, errors.Wrap(ErrUnauthorized, "token near to expiration and couldn't be refreshed")
	}

	expiresAt := time.Now().Add(c.cfg.TokenExpiration)
	claims := &AuthClaims{
		AccountID:   acs.AccountID,
		AccountType: acs.AccountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.cfg.PrivateSigningKey))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &AuthToken{
		Token:       tokenString,
		AccountID:   acs.AccountID,
		AccountType: acs.AccountType,
		ExpiresAt:   expiresAt,
	}, nil
}

func passwordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(passwordHash), nil
}

func (c *Controller) Register(ctx context.Context, rd *RegisterData) (*AuthToken, error) {
	if err := rd.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	var authToken *AuthToken

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		var err error
		if rd.Email == "" {
			_, err = c.s.TxGetAuthDataByPhone(ctx, tx, rd.Phone)
		} else {
			_, err = c.s.TxGetAuthDataByPhoneOrEmail(ctx, tx, rd.Phone, rd.Email)
		}
		switch err {
		case nil:
			return errors.Wrap(ErrAlreadyExists, "account with specified phone or email already exists")
		case authStorage.ErrNotFound:
		default:
			return errors.WithStack(err)
		}

		accountID := uuid.NewV4().String()

		at, err := c.generateToken(accountID, rd.Account)
		if err != nil {
			return errors.WithStack(err)
		}
		authToken = at

		account, err := toStorageAccount(rd.Account)
		if err != nil {
			return errors.WithStack(err)
		}

		ph, err := passwordHash(rd.Password)
		if err != nil {
			return errors.WithStack(err)
		}

		now := time.Now()

		ad := &authStorage.AuthData{
			AccountID:    accountID,
			Account:      account,
			Email:        rd.Email,
			Phone:        rd.Phone,
			PasswordHash: ph,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		return errors.WithStack(c.s.TxPutAuth(ctx, tx, ad))
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return authToken, nil
}

type LoginData struct {
	Login    string `valid:"stringlength(5|255),required"`
	Password string `valid:"stringlength(6|30),required"`
}

func (ld *LoginData) Validate() error {
	var fieldErrors = []struct {
		Field string
		Error error
	}{
		{
			Field: "Login",
			Error: ErrInvalidLogin,
		},
		{
			Field: "Password",
			Error: ErrInvalidPassword,
		},
	}

	if valid, err := govalidator.ValidateStruct(ld); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				return errors.Wrap(fe.Error, msg)
			}
		}

		return errors.New("login struct is filled with some invalid data")
	}

	return nil
}

func (c *Controller) Login(ctx context.Context, ld *LoginData) (*AuthToken, error) {
	if err := ld.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	ad, err := c.s.TxGetAuthDataByPhoneOrEmail(ctx, c.s.NoTx(), ld.Login, ld.Login)
	switch err {
	case nil:
	case authStorage.ErrNotFound:
		return nil, errors.Wrap(ErrUnauthorized, "specified login isn't registered")
	default:
		return nil, errors.WithStack(err)
	}

	ph, err := passwordHash(ld.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if ph != ad.PasswordHash {
		return nil, errors.Wrap(ErrUnauthorized, "wrong password")
	}

	account, err := fromStorageAccount(ad.Account)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	at, err := c.generateToken(ad.AccountID, account)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return at, nil
}

func (c *Controller) Refresh(ctx context.Context, tokenStr string) (*AuthToken, error) {
	at, err := c.refreshToken(tokenStr)
	return at, errors.WithStack(err)
}

func (c *Controller) GetAuthClaims(ctx context.Context, tokenStr string) (*AuthClaims, error) {
	return c.isAuthorized(tokenStr)
}

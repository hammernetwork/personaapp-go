package controller

import (
	"context"
	"personaapp/internal/controllers/auth/storage"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cockroachdb/errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/bcrypt"

	pkgtx "personaapp/pkg/tx"
)

type AccountType string

const (
	AccountTypeCompany AccountType = "company"
	AccountTypePersona AccountType = "persona"
	AccountTypeAdmin   AccountType = "admin"
)

func init() {
	govalidator.CustomTypeTagMap.Set("account_type", func(i interface{}, o interface{}) bool {
		if at, ok := i.(AccountType); ok {
			switch at {
			case AccountTypeCompany, AccountTypePersona, AccountTypeAdmin:
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

		rd, ok := o.(RegisterData)
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
	ErrAlreadyExists              = errors.New("already exists")
	ErrUnauthorized               = errors.New("unauthorized")
	ErrInvalidToken               = errors.New("invalid token")
	ErrInvalidLogin               = errors.New("invalid login")
	ErrInvalidLoginLength         = errors.New("invalid login length")
	ErrInvalidEmail               = errors.New("invalid email")
	ErrInvalidEmailFormat         = errors.New("invalid email format")
	ErrInvalidEmailLength         = errors.New("invalid email length")
	ErrInvalidPhone               = errors.New("invalid phone")
	ErrInvalidPhoneFormat         = errors.New("invalid phone format")
	ErrInvalidPhoneRequired       = errors.New("invalid phone required")
	ErrInvalidAccount             = errors.New("invalid account")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrInvalidPasswordLength      = errors.New("invalid password length")
	ErrInvalidOldPassword         = errors.New("invalid old password")
	ErrInvalidOldPasswordLength   = errors.New("invalid old password length")
	ErrInvalidOldPasswordNotMatch = errors.New("invalid old password not match")
	ErrAuthEntityNotFound         = errors.New("auth entity not found")
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
	TxPutAuth(ctx context.Context, tx pkgtx.Tx, ad *storage.AuthData) error
	TxGetAuthDataByID(ctx context.Context, tx pkgtx.Tx, accountID string) (*storage.AuthData, error)
	TxGetAuthDataByPhoneOrEmail(ctx context.Context, tx pkgtx.Tx, phone, email string) (*storage.AuthData, error)
	TxGetAuthDataByPhone(ctx context.Context, tx pkgtx.Tx, phone string) (*storage.AuthData, error)
	TxGetAuthDataByEmail(ctx context.Context, tx pkgtx.Tx, email string) (*storage.AuthData, error)

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
	Account  AccountType `valid:"account_type,required"`
	Password string      `valid:"stringlength(6|30),required"`
}

type AuthData struct {
	Email   string
	Phone   string
	Account AccountType
}

func (rd *RegisterData) Validate() error {
	var fieldErrors = []struct {
		Field        string
		Errors       map[string]error
		DefaultError error
	}{
		{
			Field: "Email",
			Errors: map[string]error{
				"stringlength": ErrInvalidEmailLength,
				"custom_email": ErrInvalidEmailFormat,
			},
			DefaultError: ErrInvalidEmail,
		},
		{
			Field: "Phone",
			Errors: map[string]error{
				"phone": ErrInvalidPhoneFormat,
			},
			DefaultError: ErrInvalidPhone,
		},
		{
			Field: "Account",
			Errors: map[string]error{
				"account_type": ErrInvalidAccount,
			},
			DefaultError: ErrInvalidAccount,
		},
		{
			Field: "Password",
			Errors: map[string]error{
				"stringlength": ErrInvalidPasswordLength,
			},
			DefaultError: ErrInvalidPassword,
		},
	}

	if valid, err := govalidator.ValidateStruct(rd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				validatorError, ok := err.(govalidator.Error)
				if !ok {
					return errors.Wrap(fe.DefaultError, msg)
				}

				if err, ok := fe.Errors[validatorError.Validator]; ok {
					return errors.Wrap(err, msg)
				}

				return errors.Wrap(fe.DefaultError, msg)
			}
		}

		return errors.New("auth struct is filled with some invalid data")
	}

	return nil
}

func toStorageAccount(at AccountType) (storage.AccountType, error) {
	switch at {
	case AccountTypeCompany:
		return storage.AccountTypeCompany, nil
	case AccountTypePersona:
		return storage.AccountTypePersona, nil
	case AccountTypeAdmin:
		return storage.AccountTypeAdmin, nil
	default:
		return "", errors.New("wrong account type")
	}
}

func fromStorageAccount(at storage.AccountType) (AccountType, error) {
	switch at {
	case storage.AccountTypeCompany:
		return AccountTypeCompany, nil
	case storage.AccountTypePersona:
		return AccountTypePersona, nil
	case storage.AccountTypeAdmin:
		return AccountTypeAdmin, nil
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
		return nil, errors.Wrap(ErrInvalidToken, "wrong token")
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
		case storage.ErrNotFound:
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

		ad := &storage.AuthData{
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

// nolint:dupl // will rework
func (ld *LoginData) Validate() error {
	var fieldErrors = []struct {
		Field        string
		Errors       map[string]error
		DefaultError error
	}{
		{
			Field: "Login",
			Errors: map[string]error{
				"stringlength": ErrInvalidLoginLength,
			},
			DefaultError: ErrInvalidLogin,
		},
		{
			Field: "Password",
			Errors: map[string]error{
				"stringlength": ErrInvalidPasswordLength,
			},
			DefaultError: ErrInvalidPassword,
		},
	}

	if valid, err := govalidator.ValidateStruct(ld); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				validatorError, ok := err.(govalidator.Error)
				if !ok {
					return errors.Wrap(fe.DefaultError, msg)
				}

				if err, ok := fe.Errors[validatorError.Validator]; ok {
					return errors.Wrap(err, msg)
				}

				return errors.Wrap(fe.DefaultError, msg)
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
	case storage.ErrNotFound:
		return nil, errors.Wrap(ErrUnauthorized, "specified login isn't registered")
	default:
		return nil, errors.WithStack(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(ad.PasswordHash), []byte(ld.Password)) != nil {
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

func (c *Controller) GetSelf(ctx context.Context, accountID string) (*AuthData, error) {
	var authData *AuthData

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		ad, err := c.s.TxGetAuthDataByID(ctx, tx, accountID)

		switch errors.Cause(err) {
		case nil:
		case storage.ErrNotFound:
			return errors.WithStack(ErrInvalidAccount)
		default:
			return errors.WithStack(err)
		}

		account, err := fromStorageAccount(ad.Account)
		if err != nil {
			return errors.WithStack(err)
		}

		authData = &AuthData{
			Email:   ad.Email,
			Phone:   ad.Phone,
			Account: account,
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return authData, nil
}

// nolint:dupl,funlen // will rework
func (c *Controller) UpdateEmail(
	ctx context.Context,
	accountID string,
	email string,
	password string,
	ac AccountType,
) (*AuthToken, error) {
	rd := RegisterData{Email: email, Account: ac}
	if valid, err := govalidator.ValidateStruct(rd); !valid {
		if msg := govalidator.ErrorByField(err, "Email"); msg != "" {
			validatorError, ok := err.(govalidator.Error)
			if !ok {
				return nil, errors.Wrap(ErrInvalidEmail, msg)
			}

			switch validatorError.Validator {
			case "stringlength":
				return nil, errors.Wrap(ErrInvalidEmailLength, msg)
			case "custom_email":
				return nil, errors.Wrap(ErrInvalidEmailFormat, msg)
			default:
				return nil, errors.Wrap(ErrInvalidEmail, msg)
			}
		}
	}

	ad, err := c.s.TxGetAuthDataByID(ctx, c.s.NoTx(), accountID)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrAuthEntityNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(ad.PasswordHash), []byte(password)) != nil {
		return nil, errors.Wrap(ErrInvalidPassword, "wrong password")
	}

	ad.Email = email
	ad.UpdatedAt = time.Now()

	if err := pkgtx.RunInTx(ctx, c.s, putAuthUpdatedMail(email, c, ad)); err != nil {
		return nil, errors.WithStack(err)
	}

	at, err := fromStorageAccount(ad.Account)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sat, err := c.generateToken(accountID, at)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sat, nil
}

func putAuthUpdatedMail(
	email string,
	c *Controller,
	ad *storage.AuthData,
) func(ctx context.Context, tx pkgtx.Tx) error {
	return func(ctx context.Context, tx pkgtx.Tx) error {
		if email != "" {
			switch _, err := c.s.TxGetAuthDataByEmail(ctx, tx, email); errors.Cause(err) {
			case storage.ErrNotFound:
			case nil:
				return errors.WithStack(ErrAlreadyExists)
			}
		}

		return errors.WithStack(c.s.TxPutAuth(ctx, tx, ad))
	}
}

// nolint:dupl // will rework
func (c *Controller) UpdatePhone(
	ctx context.Context,
	accountID string,
	phone string,
	password string,
) (*AuthToken, error) {
	rd := RegisterData{Phone: phone}
	if valid, err := govalidator.ValidateStruct(rd); !valid {
		if msg := govalidator.ErrorByField(err, "Phone"); msg != "" {
			validatorError, ok := err.(govalidator.Error)
			if !ok {
				return nil, errors.Wrap(ErrInvalidPhone, msg)
			}

			switch validatorError.Validator {
			case "phone":
				return nil, errors.Wrap(ErrInvalidPhoneFormat, msg)
			case "required":
				return nil, errors.Wrap(ErrInvalidPhoneRequired, msg)
			default:
				return nil, errors.Wrap(ErrInvalidPhone, msg)
			}
		}
	}

	ad, err := c.s.TxGetAuthDataByID(ctx, c.s.NoTx(), accountID)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrAuthEntityNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(ad.PasswordHash), []byte(password)) != nil {
		return nil, errors.Wrap(ErrInvalidPassword, "wrong password")
	}

	ad.Phone = phone
	ad.UpdatedAt = time.Now()

	if err := pkgtx.RunInTx(ctx, c.s, putAuthUpdatedPhone(c, phone, ad)); err != nil {
		return nil, errors.WithStack(err)
	}

	at, err := fromStorageAccount(ad.Account)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sat, err := c.generateToken(accountID, at)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sat, nil
}

func putAuthUpdatedPhone(
	c *Controller,
	phone string,
	ad *storage.AuthData,
) func(ctx context.Context, tx pkgtx.Tx) error {
	return func(ctx context.Context, tx pkgtx.Tx) error {
		switch _, err := c.s.TxGetAuthDataByPhone(ctx, tx, phone); errors.Cause(err) {
		case storage.ErrNotFound:
		case nil:
			return errors.WithStack(ErrAlreadyExists)
		}

		return errors.WithStack(c.s.TxPutAuth(ctx, tx, ad))
	}
}

type UpdatePasswordData struct {
	OldPassword string `valid:"stringlength(6|30),required"`
	NewPassword string `valid:"stringlength(6|30),required"`
}

// nolint:dupl // will rework
func (upd *UpdatePasswordData) Validate() error {
	var fieldErrors = []struct {
		Field        string
		Errors       map[string]error
		DefaultError error
	}{
		{
			Field:        "OldPassword",
			Errors:       map[string]error{"stringlength": ErrInvalidOldPasswordLength},
			DefaultError: ErrInvalidOldPassword,
		},
		{
			Field:        "NewPassword",
			Errors:       map[string]error{"stringlength": ErrInvalidPasswordLength},
			DefaultError: ErrInvalidPassword,
		},
	}

	if valid, err := govalidator.ValidateStruct(upd); !valid {
		for _, fe := range fieldErrors {
			if msg := govalidator.ErrorByField(err, fe.Field); msg != "" {
				validatorError, ok := err.(govalidator.Error)
				if !ok {
					return errors.Wrap(fe.DefaultError, msg)
				}

				if err, ok := fe.Errors[validatorError.Validator]; ok {
					return errors.Wrap(err, msg)
				}

				return errors.Wrap(fe.DefaultError, msg)
			}
		}

		return errors.New("update password struct is filled with some invalid data")
	}

	return nil
}

func (c *Controller) UpdatePassword(
	ctx context.Context,
	accountID string,
	upd *UpdatePasswordData,
) (*AuthToken, error) {
	if err := upd.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	ad, err := c.s.TxGetAuthDataByID(ctx, c.s.NoTx(), accountID)
	switch err {
	case nil:
	case storage.ErrNotFound:
		return nil, errors.WithStack(ErrAuthEntityNotFound)
	default:
		return nil, errors.WithStack(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(ad.PasswordHash), []byte(upd.OldPassword)) != nil {
		return nil, errors.WithStack(ErrInvalidOldPasswordNotMatch)
	}

	newPasswordHash, err := passwordHash(upd.NewPassword)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ad.PasswordHash = newPasswordHash
	ad.UpdatedAt = time.Now()

	if err := pkgtx.RunInTx(ctx, c.s, func(ctx context.Context, tx pkgtx.Tx) error {
		return errors.WithStack(c.s.TxPutAuth(ctx, tx, ad))
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	at, err := fromStorageAccount(ad.Account)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sat, err := c.generateToken(accountID, at)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sat, nil
}

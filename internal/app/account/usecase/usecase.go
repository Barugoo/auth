package usecase

import (
	"bytes"
	"context"
	"image/png"

	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"github.com/barugoo/oscillo-auth/config"

	"github.com/barugoo/oscillo-auth/internal/app/errors"
	"github.com/barugoo/oscillo-auth/internal/app/models"
	"github.com/barugoo/oscillo-auth/internal/app/service"

	"github.com/barugoo/oscillo-auth/internal/app/account/repository"
)

type AccountUsecase interface {
	RegisterWithCredentials(ctx context.Context, cred *models.Credentials) (bool, error)
	AuthByCredentials(ctx context.Context, cred *models.Credentials) (string, error)
	UpdateCredentials(ctx context.Context, cred *models.Credentials) (bool, error)
	ActivateAccount(ctx context.Context, email string) (bool, error)
	Generate2FA(ctx context.Context, email string) ([]byte, error)
	Setup2FA(ctx context.Context, email, code string) (bool, error)
	Remove2FA(ctx context.Context, email, code string) (bool, error)
	Verify2FA(ctx context.Context, email, code string) (bool, error)
}

type accountUsecase struct {
	service    service.AuthService
	config     *config.ServiceConfig
	repository repository.AccountRepository
}

func NewAccountUsecase(config *config.ServiceConfig, service service.AuthService, repository repository.AccountRepository) AccountUsecase {
	return &accountUsecase{
		config:     config,
		service:    service,
		repository: repository,
	}
}

func (uc *accountUsecase) RegisterWithCredentials(ctx context.Context, cred *models.Credentials) (bool, error) {
	span := uc.service.StartSpan(ctx, "RegisterWithCredentials")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.registerWithCredentials(ctx, cred)
}

func (uc *accountUsecase) registerWithCredentials(ctx context.Context, cred *models.Credentials) (bool, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return false, err
	}

	account := &models.Account{
		Email:        cred.Email,
		PasswordHash: hash,
		IsActive:     false,
	}

	_, err = uc.repository.CreateAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return true, err
}

func (uc *accountUsecase) AuthByCredentials(ctx context.Context, cred *models.Credentials) (string, error) {
	span := uc.service.StartSpan(ctx, "AuthByCredentials")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.authByCredentials(ctx, cred)
}

func (uc *accountUsecase) authByCredentials(ctx context.Context, cred *models.Credentials) (string, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return "", err
	}

	account, err := uc.repository.GetAccountByEmail(ctx, cred.Email)
	if err != nil {
		return "", err
	}

	if !account.IsActive {
		return "", errors.ErrInactiveAccount
	}

	if account.PasswordHash != hash {
		return "", errors.ErrWrongPassword
	}

	return uc.makeAccountToken(account)
}

func (uc *accountUsecase) UpdateCredentials(ctx context.Context, cred *models.Credentials) (bool, error) {
	span := uc.service.StartSpan(ctx, "UpdateCredentials")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.updateCredentials(ctx, cred)
}

func (uc *accountUsecase) updateCredentials(ctx context.Context, cred *models.Credentials) (bool, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return false, err
	}

	account, err := uc.repository.GetAccountByEmail(ctx, cred.Email)
	if err != nil {
		return false, err
	}
	account.PasswordHash = hash

	account, err = uc.repository.UpdateAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return true, err
}

func (uc *accountUsecase) ActivateAccount(ctx context.Context, email string) (bool, error) {
	span := uc.service.StartSpan(ctx, "ActivateAccount")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.activateAccount(ctx, email)
}

func (uc *accountUsecase) activateAccount(ctx context.Context, email string) (bool, error) {
	account, err := uc.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	account.IsActive = true

	_, err = uc.repository.UpdateAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uc *accountUsecase) Generate2FA(ctx context.Context, email string) ([]byte, error) {
	span := uc.service.StartSpan(ctx, "Generate2FA")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.generate2FA(ctx, email)
}

func (uc *accountUsecase) generate2FA(ctx context.Context, email string) ([]byte, error) {
	account, err := uc.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      uc.config.Issuer2FA,
		AccountName: account.Email,
	})
	if err != nil {
		return nil, err
	}

	ok, err := uc.service.SetKV(ctx, account.ID, key.Secret())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.ErrUnableToStoreKey
	}

	return uc.genQRCode(key)
}

func (uc *accountUsecase) Setup2FA(ctx context.Context, email, code string) (bool, error) {
	span := uc.service.StartSpan(ctx, "Setup2FA")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.setup2FA(ctx, email, code)
}

func (uc *accountUsecase) setup2FA(ctx context.Context, email, code string) (bool, error) {
	account, err := uc.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	secret, err := uc.service.GetKV(ctx, account.ID)
	if err != nil {
		return true, err
	}

	valid := totp.Validate(code, secret)
	if !valid {
		return false, errors.ErrInvalid2FACode
	}

	account.Secret2FA = secret

	_, err = uc.repository.UpdateAccount(ctx, account)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (uc *accountUsecase) Remove2FA(ctx context.Context, email, code string) (bool, error) {
	span := uc.service.StartSpan(ctx, "Remove2FA")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.remove2FA(ctx, email, code)
}

func (uc *accountUsecase) remove2FA(ctx context.Context, email, code string) (bool, error) {
	account, err := uc.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if !account.Has2FA() {
		return false, errors.Err2FADisabled
	}

	valid := totp.Validate(code, account.Secret2FA)
	if !valid {
		return false, errors.ErrInvalid2FACode
	}

	account.Secret2FA = ""

	_, err = uc.repository.UpdateAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func (uc *accountUsecase) Verify2FA(ctx context.Context, email, code string) (bool, error) {
	span := uc.service.StartSpan(ctx, "Verify2FA")
	defer span.Finish()

	ctx = uc.service.ContextWithSpan(context.Background(), span)

	return uc.verify2FA(ctx, email, code)
}

func (uc *accountUsecase) verify2FA(ctx context.Context, email, code string) (bool, error) {
	account, err := uc.repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if !account.Has2FA() {
		return false, errors.Err2FADisabled
	}

	valid := totp.Validate(code, account.Secret2FA)
	if !valid {
		return false, errors.ErrInvalid2FACode
	}
	return valid, nil
}

func (uc *accountUsecase) genQRCode(key *otp.Key) ([]byte, error) {
	var buf bytes.Buffer

	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (uc *accountUsecase) makeAccountToken(account *models.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   account.Email,
		"has_2fa": account.Has2FA(),
	})
	return token.SignedString(uc.config.AppSecret)
}

func (uc *accountUsecase) hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

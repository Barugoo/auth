package usecase

import (
	"bytes"
	"github.com/barugoo/oscillo-auth/internal/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"image/png"

	"github.com/barugoo/oscillo-auth/config"
	"github.com/barugoo/oscillo-auth/internal/errors"
	"github.com/barugoo/oscillo-auth/internal/models"
	"github.com/barugoo/oscillo-auth/internal/repository"
)

type AccountUsecase interface {
	RegisterWithCredentials(cred *models.Credentials) (bool, error)
	AuthByCredentials(cred *models.Credentials) (string, error)
	UpdateCredentials(cred *models.Credentials) (bool, error)
	ActivateAccount(email string) (bool, error)
	Generate2FA(email string) ([]byte, error)
	Setup2FA(email, code string) (bool, error)
	Remove2FA(email, code string) (bool, error)
	Verify2FA(email, code string) (bool, error)
}

type accountUsecase struct {
	service service.ServiceInterface
	config  *config.ServiceConfig
	db      repository.AccountRepository
}

func NewAccountUsecase(config *config.ServiceConfig, service service.ServiceInterface, db repository.AccountRepository) AccountUsecase {
	return accountUsecase{
		config:  config,
		service: service,
		db:      db,
	}
}

func (uc *accountUsecase) RegisterWithCredentials(cred *models.Credentials) (bool, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return false, err
	}

	account := &models.Account{
		Email:        cred.Email,
		PasswordHash: hash,
		IsActive:     false,
	}

	_, err = uc.db.CreateAccount(account)
	if err != nil {
		return false, err
	}
	return true, err
}

func (uc *accountUsecase) AuthByCredentials(cred *models.Credentials) (string, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return "", err
	}

	account, err := uc.db.GetAccountByEmail(cred.Email)
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

func (uc *accountUsecase) UpdateCredentials(cred *models.Credentials) (bool, error) {
	hash, err := uc.hash(cred.Password)
	if err != nil {
		return false, err
	}

	account, err := uc.db.GetAccountByEmail(cred.Email)
	if err != nil {
		return false, err
	}
	account.PasswordHash = hash

	account, err = uc.db.UpdateAccount(account)
	if err != nil {
		return false, err
	}
	return true, err
}

func (uc accountUsecase) ActivateAccount(email string) (bool, error) {
	account, err := uc.db.GetAccountByEmail(email)
	if err != nil {
		return false, err
	}

	account.IsActive = true

	_, err = uc.db.UpdateAccount(account)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uc *accountUsecase) Generate2FA(email string) ([]byte, error) {
	account, err := uc.db.GetAccountByEmail(email)
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

	ok, err := uc.service.SetKV(account.ID, key.Secret())
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.ErrUnableToStoreKey
	}

	return uc.genQRCode(key)
}

func (uc *accountUsecase) Setup2FA(email, code string) (bool, error) {
	account, err := uc.db.GetAccountByEmail(email)
	if err != nil {
		return false, err
	}

	secret, err := uc.service.GetKV(account.ID)
	if err != nil {
		return true, err
	}

	valid := totp.Validate(code, secret)
	if !valid {
		return false, errors.ErrInvalid2FACode
	}

	account.Secret2FA = secret

	_, err = uc.db.UpdateAccount(account)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func (uc *accountUsecase) Remove2FA(email, code string) (bool, error) {
	account, err := uc.db.GetAccountByEmail(email)
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

	_, err = uc.db.UpdateAccount(account)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func (uc *accountUsecase) Verify2FA(email, code string) (bool, error) {
	account, err := uc.db.GetAccountByEmail(email)
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

package errors

import (
	"errors"
)

var (
	ErrWrongPassword    = errors.New("wrong password")
	ErrUnableToStoreKey = errors.New("unable to store key")
	ErrInvalid2FACode   = errors.New("invalid 2FA code")
	Err2FADisabled      = errors.New("2fa disabled")
	ErrInactiveAccount  = errors.New("inactive account")
)

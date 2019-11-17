package errors

import (
	"errors"
	"fmt"
)

type DeliveryError struct {
	Email string
	Err   error
}

func (e *DeliveryError) Unwrap() error {
	return e.Err
}

func (e *DeliveryError) Error() string {
	return fmt.Sprintf("%s: %v", e.Email, e.Err)
}

var (
	ErrBrokenContext = errors.New("broken context")
)

type UsecaseError struct {
	Method string
	Err    error
}

func (e *UsecaseError) Unwrap() error {
	return e.Err
}

func (e *UsecaseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Method, e.Err)
}

var (
	ErrWrongPassword    = errors.New("wrong password")
	ErrUnableToStoreKey = errors.New("unable to store key")
	ErrInvalid2FACode   = errors.New("invalid 2FA code")
	Err2FADisabled      = errors.New("2fa disabled")
	ErrInactiveAccount  = errors.New("inactive account")
)

type RepositoryError struct {
	Impl string
	Err  error
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("%s: %v", e.Impl, e.Err)
}

var (
	ErrNotFound = errors.New("not found")
)

type ServiceError struct {
	Method string
	Err    error
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Err)
}

func (e ServiceError) Unwrap() error {
	return e.Err
}

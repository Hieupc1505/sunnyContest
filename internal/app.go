package app

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal")
	ErrNotFound            = errors.New("not found")
	ErrForbidden           = errors.New("forbidden")
	ErrTooManyRequests     = errors.New("too many requests")
)

var (
	ErrPasswordInCorrect = errors.New("password is incorrect")
	ErrAccountLocked     = errors.New("account has been locked")
	ErrAccountDisabled   = errors.New("account is disabled")
	ErrAccountDeleted    = errors.New("account is deleted")
	ErrUserAlreadyExists = errors.New("user is already in use")
)

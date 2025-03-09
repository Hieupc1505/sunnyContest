package app

import "errors"

var (
	ErrEmpty               = errors.New("")
	ErrInternalServerError = errors.New("system_error")
	ErrUserNotFound        = errors.New("user.not_found")
	ErrNotFound            = errors.New("data_not_found")
	ErrForbidden           = errors.New("forbidden")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrPermissionDenied    = errors.New("permission_denied")
	ErrInvalidData         = errors.New("invalid_data")

	ErrInvalidToken = errors.New("unauthorized.missing_or_invalid_token")
	ErrTokenExpired = errors.New("unauthorized.session_expired")

	ErrPasswordInCorrect = errors.New("password is incorrect")
	ErrAccountLocked     = errors.New("account has been locked")
	ErrAccountDisabled   = errors.New("account is disabled")
	ErrAccountDeleted    = errors.New("account is deleted")
	ErrUserAlreadyExists = errors.New("user is already in use")

	ErrContestSubmitAlready = errors.New("contest.live.submit_already")
	ErrContestNotFound      = errors.New("invalid_data.contest.notfound")
	ErrInvalidContestGameID = errors.New("invalid_data.contest.game_id")
)

const (
	SuccessCode              = 0
	SystemErrCode            = 1
	InValidCode              = 2
	ErrUserNotFoundCode      = 10
	ErrPasswordInCorrectCode = 11
	ErrAccountLockedCode     = 12
	ErrAccountDisabledCode   = 13
	ErrAccountDeletedCode    = 14
	ErrUserAlreadyExistsCode = 15
)

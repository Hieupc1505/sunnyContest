package services

import (
	"errors"
)

type Error struct {
	appErr error
	svcErr error
	code   int
}

// appErr là mã lỗi của app -> BadRequest, Fobbiden, Untholization
// svcErr là lỗi mà xử lý logic ra nó có the nil
func NewError(svcErr error, appErr error) Error {
	temp := svcErr
	if appErr != nil {
		temp = appErr
	}
	return Error{
		svcErr: svcErr,
		appErr: temp,
	}
}

func (e Error) AddCode(code int) error {
	e.code = code
	return e
}
func (e Error) Code() int { return e.code }

func (e Error) Error() string {
	if errors.Is(e.svcErr, e.appErr) {
		return e.svcErr.Error()
	}
	return e.appErr.Error()
}

func (e Error) AppError() error {
	return e.appErr
}

func (e Error) SvcError() error {
	return e.svcErr
}

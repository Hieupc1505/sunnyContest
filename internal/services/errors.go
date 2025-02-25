package services

import "errors"

type Error struct {
	appErr error
	svcErr error
}

func NewError(svcErr, appErr error) error {
	return Error{
		svcErr: svcErr,
		appErr: appErr,
	}
}

func (e Error) Error() string {
	return errors.Join(e.svcErr, e.appErr).Error()
}

func (e Error) AppError() error {
	return e.appErr
}

func (e Error) SvcError() error {
	return e.svcErr
}

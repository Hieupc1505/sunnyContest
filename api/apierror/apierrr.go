package apierror

import (
	"errors"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/internal/services"
	"net/http"
)

const defaultCodeError = 9

type APIError struct {
	Message string `json:"msg"`
	Code    int    `json:"code,omitempty"`
}

func FromError(err error) (APIError, int) {
	var apiErr APIError
	var svcErr services.Error
	var status int

	if errors.As(err, &svcErr) {
		appErr := svcErr.AppError()
		apiErr.Message = appErr.Error()
		status = errorStatus(svcErr.SvcError())
		apiErr.Code = errorCode(appErr)
	} else {
		status = errorStatus(err)
		apiErr = APIError{
			Message: "Internal Server Error",
			Code:    defaultCodeError,
		}
	}
	return apiErr, status
}

// Xác định HTTP Status từ lỗi dịch vụ
func errorStatus(err error) int {
	statusMap := map[error]int{
		app.ErrBadRequest:          http.StatusBadRequest,
		app.ErrInternalServerError: http.StatusInternalServerError,
		app.ErrNotFound:            http.StatusNotFound,
	}

	if status, exists := statusMap[err]; exists {
		return status
	}
	return http.StatusInternalServerError
}

// Xác định mã lỗi
func errorCode(err error) int {
	codeMap := map[error]int{
		app.ErrPasswordInCorrect: 11,
		app.ErrAccountLocked:     12,
		app.ErrAccountDisabled:   13,
		app.ErrAccountDeleted:    14,
		app.ErrUserAlreadyExists: 15,
	}

	if code, exists := codeMap[err]; exists {
		return code
	}
	return defaultCodeError
}

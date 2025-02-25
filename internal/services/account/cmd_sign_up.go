package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/internal/services"
	"go-rest-api-boilerplate/pkg/errsx"
)

func signupUser(email string, password string) (*User, error) {
	var errs errsx.Map
	emailAddr, err := NewEmail(email)
	if err != nil {
		fmt.Println("Email error")
		errs.Set("email", err)
	}

	pwd, err := NewPassword(password)
	if err != nil {
		errs.Set("password", err)
	}

	// ✅ Kiểm tra xem map có lỗi nào không
	if len(errs) > 0 {
		return nil, errs.ToError()
	}

	user := NewUser(emailAddr, pwd)
	return user, nil
}

func (s *Service) SignUp(ctx *gin.Context, email, password string) (*User, error) {
	user, err := signupUser(email, password)
	if err != nil {
		return nil, services.NewError(app.ErrBadRequest, err)
	}

	//todo: hash password user

	// Lưu user vào db
	if _, err = s.repo.Add(ctx, *user.ToParams()); err != nil {
		return nil, services.NewError(app.ErrInternalServerError, err)
	}

	return user, nil
}

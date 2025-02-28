package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/passwordutil"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services"
	"go-rest-api-boilerplate/pkg/errsx"
	"log/slog"
)

func signupUser(username string, password string) (*User, error) {
	var errs errsx.Map
	u, err := NewUsername(username)
	if err != nil {
		fmt.Println("Email error")
		errs.Set("username", err)
	}

	pwd, err := NewPassword(password)
	if err != nil {
		errs.Set("password", err)
	}

	// ✅ Kiểm tra xem map có lỗi nào không
	if len(errs) > 0 {
		return nil, errs.ToError()
	}

	user := NewUser(u, pwd)
	return user, nil
}

func (s *Service) SignUp(ctx *gin.Context, username, password string) (*db.AddRow, error) {
	user, err := signupUser(username, password)
	if err != nil {
		return nil, services.NewError(err, nil)
	}

	hashPass, err := passwordutil.HashPassword(user.Password)
	if err != nil {
		return nil, services.NewError(err, app.ErrInternalServerError)
	}

	user.Password = hashPass

	// Lưu user vào db
	ur, err := s.repo.Add(ctx, *user.ToParams())
	if err != nil {
		slog.Info("Error adding user", "error", err)
		return nil, services.NewError(err, app.ErrInternalServerError)
	}

	return &ur, nil
}

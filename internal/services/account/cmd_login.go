package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/passwordutil"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services"
)

func (s *Service) Login(ctx *gin.Context, username, password string) (*UserInfo, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.DB_NOT_FOUND) {
			return nil, services.NewError(app.ErrEmpty, nil).AddCode(app.ErrUserNotFoundCode)
		}
		return nil, services.NewError(app.ErrInternalServerError, nil)
	}

	if err := passwordutil.ComparePassword(password, user.Password); err != nil {
		return nil, services.NewError(app.ErrPasswordInCorrect, nil)
	}

	return ToUserInfo(user), nil
}

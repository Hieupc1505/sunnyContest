package account

import (
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services"
	"go-rest-api-boilerplate/pkg/errsx"
	"log/slog"
)

func profileUser(userID int64, nickname string, avatar string) (*db.AddProfileParams, error) {
	var errs errsx.Map
	nn, err := NewNickname(nickname)
	if err != nil {
		errs.Set("nickname", err)
	}

	avt, err := NewAvatar(avatar)
	if err != nil {
		errs.Set("avatar", err)
	}

	// ✅ Kiểm tra xem map có lỗi nào không
	if len(errs) > 0 {
		return nil, errs.ToError()
	}

	prof := NewProfile(userID, nn, avt)
	return prof, nil
}

func (s *Service) AddNickname(ctx *gin.Context, userID int64, nickname string, avatar string) (*db.AddProfileRow, error) {
	profile, err := profileUser(userID, nickname, avatar)
	if err != nil {
		return nil, services.NewError(err, nil)
	}
	userProfile, err := s.repo.AddProfile(ctx, *profile)
	if err != nil {
		slog.Info("Add profile error", "error", err)
		return nil, services.NewError(app.ErrInternalServerError, nil)
	}
	return &userProfile, nil
}

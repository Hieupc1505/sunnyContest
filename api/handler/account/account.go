package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/middleware"
	"go-rest-api-boilerplate/api/response"
	app "go-rest-api-boilerplate/internal"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/account"
	"go-rest-api-boilerplate/pkg/token"
	"log/slog"
	"time"
)

func Routes(h *handler.Handler) {
	authGroup := h.App.Group("/api/v1")
	{
		authGroup.POST("/public/register", Signup(h))
		authGroup.POST("/public/login", Login(h))
	}

	authSGroup := authGroup.Use(middleware.AuthMiddleware(h))
	{
		authSGroup.POST("/register-nickname", RegisterNickname(h))
		authSGroup.GET("/get-current-user", GetCurrentUser(h))
	}
}

func GetCurrentUser(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := contextutil.GetUser(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting user from context", "error", err)
			response.Error(ctx, app.ErrPermissionDenied)
			return
		}
		fmt.Println("userID:::", userID)
		profile, err := h.UserRepo.GetUserByID(ctx, userID)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting user by ID", "error", err)
			response.Error(ctx, nil)
			return
		}
		response.Success(ctx, profile)
	}
}

func RegisterNickname(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		avatarDefaultValue := ""

		userID, err := contextutil.GetUser(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting user from context", "error", err)
			response.Error(ctx, app.ErrPermissionDenied)
			return
		}

		var data account.RegisterNickNameParams
		if err := ctx.ShouldBindJSON(&data); err != nil {
			slog.ErrorContext(ctx, "Error binding nickname data", "error", err)
			response.Error(ctx, nil)
			return
		}
		profile, err := h.AccountService.AddNickname(ctx, userID, data.NickName, avatarDefaultValue)
		if err != nil {
			slog.ErrorContext(ctx, "Error adding nickname", "error", err)
			response.Error(ctx, nil)
			return
		}
		if err := h.UserRepo.UpdateUserRole(ctx, db.UpdateUserRoleParams{ID: userID, Role: data.Type}); err != nil {
			slog.ErrorContext(ctx, "Error adding nickname", "error", err)
			response.Error(ctx, nil)
			return
		}
		response.Success(ctx, profile)

	}
}

func Login(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data account.User
		if err := ctx.ShouldBindJSON(&data); err != nil {
			slog.ErrorContext(ctx, "Error binding login data", "error", err)
			response.Error(ctx, err)
			return
		}
		user, err := h.AccountService.Login(ctx, data.Username, data.Password)
		if err != nil {
			slog.ErrorContext(ctx, "Error logging in user", "error", err)
			response.Error(ctx, err)
			return
		}

		WriteToken(ctx, h.UserRepo, h.Token, user.ID, user.Role)

		response.Success(ctx, user)
	}
}

func Signup(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data account.User
		if err := ctx.ShouldBindJSON(&data); err != nil {
			slog.ErrorContext(ctx, "Error binding data", "error", err)
			response.Error(ctx, err)
			return
		}

		user, err := h.AccountService.SignUp(ctx, data.Username, data.Password)
		if err != nil {
			slog.ErrorContext(ctx, "Error signing up user", "error", err)
			response.Error(ctx, err)
			return
		}

		WriteToken(ctx, h.UserRepo, h.Token, user.ID, user.Role)

		response.Success(ctx, user)

	}
}

func WriteToken(ctx *gin.Context, repo user_repo.ReadWriter, maker token.Maker, userID int64, role int32) {
	timeAccessTokenExpire := 24 * time.Hour
	token, _, err := maker.CreateToken(userID, role, timeAccessTokenExpire)
	if err != nil {
		slog.ErrorContext(ctx, "Error creating token", "error", err)
		response.Error(ctx, err)
		return
	}

	params, err := account.NewUpdateTokenParam(userID, token, timeAccessTokenExpire)
	if err := repo.UpdateUserToken(ctx, *params); err != nil {
		slog.ErrorContext(ctx, "Error updating user token", "error", err)
		response.Error(ctx, app.ErrInternalServerError)
		return
	}
	ctx.Writer.Header().Set("sf-token", token)

}

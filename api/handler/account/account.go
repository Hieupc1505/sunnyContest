package account

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/response"
	"go-rest-api-boilerplate/internal/services/account"
	"log/slog"
)

func Routes(h *handler.Handler) {
	authGroup := h.App.Group("/api/auth")
	{
		authGroup.POST("/signup", SignupHandler(h))
	}
}

func SignupHandler(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data account.User
		if err := ctx.ShouldBindJSON(&data); err != nil {
			slog.ErrorContext(ctx, "Error binding data", "error", err)
			response.Error(ctx, err)
			return
		}

		user, err := h.AccountService.SignUp(ctx, data.Email, data.Password)
		if err != nil {
			slog.ErrorContext(ctx, "Error signing up user", "error", err)
			response.Error(ctx, err)
			return
		}

		response.Success(ctx, user)
	}
}

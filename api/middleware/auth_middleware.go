package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/response"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/types"
	"log/slog"
	"strings"
)

var (
	AuthorizationKey = "Authorization"
)

func AuthMiddleware(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader(AuthorizationKey)
		if tokenString == "" {
			response.Error(ctx, app.ErrInvalidToken)
			ctx.Abort()
			return
		}
		token := strings.Split(tokenString, " ")[1]
		claims, err := h.Token.VerifyToken(token)
		if err != nil {
			ctx.Abort()
			if errors.Is(err, app.ErrTokenExpired) {
				response.Error(ctx, err)
				return
			}
			slog.Info("Error to verify token", "error", err)
			response.Error(ctx, nil)
			return
		}

		ctx.Set(types.UserID, claims.UserID)
		ctx.Next()
	}
}

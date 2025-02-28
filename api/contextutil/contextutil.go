package contextutil

import (
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/types"
)

// GetUser return user from context
func GetUser(ctx *gin.Context) (int64, error) {
	userID, exists := ctx.Get(types.UserID)
	if !exists {
		return 0, app.ErrForbidden
	}
	// Convert userID sang int64
	userIDInt, ok := userID.(int64)
	if !ok {
		return 0, app.ErrInternalServerError
	}
	return userIDInt, nil
}

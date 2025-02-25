package response

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/apierror"
	"net/http"
)

type Response struct {
	Error *apierror.APIError `json:"e,omitempty"`
	Data  interface{}        `json:"d,omitempty"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{Data: data})
}

func Error(ctx *gin.Context, err error) {
	apiErr, status := apierror.FromError(err)
	ctx.JSON(status, Response{Error: &apiErr})
}

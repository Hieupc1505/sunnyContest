package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/internal/services"
	"net/http"
)

type Response struct {
	Error int         `json:"e"`
	Data  interface{} `json:"d,omitempty"`
}

type ErrorResponse struct {
	Msg string `json:"msg,omitempty"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Error: 0,
		Data:  data,
	})
}

func Error(ctx *gin.Context, err error) {
	var sys services.Error
	msg := ErrorResponse{
		Msg: app.ErrInternalServerError.Error(),
	}
	response := Response{
		Error: app.SystemErrCode,
		Data:  msg,
	}

	if err != nil {
		if errors.As(err, &sys) {
			response.Error = sys.Code()
			msg.Msg = sys.Error()
		} else {
			msg.Msg = err.Error()
		}
	}

	response.Data = msg

	ctx.JSON(http.StatusOK, response)
}

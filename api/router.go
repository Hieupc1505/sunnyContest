package api

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/handler/account"
	"go-rest-api-boilerplate/api/response"
	"time"
)

func (a *API) registerRoutes() {

	h := &handler.Handler{
		App:            a.app,
		Tasks:          a.tasks,
		UserRepo:       a.store.AccountRepo,
		AccountService: a.store.AccountAPI,
	}

	account.Routes(h)

	a.app.GET("/Ping", func(ctx *gin.Context) {
		time.Sleep(time.Second)
		response.Success(ctx, gin.H{"Status": "Pong test connect"})
	})
}

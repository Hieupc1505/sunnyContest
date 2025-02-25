package handler

import (
	"github.com/gin-gonic/gin"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	"go-rest-api-boilerplate/internal/services/account"
	"sync"
)

type Handler struct {
	App   *gin.Engine
	Tasks *sync.WaitGroup

	UserRepo       user_repo.ReadWriter
	AccountService *account.Service
}

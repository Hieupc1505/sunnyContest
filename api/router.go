package api

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/handler/account"
	"go-rest-api-boilerplate/api/handler/question"
	"go-rest-api-boilerplate/api/handler/subject"
	"go-rest-api-boilerplate/api/response"
	config "go-rest-api-boilerplate/configs/app"
	"go-rest-api-boilerplate/pkg/imgUploader"
	"go-rest-api-boilerplate/pkg/token"
	"log"
	"time"
)

func (a *API) registerRoutes() {

	//maker, err := token.NewPasetoMaker(config.Envs.SymmetricKey)
	maker, err := token.NewJWTMaker(config.Envs.SymmetricKey)
	if err != nil {
		log.Fatal("Error create token maker", err)
	}

	uploader := imgUploader.NewImgbbUpload(config.Envs.ImgbbAPIKey)

	h := &handler.Handler{
		App:             a.app,
		Tasks:           a.tasks,
		UserRepo:        a.store.AccountRepo,
		AccountService:  a.store.AccountAPI,
		SubjectRepo:     a.store.SubjectRepo,
		SubjectService:  a.store.SubjectService,
		QuestionRepo:    a.store.QuestionRepo,
		QuestionService: a.store.QuestionService,
		Token:           maker,
		ImgUploader:     uploader,
	}

	account.Routes(h)
	subject.Routes(h)
	question.Routes(h)

	a.app.GET("/Ping", func(ctx *gin.Context) {
		time.Sleep(time.Second)
		response.Success(ctx, gin.H{"Status": "Pong test connect"})
	})
}

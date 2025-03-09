package api

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/handler/account"
	"go-rest-api-boilerplate/api/handler/contest"
	"go-rest-api-boilerplate/api/handler/question"
	"go-rest-api-boilerplate/api/handler/sse"
	"go-rest-api-boilerplate/api/handler/subject"
	"go-rest-api-boilerplate/api/response"
	config "go-rest-api-boilerplate/configs/app"
	"go-rest-api-boilerplate/pkg/imgUploader"
	"go-rest-api-boilerplate/pkg/sse/channel"
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

	rooms := channel.NewRooms()

	h := &handler.Handler{
		App:                a.app,
		Tasks:              a.tasks,
		UserRepo:           a.store.AccountRepo,
		AccountService:     a.store.AccountAPI,
		SubjectRepo:        a.store.SubjectRepo,
		SubjectService:     a.store.SubjectService,
		QuestionRepo:       a.store.QuestionRepo,
		QuestionService:    a.store.QuestionService,
		ContestRepo:        a.store.ContestRepo,
		ContestService:     a.store.ContestService,
		UserContestRepo:    a.store.UserContestRepo,
		UserContestService: a.store.UserContestService,
		Token:              maker,
		ImgUploader:        uploader,
		Rooms:              rooms,
	}

	account.Routes(h)
	subject.Routes(h)
	question.Routes(h)
	contest.Routes(h)
	sse.Routes(h)

	a.app.GET("/Ping", func(ctx *gin.Context) {
		time.Sleep(time.Second)
		response.Success(ctx, gin.H{"Status": "Pong test connect"})
	})
}

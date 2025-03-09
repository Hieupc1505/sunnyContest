package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api"
	config "go-rest-api-boilerplate/configs/app"
	"go-rest-api-boilerplate/configs/logger"
	"go-rest-api-boilerplate/internal/db/repo"
	contest_repo "go-rest-api-boilerplate/internal/db/repo/contest"
	question_repo "go-rest-api-boilerplate/internal/db/repo/question"
	subject_repo "go-rest-api-boilerplate/internal/db/repo/subject"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	user_contest_repo "go-rest-api-boilerplate/internal/db/repo/user_contest"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/account"
	"go-rest-api-boilerplate/internal/services/contest"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/internal/services/subject"
	"go-rest-api-boilerplate/internal/services/user_contest"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	logger := logger.NewLogger(config.Envs.Logger)
	var longTasks sync.WaitGroup
	var engine *gin.Engine

	if config.Envs.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := db.Conn()
	defer conn.Close()

	userRepo := user_repo.NewUserRepo(conn)
	accountService, err := account.NewService(ctx, userRepo)
	if err != nil {
		log.Fatal("Error creating account service: ", err)
	}

	subRepo := subject_repo.NewSubjectRepo(conn)
	subService, err := subject.NewService(ctx, subRepo)
	if err != nil {
		log.Fatal("Error creating account service: ", err)
	}

	quesRepo := question_repo.NewQuestionService(conn)
	quesService, err := question.NewService(ctx, quesRepo)
	if err != nil {
		log.Fatal("Error creating question service: ", err)
	}

	contestRepo := contest_repo.NewContestService(conn)
	contestService, err := contest.NewService(ctx, contestRepo)
	if err != nil {
		log.Fatal("Error creating question service: ", err)
	}

	userContest := user_contest_repo.NewUserContestService(conn)
	userContestService, err := user_contest.NewService(ctx, userContest)
	if err != nil {
		log.Fatal("Error creating user contest service: ", err)
	}

	store := repo.New(userRepo, accountService, subRepo, subService, quesRepo, quesService, contestRepo, contestService, userContest, userContestService)

	envPort, err := strconv.ParseInt(config.Envs.Port, 0, 64)
	if err != nil {
		log.Fatal("Invalid port found in env: ", config.Envs.Port, err)
	}

	apiInstance := api.New(
		logger,
		&longTasks,
		engine,
		store,
		api.WithPort(int(envPort)),
	)

	go apiInstance.StartServer(ctx)
	longTasks.Wait()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	if err := apiInstance.Shutdown(ctxShutdown); err != nil {
		slog.Error("Shutdown error:", err)
	} else {
		slog.Info("Server shutdown completed.")
	}
}

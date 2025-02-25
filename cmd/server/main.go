package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api"
	config "go-rest-api-boilerplate/configs/app"
	"go-rest-api-boilerplate/configs/logger"
	"go-rest-api-boilerplate/internal/db/repo"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/account"
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
	store := repo.New(userRepo, accountService)

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

package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/middleware"
	config "go-rest-api-boilerplate/configs/app"
	"go-rest-api-boilerplate/configs/logger"
	"go-rest-api-boilerplate/internal/db/repo"
	"log"
	"net/http"
	"sync"
)

type API struct {
	secure bool
	port   int
	domain string

	server *http.Server
	app    *gin.Engine
	logger *logger.LoggerZap
	tasks  *sync.WaitGroup
	store  *repo.Store
}

// NewAPI creates a new API
func New(
	log *logger.LoggerZap,
	longTasks *sync.WaitGroup,
	router *gin.Engine,
	store *repo.Store,
	opts ...OptFunc,
) *API {

	defaultPort := 8080

	instance := &http.Server{
		Addr:    fmt.Sprintf(":%d", defaultPort),
		Handler: router,
	}

	a := &API{
		secure: false,
		port:   defaultPort,
		domain: config.Envs.PublicHost,

		server: instance,
		app:    router,
		logger: log,
		tasks:  longTasks,
		store:  store,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *API) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	return nil
}

func (a *API) StartServer(ctx context.Context) error {
	//Define routers
	a.registerMiddleware()

	a.registerRoutes()

	// service connections
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	return nil
}

func (a *API) registerMiddleware() {
	// Khởi tạo RateLimiter (5 requests mỗi giây, burst tối đa 10)
	rateLimiter := middleware.NewRateLimiter(5, 10)

	// Thêm middleware rate limit vào router
	a.app.Use(rateLimiter.RateLimitMiddleware())
}

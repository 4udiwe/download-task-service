package app

import (
	"os"

	"github.com/4udiwe/download-task-service/config"
	"github.com/4udiwe/download-task-service/internal/handler"
	"github.com/4udiwe/download-task-service/internal/service"
	"github.com/4udiwe/download-task-service/internal/storage"
	"github.com/4udiwe/download-task-service/pkg/httpserver"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// Echo
	echoHandler *echo.Echo

	// Storage
	taskStorage *storage.Storage

	// Serivce
	taskService *service.TaskService

	// Handlers
	postTaskHandler    handler.Handler
	getTasksHandler    handler.Handler
	getTaskByIDHandler handler.Handler
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	return &App{
		cfg: cfg,
	}
}

func (app *App) Start() {
	// App server
	log.Info("Starting app server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()
	log.Debugf("Server port: %s", app.cfg.HTTP.Port)

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}

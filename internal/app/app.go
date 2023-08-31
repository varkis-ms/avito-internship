package app

import (
	"avito-internship/internal/config"
	v1 "avito-internship/internal/controller/http/v1"
	"avito-internship/internal/repository"
	"avito-internship/internal/service"
	"avito-internship/internal/webapi/googledrive"
	"avito-internship/pkg/database/postgresdb"
	"avito-internship/pkg/httpserver"
	"avito-internship/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

// @title Dynamic user segmentation service
// @version 1.0

// @contact.name   Markin Sergey
// @contact.email  markin-2002@yandex.ru

// @host localhost:8000
// @BasePath /api/v1

func Run(configPath string) {
	// Config
	logger := logging.GetLogger()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.WithError(err).Fatal("no config")
	}

	// Repository
	logger.Info("Initializing postgres...")
	db, err := postgresdb.New(&cfg)
	if err != nil {
		logger.WithError(err).Fatal("app.Run - postgresdb.New")
	}
	defer db.Close()
	repositories := repository.NewRepositories(db)

	// Service
	logger.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:  repositories,
		GDrive: googledrive.New(cfg.GDriveJSONFilePath),
	}
	services := service.NewServices(deps)

	// Handler
	logger.Info("Initializing handlers and routes...")
	handler := gin.Default()
	v1.NewRouter(handler, &logger, services)

	// HTTP server
	logger.Infof("Starting http server on port :%s", cfg.PortHttp)
	httpServer := httpserver.New(handler, httpserver.Port(fmt.Sprintf(cfg.PortHttp)))

	// Waiting signal
	logger.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app.Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.WithError(err).Error("app.Run - httpServer.Notify")
	}

	// Graceful shutdown
	logger.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		logger.WithError(err).Error("app.Run - httpServer.Shutdown")
	}
}

// Package project provides the project api.
package project

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/internal/project/repository"
	"github.com/devpies/saas-core/internal/project/res"
	"github.com/devpies/saas-core/internal/project/service"
	"github.com/devpies/saas-core/pkg/log"

	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	defer logger.Sync()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	pg, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		return err
	}
	defer Close()

	// Execute latest migration in production.
	if cfg.Web.Production {
		if err = res.MigrateUp(pg.URL.String()); err != nil {
			logger.Error("error connecting to user database", zap.Error(err))
			return err
		}
	}

	// Initialize 3-layered architecture.
	taskRepo := repository.NewTaskRepository(logger, pg)
	columnRepo := repository.NewColumnRepository(logger, pg)
	projectRepo := repository.NewProjectRepository(logger, pg)

	taskService := service.NewTaskService(logger, taskRepo)
	columnService := service.NewColumnService(logger, columnRepo)
	projectService := service.NewProjectService(logger, projectRepo)

	taskHandler := handler.NewTaskHandler(logger, taskService, columnService)
	columnHandler := handler.NewColumnHandler(logger, columnService)
	projectHandler := handler.NewProjectHandler(logger, projectService, columnService, taskService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, taskHandler, columnHandler, projectHandler, cfg),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting project service on %s:%s", cfg.Web.Address, cfg.Web.Port))
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		logger.Error("error on startup", zap.Error(err))
		return err
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// Give on going tasks a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			logger.Error("error on integrity issue caused shutdown", zap.Error(err))
			return err
		case err != nil:
			logger.Error("error on gracefully shutdown", zap.Error(err))
			return err
		}
	}

	return err
}

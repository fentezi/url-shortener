package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/fentezi/url-shortener/internal/config"
	"github.com/fentezi/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/fentezi/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/fentezi/url-shortener/internal/http-server/handlers/url/save"
	"github.com/fentezi/url-shortener/internal/http-server/middleware/logger"
	"github.com/fentezi/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/fentezi/url-shortener/internal/lib/logger/sl"
	"github.com/fentezi/url-shortener/internal/storage/sqlite"
	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustConfig()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	router := gin.New()
	router.Use(logger.Logger(log), gin.Recovery())
	r := router.Group("/api")
	{
		r.Use(gin.BasicAuth(gin.Accounts{cfg.HttpServer.User: cfg.HttpServer.Password}))
		r.POST("/url", save.SaveHandlerWrapper(log, storage))
		r.DELETE("/:alias", delete.DeleteHandlerWrapper(log, storage))
	}

	router.GET("/:alias", redirect.RedirectHandlerWrapper(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

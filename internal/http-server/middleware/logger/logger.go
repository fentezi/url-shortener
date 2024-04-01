package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("logger middleware enabled")

	return func(ctx *gin.Context) {
		entry := log.With(
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("remote_addr", ctx.Request.RemoteAddr),
			slog.String("user_agent", ctx.Request.UserAgent()),
		)
		t1 := time.Now()

		entry.Info("request completed",
			slog.Int("status", ctx.Writer.Status()),
			slog.Int("bytes", ctx.Writer.Size()),
			slog.String("duration", time.Since(t1).String()),
		)
		ctx.Next()
	}
}

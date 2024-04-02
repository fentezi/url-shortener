package delete

import (
	"log/slog"
	"net/http"

	resp "github.com/fentezi/url-shortener/internal/lib/api/response"
	"github.com/fentezi/url-shortener/internal/lib/logger/sl"
	"github.com/gin-gonic/gin"
)

type URLDelete interface {
	DeleteURL(alias string) error
}

func DeleteHandlerWrapper(log *slog.Logger, urlDetele URLDelete) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		DeleteHandler(log, urlDetele, ctx)
	}
}

func DeleteHandler(log *slog.Logger, urlDelete URLDelete, c *gin.Context) {
	const op = "handlers.url.delete.DeleteHandler"

	log = log.With(
		slog.String("op", op),
	)

	alias := c.Param("alias")
	if alias == "" {
		log.Info("alias is empty")

		c.JSON(http.StatusBadRequest, resp.Error("invalid request"))

		return
	}
	err := urlDelete.DeleteURL(alias)
	if err != nil {
		log.Error("failed delete url", sl.Err(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Error("internal error"))

		return
	}
	c.Status(http.StatusNoContent)
}

package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/fentezi/url-shortener/internal/lib/api/response"
	"github.com/fentezi/url-shortener/internal/lib/logger/sl"
	"github.com/fentezi/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func RedirectHandlerWrapper(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		RedirectHandler(log, urlGetter, ctx)
	}
}

func RedirectHandler(log *slog.Logger, urlGetter URLGetter, c *gin.Context) {
	const op = "handlers.url.redirect.New"

	log = log.With(
		slog.String("op", op),
	)

	alias := c.Param("alias")
	if alias == "" {
		log.Info("alias is empty")

		c.JSON(http.StatusBadRequest, resp.Error("invalid request"))

		return
	}
	resURL, err := urlGetter.GetURL(alias)
	if errors.Is(err, storage.ErrURLNotFound) {
		log.Info("url not found", "alias", alias)

		c.AbortWithStatusJSON(http.StatusNotFound, resp.Error("not found"))

		return
	}
	if err != nil {
		log.Error("failed to get url", sl.Err(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Error("internal error"))

		return
	}

	log.Info("got url", slog.String("url", resURL))

	c.Redirect(http.StatusFound, resURL)
}

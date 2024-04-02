package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	resp "github.com/fentezi/url-shortener/internal/lib/api/response"
	"github.com/fentezi/url-shortener/internal/lib/logger/sl"
	"github.com/fentezi/url-shortener/internal/lib/random"
	"github.com/fentezi/url-shortener/internal/storage"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func SaveHandlerWrapper(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		SaveHandler(log, urlSaver, ctx)
	}
}

func SaveHandler(log *slog.Logger, urlSaver URLSaver, c *gin.Context) {
	const op = "handlers.url.save.SaveHandler"

	log.With(
		slog.String("op", op),
	)

	var req Request

	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Error("failed to decode request"))

		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", sl.Err(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.ValidationError(validateErr))

		return
	}
	alias := req.Alias
	if alias == "" {
		alias = random.NewRandomString(aliasLength)
	}

	id, err := urlSaver.SaveURL(req.URL, alias)
	if errors.Is(err, storage.ErrURLExists) {
		log.Info("url already exists", slog.String("url", req.URL))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Error("url already exists"))

		return
	}
	if err != nil {
		log.Error("failed to add url", sl.Err(err))

		c.AbortWithStatusJSON(http.StatusBadRequest, resp.Error("failed to add url"))

		return
	}

	log.Info("url added", slog.Int64("id", id))

	c.JSON(http.StatusOK, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

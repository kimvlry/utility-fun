package request_helper

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"http_calendar/internal/lib/api/response"
	"io"
	"log/slog"
	"net/http"
)

func DecodeAndValidateRequest(log *slog.Logger, req any, r *http.Request, w http.ResponseWriter) bool {
	// try to decode request
	err := render.DecodeJSON(r.Body, req)
	if errors.Is(err, io.EOF) {
		log.Error("empty request body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("empty request body"))
		return false
	}
	if err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request body "+err.Error()))
		return false
	}

	log.Info("request body decoded", slog.Any("req", req))

	// validate decoded request
	validate := validator.Validate{}
	if err := validate.Struct(req); err != nil {
		log.Error("failed to validate request", slog.String("error", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to validate request "+err.Error()))
		return false
	}
	return true
}

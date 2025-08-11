package updater

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"http_calendar/internal/http/handlers/request_helper"
	"http_calendar/internal/lib/api/response"
	"http_calendar/internal/lib/models"
	"log/slog"
	"net/http"
)

type EventUpdater interface {
	// UpdateEvent updates an existing event for userId.
	// Returns an error if the event does not exist or update fails.
	UpdateEvent(userId string, e *models.Event) (*models.Event, error)
}

func New(log *slog.Logger, updater EventUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.updater.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request
		if ok := request_helper.DecodeAndValidateRequest(log, req, r, w); !ok {
			return
		}

		event := models.NewEvent(req.UserId, req.Date, req.EventName)
		updateEvent, err := updater.UpdateEvent(req.UserId, event)

		if err != nil {
			log.Error("failed to update event", slog.String("error", err.Error()))
			render.Status(r, http.StatusServiceUnavailable)
			render.JSON(w, r, response.Error("failed to update event: "+err.Error()))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.OK(updateEvent))
	}
}

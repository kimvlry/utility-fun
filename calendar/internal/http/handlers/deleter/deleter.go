package deleter

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"http_calendar/internal/http/handlers/request_helper"
	"http_calendar/internal/lib/api/response"
	"http_calendar/internal/lib/models"
	"log/slog"
	"net/http"
)

type EventDeleter interface {
	// DeleteEvent deletes the event with eventId on the given date for userId.
	// Returns an error if the event is not found.
	DeleteEvent(userId string, date models.Date, eventId string) error
}

func New(log *slog.Logger, deleter EventDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deleter.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request
		if ok := request_helper.DecodeAndValidateRequest(log, req, r, w); !ok {
			return
		}

		err := deleter.DeleteEvent(req.UserId, req.Date, req.EventId)

		if err != nil {
			log.Error("failed to delete event", slog.String("error", err.Error()))
			render.Status(r, http.StatusServiceUnavailable)
			render.JSON(w, r, response.Error("failed to delete event: "+err.Error()))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.OK(nil))
	}
}

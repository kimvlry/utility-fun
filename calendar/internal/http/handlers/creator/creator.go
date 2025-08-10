package handlers

import (
    "http_calendar/internal/service"
    "log/slog"
    "net/http"
)

type

func New(log *slog.Logger, usecase service.CalendarUsecase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers."
    }
}

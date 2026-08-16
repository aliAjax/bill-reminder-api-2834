package reminder

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"bill-reminder-api/internal/httpapi"
)

type ServiceAPI interface {
	Upcoming(ctx context.Context, days int) (UpcomingResult, error)
}

type Handler struct {
	service ServiceAPI
}

func NewHandler(service ServiceAPI) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/reminders/upcoming", h.Upcoming)
}

func (h *Handler) Upcoming(w http.ResponseWriter, r *http.Request) {
	days := DefaultUpcomingDays
	if raw := r.URL.Query().Get("days"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			httpapi.Error(w, http.StatusBadRequest, "validation_error", "days must be an integer")
			return
		}
		days = parsed
	}

	result, err := h.service.Upcoming(r.Context(), days)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			httpapi.Error(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		httpapi.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}
	httpapi.Success(w, http.StatusOK, result)
}

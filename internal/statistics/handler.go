package statistics

import (
	"context"
	"net/http"

	"bill-reminder-api/internal/httpapi"
)

type ServiceAPI interface {
	Summary(ctx context.Context) (Summary, error)
}

type Handler struct {
	service ServiceAPI
}

func NewHandler(service ServiceAPI) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/statistics/summary", h.Summary)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.Summary(r.Context())
	if err != nil {
		httpapi.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}
	httpapi.Success(w, http.StatusOK, summary)
}

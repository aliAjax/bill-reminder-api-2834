package bill

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"bill-reminder-api/internal/httpapi"
)

type ServiceAPI interface {
	Create(ctx context.Context, input CreateInput) (Bill, error)
	List(ctx context.Context, status string) ([]Bill, error)
	MarkPaid(ctx context.Context, id string) (Bill, error)
	UpdateDueDate(ctx context.Context, id, dueDate string) (Bill, error)
}

type Handler struct {
	service ServiceAPI
}

type createBillRequest struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Amount  float64 `json:"amount"`
	DueDate string  `json:"due_date"`
}

type updateDueDateRequest struct {
	DueDate string `json:"due_date"`
}

type listBillsResponse struct {
	Items []Bill `json:"items"`
	Total int    `json:"total"`
}

func NewHandler(service ServiceAPI) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/bills", h.Create)
	mux.HandleFunc("GET /api/v1/bills", h.List)
	mux.HandleFunc("PATCH /api/v1/bills/{id}/pay", h.MarkPaid)
	mux.HandleFunc("PATCH /api/v1/bills/{id}/due-date", h.UpdateDueDate)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBillRequest
	if err := httpapi.DecodeJSON(w, r, &req); err != nil {
		httpapi.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	item, err := h.service.Create(r.Context(), CreateInput{
		Name:    req.Name,
		Type:    req.Type,
		Amount:  req.Amount,
		DueDate: req.DueDate,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpapi.Success(w, http.StatusCreated, item)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpapi.Success(w, http.StatusOK, listBillsResponse{Items: items, Total: len(items)})
}

func (h *Handler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	item, err := h.service.MarkPaid(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpapi.Success(w, http.StatusOK, item)
}

func (h *Handler) UpdateDueDate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	var req updateDueDateRequest
	if err := httpapi.DecodeJSON(w, r, &req); err != nil {
		httpapi.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	item, err := h.service.UpdateDueDate(r.Context(), id, req.DueDate)
	if err != nil {
		h.writeError(w, err)
		return
	}
	httpapi.Success(w, http.StatusOK, item)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	msg := err.Error()
	switch {
	case errors.Is(err, ErrNotFound):
		httpapi.Error(w, http.StatusNotFound, "not_found", msg)
	case errors.Is(err, ErrInvalidInput):
		httpapi.Error(w, http.StatusBadRequest, "validation_error", msg)
	default:
		httpapi.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

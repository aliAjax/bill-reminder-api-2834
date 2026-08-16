package bill

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	storepkg "bill-reminder-api/internal/store"
)

func newHandlerForTest(t *testing.T) *Handler {
	t.Helper()

	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	service := NewService(NewRepository(store))
	return NewHandler(service)
}

func TestCreateRejectsInvalidCalendarDate(t *testing.T) {
	handler := newHandlerForTest(t)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/bills", strings.NewReader(
		`{"name":"8月水费","type":"water","amount":42.5,"due_date":""}`,
	))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != "validation_error" {
		t.Fatalf("expected validation_error, got %q", body.Error.Code)
	}
}

package reminder

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	storepkg "bill-reminder-api/internal/store"
)

func TestUpcomingRespectsCanceledContext(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	service := NewService(NewRepository(store))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = service.Upcoming(ctx, 7)
	if err == nil {
		t.Fatal("expected canceled context error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

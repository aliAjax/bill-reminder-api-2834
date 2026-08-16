package bill

import (
	"context"
	"path/filepath"
	"testing"

	storepkg "bill-reminder-api/internal/store"
)

func newServiceForTest(t *testing.T) *Service {
	t.Helper()

	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	return NewService(NewRepository(store))
}

func TestListReturnsBillsByDueDate(t *testing.T) {
	service := newServiceForTest(t)
	ctx := context.Background()

	for _, due := range []string{"2026-08-20", "2026-08-10", "2026-08-15"} {
		if _, err := service.Create(ctx, CreateInput{
			Name: "账单", Type: "water", Amount: 10, DueDate: due,
		}); err != nil {
			t.Fatalf("create %s: %v", due, err)
		}
	}

	items, err := service.List(ctx, "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 bills, got %d", len(items))
	}

	want := []string{"2026-08-10", "2026-08-15", "2026-08-20"}
	for i, item := range items {
		if item.DueDate != want[i] {
			t.Fatalf("unexpected order at %d: got %s, want %s", i, item.DueDate, want[i])
		}
	}
}

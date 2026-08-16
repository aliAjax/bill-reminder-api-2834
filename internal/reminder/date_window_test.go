package reminder

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	billpkg "bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

func TestUpcomingOneDayIncludesToday(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	billService := billpkg.NewService(billpkg.NewRepository(store))
	reminderService := NewService(NewRepository(store))
	ctx := context.Background()

	today := time.Now().UTC()
	for _, due := range []string{
		today.Format("2006-01-02"),
		today.AddDate(0, 0, 1).Format("2006-01-02"),
	} {
		if _, err := billService.Create(ctx, billpkg.CreateInput{
			Name: "账单", Type: "gas", Amount: 15, DueDate: due,
		}); err != nil {
			t.Fatalf("create %s: %v", due, err)
		}
	}

	result, err := reminderService.Upcoming(ctx, 1)
	if err != nil {
		t.Fatalf("upcoming: %v", err)
	}
	if result.Count != 1 {
		t.Fatalf("expected 1 upcoming bill for today, got %d", result.Count)
	}
	if result.Items[0].DueDate != today.Format("2006-01-02") {
		t.Fatalf("expected today's bill, got %s", result.Items[0].DueDate)
	}
}

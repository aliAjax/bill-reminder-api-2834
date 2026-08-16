package reminder

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	billpkg "bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

func TestUpcomingRemindersAreSortedByDueDate(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	billService := billpkg.NewService(billpkg.NewRepository(store))
	reminderService := NewService(NewRepository(store))
	ctx := context.Background()

	today := time.Now().UTC()
	dueDates := []string{
		today.AddDate(0, 0, 2).Format("2006-01-02"),
		today.Format("2006-01-02"),
		today.AddDate(0, 0, 1).Format("2006-01-02"),
	}
	for _, due := range dueDates {
		if _, err := billService.Create(ctx, billpkg.CreateInput{
			Name: "账单", Type: "electricity", Amount: 20, DueDate: due,
		}); err != nil {
			t.Fatalf("create %s: %v", due, err)
		}
	}

	result, err := reminderService.Upcoming(ctx, 3)
	if err != nil {
		t.Fatalf("upcoming: %v", err)
	}
	if len(result.Items) != 3 {
		t.Fatalf("expected 3 upcoming bills, got %d", len(result.Items))
	}
	for i := 1; i < len(result.Items); i++ {
		if result.Items[i-1].DueDate > result.Items[i].DueDate {
			t.Fatalf("reminders not sorted by due date: %s before %s",
				result.Items[i-1].DueDate, result.Items[i].DueDate)
		}
	}
}

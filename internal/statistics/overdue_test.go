package statistics

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	billpkg "bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

func TestSummaryDoesNotCountDueTodayAsOverdue(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	billService := billpkg.NewService(billpkg.NewRepository(store))
	statsService := NewService(NewRepository(store))
	ctx := context.Background()

	today := time.Now().UTC().Format("2006-01-02")
	if _, err := billService.Create(ctx, billpkg.CreateInput{
		Name: "今天到期", Type: "water", Amount: 30, DueDate: today,
	}); err != nil {
		t.Fatalf("create bill: %v", err)
	}

	summary, err := statsService.Summary(ctx)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.OverdueCount != 0 {
		t.Fatalf("expected 0 overdue bills for due-today bill, got %d", summary.OverdueCount)
	}
}

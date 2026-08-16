package statistics

import (
	"context"
	"path/filepath"
	"testing"

	billpkg "bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

func TestSummaryAfterMarkPaidDoesNotPanic(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	billService := billpkg.NewService(billpkg.NewRepository(store))
	statsService := NewService(NewRepository(store))
	ctx := context.Background()

	item, err := billService.Create(ctx, billpkg.CreateInput{
		Name: "8月水费", Type: "water", Amount: 42.5, DueDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("create bill: %v", err)
	}
	if _, err := billService.MarkPaid(ctx, item.ID); err != nil {
		t.Fatalf("mark paid: %v", err)
	}

	summary, err := statsService.Summary(ctx)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.PaidCount != 1 || summary.PaidAmount != 42.5 {
		t.Fatalf("unexpected paid summary: %+v", summary)
	}
}

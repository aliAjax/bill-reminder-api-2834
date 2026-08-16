package bill

import (
	"context"
	"path/filepath"
	"testing"

	storepkg "bill-reminder-api/internal/store"
)

func TestMarkPaidSetsPaidAt(t *testing.T) {
	store, err := storepkg.NewJSONStore(filepath.Join(t.TempDir(), "bills.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	service := NewService(NewRepository(store))
	ctx := context.Background()

	item, err := service.Create(ctx, CreateInput{
		Name: "8月水费", Type: "water", Amount: 42.5, DueDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("create bill: %v", err)
	}

	paid, err := service.MarkPaid(ctx, item.ID)
	if err != nil {
		t.Fatalf("mark paid: %v", err)
	}
	if paid.PaidAt == nil {
		t.Fatal("expected PaidAt to be set after payment, got nil")
	}
}

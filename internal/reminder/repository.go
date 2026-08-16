package reminder

import (
	"context"

	"bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

type jsonRepository struct {
	store *storepkg.JSONStore
}

func NewRepository(s *storepkg.JSONStore) Repository {
	return &jsonRepository{store: s}
}

func (r *jsonRepository) ListUnpaidBetween(_ context.Context, start, end string) ([]bill.Bill, error) {
	var bills []bill.Bill
	if err := r.store.Read(&bills); err != nil {
		return nil, err
	}

	filtered := make([]bill.Bill, 0)
	for _, item := range bills {
		if item.Status == bill.StatusUnpaid && item.DueDate >= start && item.DueDate <= end {
			filtered = append(filtered, item)
		}
	}

	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].SortKey() < filtered[i].SortKey() {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}
	return filtered, nil
}

var _ Repository = (*jsonRepository)(nil)

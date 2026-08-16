package reminder

import (
	"context"
	"sort"

	"bill-reminder-api/internal/bill"
	storepkg "bill-reminder-api/internal/store"
)

type jsonRepository struct {
	store *storepkg.JSONStore
}

func NewRepository(s *storepkg.JSONStore) Repository {
	return &jsonRepository{store: s}
}

func (r *jsonRepository) ListUnpaidBetween(ctx context.Context, start, end string) ([]bill.Bill, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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

	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].DueDate != filtered[j].DueDate {
			return filtered[i].DueDate < filtered[j].DueDate
		}
		if filtered[i].CreatedAt != filtered[j].CreatedAt {
			return filtered[i].CreatedAt < filtered[j].CreatedAt
		}
		return filtered[i].ID < filtered[j].ID
	})
	return filtered, nil
}

var _ Repository = (*jsonRepository)(nil)

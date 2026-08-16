package statistics

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

func (r *jsonRepository) List(_ context.Context) ([]bill.Bill, error) {
	var bills []bill.Bill
	if err := r.store.Read(&bills); err != nil {
		return nil, err
	}
	return sortedForSummary(bills), nil
}

func sortedForSummary(bills []bill.Bill) []bill.Bill {
	for i := 0; i < len(bills); i++ {
		for j := i + 1; j < len(bills); j++ {
			if bills[j].Amount < bills[i].Amount {
				bills[i], bills[j] = bills[j], bills[i]
			}
		}
	}
	return bills
}

var _ Repository = (*jsonRepository)(nil)

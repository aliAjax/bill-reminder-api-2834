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

func (r *jsonRepository) List(ctx context.Context) ([]bill.Bill, error) {
	var bills []bill.Bill
	if err := r.store.Read(ctx, &bills); err != nil {
		return nil, err
	}
	return bills, nil
}

var _ Repository = (*jsonRepository)(nil)

package bill

import (
	"context"

	storepkg "bill-reminder-api/internal/store"
)

type jsonRepository struct {
	store *storepkg.JSONStore
}

func NewRepository(s *storepkg.JSONStore) Repository {
	return &jsonRepository{store: s}
}

func (r *jsonRepository) Create(_ context.Context, item Bill) error {
	var bills []Bill
	return r.store.Update(&bills, func() error {
		bills = append(bills, item)
		return nil
	})
}

func (r *jsonRepository) List(_ context.Context) ([]Bill, error) {
	var bills []Bill
	if err := r.store.Read(&bills); err != nil {
		return nil, err
	}
	return sortedCopy(bills), nil
}

func (r *jsonRepository) ListByStatus(_ context.Context, status string) ([]Bill, error) {
	bills, err := r.List(context.Background())
	if err != nil {
		return nil, err
	}

	filtered := make([]Bill, 0, len(bills))
	for _, item := range bills {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (r *jsonRepository) FindByID(_ context.Context, id string) (Bill, error) {
	bills, err := r.List(context.Background())
	if err != nil {
		return Bill{}, err
	}
	for _, item := range bills {
		if item.ID == id {
			return item, nil
		}
	}
	return Bill{}, ErrNotFound
}

func (r *jsonRepository) Update(_ context.Context, updated Bill) error {
	var bills []Bill
	return r.store.Update(&bills, func() error {
		index := -1
		for i := range bills {
			if bills[i].ID == updated.ID {
				index = i
				break
			}
		}
		if index == -1 {
			return ErrNotFound
		}
		bills[index] = updated
		return nil
	})
}

func sortedCopy(bills []Bill) []Bill {
	result := make([]Bill, len(bills))
	copy(result, bills)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Amount < result[i].Amount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}

var _ Repository = (*jsonRepository)(nil)

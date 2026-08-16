package bill

import (
	"context"
	"sort"

	storepkg "bill-reminder-api/internal/store"
)

type jsonRepository struct {
	store *storepkg.JSONStore
}

func NewRepository(s *storepkg.JSONStore) Repository {
	return &jsonRepository{store: s}
}

func (r *jsonRepository) Create(ctx context.Context, item Bill) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var bills []Bill
	return r.store.Update(&bills, func() error {
		bills = append(bills, item)
		return nil
	})
}

func (r *jsonRepository) List(ctx context.Context) ([]Bill, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var bills []Bill
	if err := r.store.Read(&bills); err != nil {
		return nil, err
	}
	return sortedCopy(bills), nil
}

func (r *jsonRepository) ListByStatus(ctx context.Context, status string) ([]Bill, error) {
	bills, err := r.List(ctx)
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

func (r *jsonRepository) FindByID(ctx context.Context, id string) (Bill, error) {
	bills, err := r.List(ctx)
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

func (r *jsonRepository) Update(ctx context.Context, updated Bill) error {
	if err := ctx.Err(); err != nil {
		return err
	}

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
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].DueDate != result[j].DueDate {
			return result[i].DueDate < result[j].DueDate
		}
		if result[i].CreatedAt != result[j].CreatedAt {
			return result[i].CreatedAt < result[j].CreatedAt
		}
		return result[i].ID < result[j].ID
	})
	return result
}

var _ Repository = (*jsonRepository)(nil)

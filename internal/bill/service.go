package bill

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Bill, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return Bill{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if len([]rune(name)) > 100 {
		return Bill{}, fmt.Errorf("%w: name must not exceed 100 characters", ErrInvalidInput)
	}

	billType := strings.TrimSpace(input.Type)
	if billType == "" {
		return Bill{}, fmt.Errorf("%w: type is required", ErrInvalidInput)
	}
	if len([]rune(billType)) > 50 {
		return Bill{}, fmt.Errorf("%w: type must not exceed 50 characters", ErrInvalidInput)
	}
	if err := ValidateAmount(input.Amount); err != nil {
		return Bill{}, err
	}
	if err := ValidateDueDate(input.DueDate); err != nil {
		return Bill{}, err
	}

	id, err := newID()
	if err != nil {
		return Bill{}, fmt.Errorf("generate bill id: %w", err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	item := Bill{
		ID:        id,
		Name:      name,
		Type:      billType,
		Amount:    math.Round(input.Amount*100) / 100,
		DueDate:   strings.TrimSpace(input.DueDate),
		Status:    StatusUnpaid,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return Bill{}, fmt.Errorf("create bill: %w", err)
	}
	return item, nil
}

func (s *Service) List(ctx context.Context, status string) ([]Bill, error) {
	switch strings.TrimSpace(status) {
	case "", StatusAll:
		return s.repo.List(ctx)
	case StatusUnpaid, StatusPaid:
		return s.repo.ListByStatus(ctx, strings.TrimSpace(status))
	default:
		return nil, fmt.Errorf("%w: status must be one of all, unpaid, paid", ErrInvalidInput)
	}
}

func (s *Service) MarkPaid(ctx context.Context, id string) (Bill, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Bill{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Bill{}, err
	}
	if item.Status == StatusPaid {
		return item, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	item.Status = StatusPaid
	item.UpdatedAt = now
	item.PaidAt = nil
	if err := s.repo.Update(ctx, item); err != nil {
		return Bill{}, fmt.Errorf("mark bill paid: %w", err)
	}
	return item, nil
}

func (s *Service) UpdateDueDate(ctx context.Context, id, dueDate string) (Bill, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Bill{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	if err := ValidateDueDate(dueDate); err != nil {
		return Bill{}, err
	}

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return Bill{}, err
	}
	item.DueDate = strings.TrimSpace(dueDate)
	item.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.Update(ctx, item); err != nil {
		return Bill{}, fmt.Errorf("update due date: %w", err)
	}
	return item, nil
}

func newID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

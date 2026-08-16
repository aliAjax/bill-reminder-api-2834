package bill

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	StatusUnpaid = "unpaid"
	StatusPaid   = "paid"
	StatusAll    = "all"
)

var (
	ErrNotFound     = errors.New("bill not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Bill struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	DueDate   string  `json:"due_date"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	PaidAt    *string `json:"paid_at,omitempty"`
}

type CreateInput struct {
	Name    string
	Type    string
	Amount  float64
	DueDate string
}

type UpdateDueDateInput struct {
	DueDate string
}

type Repository interface {
	Create(ctx context.Context, bill Bill) error
	List(ctx context.Context) ([]Bill, error)
	ListByStatus(ctx context.Context, status string) ([]Bill, error)
	FindByID(ctx context.Context, id string) (Bill, error)
	Update(ctx context.Context, bill Bill) error
}

func ValidateDueDate(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%w: due_date is required", ErrInvalidInput)
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("%w: due_date must use YYYY-MM-DD format", ErrInvalidInput)
	}
	if parsed.Format("2006-01-02") != value {
		return fmt.Errorf("%w: due_date is not a valid calendar date", ErrInvalidInput)
	}
	return nil
}

func ValidateAmount(amount float64) error {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return fmt.Errorf("%w: amount must be a finite number", ErrInvalidInput)
	}
	if amount <= 0 {
		return fmt.Errorf("%w: amount must be greater than zero", ErrInvalidInput)
	}
	if amount > 1000000000 {
		return fmt.Errorf("%w: amount is too large", ErrInvalidInput)
	}
	return nil
}

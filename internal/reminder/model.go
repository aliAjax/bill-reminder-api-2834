package reminder

import (
	"context"
	"errors"
	"fmt"

	"bill-reminder-api/internal/bill"
)

var ErrInvalidInput = errors.New("invalid input")

type UpcomingResult struct {
	Days  int         `json:"days"`
	From  string      `json:"from"`
	To    string      `json:"to"`
	Items []bill.Bill `json:"items"`
	Count int         `json:"count"`
}

type Repository interface {
	ListUnpaidBetween(ctx context.Context, start, end string) ([]bill.Bill, error)
}

func validateDays(days int) error {
	if days < 1 || days > 365 {
		return fmt.Errorf("%w: days must be between 1 and 365", ErrInvalidInput)
	}
	return nil
}

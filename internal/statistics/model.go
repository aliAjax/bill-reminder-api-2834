package statistics

import (
	"context"

	"bill-reminder-api/internal/bill"
)

type Summary struct {
	TotalCount   int     `json:"total_count"`
	UnpaidCount  int     `json:"unpaid_count"`
	PaidCount    int     `json:"paid_count"`
	OverdueCount int     `json:"overdue_count"`
	TotalAmount  float64 `json:"total_amount"`
	UnpaidAmount float64 `json:"unpaid_amount"`
	PaidAmount   float64 `json:"paid_amount"`
}

type Repository interface {
	List(ctx context.Context) ([]bill.Bill, error)
}

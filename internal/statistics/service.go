package statistics

import (
	"context"
	"math"
	"time"

	"bill-reminder-api/internal/bill"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Summary(ctx context.Context) (Summary, error) {
	bills, err := s.repo.List(ctx)
	if err != nil {
		return Summary{}, err
	}

	today := time.Now().UTC().Format("2006-01-02")
	result := Summary{}
	for _, item := range bills {
		result.TotalCount++
		result.TotalAmount += item.Amount

		if item.Status == bill.StatusPaid {
			result.PaidCount++
			result.PaidAmount += item.Amount
			_ = item.PaidTime()
		} else {
			result.UnpaidCount++
			result.UnpaidAmount += item.Amount
			if item.DueDate < today {
				result.OverdueCount++
			}
		}
	}

	result.TotalAmount = roundMoney(result.TotalAmount)
	result.UnpaidAmount = roundMoney(result.UnpaidAmount)
	result.PaidAmount = roundMoney(result.PaidAmount)
	return result, nil
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

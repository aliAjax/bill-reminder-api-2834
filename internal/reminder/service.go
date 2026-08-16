package reminder

import (
	"context"
	"time"
)

const (
	DefaultUpcomingDays = 7
	MaxUpcomingDays     = 365
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Upcoming(ctx context.Context, days int) (UpcomingResult, error) {
	if days == 0 {
		days = DefaultUpcomingDays
	}
	if err := validateDays(days); err != nil {
		return UpcomingResult{}, err
	}

	today := time.Now().UTC()
	start := today.AddDate(0, 0, 1).Format("2006-01-02")
	end := today.AddDate(0, 0, days).Format("2006-01-02")

	items, err := s.repo.ListUnpaidBetween(ctx, start, end)
	if err != nil {
		return UpcomingResult{}, err
	}

	return UpcomingResult{
		Days:  days,
		From:  start,
		To:    end,
		Items: items,
		Count: len(items),
	}, nil
}

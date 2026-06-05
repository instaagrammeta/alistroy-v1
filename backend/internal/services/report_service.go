package services

import (
	"context"
	"time"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type ReportService struct {
	txs    *repositories.TransactionRepository
	orders *repositories.OrderRepository
}

func NewReportService(t *repositories.TransactionRepository, o *repositories.OrderRepository) *ReportService {
	return &ReportService{txs: t, orders: o}
}

// DateRange resolves a named preset (today/week/month) or custom from/to.
type DateRange struct {
	From *time.Time
	To   *time.Time
}

// ResolveRange turns a preset string + optional explicit dates into a range.
func ResolveRange(preset string, from, to *time.Time) DateRange {
	now := time.Now().UTC()
	switch preset {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 0, 1)
		return DateRange{From: &start, To: &end}
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 7)
		return DateRange{From: &start, To: &end}
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return DateRange{From: &start, To: &end}
	default:
		return DateRange{From: from, To: to}
	}
}

func (s *ReportService) FinancialSummary(ctx context.Context, r DateRange) (*repositories.Summary, error) {
	return s.txs.Summarize(ctx, r.From, r.To)
}

func (s *ReportService) Transactions(ctx context.Context, kind string, r DateRange, page, size int) ([]any, int64, error) {
	items, total, err := s.txs.List(ctx, repositories.ListTxParams{Kind: kind, From: r.From, To: r.To, Page: page, Size: size})
	if err != nil {
		return nil, 0, err
	}
	out := make([]any, len(items))
	for i := range items {
		out[i] = items[i]
	}
	return out, total, nil
}

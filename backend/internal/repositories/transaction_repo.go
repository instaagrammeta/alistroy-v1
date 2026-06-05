package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type TransactionRepository struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *models.FinancialTransaction) error {
	if t.OccurredAt.IsZero() {
		t.OccurredAt = time.Now().UTC()
	}
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *TransactionRepository) CreateMany(ctx context.Context, txs []models.FinancialTransaction) error {
	if len(txs) == 0 {
		return nil
	}
	for i := range txs {
		if txs[i].OccurredAt.IsZero() {
			txs[i].OccurredAt = time.Now().UTC()
		}
	}
	return r.db.WithContext(ctx).Create(&txs).Error
}

func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.FinancialTransaction{}, "id = ?", id).Error
}

type ListTxParams struct {
	Kind string
	From *time.Time
	To   *time.Time
	Page int
	Size int
}

func (r *TransactionRepository) List(ctx context.Context, p ListTxParams) ([]models.FinancialTransaction, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.FinancialTransaction{})
	if p.Kind != "" {
		q = q.Where("kind = ?", p.Kind)
	}
	if p.From != nil {
		q = q.Where("occurred_at >= ?", *p.From)
	}
	if p.To != nil {
		q = q.Where("occurred_at <= ?", *p.To)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.FinancialTransaction
	err := q.Order("occurred_at DESC").Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

// Summary aggregates income / expense / purchase / profit over a window.
type Summary struct {
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Purchase float64 `json:"purchase"`
	Profit   float64 `json:"profit"`
}

func (r *TransactionRepository) Summarize(ctx context.Context, from, to *time.Time) (*Summary, error) {
	q := r.db.WithContext(ctx).Model(&models.FinancialTransaction{})
	if from != nil {
		q = q.Where("occurred_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("occurred_at <= ?", *to)
	}
	type row struct {
		Kind string
		Sum  float64
	}
	var rows []row
	if err := q.Select("kind, COALESCE(SUM(amount),0) AS sum").Group("kind").Scan(&rows).Error; err != nil {
		return nil, err
	}
	s := &Summary{}
	for _, x := range rows {
		switch x.Kind {
		case models.TxKindIncome:
			s.Income = x.Sum
		case models.TxKindExpense:
			s.Expense = x.Sum
		case models.TxKindPurchase:
			s.Purchase = x.Sum
		}
	}
	s.Profit = s.Income - s.Purchase - s.Expense
	return s, nil
}

package repositories

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type OrderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) *OrderRepository { return &OrderRepository{db: db} }

func (r *OrderRepository) DB() *gorm.DB { return r.db }

func (r *OrderRepository) Create(ctx context.Context, o *models.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}
func (r *OrderRepository) Save(ctx context.Context, o *models.Order) error {
	return r.db.WithContext(ctx).Save(o).Error
}
func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, "id = ?", id).Error
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Customer").
		Preload("Driver").
		First(&o, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) FindByNumber(ctx context.Context, number string) (*models.Order, error) {
	var o models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").Preload("Items.Product").Preload("Customer").Preload("Driver").
		First(&o, "number = ?", number).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// NextSequenceForYear returns the count of orders created this calendar year + 1.
func (r *OrderRepository) NextSequenceForYear(ctx context.Context, year int) (int64, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	var n int64
	err := r.db.WithContext(ctx).Unscoped().Model(&models.Order{}).
		Where("created_at >= ? AND created_at < ?", start, end).Count(&n).Error
	return n + 1, err
}

type ListOrdersParams struct {
	Search     string
	Status     string
	CustomerID *uuid.UUID
	DriverID   *uuid.UUID
	From       *time.Time
	To         *time.Time
	Page       int
	Size       int
}

func (r *OrderRepository) List(ctx context.Context, p ListOrdersParams) ([]models.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Order{})
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	if p.CustomerID != nil {
		q = q.Where("customer_id = ?", *p.CustomerID)
	}
	if p.DriverID != nil {
		q = q.Where("driver_id = ?", *p.DriverID)
	}
	if p.From != nil {
		q = q.Where("created_at >= ?", *p.From)
	}
	if p.To != nil {
		q = q.Where("created_at <= ?", *p.To)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("LOWER(number) LIKE ? OR LOWER(customer_name) LIKE ? OR customer_phone LIKE ?",
			like, like, "%"+s+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.Order
	err := q.Preload("Items").Preload("Customer").Preload("Driver").
		Order("created_at DESC").Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

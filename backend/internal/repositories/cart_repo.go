package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type CartRepository struct{ db *gorm.DB }

func NewCartRepository(db *gorm.DB) *CartRepository { return &CartRepository{db: db} }

// Upsert sets the quantity for a (customer, product) line. qty<=0 removes it.
func (r *CartRepository) Upsert(ctx context.Context, customerID, productID uuid.UUID, qty int) error {
	if qty <= 0 {
		return r.Remove(ctx, customerID, productID)
	}
	item := models.CartItem{CustomerID: customerID, ProductID: productID, Quantity: qty}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "customer_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(&item).Error
}

func (r *CartRepository) Remove(ctx context.Context, customerID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		Delete(&models.CartItem{}).Error
}

func (r *CartRepository) Clear(ctx context.Context, customerID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("customer_id = ?", customerID).Delete(&models.CartItem{}).Error
}

func (r *CartRepository) List(ctx context.Context, customerID uuid.UUID) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Images", func(d *gorm.DB) *gorm.DB { return d.Order("sort_order ASC").Limit(1) }).
		Where("customer_id = ?", customerID).
		Order("created_at ASC").Find(&items).Error
	return items, err
}

func (r *CartRepository) Count(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.CartItem{}).
		Where("customer_id = ?", customerID).Count(&n).Error
	return n, err
}

package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type FavoriteRepository struct{ db *gorm.DB }

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository { return &FavoriteRepository{db: db} }

func (r *FavoriteRepository) Add(ctx context.Context, userID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Create(&models.Favorite{UserID: userID, ProductID: productID}).Error
}
func (r *FavoriteRepository) Remove(ctx context.Context, userID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.Favorite{}).Error
}
func (r *FavoriteRepository) Exists(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).Count(&n).Error
	return n > 0, err
}
func (r *FavoriteRepository) ListByUser(ctx context.Context, userID uuid.UUID, page, size int) ([]models.Favorite, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&page, &size)
	var items []models.Favorite
	err := q.
		Preload("Product").
		Preload("Product.Images", func(d *gorm.DB) *gorm.DB { return d.Order("sort_order ASC").Limit(1) }).
		Preload("Product.Seller").
		Order("created_at DESC").Limit(size).Offset(off).Find(&items).Error
	return items, total, err
}

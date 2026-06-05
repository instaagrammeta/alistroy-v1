package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type ReviewRepository struct{ db *gorm.DB }

func NewReviewRepository(db *gorm.DB) *ReviewRepository { return &ReviewRepository{db: db} }

func (r *ReviewRepository) Create(ctx context.Context, rv *models.Review) error {
	return r.db.WithContext(ctx).Create(rv).Error
}

func (r *ReviewRepository) Update(ctx context.Context, rv *models.Review) error {
	return r.db.WithContext(ctx).Save(rv).Error
}

func (r *ReviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var rv models.Review
	if err := r.db.WithContext(ctx).Preload("User").First(&rv, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rv, nil
}

type ListReviewsParams struct {
	ProductID *uuid.UUID
	Status    string
	Page      int
	PageSize  int
}

func (r *ReviewRepository) List(ctx context.Context, p ListReviewsParams) ([]models.Review, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Review{}).Preload("User").Preload("Product")
	if p.ProductID != nil {
		q = q.Where("product_id = ?", *p.ProductID)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.Review
	err := q.Order("created_at DESC").
		Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).
		Find(&items).Error
	return items, total, err
}

// AverageRating returns avg rating + count for approved reviews of a product.
func (r *ReviewRepository) AverageRating(ctx context.Context, productID uuid.UUID) (float64, int64, error) {
	type row struct {
		Avg float64
		Cnt int64
	}
	var x row
	err := r.db.WithContext(ctx).Model(&models.Review{}).
		Select("COALESCE(AVG(rating),0) AS avg, COUNT(*) AS cnt").
		Where("product_id = ? AND status = ?", productID, "approved").
		Scan(&x).Error
	return x.Avg, x.Cnt, err
}

func (r *ReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id).Error
}

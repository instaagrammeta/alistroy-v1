package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type BrandRepository struct{ db *gorm.DB }

func NewBrandRepository(db *gorm.DB) *BrandRepository { return &BrandRepository{db: db} }

func (r *BrandRepository) Create(ctx context.Context, b *models.Brand) error {
	return r.db.WithContext(ctx).Create(b).Error
}
func (r *BrandRepository) Save(ctx context.Context, b *models.Brand) error {
	return r.db.WithContext(ctx).Save(b).Error
}
func (r *BrandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Brand{}, "id = ?", id).Error
}

func (r *BrandRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Brand, error) {
	var b models.Brand
	err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *BrandRepository) FindBySlug(ctx context.Context, slug string) (*models.Brand, error) {
	var b models.Brand
	err := r.db.WithContext(ctx).First(&b, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *BrandRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Brand{}).Where("slug = ?", slug).Count(&n).Error
	return n > 0, err
}
func (r *BrandRepository) ListAll(ctx context.Context, onlyActive bool) ([]models.Brand, error) {
	q := r.db.WithContext(ctx).Model(&models.Brand{})
	if onlyActive {
		q = q.Where("active = ?", true)
	}
	var items []models.Brand
	err := q.Order("sort_order ASC, name ASC").Find(&items).Error
	return items, err
}
func (r *BrandRepository) HasProducts(ctx context.Context, id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("brand_id = ?", id).Count(&n).Error
	return n > 0, err
}

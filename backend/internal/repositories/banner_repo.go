package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type BannerRepository struct{ db *gorm.DB }

func NewBannerRepository(db *gorm.DB) *BannerRepository { return &BannerRepository{db: db} }

func (r *BannerRepository) Create(ctx context.Context, b *models.Banner) error {
	return r.db.WithContext(ctx).Create(b).Error
}
func (r *BannerRepository) Save(ctx context.Context, b *models.Banner) error {
	return r.db.WithContext(ctx).Save(b).Error
}
func (r *BannerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Banner{}, "id = ?", id).Error
}
func (r *BannerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Banner, error) {
	var b models.Banner
	err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *BannerRepository) ListActive(ctx context.Context) ([]models.Banner, error) {
	var items []models.Banner
	err := r.db.WithContext(ctx).Where("active = ?", true).
		Order("position ASC, sort_order ASC, created_at DESC").Find(&items).Error
	return items, err
}
func (r *BannerRepository) ListAll(ctx context.Context, position string) ([]models.Banner, error) {
	q := r.db.WithContext(ctx).Model(&models.Banner{})
	if position != "" {
		q = q.Where("position = ?", position)
	}
	var items []models.Banner
	err := q.Order("position ASC, sort_order ASC, created_at DESC").Find(&items).Error
	return items, err
}

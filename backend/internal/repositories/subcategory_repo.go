package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type SubcategoryRepository struct{ db *gorm.DB }

func NewSubcategoryRepository(db *gorm.DB) *SubcategoryRepository {
	return &SubcategoryRepository{db: db}
}

func (r *SubcategoryRepository) Create(ctx context.Context, s *models.Subcategory) error {
	return r.db.WithContext(ctx).Create(s).Error
}
func (r *SubcategoryRepository) Save(ctx context.Context, s *models.Subcategory) error {
	return r.db.WithContext(ctx).Save(s).Error
}
func (r *SubcategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Subcategory{}, "id = ?", id).Error
}

func (r *SubcategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Subcategory, error) {
	var s models.Subcategory
	err := r.db.WithContext(ctx).Preload("Category").First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SubcategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Subcategory, error) {
	var s models.Subcategory
	err := r.db.WithContext(ctx).Preload("Category").First(&s, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SubcategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Subcategory{}).Where("slug = ?", slug).Count(&n).Error
	return n > 0, err
}
func (r *SubcategoryRepository) ListByCategory(ctx context.Context, categoryID uuid.UUID, onlyActive bool) ([]models.Subcategory, error) {
	q := r.db.WithContext(ctx).Where("category_id = ?", categoryID)
	if onlyActive {
		q = q.Where("active = ?", true)
	}
	var items []models.Subcategory
	err := q.Order("sort_order ASC, name_tj ASC").Find(&items).Error
	return items, err
}
func (r *SubcategoryRepository) HasProducts(ctx context.Context, id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("subcategory_id = ?", id).Count(&n).Error
	return n > 0, err
}

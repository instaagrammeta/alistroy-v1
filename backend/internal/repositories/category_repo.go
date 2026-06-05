package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type CategoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}
func (r *CategoryRepository) Save(ctx context.Context, c *models.Category) error {
	return r.db.WithContext(ctx).Save(c).Error
}
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var c models.Category
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var c models.Category
	err := r.db.WithContext(ctx).First(&c, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Category{}).Where("slug = ?", slug).Count(&n).Error
	return n > 0, err
}

func (r *CategoryRepository) ListAll(ctx context.Context, onlyActive bool) ([]models.Category, error) {
	q := r.db.WithContext(ctx).Model(&models.Category{})
	if onlyActive {
		q = q.Where("active = ?", true)
	}
	var items []models.Category
	err := q.Preload("Subcategories", func(db *gorm.DB) *gorm.DB {
		return db.Where("active = ?", true).Order("sort_order ASC, name_tj ASC")
	}).Order("sort_order ASC, name_tj ASC").Find(&items).Error
	return items, err
}

func (r *CategoryRepository) HasChildren(ctx context.Context, id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Subcategory{}).Where("category_id = ?", id).Count(&n).Error
	return n > 0, err
}

func (r *CategoryRepository) HasProducts(ctx context.Context, id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("category_id = ?", id).Count(&n).Error
	return n > 0, err
}

// Popular returns categories sorted by approved-product count desc.
func (r *CategoryRepository) Popular(ctx context.Context, limit int) ([]models.Category, error) {
	if limit <= 0 || limit > 50 {
		limit = 8
	}
	type row struct {
		CategoryID uuid.UUID
		Cnt        int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT category_id, COUNT(*) AS cnt
		FROM products
		WHERE status = ? AND deleted_at IS NULL AND is_available = TRUE
		GROUP BY category_id
		ORDER BY cnt DESC
		LIMIT ?`, models.ProductStatusApproved, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.CategoryID)
	}
	var cats []models.Category
	if err := r.db.WithContext(ctx).Where("id IN ? AND active = TRUE", ids).Find(&cats).Error; err != nil {
		return nil, err
	}
	by := make(map[uuid.UUID]models.Category, len(cats))
	for _, c := range cats {
		by[c.ID] = c
	}
	out := make([]models.Category, 0, len(rows))
	for _, r := range rows {
		if c, ok := by[r.CategoryID]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

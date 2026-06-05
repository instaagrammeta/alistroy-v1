package repositories

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type ProductRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) *ProductRepository { return &ProductRepository{db: db} }

func (r *ProductRepository) DB() *gorm.DB { return r.db }

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProductRepository) UpdateColumns(ctx context.Context, id uuid.UUID, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).
		Preload("Seller").
		Preload("Seller.User").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).
		Preload("Seller").
		Preload("Seller.User").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		First(&p, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("slug = ?", slug).Count(&n).Error
	return n > 0, err
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

func (r *ProductRepository) IncrCounter(ctx context.Context, id uuid.UUID, column string) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn(column, gorm.Expr(column+" + 1")).Error
}

// ListParams drives the public/admin product listings.
type ListParams struct {
	Search        string
	CategoryID    *uuid.UUID
	CategorySlug  string
	SellerID      *uuid.UUID
	SellerSlug    string
	Status        string // empty = no filter; "approved" for public listing
	OnlyAvailable bool
	IsFeatured    *bool
	MinPrice      *float64
	MaxPrice      *float64
	Sort          string // newest|oldest|price_asc|price_desc|popular
	Page          int
	PageSize      int
}

func (r *ProductRepository) List(ctx context.Context, p ListParams) ([]models.Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Product{}).
		Preload("Seller").
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC").Limit(5) })

	if p.Status != "" {
		q = q.Where("products.status = ?", p.Status)
	}
	if p.OnlyAvailable {
		q = q.Where("products.is_available = ?", true)
	}
	if p.IsFeatured != nil {
		q = q.Where("products.is_featured = ?", *p.IsFeatured)
	}
	if p.CategoryID != nil {
		q = q.Where("products.category_id = ?", *p.CategoryID)
	}
	if p.SellerID != nil {
		q = q.Where("products.seller_id = ?", *p.SellerID)
	}
	if p.CategorySlug != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", p.CategorySlug)
	}
	if p.SellerSlug != "" {
		q = q.Joins("JOIN sellers ON sellers.id = products.seller_id").
			Where("sellers.slug = ?", p.SellerSlug)
	}
	if p.MinPrice != nil {
		q = q.Where("products.price >= ?", *p.MinPrice)
	}
	if p.MaxPrice != nil {
		q = q.Where("products.price <= ?", *p.MaxPrice)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where(`LOWER(products.title_tj) LIKE ?
			OR LOWER(products.title_ru) LIKE ?
			OR LOWER(products.description_tj) LIKE ?
			OR LOWER(products.description_ru) LIKE ?
			OR LOWER(products.sku) LIKE ?`, like, like, like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch p.Sort {
	case "oldest":
		q = q.Order("products.created_at ASC")
	case "price_asc":
		q = q.Order("products.price ASC")
	case "price_desc":
		q = q.Order("products.price DESC")
	case "popular":
		q = q.Order("products.views_count DESC")
	default: // "newest"
		q = q.Order("products.is_featured DESC, products.created_at DESC")
	}

	var items []models.Product
	err := q.Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).Find(&items).Error
	return items, total, err
}

// Related returns products in the same category (excluding the source product).
func (r *ProductRepository) Related(ctx context.Context, p *models.Product, limit int) ([]models.Product, error) {
	var items []models.Product
	err := r.db.WithContext(ctx).
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC").Limit(1) }).
		Preload("Seller").
		Where("category_id = ? AND id <> ? AND status = ? AND is_available = ?",
			p.CategoryID, p.ID, models.ProductStatusApproved, true).
		Order("created_at DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// CountByStatus returns counts grouped by status (used in admin dashboard).
func (r *ProductRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).Model(&models.Product{}).
		Select("status, COUNT(*) AS cnt").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, r := range rows {
		out[r.Status] = r.Cnt
	}
	return out, nil
}

// ReplaceImages replaces all images for a product within a transaction.
func (r *ProductRepository) ReplaceImages(ctx context.Context, productID uuid.UUID, images []models.ProductImage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error; err != nil {
			return err
		}
		if len(images) == 0 {
			return nil
		}
		for i := range images {
			images[i].ProductID = productID
		}
		return tx.Create(&images).Error
	})
}

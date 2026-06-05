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
func (r *ProductRepository) Save(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}
func (r *ProductRepository) UpdateColumns(ctx context.Context, id uuid.UUID, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ProductRepository) preloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Seller").
		Preload("Category").
		Preload("Subcategory").
		Preload("Brand").
		Preload("Images", func(d *gorm.DB) *gorm.DB { return d.Order("sort_order ASC") })
}

func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	err := r.preloads(r.db.WithContext(ctx)).First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var p models.Product
	err := r.preloads(r.db.WithContext(ctx)).First(&p, "slug = ?", slug).Error
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

func (r *ProductRepository) IncrCounter(ctx context.Context, id uuid.UUID, column string) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn(column, gorm.Expr(column+" + 1")).Error
}

func (r *ProductRepository) AdjustStock(ctx context.Context, id uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn("stock_quantity", gorm.Expr("GREATEST(stock_quantity + ?, 0)", delta)).Error
}

type ListProductsParams struct {
	Search        string
	CategoryID    *uuid.UUID
	CategorySlug  string
	SubcategoryID *uuid.UUID
	SubcatSlug    string
	BrandID       *uuid.UUID
	BrandIDs      []uuid.UUID
	SellerID      *uuid.UUID
	SellerSlug    string
	Status        string
	OnlyAvailable bool
	IsFeatured    *bool
	LowStockOnly  bool
	MinPrice      *float64
	MaxPrice      *float64
	Sort          string
	Page          int
	Size          int
}

func (r *ProductRepository) List(ctx context.Context, p ListProductsParams) ([]models.Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Product{})

	if p.Status != "" {
		q = q.Where("products.status = ?", p.Status)
	}
	if p.OnlyAvailable {
		q = q.Where("products.is_available = ?", true)
	}
	if p.IsFeatured != nil {
		q = q.Where("products.is_featured = ?", *p.IsFeatured)
	}
	if p.LowStockOnly {
		q = q.Where("products.stock_quantity <= products.minimum_stock")
	}
	if p.CategoryID != nil {
		q = q.Where("products.category_id = ?", *p.CategoryID)
	}
	if p.SubcategoryID != nil {
		q = q.Where("products.subcategory_id = ?", *p.SubcategoryID)
	}
	if p.BrandID != nil {
		q = q.Where("products.brand_id = ?", *p.BrandID)
	}
	if len(p.BrandIDs) > 0 {
		q = q.Where("products.brand_id IN ?", p.BrandIDs)
	}
	if p.SellerID != nil {
		q = q.Where("products.seller_id = ?", *p.SellerID)
	}
	if p.CategorySlug != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", p.CategorySlug)
	}
	if p.SubcatSlug != "" {
		q = q.Joins("JOIN subcategories ON subcategories.id = products.subcategory_id").
			Where("subcategories.slug = ?", p.SubcatSlug)
	}
	if p.SellerSlug != "" {
		q = q.Joins("JOIN sellers ON sellers.id = products.seller_id").
			Where("sellers.slug = ?", p.SellerSlug)
	}
	if p.MinPrice != nil {
		q = q.Where("products.sale_price >= ?", *p.MinPrice)
	}
	if p.MaxPrice != nil {
		q = q.Where("products.sale_price <= ?", *p.MaxPrice)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where(`LOWER(products.name_tj) LIKE ? OR LOWER(products.name_ru) LIKE ?
			OR LOWER(products.sku) LIKE ? OR LOWER(products.description_tj) LIKE ?
			OR LOWER(products.description_ru) LIKE ?`, like, like, like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch p.Sort {
	case "oldest":
		q = q.Order("products.created_at ASC")
	case "price_asc":
		q = q.Order("products.sale_price ASC")
	case "price_desc":
		q = q.Order("products.sale_price DESC")
	case "popular":
		q = q.Order("products.views_count DESC")
	default:
		q = q.Order("products.is_featured DESC, products.created_at DESC")
	}

	off := applyPaging(&p.Page, &p.Size)
	var items []models.Product
	err := r.preloads(q).Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

func (r *ProductRepository) Related(ctx context.Context, p *models.Product, limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 24 {
		limit = 8
	}
	var items []models.Product
	err := r.db.WithContext(ctx).
		Preload("Images", func(d *gorm.DB) *gorm.DB { return d.Order("sort_order ASC").Limit(1) }).
		Preload("Seller").
		Where("category_id = ? AND id <> ? AND status = ? AND is_available = TRUE",
			p.CategoryID, p.ID, models.ProductStatusApproved).
		Order("created_at DESC").Limit(limit).Find(&items).Error
	return items, err
}

func (r *ProductRepository) PriceBounds(ctx context.Context, categorySlug string) (min, max float64, err error) {
	type row struct{ Min, Max float64 }
	var x row
	q := r.db.WithContext(ctx).Model(&models.Product{}).
		Where("status = ? AND is_available = TRUE", models.ProductStatusApproved)
	if categorySlug != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", categorySlug)
	}
	err = q.Select("COALESCE(MIN(sale_price),0) AS min, COALESCE(MAX(sale_price),0) AS max").Scan(&x).Error
	return x.Min, x.Max, err
}

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

func (r *ProductRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).Model(&models.Product{}).
		Select("status, COUNT(*) AS cnt").Group("status").Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, x := range rows {
		out[x.Status] = x.Cnt
	}
	return out, nil
}

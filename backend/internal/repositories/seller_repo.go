package repositories

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type SellerRepository struct{ db *gorm.DB }

func NewSellerRepository(db *gorm.DB) *SellerRepository { return &SellerRepository{db: db} }

func (r *SellerRepository) Create(ctx context.Context, s *models.Seller) error {
	return r.db.WithContext(ctx).Create(s).Error
}
func (r *SellerRepository) Save(ctx context.Context, s *models.Seller) error {
	return r.db.WithContext(ctx).Save(s).Error
}
func (r *SellerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Seller{}, "id = ?", id).Error
}

func (r *SellerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	var s models.Seller
	err := r.db.WithContext(ctx).Preload("BusinessCat").First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SellerRepository) FindByUserID(ctx context.Context, uid uuid.UUID) (*models.Seller, error) {
	var s models.Seller
	err := r.db.WithContext(ctx).First(&s, "user_id = ?", uid).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SellerRepository) FindBySlug(ctx context.Context, slug string) (*models.Seller, error) {
	var s models.Seller
	err := r.db.WithContext(ctx).Preload("BusinessCat").First(&s, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SellerRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Seller{}).Where("slug = ?", slug).Count(&n).Error
	return n > 0, err
}

type ListSellersParams struct {
	Search   string
	Active   *bool
	Featured *bool
	Page     int
	Size     int
}

func (r *SellerRepository) List(ctx context.Context, p ListSellersParams) ([]models.Seller, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Seller{})
	if p.Active != nil {
		q = q.Where("active = ?", *p.Active)
	}
	if p.Featured != nil {
		q = q.Where("is_featured = ?", *p.Featured)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where(`LOWER(full_name) LIKE ? OR LOWER(company_name) LIKE ?
			OR LOWER(market_name) LIKE ? OR LOWER(city) LIKE ? OR phone LIKE ?`,
			like, like, like, like, "%"+s+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.Seller
	err := q.Preload("BusinessCat").
		Order("is_featured DESC, created_at DESC").
		Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

// Top returns sellers with the most approved products.
func (r *SellerRepository) Top(ctx context.Context, limit int) ([]models.Seller, error) {
	if limit <= 0 || limit > 50 {
		limit = 8
	}
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Table("products").
		Select("seller_id").
		Where("status = ? AND deleted_at IS NULL", models.ProductStatusApproved).
		Group("seller_id").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("seller_id", &ids).Error
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	var sellers []models.Seller
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND active = TRUE", ids).Find(&sellers).Error; err != nil {
		return nil, err
	}
	by := make(map[uuid.UUID]models.Seller, len(sellers))
	for _, s := range sellers {
		by[s.ID] = s
	}
	out := make([]models.Seller, 0, len(ids))
	for _, id := range ids {
		if s, ok := by[id]; ok {
			out = append(out, s)
		}
	}
	return out, nil
}

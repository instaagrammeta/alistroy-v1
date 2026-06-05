package repositories

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type DriverRepository struct{ db *gorm.DB }

func NewDriverRepository(db *gorm.DB) *DriverRepository { return &DriverRepository{db: db} }

func (r *DriverRepository) Create(ctx context.Context, d *models.Driver) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *DriverRepository) Save(ctx context.Context, d *models.Driver) error {
	return r.db.WithContext(ctx).Save(d).Error
}
func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Driver{}, "id = ?", id).Error
}

func (r *DriverRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	var d models.Driver
	err := r.db.WithContext(ctx).First(&d, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *DriverRepository) FindByUserID(ctx context.Context, uid uuid.UUID) (*models.Driver, error) {
	var d models.Driver
	err := r.db.WithContext(ctx).First(&d, "user_id = ?", uid).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

type ListDriversParams struct {
	Search string
	Active *bool
	Page   int
	Size   int
}

func (r *DriverRepository) List(ctx context.Context, p ListDriversParams) ([]models.Driver, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Driver{})
	if p.Active != nil {
		q = q.Where("active = ?", *p.Active)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("LOWER(full_name) LIKE ? OR phone LIKE ? OR LOWER(vehicle) LIKE ?",
			like, "%"+s+"%", like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.Driver
	err := q.Order("created_at DESC").Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

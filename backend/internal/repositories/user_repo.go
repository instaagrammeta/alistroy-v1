package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).Preload("Seller").First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	email = strings.ToLower(strings.TrimSpace(email))
	if err := r.db.WithContext(ctx).Preload("Seller").First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByResetToken(ctx context.Context, token string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, "reset_token = ?", token).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	email = strings.ToLower(strings.TrimSpace(email))
	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

type ListUsersParams struct {
	Role     string
	Search   string
	Page     int
	PageSize int
}

func (r *UserRepository) List(ctx context.Context, p ListUsersParams) ([]models.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.User{})
	if p.Role != "" {
		q = q.Where("role = ?", p.Role)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("LOWER(email) LIKE ? OR LOWER(name) LIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.User
	err := q.Order("created_at DESC").
		Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).
		Find(&items).Error
	return items, total, err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// IsNotFound returns true for the GORM "record not found" sentinel.
func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

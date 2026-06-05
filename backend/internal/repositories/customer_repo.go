package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type CustomerRepository struct{ db *gorm.DB }

func NewCustomerRepository(db *gorm.DB) *CustomerRepository { return &CustomerRepository{db: db} }

func (r *CustomerRepository) Create(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CustomerRepository) Save(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var c models.Customer
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) FindByUserID(ctx context.Context, uid uuid.UUID) (*models.Customer, error) {
	var c models.Customer
	err := r.db.WithContext(ctx).First(&c, "user_id = ?", uid).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type ListCustomersParams struct {
	Search string
	Page   int
	Size   int
}

func (r *CustomerRepository) List(ctx context.Context, p ListCustomersParams) ([]models.Customer, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Customer{})
	if p.Search != "" {
		like := "%" + p.Search + "%"
		q = q.Joins("JOIN users ON users.id = customers.user_id").
			Where("users.name ILIKE ? OR users.phone LIKE ? OR users.email ILIKE ? OR customers.company ILIKE ?",
				like, like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.Customer
	err := q.Order("customers.created_at DESC").
		Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, "id = ?", id).Error
}

// --- Addresses ---

func (r *CustomerRepository) AddAddress(ctx context.Context, a *models.Address) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *CustomerRepository) Addresses(ctx context.Context, customerID uuid.UUID) ([]models.Address, error) {
	var items []models.Address
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).
		Order("is_primary DESC, created_at DESC").Find(&items).Error
	return items, err
}

func (r *CustomerRepository) DeleteAddress(ctx context.Context, customerID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND customer_id = ?", addressID, customerID).
		Delete(&models.Address{}).Error
}

package repositories

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) DB() *gorm.DB { return r.db }

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Login = strings.ToLower(strings.TrimSpace(u.Login))
	u.Phone = strings.TrimSpace(u.Phone)
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) Save(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Preload("Customer").Preload("Seller").Preload("Driver").
		First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var u models.User
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Seller").Preload("Driver").
		First(&u, "email = ?", email).Error
	return ptrIfOK(&u, err)
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var u models.User
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Seller").Preload("Driver").
		First(&u, "phone = ?", phone).Error
	return ptrIfOK(&u, err)
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	if login == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var u models.User
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Seller").Preload("Driver").
		First(&u, "login = ?", login).Error
	return ptrIfOK(&u, err)
}

func (r *UserRepository) FindByGoogleID(ctx context.Context, gid string) (*models.User, error) {
	if strings.TrimSpace(gid) == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var u models.User
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Seller").Preload("Driver").
		First(&u, "google_id = ?", gid).Error
	return ptrIfOK(&u, err)
}

// FindByIdentifier accepts email/phone/login (whichever the user typed).
func (r *UserRepository) FindByIdentifier(ctx context.Context, ident string) (*models.User, error) {
	ident = strings.TrimSpace(ident)
	if ident == "" {
		return nil, gorm.ErrRecordNotFound
	}
	lower := strings.ToLower(ident)
	var u models.User
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Seller").Preload("Driver").
		Where("email = ? OR login = ? OR phone = ?", lower, lower, ident).
		First(&u).Error
	return ptrIfOK(&u, err)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return false, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&n).Error
	return n > 0, err
}

func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false, nil
	}
	var n int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("phone = ?", phone).Count(&n).Error
	return n > 0, err
}

type ListUsersParams struct {
	Role   string
	Status string
	Search string
	Page   int
	Size   int
}

func (r *UserRepository) List(ctx context.Context, p ListUsersParams) ([]models.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.User{})
	if p.Role != "" {
		q = q.Where("role = ?", p.Role)
	}
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	if s := strings.TrimSpace(p.Search); s != "" {
		like := "%" + strings.ToLower(s) + "%"
		q = q.Where("LOWER(email) LIKE ? OR LOWER(name) LIKE ? OR phone LIKE ? OR LOWER(login) LIKE ?",
			like, like, "%"+s+"%", like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&p.Page, &p.Size)
	var items []models.User
	err := q.Preload("Seller").Preload("Customer").Preload("Driver").
		Order("created_at DESC").Limit(p.Size).Offset(off).Find(&items).Error
	return items, total, err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func ptrIfOK(u *models.User, err error) (*models.User, error) {
	if err != nil {
		return nil, err
	}
	return u, nil
}

package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type DriverService struct {
	drivers *repositories.DriverRepository
	users   *repositories.UserRepository
}

func NewDriverService(d *repositories.DriverRepository, u *repositories.UserRepository) *DriverService {
	return &DriverService{drivers: d, users: u}
}

type DriverInput struct {
	FullName string
	Age      int
	Phone    string
	PhoneAlt string
	WhatsApp string
	Telegram string
	Vehicle  string
	PhotoURL string
	Notes    string
	Login    string
	Password string
	Active   *bool
	OnDuty   *bool
}

func (s *DriverService) CreateByAdmin(ctx context.Context, in DriverInput) (*models.Driver, error) {
	if strings.TrimSpace(in.FullName) == "" || strings.TrimSpace(in.Login) == "" || in.Password == "" {
		return nil, ErrValidation
	}
	if exists, _ := s.users.FindByLogin(ctx, in.Login); exists != nil {
		return nil, ErrConflict
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:         in.FullName,
		Login:        in.Login,
		Phone:        in.Phone,
		PasswordHash: string(hash),
		Role:         models.RoleDriver,
		Status:       models.UserStatusActive,
		Locale:       "tg",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, ErrConflict
	}
	d := &models.Driver{
		UserID:   user.ID,
		FullName: in.FullName,
		Age:      in.Age,
		Phone:    in.Phone,
		PhoneAlt: in.PhoneAlt,
		WhatsApp: in.WhatsApp,
		Telegram: in.Telegram,
		Vehicle:  in.Vehicle,
		PhotoURL: in.PhotoURL,
		Notes:    in.Notes,
		Active:   boolOr(in.Active, true),
		OnDuty:   boolOr(in.OnDuty, true),
	}
	if err := s.drivers.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *DriverService) UpdateByAdmin(ctx context.Context, id uuid.UUID, in DriverInput) (*models.Driver, error) {
	d, err := s.drivers.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.FullName != "" {
		d.FullName = in.FullName
	}
	if in.Age > 0 {
		d.Age = in.Age
	}
	d.Phone = firstNonEmpty(in.Phone, d.Phone)
	d.PhoneAlt = firstNonEmpty(in.PhoneAlt, d.PhoneAlt)
	d.WhatsApp = firstNonEmpty(in.WhatsApp, d.WhatsApp)
	d.Telegram = firstNonEmpty(in.Telegram, d.Telegram)
	d.Vehicle = firstNonEmpty(in.Vehicle, d.Vehicle)
	if in.PhotoURL != "" {
		d.PhotoURL = in.PhotoURL
	}
	d.Notes = firstNonEmpty(in.Notes, d.Notes)
	if in.Active != nil {
		d.Active = *in.Active
	}
	if in.OnDuty != nil {
		d.OnDuty = *in.OnDuty
	}
	if err := s.drivers.Save(ctx, d); err != nil {
		return nil, err
	}
	if in.Password != "" || in.Login != "" {
		if u, err := s.users.FindByID(ctx, d.UserID); err == nil {
			if in.Login != "" {
				u.Login = in.Login
			}
			if in.Password != "" {
				if h, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost); err == nil {
					u.PasswordHash = string(h)
				}
			}
			_ = s.users.Save(ctx, u)
		}
	}
	return d, nil
}

func (s *DriverService) Delete(ctx context.Context, id uuid.UUID) error {
	d, err := s.drivers.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	return s.users.Delete(ctx, d.UserID)
}

func (s *DriverService) GetByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	d, err := s.drivers.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return d, nil
}
func (s *DriverService) GetByUserID(ctx context.Context, uid uuid.UUID) (*models.Driver, error) {
	d, err := s.drivers.FindByUserID(ctx, uid)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return d, nil
}

func (s *DriverService) List(ctx context.Context, search string, active *bool, page, size int) ([]models.Driver, int64, error) {
	return s.drivers.List(ctx, repositories.ListDriversParams{Search: search, Active: active, Page: page, Size: size})
}

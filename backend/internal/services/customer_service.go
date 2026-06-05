package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type CustomerService struct {
	customers *repositories.CustomerRepository
	users     *repositories.UserRepository
}

func NewCustomerService(c *repositories.CustomerRepository, u *repositories.UserRepository) *CustomerService {
	return &CustomerService{customers: c, users: u}
}

type CustomerInput struct {
	Name     string
	Company  string
	Status   string
	Phone    string
	PhoneAlt string
	Address  string
	City     string
	Notes    string
	Password string
}

// CreateByAdmin creates a customer User + profile from the admin panel.
func (s *CustomerService) CreateByAdmin(ctx context.Context, in CustomerInput) (*models.Customer, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Phone) == "" {
		return nil, ErrValidation
	}
	if ex, _ := s.users.ExistsByPhone(ctx, in.Phone); ex {
		return nil, ErrConflict
	}
	hashStr := ""
	if in.Password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashStr = string(h)
	}
	status := models.UserStatusActive
	if in.Status != "" {
		status = in.Status
	}
	user := &models.User{
		Name: in.Name, Phone: in.Phone, PasswordHash: hashStr,
		Role: models.RoleCustomer, Status: status, Locale: "tg",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, ErrConflict
	}
	cust := &models.Customer{
		UserID: user.ID, Company: in.Company, Address: in.Address, City: in.City, Notes: in.Notes,
	}
	if err := s.customers.Create(ctx, cust); err != nil {
		return nil, err
	}
	cust.UserID = user.ID
	return cust, nil
}

func (s *CustomerService) UpdateByAdmin(ctx context.Context, id uuid.UUID, in CustomerInput) (*models.Customer, error) {
	cust, err := s.customers.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	cust.Company = firstNonEmpty(in.Company, cust.Company)
	cust.Address = firstNonEmpty(in.Address, cust.Address)
	cust.City = firstNonEmpty(in.City, cust.City)
	cust.Notes = firstNonEmpty(in.Notes, cust.Notes)
	if err := s.customers.Save(ctx, cust); err != nil {
		return nil, err
	}
	if u, err := s.users.FindByID(ctx, cust.UserID); err == nil {
		if in.Name != "" {
			u.Name = in.Name
		}
		if in.Phone != "" {
			u.Phone = in.Phone
		}
		if in.Status != "" {
			u.Status = in.Status
		}
		if in.Password != "" {
			if h, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost); err == nil {
				u.PasswordHash = string(h)
			}
		}
		_ = s.users.Save(ctx, u)
	}
	return cust, nil
}

func (s *CustomerService) Get(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	c, err := s.customers.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}
func (s *CustomerService) GetByUserID(ctx context.Context, uid uuid.UUID) (*models.Customer, error) {
	c, err := s.customers.FindByUserID(ctx, uid)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}
func (s *CustomerService) List(ctx context.Context, search string, page, size int) ([]models.Customer, int64, error) {
	return s.customers.List(ctx, repositories.ListCustomersParams{Search: search, Page: page, Size: size})
}
func (s *CustomerService) Delete(ctx context.Context, id uuid.UUID) error {
	c, err := s.customers.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	return s.users.Delete(ctx, c.UserID)
}

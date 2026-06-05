package services

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type UserService struct {
	users *repositories.UserRepository
}

func NewUserService(u *repositories.UserRepository) *UserService { return &UserService{users: u} }

func (s *UserService) List(ctx context.Context, role, status, search string, page, size int) ([]models.User, int64, error) {
	return s.users.List(ctx, repositories.ListUsersParams{Role: role, Status: status, Search: search, Page: page, Size: size})
}
func (s *UserService) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

type AdminUpdateUserInput struct {
	Name     string
	Phone    string
	Status   string
	Password string
}

func (s *UserService) AdminUpdate(ctx context.Context, id uuid.UUID, in AdminUpdateUserInput) (*models.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
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
		if len(in.Password) < 8 {
			return nil, ErrValidation
		}
		if h, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost); err == nil {
			u.PasswordHash = string(h)
		}
	}
	if err := s.users.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.users.Delete(ctx, id)
}

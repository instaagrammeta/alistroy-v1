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

type ListUsersInput struct {
	Role     string
	Search   string
	Page     int
	PageSize int
}

func (s *UserService) List(ctx context.Context, in ListUsersInput) ([]models.User, int64, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 20
	}
	return s.users.List(ctx, repositories.ListUsersParams{
		Role:     in.Role,
		Search:   in.Search,
		Page:     in.Page,
		PageSize: in.PageSize,
	})
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
	Role     string
	IsActive *bool
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
	if in.Role != "" {
		switch in.Role {
		case models.RoleAdmin, models.RoleSeller, models.RoleCustomer:
			u.Role = in.Role
		default:
			return nil, ErrValidation
		}
	}
	if in.IsActive != nil {
		u.IsActive = *in.IsActive
	}
	if in.Password != "" {
		if len(in.Password) < 8 {
			return nil, ErrValidation
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = string(hash)
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.users.Delete(ctx, id)
}

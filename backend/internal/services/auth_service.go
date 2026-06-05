package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/jwt"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type AuthService struct {
	users     *repositories.UserRepository
	customers *repositories.CustomerRepository
	jwt       *jwt.Manager
}

func NewAuthService(u *repositories.UserRepository, c *repositories.CustomerRepository, jm *jwt.Manager) *AuthService {
	return &AuthService{users: u, customers: c, jwt: jm}
}

type TokenPair struct {
	AccessToken      string       `json:"access_token"`
	RefreshToken     string       `json:"refresh_token"`
	AccessExpiresAt  time.Time    `json:"access_expires_at"`
	RefreshExpiresAt time.Time    `json:"refresh_expires_at"`
	User             *models.User `json:"user"`
}

// RegisterCustomerInput — phone + password registration (email optional).
type RegisterCustomerInput struct {
	Name     string
	Phone    string
	Email    string
	Password string
	Address  string
	City     string
	Locale   string
}

func (s *AuthService) RegisterCustomer(ctx context.Context, in RegisterCustomerInput) (*TokenPair, error) {
	in.Phone = strings.TrimSpace(in.Phone)
	if in.Name == "" || in.Phone == "" || in.Password == "" {
		return nil, ErrValidation
	}
	if exists, _ := s.users.ExistsByPhone(ctx, in.Phone); exists {
		return nil, ErrConflict
	}
	if in.Email != "" {
		if exists, _ := s.users.ExistsByEmail(ctx, in.Email); exists {
			return nil, ErrConflict
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:         in.Name,
		Phone:        in.Phone,
		Email:        strings.ToLower(strings.TrimSpace(in.Email)),
		PasswordHash: string(hash),
		Role:         models.RoleCustomer,
		Status:       models.UserStatusActive,
		Locale:       defaultStr(in.Locale, "tg"),
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	cust := &models.Customer{UserID: user.ID, Address: in.Address, City: in.City}
	if err := s.customers.Create(ctx, cust); err != nil {
		return nil, err
	}
	user.Customer = cust
	return s.issue(user)
}

// LoginInput accepts phone / email / login + password.
func (s *AuthService) Login(ctx context.Context, identifier, password string) (*TokenPair, error) {
	user, err := s.users.FindByIdentifier(ctx, identifier)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if user.Status != models.UserStatusActive {
		return nil, ErrForbidden
	}
	if user.PasswordHash == "" ||
		bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	now := time.Now().UTC()
	user.LastLoginAt = &now
	_ = s.users.Save(ctx, user)
	return s.issue(user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwt.Parse(refreshToken)
	if err != nil || claims.Type != jwt.RefreshToken {
		return nil, ErrUnauthorized
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if user.Status != models.UserStatusActive {
		return nil, ErrForbidden
	}
	return s.issue(user)
}

// GoogleUpsert finds or creates a customer from a verified Google profile.
// New Google users still must complete phone+address (enforced by handler).
func (s *AuthService) GoogleUpsert(ctx context.Context, googleID, email, name, avatar string) (*TokenPair, bool, error) {
	if googleID == "" {
		return nil, false, ErrValidation
	}
	if u, err := s.users.FindByGoogleID(ctx, googleID); err == nil {
		pair, err := s.issue(u)
		return pair, false, err
	}
	if email != "" {
		if u, err := s.users.FindByEmail(ctx, email); err == nil {
			u.GoogleID = googleID
			if u.AvatarURL == "" {
				u.AvatarURL = avatar
			}
			_ = s.users.Save(ctx, u)
			pair, err := s.issue(u)
			return pair, false, err
		}
	}
	user := &models.User{
		Name:      defaultStr(name, "Customer"),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		GoogleID:  googleID,
		AvatarURL: avatar,
		Role:      models.RoleCustomer,
		Status:    models.UserStatusActive,
		Locale:    "tg",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, false, err
	}
	cust := &models.Customer{UserID: user.ID}
	if err := s.customers.Create(ctx, cust); err != nil {
		return nil, false, err
	}
	user.Customer = cust
	pair, err := s.issue(user)
	return pair, true, err // needsProfile = true
}

func (s *AuthService) Me(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

type UpdateProfileInput struct {
	Name    string
	Phone   string
	Locale  string
	Address string
	City    string
	Company string
}

func (s *AuthService) UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateProfileInput) (*models.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.Name != "" {
		u.Name = in.Name
	}
	if in.Phone != "" && in.Phone != u.Phone {
		if exists, _ := s.users.ExistsByPhone(ctx, in.Phone); exists {
			return nil, ErrConflict
		}
		u.Phone = in.Phone
	}
	if in.Locale != "" {
		u.Locale = in.Locale
	}
	if err := s.users.Save(ctx, u); err != nil {
		return nil, err
	}
	// Update customer profile if present.
	if u.Role == models.RoleCustomer {
		if cust, err := s.customers.FindByUserID(ctx, u.ID); err == nil {
			if in.Address != "" {
				cust.Address = in.Address
			}
			if in.City != "" {
				cust.City = in.City
			}
			if in.Company != "" {
				cust.Company = in.Company
			}
			_ = s.customers.Save(ctx, cust)
			u.Customer = cust
		}
	}
	return u, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, id uuid.UUID, oldPw, newPw string) error {
	if len(newPw) < 8 {
		return ErrValidation
	}
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if u.PasswordHash != "" {
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPw)) != nil {
			return ErrInvalidCredentials
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return s.users.Save(ctx, u)
}

func (s *AuthService) issue(u *models.User) (*TokenPair, error) {
	at, atExp, err := s.jwt.Generate(u.ID, u.Role, jwt.AccessToken)
	if err != nil {
		return nil, err
	}
	rt, rtExp, err := s.jwt.Generate(u.ID, u.Role, jwt.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:      at,
		RefreshToken:     rt,
		AccessExpiresAt:  atExp,
		RefreshExpiresAt: rtExp,
		User:             u,
	}, nil
}

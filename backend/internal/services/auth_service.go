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
	users   *repositories.UserRepository
	sellers *repositories.SellerRepository
	jwt     *jwt.Manager
}

func NewAuthService(users *repositories.UserRepository, sellers *repositories.SellerRepository, jm *jwt.Manager) *AuthService {
	return &AuthService{users: users, sellers: sellers, jwt: jm}
}

type RegisterInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	Role     string // "customer" (default) or "seller"
	Locale   string
	// Used only when Role == "seller":
	SellerName string
	City       string
}

type TokenPair struct {
	AccessToken      string       `json:"access_token"`
	RefreshToken     string       `json:"refresh_token"`
	AccessExpiresAt  time.Time    `json:"access_expires_at"`
	RefreshExpiresAt time.Time    `json:"refresh_expires_at"`
	User             *models.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*TokenPair, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.Email == "" || in.Password == "" || in.Name == "" {
		return nil, ErrValidation
	}
	if in.Role == "" {
		in.Role = models.RoleCustomer
	}
	if in.Role != models.RoleCustomer && in.Role != models.RoleSeller {
		return nil, ErrValidation
	}
	if in.Locale == "" {
		in.Locale = "tg"
	}

	exists, err := s.users.ExistsByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		Name:         in.Name,
		Phone:        in.Phone,
		Role:         in.Role,
		Locale:       in.Locale,
		IsActive:     true,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	if in.Role == models.RoleSeller {
		seller, err := s.createSellerFor(ctx, user, in.SellerName, in.City)
		if err != nil {
			return nil, err
		}
		user.Seller = seller
	}

	return s.issueTokens(user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, ErrForbidden
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	now := time.Now().UTC()
	user.LastLoginAt = &now
	_ = s.users.Update(ctx, user)
	return s.issueTokens(user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwt.Parse(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if claims.Type != jwt.RefreshToken {
		return nil, ErrUnauthorized
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if !user.IsActive {
		return nil, ErrForbidden
	}
	return s.issueTokens(user)
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
	Name   string
	Phone  string
	Locale string
}

func (s *AuthService) UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateProfileInput) (*models.User, error) {
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
	if in.Locale != "" {
		u.Locale = in.Locale
	}
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	if newPassword == "" || len(newPassword) < 8 {
		return ErrValidation
	}
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return s.users.Update(ctx, u)
}

// RequestPasswordReset issues a one-time reset token. The token is returned
// here so it can be surfaced to the integrator (email/SMS gateway) — in this
// build we expose it via the API response (admin can also use it manually).
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		// Do not leak existence; return empty token.
		return "", nil
	}
	tok, err := RandomToken(24)
	if err != nil {
		return "", err
	}
	exp := time.Now().UTC().Add(2 * time.Hour)
	u.ResetToken = tok
	u.ResetExpiresAt = &exp
	if err := s.users.Update(ctx, u); err != nil {
		return "", err
	}
	return tok, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrValidation
	}
	u, err := s.users.FindByResetToken(ctx, token)
	if err != nil || u.ResetExpiresAt == nil || time.Now().UTC().After(*u.ResetExpiresAt) {
		return ErrUnauthorized
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	u.ResetToken = ""
	u.ResetExpiresAt = nil
	return s.users.Update(ctx, u)
}

// ----- helpers -----

func (s *AuthService) issueTokens(u *models.User) (*TokenPair, error) {
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

func (s *AuthService) createSellerFor(ctx context.Context, user *models.User, sellerName, city string) (*models.Seller, error) {
	name := strings.TrimSpace(sellerName)
	if name == "" {
		name = user.Name
	}
	slug := Slugify(name)
	finalSlug := slug
	for i := 1; i < 1000; i++ {
		exists, err := s.sellers.ExistsBySlug(ctx, finalSlug)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		finalSlug = slug + "-" + intToStr(i)
	}
	seller := &models.Seller{
		UserID:   user.ID,
		Name:     name,
		Slug:     finalSlug,
		Phone:    user.Phone,
		WhatsApp: user.Phone,
		City:     city,
		Status:   models.SellerStatusPending,
	}
	if err := s.sellers.Create(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}

func intToStr(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	out := make([]byte, 0, 10)
	for n > 0 {
		out = append([]byte{digits[n%10]}, out...)
		n /= 10
	}
	return string(out)
}

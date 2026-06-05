package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type SellerService struct {
	sellers *repositories.SellerRepository
	users   *repositories.UserRepository
}

func NewSellerService(s *repositories.SellerRepository, u *repositories.UserRepository) *SellerService {
	return &SellerService{sellers: s, users: u}
}

// SellerInput holds everything admin enters to create a seller (incl. login).
type SellerInput struct {
	FullName         string
	CompanyName      string
	MarketName       string
	Phone            string
	PhoneAlt         string
	WhatsApp         string
	Telegram         string
	TelegramUsername string
	Address          string
	City             string
	BusinessCategory *uuid.UUID
	Notes            string
	LogoURL          string
	Login            string
	Password         string
	Active           *bool
	IsFeatured       *bool
}

// CreateByAdmin creates a seller User (role=seller) + Seller profile.
func (s *SellerService) CreateByAdmin(ctx context.Context, in SellerInput) (*models.Seller, error) {
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
		Role:         models.RoleSeller,
		Status:       models.UserStatusActive,
		Locale:       "tg",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, ErrConflict
	}
	slug, err := uniqueSlug(ctx, firstNonEmpty(in.MarketName, in.CompanyName, in.FullName), s.sellers.ExistsBySlug)
	if err != nil {
		return nil, err
	}
	seller := &models.Seller{
		UserID:           user.ID,
		FullName:         in.FullName,
		CompanyName:      in.CompanyName,
		MarketName:       in.MarketName,
		Slug:             slug,
		Phone:            in.Phone,
		PhoneAlt:         in.PhoneAlt,
		WhatsApp:         in.WhatsApp,
		Telegram:         in.Telegram,
		TelegramUsername: in.TelegramUsername,
		Address:          in.Address,
		City:             in.City,
		BusinessCategory: in.BusinessCategory,
		Notes:            in.Notes,
		LogoURL:          in.LogoURL,
		Active:           boolOr(in.Active, true),
		IsFeatured:       boolOr(in.IsFeatured, false),
	}
	if err := s.sellers.Create(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}

func (s *SellerService) UpdateByAdmin(ctx context.Context, id uuid.UUID, in SellerInput) (*models.Seller, error) {
	seller, err := s.sellers.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	applySellerFields(seller, in)
	if in.Active != nil {
		seller.Active = *in.Active
	}
	if in.IsFeatured != nil {
		seller.IsFeatured = *in.IsFeatured
	}
	if err := s.sellers.Save(ctx, seller); err != nil {
		return nil, err
	}
	// Optional password reset / login update on the user.
	if in.Password != "" || in.Login != "" {
		if u, err := s.users.FindByID(ctx, seller.UserID); err == nil {
			if in.Login != "" {
				u.Login = in.Login
			}
			if in.Password != "" {
				if h, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost); err == nil {
					u.PasswordHash = string(h)
				}
			}
			if in.FullName != "" {
				u.Name = in.FullName
			}
			_ = s.users.Save(ctx, u)
		}
	}
	return seller, nil
}

// UpdateOwn lets the seller edit their own profile (no login/active changes).
func (s *SellerService) UpdateOwn(ctx context.Context, userID uuid.UUID, in SellerInput) (*models.Seller, error) {
	seller, err := s.sellers.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	applySellerFields(seller, in)
	if err := s.sellers.Save(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}

func applySellerFields(seller *models.Seller, in SellerInput) {
	if in.FullName != "" {
		seller.FullName = in.FullName
	}
	seller.CompanyName = firstNonEmpty(in.CompanyName, seller.CompanyName)
	seller.MarketName = firstNonEmpty(in.MarketName, seller.MarketName)
	if in.Phone != "" {
		seller.Phone = in.Phone
	}
	seller.PhoneAlt = firstNonEmpty(in.PhoneAlt, seller.PhoneAlt)
	seller.WhatsApp = firstNonEmpty(in.WhatsApp, seller.WhatsApp)
	seller.Telegram = firstNonEmpty(in.Telegram, seller.Telegram)
	seller.TelegramUsername = firstNonEmpty(in.TelegramUsername, seller.TelegramUsername)
	seller.Address = firstNonEmpty(in.Address, seller.Address)
	seller.City = firstNonEmpty(in.City, seller.City)
	if in.BusinessCategory != nil {
		seller.BusinessCategory = in.BusinessCategory
	}
	seller.Notes = firstNonEmpty(in.Notes, seller.Notes)
	if in.LogoURL != "" {
		seller.LogoURL = in.LogoURL
	}
}

func (s *SellerService) Delete(ctx context.Context, id uuid.UUID) error {
	seller, err := s.sellers.FindByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	// Removing the user cascades to the seller + products.
	return s.users.Delete(ctx, seller.UserID)
}

func (s *SellerService) GetByID(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	return wrapSeller(s.sellers.FindByID(ctx, id))
}
func (s *SellerService) GetBySlug(ctx context.Context, slug string) (*models.Seller, error) {
	return wrapSeller(s.sellers.FindBySlug(ctx, slug))
}
func (s *SellerService) GetByUserID(ctx context.Context, uid uuid.UUID) (*models.Seller, error) {
	return wrapSeller(s.sellers.FindByUserID(ctx, uid))
}

type ListSellersInput struct {
	Search   string
	Active   *bool
	Featured *bool
	Page     int
	Size     int
}

func (s *SellerService) List(ctx context.Context, in ListSellersInput) ([]models.Seller, int64, error) {
	return s.sellers.List(ctx, repositories.ListSellersParams{
		Search: in.Search, Active: in.Active, Featured: in.Featured, Page: in.Page, Size: in.Size,
	})
}
func (s *SellerService) Top(ctx context.Context, limit int) ([]models.Seller, error) {
	return s.sellers.Top(ctx, limit)
}

func wrapSeller(s *models.Seller, err error) (*models.Seller, error) {
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

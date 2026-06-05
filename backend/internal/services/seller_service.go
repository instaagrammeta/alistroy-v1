package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type SellerService struct {
	repo *repositories.SellerRepository
}

func NewSellerService(repo *repositories.SellerRepository) *SellerService {
	return &SellerService{repo: repo}
}

type SellerProfileInput struct {
	Name          string
	DescriptionTJ string
	DescriptionRU string
	LogoURL       string
	Phone         string
	WhatsApp      string
	Address       string
	City          string
}

func (s *SellerService) GetByID(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	out, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return out, nil
}

func (s *SellerService) GetBySlug(ctx context.Context, slug string) (*models.Seller, error) {
	out, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return out, nil
}

func (s *SellerService) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Seller, error) {
	out, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return out, nil
}

func (s *SellerService) UpdateOwn(ctx context.Context, userID uuid.UUID, in SellerProfileInput) (*models.Seller, error) {
	seller, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	applyProfile(seller, in)
	if err := s.repo.Update(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}

func (s *SellerService) AdminUpdate(ctx context.Context, id uuid.UUID, in SellerProfileInput, status string, isFeatured *bool) (*models.Seller, error) {
	seller, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	applyProfile(seller, in)
	if status != "" {
		seller.Status = status
	}
	if isFeatured != nil {
		seller.IsFeatured = *isFeatured
	}
	if err := s.repo.Update(ctx, seller); err != nil {
		return nil, err
	}
	return seller, nil
}

func applyProfile(seller *models.Seller, in SellerProfileInput) {
	if in.Name != "" {
		seller.Name = in.Name
	}
	if in.DescriptionTJ != "" {
		seller.DescriptionTJ = in.DescriptionTJ
	}
	if in.DescriptionRU != "" {
		seller.DescriptionRU = in.DescriptionRU
	}
	if in.LogoURL != "" {
		seller.LogoURL = in.LogoURL
	}
	if in.Phone != "" {
		seller.Phone = strings.TrimSpace(in.Phone)
	}
	if in.WhatsApp != "" {
		seller.WhatsApp = strings.TrimSpace(in.WhatsApp)
	}
	if in.Address != "" {
		seller.Address = in.Address
	}
	if in.City != "" {
		seller.City = in.City
	}
}

type ListSellersInput struct {
	Search   string
	Status   string
	Featured *bool
	Page     int
	PageSize int
}

func (s *SellerService) List(ctx context.Context, in ListSellersInput) ([]models.Seller, int64, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 20
	}
	return s.repo.List(ctx, repositories.ListSellersParams{
		Search:   in.Search,
		Status:   in.Status,
		Featured: in.Featured,
		Page:     in.Page,
		PageSize: in.PageSize,
	})
}

func (s *SellerService) Top(ctx context.Context, limit int) ([]models.Seller, error) {
	if limit <= 0 || limit > 50 {
		limit = 8
	}
	return s.repo.Top(ctx, limit)
}

func (s *SellerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

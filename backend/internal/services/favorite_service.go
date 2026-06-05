package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type FavoriteService struct {
	repo     *repositories.FavoriteRepository
	products *repositories.ProductRepository
}

func NewFavoriteService(r *repositories.FavoriteRepository, p *repositories.ProductRepository) *FavoriteService {
	return &FavoriteService{repo: r, products: p}
}

func (s *FavoriteService) Add(ctx context.Context, userID, productID uuid.UUID) error {
	if _, err := s.products.FindByID(ctx, productID); err != nil {
		return ErrNotFound
	}
	exists, err := s.repo.Exists(ctx, userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.repo.Add(ctx, userID, productID)
}

func (s *FavoriteService) Remove(ctx context.Context, userID, productID uuid.UUID) error {
	return s.repo.Remove(ctx, userID, productID)
}

func (s *FavoriteService) Has(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, userID, productID)
}

func (s *FavoriteService) List(ctx context.Context, userID uuid.UUID, page, size int) ([]models.Favorite, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.repo.ListByUser(ctx, userID, page, size)
}

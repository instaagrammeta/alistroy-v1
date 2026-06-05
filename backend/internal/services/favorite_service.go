package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type FavoriteService struct {
	favs     *repositories.FavoriteRepository
	products *repositories.ProductRepository
}

func NewFavoriteService(f *repositories.FavoriteRepository, p *repositories.ProductRepository) *FavoriteService {
	return &FavoriteService{favs: f, products: p}
}

func (s *FavoriteService) Add(ctx context.Context, userID, productID uuid.UUID) error {
	if _, err := s.products.FindByID(ctx, productID); err != nil {
		return ErrNotFound
	}
	if ex, _ := s.favs.Exists(ctx, userID, productID); ex {
		return nil
	}
	return s.favs.Add(ctx, userID, productID)
}
func (s *FavoriteService) Remove(ctx context.Context, userID, productID uuid.UUID) error {
	return s.favs.Remove(ctx, userID, productID)
}
func (s *FavoriteService) Has(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return s.favs.Exists(ctx, userID, productID)
}
func (s *FavoriteService) List(ctx context.Context, userID uuid.UUID, page, size int) ([]models.Favorite, int64, error) {
	return s.favs.ListByUser(ctx, userID, page, size)
}

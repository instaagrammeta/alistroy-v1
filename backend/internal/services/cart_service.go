package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type CartService struct {
	cart     *repositories.CartRepository
	products *repositories.ProductRepository
}

func NewCartService(c *repositories.CartRepository, p *repositories.ProductRepository) *CartService {
	return &CartService{cart: c, products: p}
}

func (s *CartService) Set(ctx context.Context, customerID, productID uuid.UUID, qty int) error {
	if _, err := s.products.FindByID(ctx, productID); err != nil {
		return ErrNotFound
	}
	return s.cart.Upsert(ctx, customerID, productID, qty)
}
func (s *CartService) Remove(ctx context.Context, customerID, productID uuid.UUID) error {
	return s.cart.Remove(ctx, customerID, productID)
}
func (s *CartService) Clear(ctx context.Context, customerID uuid.UUID) error {
	return s.cart.Clear(ctx, customerID)
}
func (s *CartService) List(ctx context.Context, customerID uuid.UUID) ([]models.CartItem, error) {
	return s.cart.List(ctx, customerID)
}
func (s *CartService) Count(ctx context.Context, customerID uuid.UUID) (int64, error) {
	return s.cart.Count(ctx, customerID)
}

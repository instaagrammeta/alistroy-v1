package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type ReviewService struct {
	reviews  *repositories.ReviewRepository
	products *repositories.ProductRepository
}

func NewReviewService(r *repositories.ReviewRepository, p *repositories.ProductRepository) *ReviewService {
	return &ReviewService{reviews: r, products: p}
}

type ReviewInput struct {
	ProductID uuid.UUID
	Rating    int
	Comment   string
}

func (s *ReviewService) Create(ctx context.Context, userID uuid.UUID, in ReviewInput) (*models.Review, error) {
	if in.Rating < 1 || in.Rating > 5 || in.ProductID == uuid.Nil {
		return nil, ErrValidation
	}
	if _, err := s.products.FindByID(ctx, in.ProductID); err != nil {
		return nil, ErrNotFound
	}
	r := &models.Review{
		ProductID: in.ProductID,
		UserID:    userID,
		Rating:    in.Rating,
		Comment:   in.Comment,
		Status:    "pending",
	}
	if err := s.reviews.Create(ctx, r); err != nil {
		return nil, err
	}
	return s.reviews.FindByID(ctx, r.ID)
}

func (s *ReviewService) ListForProduct(ctx context.Context, productID uuid.UUID, page, size int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.reviews.List(ctx, repositories.ListReviewsParams{
		ProductID: &productID,
		Status:    "approved",
		Page:      page,
		PageSize:  size,
	})
}

func (s *ReviewService) AdminList(ctx context.Context, status string, page, size int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.reviews.List(ctx, repositories.ListReviewsParams{
		Status:   status,
		Page:     page,
		PageSize: size,
	})
}

func (s *ReviewService) Moderate(ctx context.Context, id uuid.UUID, status string) (*models.Review, error) {
	if status != "approved" && status != "rejected" && status != "pending" {
		return nil, ErrValidation
	}
	r, err := s.reviews.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	r.Status = status
	if err := s.reviews.Update(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *ReviewService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.reviews.Delete(ctx, id)
}

func (s *ReviewService) Average(ctx context.Context, productID uuid.UUID) (float64, int64, error) {
	return s.reviews.AverageRating(ctx, productID)
}

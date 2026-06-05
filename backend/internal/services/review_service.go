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
	rv := &models.Review{ProductID: in.ProductID, UserID: userID, Rating: in.Rating, Comment: in.Comment, Status: models.ReviewStatusPending}
	if err := s.reviews.Create(ctx, rv); err != nil {
		return nil, err
	}
	return s.reviews.FindByID(ctx, rv.ID)
}

func (s *ReviewService) ListForProduct(ctx context.Context, productID uuid.UUID, page, size int) ([]models.Review, int64, error) {
	return s.reviews.List(ctx, repositories.ListReviewsParams{ProductID: &productID, Status: models.ReviewStatusApproved, Page: page, Size: size})
}
func (s *ReviewService) AdminList(ctx context.Context, status string, page, size int) ([]models.Review, int64, error) {
	return s.reviews.List(ctx, repositories.ListReviewsParams{Status: status, Page: page, Size: size})
}
func (s *ReviewService) Moderate(ctx context.Context, id uuid.UUID, status string) (*models.Review, error) {
	if status != models.ReviewStatusApproved && status != models.ReviewStatusRejected && status != models.ReviewStatusPending {
		return nil, ErrValidation
	}
	rv, err := s.reviews.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	rv.Status = status
	if err := s.reviews.Save(ctx, rv); err != nil {
		return nil, err
	}
	return rv, nil
}
func (s *ReviewService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.reviews.Delete(ctx, id)
}
func (s *ReviewService) Average(ctx context.Context, productID uuid.UUID) (float64, int64, error) {
	return s.reviews.Average(ctx, productID)
}

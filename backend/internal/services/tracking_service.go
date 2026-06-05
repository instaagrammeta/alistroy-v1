package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type TrackingService struct {
	products *repositories.ProductRepository
	repo     *repositories.TrackingRepository
}

func NewTrackingService(p *repositories.ProductRepository, t *repositories.TrackingRepository) *TrackingService {
	return &TrackingService{products: p, repo: t}
}

type TrackInput struct {
	ProductID uuid.UUID
	Event     string
	UserID    *uuid.UUID
	IP        string
	UserAgent string
}

// Record persists a tracking event AND atomically increments the per-product counter.
func (s *TrackingService) Record(ctx context.Context, in TrackInput) error {
	p, err := s.products.FindByID(ctx, in.ProductID)
	if err != nil {
		return ErrNotFound
	}
	column := ""
	switch in.Event {
	case models.EventView:
		column = "views_count"
	case models.EventPhoneClick:
		column = "phone_clicks"
	case models.EventWhatsAppClick:
		column = "whats_app_clicks"
	default:
		return ErrValidation
	}
	if err := s.products.IncrCounter(ctx, p.ID, column); err != nil {
		return err
	}
	ev := &models.TrackingEvent{
		ProductID: p.ID,
		SellerID:  p.SellerID,
		UserID:    in.UserID,
		Event:     in.Event,
		IP:        in.IP,
		UserAgent: in.UserAgent,
	}
	return s.repo.Insert(ctx, ev)
}

func (s *TrackingService) GlobalTotals(ctx context.Context) (*repositories.Totals, error) {
	return s.repo.GlobalTotals(ctx)
}

func (s *TrackingService) SellerTotals(ctx context.Context, sellerID uuid.UUID) (*repositories.SellerTotals, error) {
	return s.repo.SellerTotals(ctx, sellerID)
}

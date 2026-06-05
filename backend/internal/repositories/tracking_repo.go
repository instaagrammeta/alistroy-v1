package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type TrackingRepository struct{ db *gorm.DB }

func NewTrackingRepository(db *gorm.DB) *TrackingRepository { return &TrackingRepository{db: db} }

func (r *TrackingRepository) Insert(ctx context.Context, e *models.TrackingEvent) error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now().UTC()
	}
	return r.db.WithContext(ctx).Create(e).Error
}

// Totals returns global totals for the admin dashboard.
type Totals struct {
	TotalProducts         int64 `json:"total_products"`
	TotalApprovedProducts int64 `json:"total_approved_products"`
	TotalPendingProducts  int64 `json:"total_pending_products"`
	TotalSellers          int64 `json:"total_sellers"`
	TotalUsers            int64 `json:"total_users"`
	TotalViews            int64 `json:"total_views"`
	TotalPhoneClicks      int64 `json:"total_phone_clicks"`
	TotalWhatsAppClicks   int64 `json:"total_whatsapp_clicks"`
	TotalReviews          int64 `json:"total_reviews"`
}

func (r *TrackingRepository) GlobalTotals(ctx context.Context) (*Totals, error) {
	t := &Totals{}
	db := r.db.WithContext(ctx)

	if err := db.Model(&models.Product{}).Count(&t.TotalProducts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Product{}).Where("status = ?", models.ProductStatusApproved).Count(&t.TotalApprovedProducts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Product{}).Where("status = ?", models.ProductStatusPending).Count(&t.TotalPendingProducts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Seller{}).Count(&t.TotalSellers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.User{}).Count(&t.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Review{}).Count(&t.TotalReviews).Error; err != nil {
		return nil, err
	}

	type sums struct {
		Views   int64
		PhoneCl int64
		WaCl    int64
	}
	var s sums
	if err := db.Model(&models.Product{}).
		Select(`COALESCE(SUM(views_count),0) AS views,
			COALESCE(SUM(phone_clicks),0) AS phone_cl,
			COALESCE(SUM(whats_app_clicks),0) AS wa_cl`).
		Scan(&s).Error; err != nil {
		return nil, err
	}
	t.TotalViews = s.Views
	t.TotalPhoneClicks = s.PhoneCl
	t.TotalWhatsAppClicks = s.WaCl
	return t, nil
}

// SellerTotals returns aggregate stats for a single seller.
type SellerTotals struct {
	TotalProducts       int64 `json:"total_products"`
	ApprovedProducts    int64 `json:"approved_products"`
	PendingProducts     int64 `json:"pending_products"`
	TotalViews          int64 `json:"total_views"`
	TotalPhoneClicks    int64 `json:"total_phone_clicks"`
	TotalWhatsAppClicks int64 `json:"total_whatsapp_clicks"`
}

func (r *TrackingRepository) SellerTotals(ctx context.Context, sellerID uuid.UUID) (*SellerTotals, error) {
	t := &SellerTotals{}
	db := r.db.WithContext(ctx)

	if err := db.Model(&models.Product{}).Where("seller_id = ?", sellerID).Count(&t.TotalProducts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Product{}).Where("seller_id = ? AND status = ?", sellerID, models.ProductStatusApproved).Count(&t.ApprovedProducts).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&models.Product{}).Where("seller_id = ? AND status = ?", sellerID, models.ProductStatusPending).Count(&t.PendingProducts).Error; err != nil {
		return nil, err
	}
	type sums struct {
		Views   int64
		PhoneCl int64
		WaCl    int64
	}
	var s sums
	if err := db.Model(&models.Product{}).
		Where("seller_id = ?", sellerID).
		Select(`COALESCE(SUM(views_count),0) AS views,
			COALESCE(SUM(phone_clicks),0) AS phone_cl,
			COALESCE(SUM(whats_app_clicks),0) AS wa_cl`).
		Scan(&s).Error; err != nil {
		return nil, err
	}
	t.TotalViews = s.Views
	t.TotalPhoneClicks = s.PhoneCl
	t.TotalWhatsAppClicks = s.WaCl
	return t, nil
}

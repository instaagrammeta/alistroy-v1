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

// Totals are the global dashboard counters.
type Totals struct {
	TotalProducts       int64   `json:"total_products"`
	ApprovedProducts    int64   `json:"approved_products"`
	PendingProducts     int64   `json:"pending_products"`
	TotalCategories     int64   `json:"total_categories"`
	TotalBrands         int64   `json:"total_brands"`
	TotalSellers        int64   `json:"total_sellers"`
	TotalCustomers      int64   `json:"total_customers"`
	TotalDrivers        int64   `json:"total_drivers"`
	TotalOrders         int64   `json:"total_orders"`
	NewOrders           int64   `json:"new_orders"`
	TotalViews          int64   `json:"total_views"`
	TotalPhoneClicks    int64   `json:"total_phone_clicks"`
	TotalWhatsAppClicks int64   `json:"total_whatsapp_clicks"`
	TotalTelegramClicks int64   `json:"total_telegram_clicks"`
	TotalRevenue        float64 `json:"total_revenue"`
	TotalProfit         float64 `json:"total_profit"`
}

func (r *TrackingRepository) GlobalTotals(ctx context.Context) (*Totals, error) {
	t := &Totals{}
	db := r.db.WithContext(ctx)
	count := func(model any, where string, args ...any) (int64, error) {
		var n int64
		q := db.Model(model)
		if where != "" {
			q = q.Where(where, args...)
		}
		return n, q.Count(&n).Error
	}
	var err error
	if t.TotalProducts, err = count(&models.Product{}, ""); err != nil {
		return nil, err
	}
	if t.ApprovedProducts, err = count(&models.Product{}, "status = ?", models.ProductStatusApproved); err != nil {
		return nil, err
	}
	if t.PendingProducts, err = count(&models.Product{}, "status = ?", models.ProductStatusPending); err != nil {
		return nil, err
	}
	if t.TotalCategories, err = count(&models.Category{}, ""); err != nil {
		return nil, err
	}
	if t.TotalBrands, err = count(&models.Brand{}, ""); err != nil {
		return nil, err
	}
	if t.TotalSellers, err = count(&models.Seller{}, ""); err != nil {
		return nil, err
	}
	if t.TotalCustomers, err = count(&models.Customer{}, ""); err != nil {
		return nil, err
	}
	if t.TotalDrivers, err = count(&models.Driver{}, ""); err != nil {
		return nil, err
	}
	if t.TotalOrders, err = count(&models.Order{}, ""); err != nil {
		return nil, err
	}
	if t.NewOrders, err = count(&models.Order{}, "status = ?", models.OrderStatusNew); err != nil {
		return nil, err
	}

	type sums struct {
		Views, Phone, Wa, Tg int64
	}
	var s sums
	if err = db.Model(&models.Product{}).
		Select(`COALESCE(SUM(views_count),0) AS views,
			COALESCE(SUM(phone_clicks),0) AS phone,
			COALESCE(SUM(whats_app_clicks),0) AS wa,
			COALESCE(SUM(telegram_clicks),0) AS tg`).Scan(&s).Error; err != nil {
		return nil, err
	}
	t.TotalViews = s.Views
	t.TotalPhoneClicks = s.Phone
	t.TotalWhatsAppClicks = s.Wa
	t.TotalTelegramClicks = s.Tg

	type money struct {
		Revenue, Profit float64
	}
	var m money
	if err = db.Model(&models.Order{}).
		Where("status = ?", models.OrderStatusCompleted).
		Select("COALESCE(SUM(total),0) AS revenue, COALESCE(SUM(profit),0) AS profit").
		Scan(&m).Error; err != nil {
		return nil, err
	}
	t.TotalRevenue = m.Revenue
	t.TotalProfit = m.Profit
	return t, nil
}

// SellerTotals for the seller dashboard.
type SellerTotals struct {
	TotalProducts       int64 `json:"total_products"`
	ApprovedProducts    int64 `json:"approved_products"`
	PendingProducts     int64 `json:"pending_products"`
	LowStock            int64 `json:"low_stock"`
	TotalViews          int64 `json:"total_views"`
	TotalPhoneClicks    int64 `json:"total_phone_clicks"`
	TotalWhatsAppClicks int64 `json:"total_whatsapp_clicks"`
	TotalTelegramClicks int64 `json:"total_telegram_clicks"`
}

func (r *TrackingRepository) SellerTotals(ctx context.Context, sellerID uuid.UUID) (*SellerTotals, error) {
	t := &SellerTotals{}
	db := r.db.WithContext(ctx)
	count := func(where string, args ...any) (int64, error) {
		var n int64
		return n, db.Model(&models.Product{}).Where("seller_id = ?", sellerID).Where(where, args...).Count(&n).Error
	}
	var err error
	if err = db.Model(&models.Product{}).Where("seller_id = ?", sellerID).Count(&t.TotalProducts).Error; err != nil {
		return nil, err
	}
	if t.ApprovedProducts, err = count("status = ?", models.ProductStatusApproved); err != nil {
		return nil, err
	}
	if t.PendingProducts, err = count("status = ?", models.ProductStatusPending); err != nil {
		return nil, err
	}
	var low int64
	if err = db.Model(&models.Product{}).
		Where("seller_id = ? AND stock_quantity <= minimum_stock", sellerID).Count(&low).Error; err != nil {
		return nil, err
	}
	t.LowStock = low

	type sums struct {
		Views, Phone, Wa, Tg int64
	}
	var s sums
	if err = db.Model(&models.Product{}).Where("seller_id = ?", sellerID).
		Select(`COALESCE(SUM(views_count),0) AS views,
			COALESCE(SUM(phone_clicks),0) AS phone,
			COALESCE(SUM(whats_app_clicks),0) AS wa,
			COALESCE(SUM(telegram_clicks),0) AS tg`).Scan(&s).Error; err != nil {
		return nil, err
	}
	t.TotalViews = s.Views
	t.TotalPhoneClicks = s.Phone
	t.TotalWhatsAppClicks = s.Wa
	t.TotalTelegramClicks = s.Tg
	return t, nil
}

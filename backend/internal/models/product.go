package models

import "github.com/google/uuid"

// Product is the central marketplace entity. Bilingual TJ/RU. Goes through a
// pending → approved flow; admin sets the sale price + contact owner.
type Product struct {
	BaseModel
	SellerID      uuid.UUID  `gorm:"type:uuid;index;not null" json:"seller_id"`
	CategoryID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"category_id"`
	SubcategoryID *uuid.UUID `gorm:"type:uuid;index" json:"subcategory_id,omitempty"`
	BrandID       *uuid.UUID `gorm:"type:uuid;index" json:"brand_id,omitempty"`

	Slug          string `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	SKU           string `gorm:"size:64;index" json:"sku"`
	NameTJ        string `gorm:"size:255;not null" json:"name_tj"`
	NameRU        string `gorm:"size:255;not null" json:"name_ru"`
	DescriptionTJ string `gorm:"type:text" json:"description_tj"`
	DescriptionRU string `gorm:"type:text" json:"description_ru"`

	Unit            string  `gorm:"size:20;not null;default:pcs" json:"unit"`
	Currency        string  `gorm:"size:8;not null;default:TJS" json:"currency"`
	CostPrice       float64 `gorm:"type:numeric(14,2);not null;default:0" json:"cost_price"` // seller's price
	SalePrice       float64 `gorm:"type:numeric(14,2);not null;default:0" json:"sale_price"` // public price (admin sets)
	DiscountPercent float64 `gorm:"type:numeric(5,2);not null;default:0" json:"discount_percent"`
	StockQuantity   int     `gorm:"not null;default:0" json:"stock_quantity"`
	MinimumStock    int     `gorm:"not null;default:0" json:"minimum_stock"`
	IsAvailable     bool    `gorm:"not null;default:true;index" json:"is_available"`
	IsFeatured      bool    `gorm:"not null;default:false;index" json:"is_featured"`

	// Contact routing chosen by admin during moderation.
	ContactOwner    string `gorm:"size:10;not null;default:admin" json:"contact_owner"` // admin|seller
	ContactPhone    string `gorm:"size:32" json:"contact_phone"`
	ContactWhatsApp string `gorm:"size:32" json:"contact_whatsapp"`
	ContactTelegram string `gorm:"size:32" json:"contact_telegram"`

	Status        string `gorm:"size:20;not null;default:draft;index" json:"status"`
	RejectionNote string `gorm:"type:text" json:"rejection_note"`

	// Aggregated tracking counters (atomically incremented).
	ViewsCount     int64 `gorm:"not null;default:0" json:"views_count"`
	PhoneClicks    int64 `gorm:"not null;default:0" json:"phone_clicks"`
	WhatsAppClicks int64 `gorm:"not null;default:0" json:"whatsapp_clicks"`
	TelegramClicks int64 `gorm:"not null;default:0" json:"telegram_clicks"`

	Seller      *Seller        `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Subcategory *Subcategory   `gorm:"foreignKey:SubcategoryID" json:"subcategory,omitempty"`
	Brand       *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Images      []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// ProductImage stores one image per row, ordered by SortOrder.
type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	Alt       string    `gorm:"size:255" json:"alt"`
	SortOrder int       `gorm:"not null;default:0;index" json:"sort_order"`
	IsCover   bool      `gorm:"not null;default:false" json:"is_cover"`
}

// Profit returns SalePrice - CostPrice.
func (p *Product) Profit() float64 { return p.SalePrice - p.CostPrice }

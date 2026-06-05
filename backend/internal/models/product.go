package models

import "github.com/google/uuid"

// Product is the central marketplace entity. All public-facing fields are bilingual.
type Product struct {
	BaseModel
	SellerID   uuid.UUID `gorm:"type:uuid;index;not null" json:"seller_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;index;not null" json:"category_id"`

	Slug string `gorm:"size:220;uniqueIndex;not null" json:"slug"`
	SKU  string `gorm:"size:64;index" json:"sku"`

	TitleTJ       string `gorm:"size:255;not null" json:"title_tj"`
	TitleRU       string `gorm:"size:255;not null" json:"title_ru"`
	DescriptionTJ string `gorm:"type:text" json:"description_tj"`
	DescriptionRU string `gorm:"type:text" json:"description_ru"`

	Price         float64 `gorm:"type:numeric(14,2);not null;default:0" json:"price"`
	Currency      string  `gorm:"size:8;not null;default:TJS" json:"currency"`
	Unit          string  `gorm:"size:20;not null;default:pcs" json:"unit"`
	StockQuantity int     `gorm:"not null;default:0" json:"stock_quantity"`
	IsAvailable   bool    `gorm:"not null;default:true;index" json:"is_available"`

	// Contact-routing fields configured by admin during moderation.
	ContactType    string `gorm:"size:10;not null;default:admin" json:"contact_type"` // "admin" | "seller"
	PhoneNumber    string `gorm:"size:32" json:"phone_number"`
	WhatsAppNumber string `gorm:"size:32" json:"whatsapp_number"`

	Status        string `gorm:"size:20;not null;default:draft;index" json:"status"`
	RejectionNote string `gorm:"type:text" json:"rejection_note"`
	IsFeatured    bool   `gorm:"not null;default:false;index" json:"is_featured"`

	// Aggregated tracking counters (updated atomically).
	ViewsCount     int64 `gorm:"not null;default:0" json:"views_count"`
	PhoneClicks    int64 `gorm:"not null;default:0" json:"phone_clicks"`
	WhatsAppClicks int64 `gorm:"not null;default:0" json:"whatsapp_clicks"`

	Seller   *Seller        `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Category *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Images   []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// ProductImage is a single image associated with a product.
type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	Alt       string    `gorm:"size:255" json:"alt"`
	SortOrder int       `gorm:"not null;default:0;index" json:"sort_order"`
	IsCover   bool      `gorm:"not null;default:false" json:"is_cover"`
}

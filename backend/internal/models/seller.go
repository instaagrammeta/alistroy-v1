package models

import "github.com/google/uuid"

// Seller is the public storefront / inventory owner. Created only by admin.
type Seller struct {
	BaseModel
	UserID           uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	FullName         string     `gorm:"size:160;not null" json:"full_name"`
	CompanyName      string     `gorm:"size:200" json:"company_name"`
	MarketName       string     `gorm:"size:160" json:"market_name"`
	Slug             string     `gorm:"size:200;uniqueIndex" json:"slug"`
	Phone            string     `gorm:"size:32" json:"phone"`
	PhoneAlt         string     `gorm:"size:32" json:"phone_alt"`
	WhatsApp         string     `gorm:"size:32" json:"whatsapp"`
	Telegram         string     `gorm:"size:32" json:"telegram"`
	TelegramUsername string     `gorm:"size:64" json:"telegram_username"`
	Address          string     `gorm:"size:500" json:"address"`
	City             string     `gorm:"size:100" json:"city"`
	BusinessCategory *uuid.UUID `gorm:"type:uuid;index" json:"business_category_id"` // optional category link
	LogoURL          string     `gorm:"size:500" json:"logo_url"`
	Notes            string     `gorm:"type:text" json:"notes"`
	Active           bool       `gorm:"not null;default:true;index" json:"active"`
	IsFeatured       bool       `gorm:"not null;default:false;index" json:"is_featured"`

	BusinessCat *Category `gorm:"foreignKey:BusinessCategory" json:"business_category,omitempty"`
}

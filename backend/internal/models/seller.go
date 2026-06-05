package models

import "github.com/google/uuid"

// Seller stores public storefront data attached to a user.
type Seller struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Name          string    `gorm:"size:160;not null" json:"name"`
	Slug          string    `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	DescriptionTJ string    `gorm:"type:text" json:"description_tj"`
	DescriptionRU string    `gorm:"type:text" json:"description_ru"`
	LogoURL       string    `gorm:"size:500" json:"logo_url"`
	Phone         string    `gorm:"size:32" json:"phone"`
	WhatsApp      string    `gorm:"size:32" json:"whatsapp"`
	Address       string    `gorm:"size:255" json:"address"`
	City          string    `gorm:"size:100" json:"city"`
	Status        string    `gorm:"size:20;not null;default:pending;index" json:"status"`
	IsFeatured    bool      `gorm:"not null;default:false;index" json:"is_featured"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

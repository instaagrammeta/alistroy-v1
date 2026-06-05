package models

import "time"

// User is the authentication root for any role: customer, seller, driver, admin.
// Role-specific profile data lives on Customer / Seller / Driver records that
// reference UserID 1:1.
type User struct {
	BaseModel
	Email        string     `gorm:"size:255;uniqueIndex" json:"email"`
	Phone        string     `gorm:"size:32;uniqueIndex" json:"phone"`
	Login        string     `gorm:"size:64;uniqueIndex" json:"login"` // optional, for sellers/drivers
	PasswordHash string     `gorm:"size:255" json:"-"`
	Name         string     `gorm:"size:160;not null" json:"name"`
	Role         string     `gorm:"size:20;not null;default:customer;index" json:"role"`
	Status       string     `gorm:"size:20;not null;default:active;index" json:"status"`
	Locale       string     `gorm:"size:5;not null;default:tg" json:"locale"`
	AvatarURL    string     `gorm:"size:500" json:"avatar_url"`
	GoogleID     string     `gorm:"size:128;index" json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`

	Customer *Customer `gorm:"foreignKey:UserID" json:"customer,omitempty"`
	Seller   *Seller   `gorm:"foreignKey:UserID" json:"seller,omitempty"`
	Driver   *Driver   `gorm:"foreignKey:UserID" json:"driver,omitempty"`
}

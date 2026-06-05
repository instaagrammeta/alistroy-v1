package models

import "time"

// User is the authentication entity. It can be customer / seller / admin.
type User struct {
	BaseModel
	Email          string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash   string     `gorm:"size:255;not null" json:"-"`
	Name           string     `gorm:"size:120;not null" json:"name"`
	Phone          string     `gorm:"size:32" json:"phone"`
	Role           string     `gorm:"size:20;not null;default:customer;index" json:"role"`
	Locale         string     `gorm:"size:5;not null;default:tg" json:"locale"`
	IsActive       bool       `gorm:"not null;default:true" json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	ResetToken     string     `gorm:"size:128;index" json:"-"`
	ResetExpiresAt *time.Time `json:"-"`

	Seller *Seller `gorm:"foreignKey:UserID" json:"seller,omitempty"`
}

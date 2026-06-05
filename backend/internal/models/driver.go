package models

import "github.com/google/uuid"

// Driver is a delivery operator. Created only by admin.
type Driver struct {
	BaseModel
	UserID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	FullName string    `gorm:"size:160;not null" json:"full_name"`
	Age      int       `gorm:"default:0" json:"age"`
	Phone    string    `gorm:"size:32" json:"phone"`
	PhoneAlt string    `gorm:"size:32" json:"phone_alt"`
	WhatsApp string    `gorm:"size:32" json:"whatsapp"`
	Telegram string    `gorm:"size:32" json:"telegram"`
	Vehicle  string    `gorm:"size:160" json:"vehicle"` // car / scooter / on foot
	PhotoURL string    `gorm:"size:500" json:"photo_url"`
	Notes    string    `gorm:"type:text" json:"notes"`
	Active   bool      `gorm:"not null;default:true;index" json:"active"`
	OnDuty   bool      `gorm:"not null;default:true" json:"on_duty"`
}

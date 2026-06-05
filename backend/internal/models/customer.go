package models

import "github.com/google/uuid"

// Customer holds the buyer profile attached to a User (1:1).
type Customer struct {
	BaseModel
	UserID  uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Company string    `gorm:"size:160" json:"company"`
	Address string    `gorm:"size:500" json:"address"`
	City    string    `gorm:"size:100" json:"city"`
	Notes   string    `gorm:"type:text" json:"notes"`
}

// Address is an additional shipping/contact address per customer.
type Address struct {
	BaseModel
	CustomerID uuid.UUID `gorm:"type:uuid;index;not null" json:"customer_id"`
	Label      string    `gorm:"size:80" json:"label"`
	Phone      string    `gorm:"size:32" json:"phone"`
	City       string    `gorm:"size:100" json:"city"`
	Street     string    `gorm:"size:500" json:"street"`
	IsPrimary  bool      `gorm:"not null;default:false" json:"is_primary"`
}

package models

import "github.com/google/uuid"

// Review is a rating + comment authored by a customer on a product.
type Review struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Rating    int       `gorm:"not null" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	Status    string    `gorm:"size:20;not null;default:pending;index" json:"status"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

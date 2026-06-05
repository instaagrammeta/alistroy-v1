package models

import "github.com/google/uuid"

// Review is a customer-authored rating + comment on a product.
type Review struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Rating    int       `gorm:"not null" json:"rating"` // 1..5
	Comment   string    `gorm:"type:text" json:"comment"`
	Status    string    `gorm:"size:20;not null;default:pending;index" json:"status"` // pending|approved|rejected

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

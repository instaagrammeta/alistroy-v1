package models

import "github.com/google/uuid"

// Favorite represents a "saved product" entry for a user.
type Favorite struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;index:idx_user_product,unique;not null" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;index:idx_user_product,unique;not null" json:"product_id"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

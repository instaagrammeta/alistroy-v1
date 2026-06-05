package models

import "github.com/google/uuid"

// CartItem is a per-customer pending line. We use a flat table (no Cart row)
// because each customer has at most one active cart and we never need cart-
// level metadata. The unique index (customer_id, product_id) guarantees no
// duplicates.
type CartItem struct {
	BaseModel
	CustomerID uuid.UUID `gorm:"type:uuid;index:idx_cart_customer_product,unique;not null" json:"customer_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;index:idx_cart_customer_product,unique;not null" json:"product_id"`
	Quantity   int       `gorm:"not null;default:1" json:"quantity"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

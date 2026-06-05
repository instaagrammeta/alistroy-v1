package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role constants used across the application.
const (
	RoleCustomer = "customer"
	RoleSeller   = "seller"
	RoleAdmin    = "admin"
)

// Product status constants.
const (
	ProductStatusDraft    = "draft"
	ProductStatusPending  = "pending"
	ProductStatusApproved = "approved"
	ProductStatusRejected = "rejected"
)

// Contact-routing type for a product.
const (
	ContactTypeAdmin  = "admin"
	ContactTypeSeller = "seller"
)

// Seller status constants.
const (
	SellerStatusPending  = "pending"
	SellerStatusApproved = "approved"
	SellerStatusBlocked  = "blocked"
)

// BaseModel embeds the common ID + timestamps + soft-delete columns.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

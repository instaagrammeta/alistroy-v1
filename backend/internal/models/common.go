package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role / status / type constants used across models.
const (
	RoleCustomer = "customer"
	RoleSeller   = "seller"
	RoleDriver   = "driver"
	RoleAdmin    = "admin"

	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBlocked  = "blocked"

	ProductStatusDraft    = "draft"
	ProductStatusPending  = "pending"
	ProductStatusApproved = "approved"
	ProductStatusRejected = "rejected"

	ContactOwnerAdmin  = "admin"
	ContactOwnerSeller = "seller"

	OrderStatusNew        = "new"
	OrderStatusProcessing = "processing"
	OrderStatusAssigned   = "assigned"
	OrderStatusOnDelivery = "on_delivery"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"

	BannerPositionHero        = "hero"
	BannerPositionFeatures    = "features" // sub-banner row with icons (free shipping etc)
	BannerPositionSide        = "side"
	BannerPositionMidLarge    = "mid_large"
	BannerPositionMidSmall    = "mid_small"
	BannerPositionFooter      = "footer"
	BannerPositionCategoryRow = "category_row" // square category tiles 1000x1000

	TxKindIncome   = "income"
	TxKindExpense  = "expense"
	TxKindPurchase = "purchase"

	NotifKindOrder   = "order"
	NotifKindChat    = "chat"
	NotifKindProduct = "product"
	NotifKindSystem  = "system"

	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
)

// Roles set for validation.
func ValidRoles() []string { return []string{RoleCustomer, RoleSeller, RoleDriver, RoleAdmin} }

// Time-stamped + soft-delete base.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

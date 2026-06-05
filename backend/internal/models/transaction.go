package models

import (
	"time"

	"github.com/google/uuid"
)

// FinancialTransaction is the accounting entry. Created automatically when an
// order completes (one income + one matching purchase) and may also be created
// manually by admin for ad-hoc expenses/income.
type FinancialTransaction struct {
	BaseModel
	Kind        string     `gorm:"size:20;not null;index" json:"kind"` // income|expense|purchase
	OrderID     *uuid.UUID `gorm:"type:uuid;index" json:"order_id,omitempty"`
	ProductID   *uuid.UUID `gorm:"type:uuid;index" json:"product_id,omitempty"`
	SellerID    *uuid.UUID `gorm:"type:uuid;index" json:"seller_id,omitempty"`
	Amount      float64    `gorm:"type:numeric(14,2);not null;default:0" json:"amount"`
	Currency    string     `gorm:"size:8;not null;default:TJS" json:"currency"`
	Description string     `gorm:"size:500" json:"description"`
	OccurredAt  time.Time  `gorm:"not null;index" json:"occurred_at"`
}

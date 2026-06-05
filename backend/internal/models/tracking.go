package models

import (
	"time"

	"github.com/google/uuid"
)

// Tracking event types.
const (
	EventView          = "view"
	EventPhoneClick    = "phone_click"
	EventWhatsAppClick = "whatsapp_click"
	EventTelegramClick = "telegram_click"
)

// TrackingEvent is an immutable analytics row. Aggregated counters are also
// kept on the Product row for O(1) reads.
type TrackingEvent struct {
	BaseModel
	ProductID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"product_id"`
	SellerID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"seller_id"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Event      string     `gorm:"size:32;not null;index" json:"event"`
	IP         string     `gorm:"size:64" json:"ip"`
	UserAgent  string     `gorm:"size:500" json:"user_agent"`
	Referrer   string     `gorm:"size:500" json:"referrer"`
	OccurredAt time.Time  `gorm:"not null;index" json:"occurred_at"`
}

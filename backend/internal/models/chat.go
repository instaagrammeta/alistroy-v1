package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatRoom is a conversation between a customer and the admin team.
// Each customer has exactly one room (lazy-created on first message).
type ChatRoom struct {
	BaseModel
	CustomerID     uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"customer_id"`
	LastMessageAt  *time.Time `gorm:"index" json:"last_message_at,omitempty"`
	UnreadAdmin    int        `gorm:"not null;default:0" json:"unread_admin"`
	UnreadCustomer int        `gorm:"not null;default:0" json:"unread_customer"`

	Customer *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Messages []ChatMessage `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
}

// ChatMessage is a single message inside a ChatRoom.
type ChatMessage struct {
	BaseModel
	RoomID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"room_id"`
	SenderID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"sender_id"`
	SenderRole string     `gorm:"size:20;not null" json:"sender_role"` // customer|admin
	Body       string     `gorm:"type:text" json:"body"`
	ReadAt     *time.Time `json:"read_at,omitempty"`

	Attachments []ChatAttachment `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
}

// ChatAttachment is a file (image/video) bound to a message.
type ChatAttachment struct {
	BaseModel
	MessageID uuid.UUID `gorm:"type:uuid;index;not null" json:"message_id"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	MimeType  string    `gorm:"size:100" json:"mime_type"`
	SizeBytes int64     `gorm:"not null;default:0" json:"size_bytes"`
}

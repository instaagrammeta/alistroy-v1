package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification is a user-targeted ping (order updated, new chat message, etc).
type Notification struct {
	BaseModel
	UserID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"user_id"`
	Kind    string     `gorm:"size:32;not null;index" json:"kind"`
	TitleTJ string     `gorm:"size:255" json:"title_tj"`
	TitleRU string     `gorm:"size:255" json:"title_ru"`
	BodyTJ  string     `gorm:"type:text" json:"body_tj"`
	BodyRU  string     `gorm:"type:text" json:"body_ru"`
	LinkURL string     `gorm:"size:500" json:"link_url"`
	ReadAt  *time.Time `json:"read_at,omitempty"`
}

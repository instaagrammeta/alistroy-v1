package models

import "github.com/google/uuid"

// Subcategory belongs to a Category. Bilingual + optional icon.
type Subcategory struct {
	BaseModel
	CategoryID uuid.UUID `gorm:"type:uuid;index;not null" json:"category_id"`
	Slug       string    `gorm:"size:160;uniqueIndex;not null" json:"slug"`
	NameTJ     string    `gorm:"size:200;not null" json:"name_tj"`
	NameRU     string    `gorm:"size:200;not null" json:"name_ru"`
	IconURL    string    `gorm:"size:500" json:"icon_url"`
	SortOrder  int       `gorm:"not null;default:0;index" json:"sort_order"`
	Active     bool      `gorm:"not null;default:true;index" json:"active"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

package models

// Category is bilingual (TJ / RU). It supports a single level of nesting via ParentID.
type Category struct {
	BaseModel
	Slug      string  `gorm:"size:120;uniqueIndex;not null" json:"slug"`
	TitleTJ   string  `gorm:"size:160;not null" json:"title_tj"`
	TitleRU   string  `gorm:"size:160;not null" json:"title_ru"`
	IconURL   string  `gorm:"size:500" json:"icon_url"`
	SortOrder int     `gorm:"not null;default:0;index" json:"sort_order"`
	IsActive  bool    `gorm:"not null;default:true;index" json:"is_active"`
	ParentID  *string `gorm:"type:uuid;index" json:"parent_id,omitempty"`

	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

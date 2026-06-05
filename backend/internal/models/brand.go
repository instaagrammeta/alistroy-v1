package models

// Brand is a supplier/manufacturer label. Logo should be at least 500x500.
type Brand struct {
	BaseModel
	Slug      string `gorm:"size:160;uniqueIndex;not null" json:"slug"`
	Name      string `gorm:"size:200;not null" json:"name"`
	LogoURL   string `gorm:"size:500" json:"logo_url"`
	Active    bool   `gorm:"not null;default:true;index" json:"active"`
	SortOrder int    `gorm:"not null;default:0" json:"sort_order"`
}

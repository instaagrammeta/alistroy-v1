package models

// Category is bilingual (TJ / RU). Holds an icon (500x500), description and
// SEO slug. Subcategories belong to it via Subcategory.CategoryID.
type Category struct {
	BaseModel
	Slug          string `gorm:"size:160;uniqueIndex;not null" json:"slug"`
	NameTJ        string `gorm:"size:200;not null" json:"name_tj"`
	NameRU        string `gorm:"size:200;not null" json:"name_ru"`
	DescriptionTJ string `gorm:"type:text" json:"description_tj"`
	DescriptionRU string `gorm:"type:text" json:"description_ru"`
	IconURL       string `gorm:"size:500" json:"icon_url"`
	BannerURL     string `gorm:"size:500" json:"banner_url"`
	SEOTitleTJ    string `gorm:"size:255" json:"seo_title_tj"`
	SEOTitleRU    string `gorm:"size:255" json:"seo_title_ru"`
	SEODescTJ     string `gorm:"type:text" json:"seo_description_tj"`
	SEODescRU     string `gorm:"type:text" json:"seo_description_ru"`
	SortOrder     int    `gorm:"not null;default:0;index" json:"sort_order"`
	Active        bool   `gorm:"not null;default:true;index" json:"active"`

	Subcategories []Subcategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"subcategories,omitempty"`
}

package models

// Banner is a CMS-controlled promotional unit on the public site. Admin uploads
// device-specific images (desktop / tablet / mobile), a target link and a
// sort order. Position determines where it renders (hero, side, mid_*, etc.).
type Banner struct {
	BaseModel
	Position      string `gorm:"size:32;not null;index" json:"position"`
	TitleTJ       string `gorm:"size:255" json:"title_tj"`
	TitleRU       string `gorm:"size:255" json:"title_ru"`
	DescriptionTJ string `gorm:"type:text" json:"description_tj"`
	DescriptionRU string `gorm:"type:text" json:"description_ru"`
	DesktopURL    string `gorm:"size:500" json:"desktop_url"`
	TabletURL     string `gorm:"size:500" json:"tablet_url"`
	MobileURL     string `gorm:"size:500" json:"mobile_url"`
	LinkURL       string `gorm:"size:500" json:"link_url"`
	SortOrder     int    `gorm:"not null;default:0;index" json:"sort_order"`
	Active        bool   `gorm:"not null;default:true;index" json:"active"`
}

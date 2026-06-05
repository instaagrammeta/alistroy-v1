package models

// Setting is a key/value store used for site-wide configuration (logo, favicon,
// marketplace contacts, SEO defaults). Editable from the admin panel.
type Setting struct {
	BaseModel
	Key   string `gorm:"size:80;uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

// Well-known setting keys.
const (
	SettingLogoURL          = "logo_url"
	SettingFaviconURL       = "favicon_url"
	SettingSiteNameTJ       = "site_name_tj"
	SettingSiteNameRU       = "site_name_ru"
	SettingTaglineTJ        = "tagline_tj"
	SettingTaglineRU        = "tagline_ru"
	SettingMarketplacePhone = "marketplace_phone"
	SettingMarketplaceWA    = "marketplace_whatsapp"
	SettingSEODescriptionTJ = "seo_description_tj"
	SettingSEODescriptionRU = "seo_description_ru"
	SettingHeroTitleTJ      = "hero_title_tj"
	SettingHeroTitleRU      = "hero_title_ru"
	SettingHeroSubtitleTJ   = "hero_subtitle_tj"
	SettingHeroSubtitleRU   = "hero_subtitle_ru"
	SettingHeroImageURL     = "hero_image_url"
	SettingFooterAddress    = "footer_address"
	SettingFooterEmail      = "footer_email"
)

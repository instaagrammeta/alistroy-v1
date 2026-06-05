package models

// Setting is a flat key/value store editable from the admin panel. Used for
// site logo/favicon, marketplace contacts, hero copy, footer info, etc.
type Setting struct {
	BaseModel
	Key   string `gorm:"size:80;uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

// Well-known keys.
const (
	SettingLogoURL          = "logo_url"
	SettingFaviconURL       = "favicon_url"
	SettingSiteNameTJ       = "site_name_tj"
	SettingSiteNameRU       = "site_name_ru"
	SettingTaglineTJ        = "tagline_tj"
	SettingTaglineRU        = "tagline_ru"
	SettingSEODescriptionTJ = "seo_description_tj"
	SettingSEODescriptionRU = "seo_description_ru"
	SettingMarketplacePhone = "marketplace_phone"
	SettingMarketplaceWA    = "marketplace_whatsapp"
	SettingMarketplaceTG    = "marketplace_telegram"
	SettingMarketplaceTGUN  = "marketplace_telegram_username"
	SettingFooterAddress    = "footer_address"
	SettingFooterEmail      = "footer_email"
	SettingHeroTitleTJ      = "hero_title_tj"
	SettingHeroTitleRU      = "hero_title_ru"
	SettingHeroSubtitleTJ   = "hero_subtitle_tj"
	SettingHeroSubtitleRU   = "hero_subtitle_ru"
)

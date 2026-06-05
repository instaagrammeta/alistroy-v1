package handlers

// Request DTOs with validation tags (go-playground/validator via Gin binding).

// ---- Auth ----

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=160"`
	Phone    string `json:"phone" binding:"required,min=5,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Address  string `json:"address" binding:"omitempty,max=500"`
	City     string `json:"city" binding:"omitempty,max=100"`
	Locale   string `json:"locale" binding:"omitempty,oneof=tg ru"`
}

type loginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // phone / email / login
	Password   string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type updateProfileRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=160"`
	Phone   string `json:"phone" binding:"omitempty,max=32"`
	Locale  string `json:"locale" binding:"omitempty,oneof=tg ru"`
	Address string `json:"address" binding:"omitempty,max=500"`
	City    string `json:"city" binding:"omitempty,max=100"`
	Company string `json:"company" binding:"omitempty,max=160"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"omitempty"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

type googleCallbackRequest struct {
	Code string `json:"code" binding:"required"`
}

// ---- Category / Subcategory / Brand ----

type categoryRequest struct {
	NameTJ        string `json:"name_tj" binding:"required,min=1,max=200"`
	NameRU        string `json:"name_ru" binding:"required,min=1,max=200"`
	DescriptionTJ string `json:"description_tj" binding:"omitempty"`
	DescriptionRU string `json:"description_ru" binding:"omitempty"`
	Slug          string `json:"slug" binding:"omitempty,max=160"`
	IconURL       string `json:"icon_url" binding:"omitempty,max=500"`
	BannerURL     string `json:"banner_url" binding:"omitempty,max=500"`
	SEOTitleTJ    string `json:"seo_title_tj" binding:"omitempty,max=255"`
	SEOTitleRU    string `json:"seo_title_ru" binding:"omitempty,max=255"`
	SEODescTJ     string `json:"seo_description_tj" binding:"omitempty"`
	SEODescRU     string `json:"seo_description_ru" binding:"omitempty"`
	SortOrder     int    `json:"sort_order"`
	Active        *bool  `json:"active"`
}

type subcategoryRequest struct {
	CategoryID string `json:"category_id" binding:"required,uuid"`
	NameTJ     string `json:"name_tj" binding:"required,min=1,max=200"`
	NameRU     string `json:"name_ru" binding:"required,min=1,max=200"`
	Slug       string `json:"slug" binding:"omitempty,max=160"`
	IconURL    string `json:"icon_url" binding:"omitempty,max=500"`
	SortOrder  int    `json:"sort_order"`
	Active     *bool  `json:"active"`
}

type brandRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=200"`
	Slug      string `json:"slug" binding:"omitempty,max=160"`
	LogoURL   string `json:"logo_url" binding:"omitempty,max=500"`
	SortOrder int    `json:"sort_order"`
	Active    *bool  `json:"active"`
}

// ---- Seller ----

type sellerRequest struct {
	FullName         string `json:"full_name" binding:"required,min=2,max=160"`
	CompanyName      string `json:"company_name" binding:"omitempty,max=200"`
	MarketName       string `json:"market_name" binding:"omitempty,max=160"`
	Phone            string `json:"phone" binding:"omitempty,max=32"`
	PhoneAlt         string `json:"phone_alt" binding:"omitempty,max=32"`
	WhatsApp         string `json:"whatsapp" binding:"omitempty,max=32"`
	Telegram         string `json:"telegram" binding:"omitempty,max=32"`
	TelegramUsername string `json:"telegram_username" binding:"omitempty,max=64"`
	Address          string `json:"address" binding:"omitempty,max=500"`
	City             string `json:"city" binding:"omitempty,max=100"`
	BusinessCategory string `json:"business_category_id" binding:"omitempty,uuid"`
	Notes            string `json:"notes" binding:"omitempty"`
	LogoURL          string `json:"logo_url" binding:"omitempty,max=500"`
	Login            string `json:"login" binding:"omitempty,max=64"`
	Password         string `json:"password" binding:"omitempty,min=6,max=128"`
	Active           *bool  `json:"active"`
	IsFeatured       *bool  `json:"is_featured"`
}

type sellerProfileRequest struct {
	FullName         string `json:"full_name" binding:"omitempty,max=160"`
	CompanyName      string `json:"company_name" binding:"omitempty,max=200"`
	MarketName       string `json:"market_name" binding:"omitempty,max=160"`
	Phone            string `json:"phone" binding:"omitempty,max=32"`
	PhoneAlt         string `json:"phone_alt" binding:"omitempty,max=32"`
	WhatsApp         string `json:"whatsapp" binding:"omitempty,max=32"`
	Telegram         string `json:"telegram" binding:"omitempty,max=32"`
	TelegramUsername string `json:"telegram_username" binding:"omitempty,max=64"`
	Address          string `json:"address" binding:"omitempty,max=500"`
	City             string `json:"city" binding:"omitempty,max=100"`
	Notes            string `json:"notes" binding:"omitempty"`
	LogoURL          string `json:"logo_url" binding:"omitempty,max=500"`
}

// ---- Driver ----

type driverRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=160"`
	Age      int    `json:"age" binding:"omitempty,min=0,max=120"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	PhoneAlt string `json:"phone_alt" binding:"omitempty,max=32"`
	WhatsApp string `json:"whatsapp" binding:"omitempty,max=32"`
	Telegram string `json:"telegram" binding:"omitempty,max=32"`
	Vehicle  string `json:"vehicle" binding:"omitempty,max=160"`
	PhotoURL string `json:"photo_url" binding:"omitempty,max=500"`
	Notes    string `json:"notes" binding:"omitempty"`
	Login    string `json:"login" binding:"omitempty,max=64"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
	Active   *bool  `json:"active"`
	OnDuty   *bool  `json:"on_duty"`
}

// ---- Product ----

type productImageDTO struct {
	URL string `json:"url" binding:"required,max=500"`
	Alt string `json:"alt" binding:"omitempty,max=255"`
}

type productCreateRequest struct {
	SellerID        string            `json:"seller_id" binding:"omitempty,uuid"`
	CategoryID      string            `json:"category_id" binding:"required,uuid"`
	SubcategoryID   string            `json:"subcategory_id" binding:"omitempty,uuid"`
	BrandID         string            `json:"brand_id" binding:"omitempty,uuid"`
	SKU             string            `json:"sku" binding:"omitempty,max=64"`
	NameTJ          string            `json:"name_tj" binding:"required,min=2,max=255"`
	NameRU          string            `json:"name_ru" binding:"required,min=2,max=255"`
	DescriptionTJ   string            `json:"description_tj" binding:"omitempty"`
	DescriptionRU   string            `json:"description_ru" binding:"omitempty"`
	Unit            string            `json:"unit" binding:"omitempty,max=20"`
	Currency        string            `json:"currency" binding:"omitempty,max=8"`
	CostPrice       float64           `json:"cost_price" binding:"omitempty,min=0"`
	SalePrice       float64           `json:"sale_price" binding:"omitempty,min=0"`
	DiscountPercent float64           `json:"discount_percent" binding:"omitempty,min=0,max=100"`
	StockQuantity   int               `json:"stock_quantity" binding:"omitempty,min=0"`
	MinimumStock    int               `json:"minimum_stock" binding:"omitempty,min=0"`
	IsAvailable     *bool             `json:"is_available"`
	IsFeatured      *bool             `json:"is_featured"`
	ContactOwner    string            `json:"contact_owner" binding:"omitempty,oneof=admin seller"`
	ContactPhone    string            `json:"contact_phone" binding:"omitempty,max=32"`
	ContactWhatsApp string            `json:"contact_whatsapp" binding:"omitempty,max=32"`
	ContactTelegram string            `json:"contact_telegram" binding:"omitempty,max=32"`
	Status          string            `json:"status" binding:"omitempty,oneof=draft pending approved rejected"`
	Images          []productImageDTO `json:"images" binding:"omitempty,dive"`
}

type productUpdateRequest struct {
	SellerID        string             `json:"seller_id" binding:"omitempty,uuid"`
	CategoryID      string             `json:"category_id" binding:"omitempty,uuid"`
	SubcategoryID   string             `json:"subcategory_id" binding:"omitempty,uuid"`
	BrandID         string             `json:"brand_id" binding:"omitempty,uuid"`
	SKU             string             `json:"sku" binding:"omitempty,max=64"`
	NameTJ          string             `json:"name_tj" binding:"omitempty,min=2,max=255"`
	NameRU          string             `json:"name_ru" binding:"omitempty,min=2,max=255"`
	DescriptionTJ   string             `json:"description_tj" binding:"omitempty"`
	DescriptionRU   string             `json:"description_ru" binding:"omitempty"`
	Unit            string             `json:"unit" binding:"omitempty,max=20"`
	Currency        string             `json:"currency" binding:"omitempty,max=8"`
	CostPrice       float64            `json:"cost_price" binding:"omitempty,min=0"`
	SalePrice       float64            `json:"sale_price" binding:"omitempty,min=0"`
	DiscountPercent float64            `json:"discount_percent" binding:"omitempty,min=0,max=100"`
	StockQuantity   int                `json:"stock_quantity" binding:"omitempty,min=0"`
	MinimumStock    int                `json:"minimum_stock" binding:"omitempty,min=0"`
	IsAvailable     *bool              `json:"is_available"`
	IsFeatured      *bool              `json:"is_featured"`
	ContactOwner    string             `json:"contact_owner" binding:"omitempty,oneof=admin seller"`
	ContactPhone    string             `json:"contact_phone" binding:"omitempty,max=32"`
	ContactWhatsApp string             `json:"contact_whatsapp" binding:"omitempty,max=32"`
	ContactTelegram string             `json:"contact_telegram" binding:"omitempty,max=32"`
	Status          string             `json:"status" binding:"omitempty,oneof=draft pending approved rejected"`
	Images          *[]productImageDTO `json:"images" binding:"omitempty,dive"`
}

type moderationRequest struct {
	Status          string  `json:"status" binding:"required,oneof=approved rejected pending"`
	SalePrice       float64 `json:"sale_price" binding:"omitempty,min=0"`
	ContactOwner    string  `json:"contact_owner" binding:"omitempty,oneof=admin seller"`
	ContactPhone    string  `json:"contact_phone" binding:"omitempty,max=32"`
	ContactWhatsApp string  `json:"contact_whatsapp" binding:"omitempty,max=32"`
	ContactTelegram string  `json:"contact_telegram" binding:"omitempty,max=32"`
	RejectionNote   string  `json:"rejection_note" binding:"omitempty,max=1000"`
	IsFeatured      *bool   `json:"is_featured"`
}

// ---- Customer ----

type customerRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=160"`
	Company  string `json:"company" binding:"omitempty,max=160"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive blocked"`
	Phone    string `json:"phone" binding:"required,max=32"`
	PhoneAlt string `json:"phone_alt" binding:"omitempty,max=32"`
	Address  string `json:"address" binding:"omitempty,max=500"`
	City     string `json:"city" binding:"omitempty,max=100"`
	Notes    string `json:"notes" binding:"omitempty"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// ---- Order ----

type orderItemDTO struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type orderCreateRequest struct {
	CustomerID      string         `json:"customer_id" binding:"omitempty,uuid"`
	CustomerName    string         `json:"customer_name" binding:"omitempty,max=160"`
	CustomerPhone   string         `json:"customer_phone" binding:"omitempty,max=32"`
	DeliveryAddress string         `json:"delivery_address" binding:"omitempty,max=500"`
	DeliveryDate    string         `json:"delivery_date" binding:"omitempty"`
	DiscountPercent float64        `json:"discount_percent" binding:"omitempty,min=0,max=100"`
	Notes           string         `json:"notes" binding:"omitempty"`
	Status          string         `json:"status" binding:"omitempty,oneof=new processing assigned on_delivery completed cancelled"`
	Items           []orderItemDTO `json:"items" binding:"required,min=1,dive"`
}

type orderStatusRequest struct {
	Status   string `json:"status" binding:"required,oneof=new processing assigned on_delivery completed cancelled"`
	DriverID string `json:"driver_id" binding:"omitempty,uuid"`
}

type checkoutRequest struct {
	DeliveryAddress string `json:"delivery_address" binding:"omitempty,max=500"`
	DeliveryDate    string `json:"delivery_date" binding:"omitempty"`
	Notes           string `json:"notes" binding:"omitempty"`
}

// ---- Cart ----

type cartSetRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=0"`
}

// ---- Reviews ----

type reviewCreateRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=2000"`
}

type reviewModerateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`
}

// ---- Banner ----

type bannerRequest struct {
	Position      string `json:"position" binding:"required,max=32"`
	TitleTJ       string `json:"title_tj" binding:"omitempty,max=255"`
	TitleRU       string `json:"title_ru" binding:"omitempty,max=255"`
	DescriptionTJ string `json:"description_tj" binding:"omitempty"`
	DescriptionRU string `json:"description_ru" binding:"omitempty"`
	DesktopURL    string `json:"desktop_url" binding:"omitempty,max=500"`
	TabletURL     string `json:"tablet_url" binding:"omitempty,max=500"`
	MobileURL     string `json:"mobile_url" binding:"omitempty,max=500"`
	LinkURL       string `json:"link_url" binding:"omitempty,max=500"`
	SortOrder     int    `json:"sort_order"`
	Active        *bool  `json:"active"`
}

// ---- Settings ----

type settingsUpdateRequest struct {
	Items map[string]string `json:"items" binding:"required"`
}

// ---- Tracking ----

type trackRequest struct {
	Event string `json:"event" binding:"required,oneof=view phone_click whatsapp_click telegram_click"`
}

// ---- Chat ----

type chatSendRequest struct {
	Body        string `json:"body" binding:"omitempty,max=4000"`
	Attachments []struct {
		URL       string `json:"url" binding:"required,max=500"`
		MimeType  string `json:"mime_type" binding:"omitempty,max=100"`
		SizeBytes int64  `json:"size_bytes" binding:"omitempty"`
	} `json:"attachments" binding:"omitempty,dive"`
}

// ---- Admin user update ----

type adminUserUpdateRequest struct {
	Name     string `json:"name" binding:"omitempty,max=160"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive blocked"`
	Password string `json:"password" binding:"omitempty,min=8,max=128"`
}

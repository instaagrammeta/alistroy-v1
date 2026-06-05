package handlers

// DTOs for incoming JSON payloads. Validation tags use go-playground/validator.

// ----- Auth -----

type registerRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8,max=128"`
	Name       string `json:"name" binding:"required,min=2,max=120"`
	Phone      string `json:"phone" binding:"omitempty,max=32"`
	Role       string `json:"role" binding:"omitempty,oneof=customer seller"`
	Locale     string `json:"locale" binding:"omitempty,oneof=tg ru"`
	SellerName string `json:"seller_name" binding:"omitempty,max=160"`
	City       string `json:"city" binding:"omitempty,max=100"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type updateProfileRequest struct {
	Name   string `json:"name" binding:"omitempty,min=2,max=120"`
	Phone  string `json:"phone" binding:"omitempty,max=32"`
	Locale string `json:"locale" binding:"omitempty,oneof=tg ru"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

// ----- Categories -----

type categoryRequest struct {
	TitleTJ   string  `json:"title_tj" binding:"required,min=1,max=160"`
	TitleRU   string  `json:"title_ru" binding:"required,min=1,max=160"`
	Slug      string  `json:"slug" binding:"omitempty,max=120"`
	IconURL   string  `json:"icon_url" binding:"omitempty,max=500"`
	SortOrder int     `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
	ParentID  *string `json:"parent_id"`
}

// ----- Seller profile -----

type sellerProfileRequest struct {
	Name          string `json:"name" binding:"omitempty,max=160"`
	DescriptionTJ string `json:"description_tj" binding:"omitempty"`
	DescriptionRU string `json:"description_ru" binding:"omitempty"`
	LogoURL       string `json:"logo_url" binding:"omitempty,max=500"`
	Phone         string `json:"phone" binding:"omitempty,max=32"`
	WhatsApp      string `json:"whatsapp" binding:"omitempty,max=32"`
	Address       string `json:"address" binding:"omitempty,max=255"`
	City          string `json:"city" binding:"omitempty,max=100"`
}

type adminSellerRequest struct {
	sellerProfileRequest
	Status     string `json:"status" binding:"omitempty,oneof=pending approved blocked"`
	IsFeatured *bool  `json:"is_featured"`
}

// ----- Products -----

type productImageDTO struct {
	URL string `json:"url" binding:"required,max=500"`
	Alt string `json:"alt" binding:"omitempty,max=255"`
}

type productCreateRequest struct {
	CategoryID    string            `json:"category_id" binding:"required,uuid"`
	TitleTJ       string            `json:"title_tj" binding:"required,min=2,max=255"`
	TitleRU       string            `json:"title_ru" binding:"required,min=2,max=255"`
	DescriptionTJ string            `json:"description_tj" binding:"omitempty"`
	DescriptionRU string            `json:"description_ru" binding:"omitempty"`
	Price         float64           `json:"price" binding:"min=0"`
	Currency      string            `json:"currency" binding:"omitempty,max=8"`
	Unit          string            `json:"unit" binding:"omitempty,max=20"`
	SKU           string            `json:"sku" binding:"omitempty,max=64"`
	StockQuantity int               `json:"stock_quantity" binding:"omitempty,min=0"`
	IsAvailable   *bool             `json:"is_available"`
	Images        []productImageDTO `json:"images" binding:"omitempty,dive"`
}

type productUpdateRequest struct {
	CategoryID    string             `json:"category_id" binding:"omitempty,uuid"`
	TitleTJ       string             `json:"title_tj" binding:"omitempty,min=2,max=255"`
	TitleRU       string             `json:"title_ru" binding:"omitempty,min=2,max=255"`
	DescriptionTJ string             `json:"description_tj" binding:"omitempty"`
	DescriptionRU string             `json:"description_ru" binding:"omitempty"`
	Price         float64            `json:"price" binding:"omitempty,min=0"`
	Currency      string             `json:"currency" binding:"omitempty,max=8"`
	Unit          string             `json:"unit" binding:"omitempty,max=20"`
	SKU           string             `json:"sku" binding:"omitempty,max=64"`
	StockQuantity int                `json:"stock_quantity" binding:"omitempty,min=0"`
	IsAvailable   *bool              `json:"is_available"`
	Images        *[]productImageDTO `json:"images" binding:"omitempty,dive"`
}

type adminProductUpdateRequest struct {
	productUpdateRequest
	ContactType    string `json:"contact_type" binding:"omitempty,oneof=admin seller"`
	PhoneNumber    string `json:"phone_number" binding:"omitempty,max=32"`
	WhatsAppNumber string `json:"whatsapp_number" binding:"omitempty,max=32"`
	IsFeatured     *bool  `json:"is_featured"`
	Status         string `json:"status" binding:"omitempty,oneof=draft pending approved rejected"`
	RejectionNote  string `json:"rejection_note" binding:"omitempty,max=1000"`
}

type moderationRequest struct {
	Status         string `json:"status" binding:"required,oneof=approved rejected pending"`
	ContactType    string `json:"contact_type" binding:"omitempty,oneof=admin seller"`
	PhoneNumber    string `json:"phone_number" binding:"omitempty,max=32"`
	WhatsAppNumber string `json:"whatsapp_number" binding:"omitempty,max=32"`
	RejectionNote  string `json:"rejection_note" binding:"omitempty,max=1000"`
}

// ----- Reviews -----

type reviewCreateRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=2000"`
}

type reviewModerateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`
}

// ----- Tracking -----

type trackRequest struct {
	Event string `json:"event" binding:"required,oneof=view phone_click whatsapp_click"`
}

// ----- Settings -----

type settingsUpdateRequest struct {
	Items map[string]string `json:"items" binding:"required"`
}

// ----- Admin user update -----

type adminUserUpdateRequest struct {
	Name     string `json:"name" binding:"omitempty,max=120"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Role     string `json:"role" binding:"omitempty,oneof=customer seller admin"`
	IsActive *bool  `json:"is_active"`
	Password string `json:"password" binding:"omitempty,min=8,max=128"`
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type SellerHandler struct {
	sellers  *services.SellerService
	tracking *services.TrackingService
}

func NewSellerHandler(s *services.SellerService, t *services.TrackingService) *SellerHandler {
	return &SellerHandler{sellers: s, tracking: t}
}

// ---- Public ----

func (h *SellerHandler) List(c *gin.Context) {
	page, size := paginate(c)
	active := true
	items, total, err := h.sellers.List(c.Request.Context(), services.ListSellersInput{
		Search: c.Query("q"), Active: &active, Featured: boolQuery(c, "featured"),
		Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *SellerHandler) Top(c *gin.Context) {
	items, err := h.sellers.Top(c.Request.Context(), intQuery(c, "limit", 6))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *SellerHandler) GetBySlug(c *gin.Context) {
	s, err := h.sellers.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

// ---- Seller self ----

func (h *SellerHandler) Me(c *gin.Context) {
	s, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

func (h *SellerHandler) UpdateMe(c *gin.Context) {
	var req sellerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	s, err := h.sellers.UpdateOwn(c.Request.Context(), middleware.MustUserID(c), services.SellerInput{
		FullName: req.FullName, CompanyName: req.CompanyName, MarketName: req.MarketName,
		Phone: req.Phone, PhoneAlt: req.PhoneAlt, WhatsApp: req.WhatsApp,
		Telegram: req.Telegram, TelegramUsername: req.TelegramUsername,
		Address: req.Address, City: req.City, Notes: req.Notes, LogoURL: req.LogoURL,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

func (h *SellerHandler) MyStats(c *gin.Context) {
	s, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	totals, err := h.tracking.SellerTotals(c.Request.Context(), s.ID)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, totals)
}

// ---- Admin ----

func (h *SellerHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.sellers.List(c.Request.Context(), services.ListSellersInput{
		Search: c.Query("q"), Active: boolQuery(c, "active"), Featured: boolQuery(c, "featured"),
		Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *SellerHandler) AdminGet(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	s, err := h.sellers.GetByID(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

func (h *SellerHandler) AdminCreate(c *gin.Context) {
	var req sellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	s, err := h.sellers.CreateByAdmin(c.Request.Context(), toSellerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, s)
}

func (h *SellerHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req sellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	s, err := h.sellers.UpdateByAdmin(c.Request.Context(), id, toSellerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

func (h *SellerHandler) AdminDelete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.sellers.Delete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func toSellerInput(req sellerRequest) services.SellerInput {
	in := services.SellerInput{
		FullName: req.FullName, CompanyName: req.CompanyName, MarketName: req.MarketName,
		Phone: req.Phone, PhoneAlt: req.PhoneAlt, WhatsApp: req.WhatsApp,
		Telegram: req.Telegram, TelegramUsername: req.TelegramUsername,
		Address: req.Address, City: req.City, Notes: req.Notes, LogoURL: req.LogoURL,
		Login: req.Login, Password: req.Password, Active: req.Active, IsFeatured: req.IsFeatured,
	}
	if req.BusinessCategory != "" {
		if id, err := uuid.Parse(req.BusinessCategory); err == nil {
			in.BusinessCategory = &id
		}
	}
	return in
}

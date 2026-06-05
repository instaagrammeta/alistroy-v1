package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type SellerHandler struct {
	sellers  *services.SellerService
	products *services.ProductService
	tracking *services.TrackingService
}

func NewSellerHandler(s *services.SellerService, p *services.ProductService, t *services.TrackingService) *SellerHandler {
	return &SellerHandler{sellers: s, products: p, tracking: t}
}

// Public ---

func (h *SellerHandler) List(c *gin.Context) {
	page, size := paginate(c)
	feat := boolQuery(c, "featured")
	items, total, err := h.sellers.List(c.Request.Context(), services.ListSellersInput{
		Search:   c.Query("q"),
		Status:   "approved",
		Featured: feat,
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
}

func (h *SellerHandler) Top(c *gin.Context) {
	limit := intQuery(c, "limit", 6)
	items, err := h.sellers.Top(c.Request.Context(), limit)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *SellerHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	s, err := h.sellers.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

// Authenticated seller ---

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
	s, err := h.sellers.UpdateOwn(c.Request.Context(), middleware.MustUserID(c), services.SellerProfileInput{
		Name:          req.Name,
		DescriptionTJ: req.DescriptionTJ,
		DescriptionRU: req.DescriptionRU,
		LogoURL:       req.LogoURL,
		Phone:         req.Phone,
		WhatsApp:      req.WhatsApp,
		Address:       req.Address,
		City:          req.City,
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

// Admin ---

func (h *SellerHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.sellers.List(c.Request.Context(), services.ListSellersInput{
		Search:   c.Query("q"),
		Status:   c.Query("status"),
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
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

func (h *SellerHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req adminSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	s, err := h.sellers.AdminUpdate(c.Request.Context(), id, services.SellerProfileInput{
		Name:          req.Name,
		DescriptionTJ: req.DescriptionTJ,
		DescriptionRU: req.DescriptionRU,
		LogoURL:       req.LogoURL,
		Phone:         req.Phone,
		WhatsApp:      req.WhatsApp,
		Address:       req.Address,
		City:          req.City,
	}, req.Status, req.IsFeatured)
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

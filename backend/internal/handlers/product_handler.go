package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type ProductHandler struct {
	products *services.ProductService
	sellers  *services.SellerService
	reviews  *services.ReviewService
}

func NewProductHandler(p *services.ProductService, s *services.SellerService, r *services.ReviewService) *ProductHandler {
	return &ProductHandler{products: p, sellers: s, reviews: r}
}

// ---- Public ----

func (h *ProductHandler) List(c *gin.Context) {
	page, size := paginate(c)
	var brandIDs []uuid.UUID
	for _, raw := range c.QueryArray("brand") {
		if id, err := uuid.Parse(raw); err == nil {
			brandIDs = append(brandIDs, id)
		}
	}
	items, total, err := h.products.ListPublic(c.Request.Context(), services.ListInput{
		Search:       c.Query("q"),
		CategorySlug: c.Query("category"),
		SubcatSlug:   c.Query("subcategory"),
		SellerSlug:   c.Query("seller"),
		BrandIDs:     brandIDs,
		IsFeatured:   boolQuery(c, "featured"),
		MinPrice:     float64Query(c, "min_price"),
		MaxPrice:     float64Query(c, "max_price"),
		Sort:         c.Query("sort"),
		Page:         page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	p, err := h.products.GetPublicBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	avg, cnt, _ := h.reviews.Average(c.Request.Context(), p.ID)
	httpx.OK(c, gin.H{"product": p, "avg_rating": avg, "review_count": cnt})
}

func (h *ProductHandler) Related(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	items, err := h.products.Related(c.Request.Context(), id, intQuery(c, "limit", 8))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *ProductHandler) PriceBounds(c *gin.Context) {
	min, max, err := h.products.PriceBounds(c.Request.Context(), c.Query("category"))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"min": min, "max": max})
}

// ---- Seller ----

func (h *ProductHandler) MyList(c *gin.Context) {
	page, size := paginate(c)
	seller, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	items, total, err := h.products.ListForSeller(c.Request.Context(), seller.ID, services.ListInput{
		Search: c.Query("q"), Status: c.Query("status"),
		LowStockOnly: boolPtrTrue(boolQuery(c, "low_stock")),
		Sort:         c.Query("sort"), Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *ProductHandler) MyCreate(c *gin.Context) {
	var req productCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	seller, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	in := toProductInput(req)
	p, err := h.products.CreateBySeller(c.Request.Context(), seller.ID, in)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, p)
}

func (h *ProductHandler) MyUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req productUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	seller, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	p, err := h.products.UpdateBySeller(c.Request.Context(), seller.ID, id, toProductUpdateInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProductHandler) MyDelete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	seller, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	if err := h.products.DeleteBySeller(c.Request.Context(), seller.ID, id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// ---- Admin ----

func (h *ProductHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.products.ListAdmin(c.Request.Context(), services.ListInput{
		Search: c.Query("q"), Status: c.Query("status"),
		CategorySlug: c.Query("category"), SellerSlug: c.Query("seller"),
		SellerID:     optionalUUIDQuery(c, "seller_id"),
		LowStockOnly: boolPtrTrue(boolQuery(c, "low_stock")),
		MinPrice:     float64Query(c, "min_price"), MaxPrice: float64Query(c, "max_price"),
		Sort: c.Query("sort"), Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *ProductHandler) AdminGet(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	p, err := h.products.GetByID(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProductHandler) AdminCreate(c *gin.Context) {
	var req productCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	in := toProductInput(req)
	if req.SellerID != "" {
		if sid, err := uuid.Parse(req.SellerID); err == nil {
			in.SellerID = sid
		}
	}
	p, err := h.products.CreateByAdmin(c.Request.Context(), in)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, p)
}

func (h *ProductHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req productUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	p, err := h.products.UpdateByAdmin(c.Request.Context(), id, toProductUpdateInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProductHandler) AdminModerate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req moderationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	p, err := h.products.Moderate(c.Request.Context(), id, services.ModerationDecision{
		Status: req.Status, SalePrice: req.SalePrice,
		ContactOwner: req.ContactOwner, ContactPhone: req.ContactPhone,
		ContactWhatsApp: req.ContactWhatsApp, ContactTelegram: req.ContactTelegram,
		RejectionNote: req.RejectionNote, IsFeatured: req.IsFeatured,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, p)
}

func (h *ProductHandler) AdminDelete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.products.DeleteByAdmin(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// ---- mapping helpers ----

func toProductInput(req productCreateRequest) services.ProductInput {
	in := services.ProductInput{
		SKU: req.SKU, NameTJ: req.NameTJ, NameRU: req.NameRU,
		DescriptionTJ: req.DescriptionTJ, DescriptionRU: req.DescriptionRU,
		Unit: req.Unit, Currency: req.Currency,
		CostPrice: req.CostPrice, SalePrice: req.SalePrice, DiscountPercent: req.DiscountPercent,
		StockQuantity: req.StockQuantity, MinimumStock: req.MinimumStock,
		IsAvailable: req.IsAvailable, IsFeatured: req.IsFeatured,
		ContactOwner: req.ContactOwner, ContactPhone: req.ContactPhone,
		ContactWhatsApp: req.ContactWhatsApp, ContactTelegram: req.ContactTelegram,
		Status: req.Status,
		Images: toServiceImages(req.Images),
	}
	if req.CategoryID != "" {
		if id, err := uuid.Parse(req.CategoryID); err == nil {
			in.CategoryID = id
		}
	}
	if req.SubcategoryID != "" {
		if id, err := uuid.Parse(req.SubcategoryID); err == nil {
			in.SubcategoryID = &id
		}
	}
	if req.BrandID != "" {
		if id, err := uuid.Parse(req.BrandID); err == nil {
			in.BrandID = &id
		}
	}
	return in
}

func toProductUpdateInput(req productUpdateRequest) services.ProductInput {
	in := services.ProductInput{
		SKU: req.SKU, NameTJ: req.NameTJ, NameRU: req.NameRU,
		DescriptionTJ: req.DescriptionTJ, DescriptionRU: req.DescriptionRU,
		Unit: req.Unit, Currency: req.Currency,
		CostPrice: req.CostPrice, SalePrice: req.SalePrice, DiscountPercent: req.DiscountPercent,
		StockQuantity: req.StockQuantity, MinimumStock: req.MinimumStock,
		IsAvailable: req.IsAvailable, IsFeatured: req.IsFeatured,
		ContactOwner: req.ContactOwner, ContactPhone: req.ContactPhone,
		ContactWhatsApp: req.ContactWhatsApp, ContactTelegram: req.ContactTelegram,
		Status: req.Status,
	}
	if req.CategoryID != "" {
		if id, err := uuid.Parse(req.CategoryID); err == nil {
			in.CategoryID = id
		}
	}
	if req.SubcategoryID != "" {
		if id, err := uuid.Parse(req.SubcategoryID); err == nil {
			in.SubcategoryID = &id
		}
	}
	if req.BrandID != "" {
		if id, err := uuid.Parse(req.BrandID); err == nil {
			in.BrandID = &id
		}
	}
	if req.SellerID != "" {
		if id, err := uuid.Parse(req.SellerID); err == nil {
			in.SellerID = id
		}
	}
	if req.Images != nil {
		in.Images = toServiceImages(*req.Images)
	}
	return in
}

func toServiceImages(in []productImageDTO) []services.ProductImageInput {
	out := make([]services.ProductImageInput, 0, len(in))
	for _, im := range in {
		out = append(out, services.ProductImageInput{URL: im.URL, Alt: im.Alt})
	}
	return out
}

func boolPtrTrue(p *bool) bool { return p != nil && *p }

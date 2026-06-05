package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type ProductHandler struct {
	products  *services.ProductService
	sellers   *services.SellerService
	favorites *services.FavoriteService
	reviews   *services.ReviewService
}

func NewProductHandler(p *services.ProductService, s *services.SellerService, f *services.FavoriteService, r *services.ReviewService) *ProductHandler {
	return &ProductHandler{products: p, sellers: s, favorites: f, reviews: r}
}

// ----- Public -----

func (h *ProductHandler) List(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.products.ListPublic(c.Request.Context(), services.ListProductsInput{
		Search:       c.Query("q"),
		CategorySlug: c.Query("category"),
		SellerSlug:   c.Query("seller"),
		IsFeatured:   boolQuery(c, "featured"),
		MinPrice:     float64Query(c, "min_price"),
		MaxPrice:     float64Query(c, "max_price"),
		Sort:         c.Query("sort"),
		Page:         page,
		PageSize:     size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	p, err := h.products.GetPublicBySlug(c.Request.Context(), slug)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	avg, cnt, _ := h.reviews.Average(c.Request.Context(), p.ID)
	httpx.OK(c, gin.H{
		"product":      p,
		"avg_rating":   avg,
		"review_count": cnt,
	})
}

func (h *ProductHandler) Related(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	limit := intQuery(c, "limit", 8)
	items, err := h.products.Related(c.Request.Context(), id, limit)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

// ----- Seller (auth: role=seller) -----

func (h *ProductHandler) MyList(c *gin.Context) {
	page, size := paginate(c)
	seller, err := h.sellers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	items, total, err := h.products.ListForSeller(c.Request.Context(), seller.ID, services.ListProductsInput{
		Search:   c.Query("q"),
		Status:   c.Query("status"),
		Sort:     c.Query("sort"),
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
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
	catID, _ := uuid.Parse(req.CategoryID)
	p, err := h.products.CreateBySeller(c.Request.Context(), seller.ID, services.ProductInput{
		CategoryID:    catID,
		TitleTJ:       req.TitleTJ,
		TitleRU:       req.TitleRU,
		DescriptionTJ: req.DescriptionTJ,
		DescriptionRU: req.DescriptionRU,
		Price:         req.Price,
		Currency:      req.Currency,
		Unit:          req.Unit,
		SKU:           req.SKU,
		StockQuantity: req.StockQuantity,
		IsAvailable:   req.IsAvailable,
		Images:        toServiceImages(req.Images),
	})
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
	in := services.ProductInput{
		TitleTJ:       req.TitleTJ,
		TitleRU:       req.TitleRU,
		DescriptionTJ: req.DescriptionTJ,
		DescriptionRU: req.DescriptionRU,
		Price:         req.Price,
		Currency:      req.Currency,
		Unit:          req.Unit,
		SKU:           req.SKU,
		StockQuantity: req.StockQuantity,
		IsAvailable:   req.IsAvailable,
	}
	if req.CategoryID != "" {
		if cid, err := uuid.Parse(req.CategoryID); err == nil {
			in.CategoryID = cid
		}
	}
	if req.Images != nil {
		in.Images = toServiceImages(*req.Images)
	}
	p, err := h.products.UpdateBySeller(c.Request.Context(), seller.ID, id, in)
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

// ----- Admin -----

func (h *ProductHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.products.ListAdmin(c.Request.Context(), services.ListProductsInput{
		Search:       c.Query("q"),
		Status:       c.Query("status"),
		CategorySlug: c.Query("category"),
		SellerSlug:   c.Query("seller"),
		MinPrice:     float64Query(c, "min_price"),
		MaxPrice:     float64Query(c, "max_price"),
		Sort:         c.Query("sort"),
		Page:         page,
		PageSize:     size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
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

func (h *ProductHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req adminProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	in := services.ProductInput{
		TitleTJ:        req.TitleTJ,
		TitleRU:        req.TitleRU,
		DescriptionTJ:  req.DescriptionTJ,
		DescriptionRU:  req.DescriptionRU,
		Price:          req.Price,
		Currency:       req.Currency,
		Unit:           req.Unit,
		SKU:            req.SKU,
		StockQuantity:  req.StockQuantity,
		IsAvailable:    req.IsAvailable,
		ContactType:    req.ContactType,
		PhoneNumber:    req.PhoneNumber,
		WhatsAppNumber: req.WhatsAppNumber,
		IsFeatured:     req.IsFeatured,
		Status:         req.Status,
		RejectionNote:  req.RejectionNote,
	}
	if req.CategoryID != "" {
		if cid, err := uuid.Parse(req.CategoryID); err == nil {
			in.CategoryID = cid
		}
	}
	if req.Images != nil {
		in.Images = toServiceImages(*req.Images)
	}
	p, err := h.products.AdminUpdate(c.Request.Context(), id, in)
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
		Status:         req.Status,
		ContactType:    req.ContactType,
		PhoneNumber:    req.PhoneNumber,
		WhatsAppNumber: req.WhatsAppNumber,
		RejectionNote:  req.RejectionNote,
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
	if err := h.products.AdminDelete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// helpers

func toServiceImages(in []productImageDTO) []services.ProductImageInput {
	out := make([]services.ProductImageInput, 0, len(in))
	for _, im := range in {
		out = append(out, services.ProductImageInput{URL: im.URL, Alt: im.Alt})
	}
	return out
}

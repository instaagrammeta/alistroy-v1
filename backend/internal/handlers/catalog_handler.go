package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type CatalogHandler struct {
	svc *services.CatalogService
}

func NewCatalogHandler(s *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: s}
}

// ---- Categories (public) ----

func (h *CatalogHandler) ListCategories(c *gin.Context) {
	items, err := h.svc.ListCategories(c.Request.Context(), true)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *CatalogHandler) PopularCategories(c *gin.Context) {
	items, err := h.svc.PopularCategories(c.Request.Context(), intQuery(c, "limit", 8))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *CatalogHandler) GetCategoryBySlug(c *gin.Context) {
	cat, err := h.svc.GetCategoryBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cat)
}

func (h *CatalogHandler) ListSubcategories(c *gin.Context) {
	catID, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.ListSubcategories(c.Request.Context(), catID, true)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

// ---- Categories (admin) ----

func (h *CatalogHandler) AdminListCategories(c *gin.Context) {
	items, err := h.svc.ListCategories(c.Request.Context(), false)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.CreateCategory(c.Request.Context(), toCategoryInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, cat)
}

func (h *CatalogHandler) UpdateCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.UpdateCategory(c.Request.Context(), id, toCategoryInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cat)
}

func (h *CatalogHandler) DeleteCategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// ---- Subcategories (admin) ----

func (h *CatalogHandler) CreateSubcategory(c *gin.Context) {
	var req subcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	catID, _ := uuid.Parse(req.CategoryID)
	sc, err := h.svc.CreateSubcategory(c.Request.Context(), services.SubcategoryInput{
		CategoryID: catID, NameTJ: req.NameTJ, NameRU: req.NameRU,
		Slug: req.Slug, IconURL: req.IconURL, SortOrder: req.SortOrder, Active: req.Active,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, sc)
}

func (h *CatalogHandler) UpdateSubcategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req subcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	in := services.SubcategoryInput{
		NameTJ: req.NameTJ, NameRU: req.NameRU, Slug: req.Slug,
		IconURL: req.IconURL, SortOrder: req.SortOrder, Active: req.Active,
	}
	if req.CategoryID != "" {
		if cid, err := uuid.Parse(req.CategoryID); err == nil {
			in.CategoryID = cid
		}
	}
	sc, err := h.svc.UpdateSubcategory(c.Request.Context(), id, in)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, sc)
}

func (h *CatalogHandler) DeleteSubcategory(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteSubcategory(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// ---- Brands ----

func (h *CatalogHandler) ListBrands(c *gin.Context) {
	onlyActive := true
	if r := c.Query("all"); r == "true" {
		onlyActive = false
	}
	items, err := h.svc.ListBrands(c.Request.Context(), onlyActive)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *CatalogHandler) CreateBrand(c *gin.Context) {
	var req brandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	b, err := h.svc.CreateBrand(c.Request.Context(), services.BrandInput{
		Name: req.Name, Slug: req.Slug, LogoURL: req.LogoURL, SortOrder: req.SortOrder, Active: req.Active,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, b)
}

func (h *CatalogHandler) UpdateBrand(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req brandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	b, err := h.svc.UpdateBrand(c.Request.Context(), id, services.BrandInput{
		Name: req.Name, Slug: req.Slug, LogoURL: req.LogoURL, SortOrder: req.SortOrder, Active: req.Active,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, b)
}

func (h *CatalogHandler) DeleteBrand(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteBrand(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func toCategoryInput(req categoryRequest) services.CategoryInput {
	return services.CategoryInput{
		NameTJ: req.NameTJ, NameRU: req.NameRU,
		DescriptionTJ: req.DescriptionTJ, DescriptionRU: req.DescriptionRU,
		Slug: req.Slug, IconURL: req.IconURL, BannerURL: req.BannerURL,
		SEOTitleTJ: req.SEOTitleTJ, SEOTitleRU: req.SEOTitleRU,
		SEODescTJ: req.SEODescTJ, SEODescRU: req.SEODescRU,
		SortOrder: req.SortOrder, Active: req.Active,
	}
}

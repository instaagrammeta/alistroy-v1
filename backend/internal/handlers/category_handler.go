package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type CategoryHandler struct {
	svc *services.CategoryService
}

func NewCategoryHandler(s *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

// Public
func (h *CategoryHandler) List(c *gin.Context) {
	items, err := h.svc.ListAll(c.Request.Context(), true)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

// Public
func (h *CategoryHandler) Popular(c *gin.Context) {
	limit := intQuery(c, "limit", 8)
	items, err := h.svc.Popular(c.Request.Context(), limit)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

// Public
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	cat, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cat)
}

// Admin
func (h *CategoryHandler) AdminList(c *gin.Context) {
	items, err := h.svc.ListAll(c.Request.Context(), false)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

// Admin
func (h *CategoryHandler) Create(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.Create(c.Request.Context(), services.CategoryInput{
		TitleTJ:   req.TitleTJ,
		TitleRU:   req.TitleRU,
		Slug:      req.Slug,
		IconURL:   req.IconURL,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		ParentID:  req.ParentID,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, cat)
}

// Admin
func (h *CategoryHandler) Update(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cat, err := h.svc.Update(c.Request.Context(), id, services.CategoryInput{
		TitleTJ:   req.TitleTJ,
		TitleRU:   req.TitleRU,
		Slug:      req.Slug,
		IconURL:   req.IconURL,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		ParentID:  req.ParentID,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cat)
}

// Admin
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

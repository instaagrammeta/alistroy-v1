package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type FavoriteHandler struct {
	svc *services.FavoriteService
}

func NewFavoriteHandler(s *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: s}
}

func (h *FavoriteHandler) List(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.svc.List(c.Request.Context(), middleware.MustUserID(c), page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Add(c.Request.Context(), middleware.MustUserID(c), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Remove(c.Request.Context(), middleware.MustUserID(c), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *FavoriteHandler) Has(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	has, err := h.svc.Has(c.Request.Context(), middleware.MustUserID(c), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"favorited": has})
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type ReviewHandler struct {
	svc *services.ReviewService
}

func NewReviewHandler(s *services.ReviewService) *ReviewHandler { return &ReviewHandler{svc: s} }

func (h *ReviewHandler) ListForProduct(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	page, size := paginate(c)
	items, total, err := h.svc.ListForProduct(c.Request.Context(), id, page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *ReviewHandler) Create(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req reviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	r, err := h.svc.Create(c.Request.Context(), middleware.MustUserID(c), services.ReviewInput{
		ProductID: id, Rating: req.Rating, Comment: req.Comment,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, r)
}

func (h *ReviewHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.svc.AdminList(c.Request.Context(), c.Query("status"), page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *ReviewHandler) AdminModerate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req reviewModerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	r, err := h.svc.Moderate(c.Request.Context(), id, req.Status)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, r)
}

func (h *ReviewHandler) AdminDelete(c *gin.Context) {
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

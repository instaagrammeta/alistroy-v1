package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type BannerHandler struct {
	svc *services.BannerService
}

func NewBannerHandler(s *services.BannerService) *BannerHandler { return &BannerHandler{svc: s} }

// Public — grouped by position for the homepage.
func (h *BannerHandler) Public(c *gin.Context) {
	grouped, err := h.svc.PublicGrouped(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, grouped)
}

func (h *BannerHandler) AdminList(c *gin.Context) {
	items, err := h.svc.AdminList(c.Request.Context(), c.Query("position"))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *BannerHandler) Create(c *gin.Context) {
	var req bannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	b, err := h.svc.Create(c.Request.Context(), toBannerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, b)
}

func (h *BannerHandler) Update(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req bannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	b, err := h.svc.Update(c.Request.Context(), id, toBannerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, b)
}

func (h *BannerHandler) Delete(c *gin.Context) {
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

func toBannerInput(req bannerRequest) services.BannerInput {
	return services.BannerInput{
		Position: req.Position, TitleTJ: req.TitleTJ, TitleRU: req.TitleRU,
		DescriptionTJ: req.DescriptionTJ, DescriptionRU: req.DescriptionRU,
		DesktopURL: req.DesktopURL, TabletURL: req.TabletURL, MobileURL: req.MobileURL,
		LinkURL: req.LinkURL, SortOrder: req.SortOrder, Active: req.Active,
	}
}

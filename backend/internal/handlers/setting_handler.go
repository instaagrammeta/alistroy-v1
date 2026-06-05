package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type SettingHandler struct {
	svc *services.SettingService
}

func NewSettingHandler(s *services.SettingService) *SettingHandler {
	return &SettingHandler{svc: s}
}

// Public — read-only view used by the frontend (logo, names, hero, etc.)
func (h *SettingHandler) Public(c *gin.Context) {
	all, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, all)
}

// Admin — read all (same data, kept separate to allow future filtering)
func (h *SettingHandler) AdminGet(c *gin.Context) {
	all, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, all)
}

// Admin — bulk update
func (h *SettingHandler) AdminUpdate(c *gin.Context) {
	var req settingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(c.Request.Context(), req.Items); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

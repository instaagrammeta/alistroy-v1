package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type TrackingHandler struct {
	svc *services.TrackingService
}

func NewTrackingHandler(s *services.TrackingService) *TrackingHandler {
	return &TrackingHandler{svc: s}
}

// Public — anonymous users may also send tracking events.
func (h *TrackingHandler) Track(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req trackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	var userID *uuid.UUID
	if u := middleware.MustUserID(c); u != uuid.Nil {
		userID = &u
	}
	if err := h.svc.Record(c.Request.Context(), services.TrackInput{
		ProductID: id,
		Event:     mapEvent(req.Event),
		UserID:    userID,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// Admin — global totals
func (h *TrackingHandler) AdminTotals(c *gin.Context) {
	t, err := h.svc.GlobalTotals(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, t)
}

func mapEvent(e string) string {
	// front-end uses friendly names; map to model constants
	switch e {
	case "view":
		return "view"
	case "phone_click":
		return "phone_click"
	case "whatsapp_click":
		return "whatsapp_click"
	}
	return ""
}

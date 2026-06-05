package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/ws"
)

// NotifySocketHandler upgrades to a websocket delivering the user's
// personal notification stream (notify:<userID>).
type NotifySocketHandler struct {
	hub *ws.Hub
}

func NewNotifySocketHandler(hub *ws.Hub) *NotifySocketHandler {
	return &NotifySocketHandler{hub: hub}
}

func (h *NotifySocketHandler) Socket(c *gin.Context) {
	uid := middleware.MustUserID(c)
	h.hub.Serve(c.Writer, c.Request, "notify:"+uid.String())
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/ws"
)

type ChatHandler struct {
	svc *services.ChatService
	hub *ws.Hub
}

func NewChatHandler(s *services.ChatService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{svc: s, hub: hub}
}

// ---- Customer ----

func (h *ChatHandler) MyRoom(c *gin.Context) {
	room, err := h.svc.CustomerRoom(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, room)
}

func (h *ChatHandler) MyMessages(c *gin.Context) {
	room, err := h.svc.CustomerRoom(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	page, size := paginate(c)
	msgs, total, err := h.svc.Messages(c.Request.Context(), room.ID, page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	_ = h.svc.MarkRead(c.Request.Context(), room.ID, models.RoleCustomer)
	httpx.List(c, msgs, httpx.NewPagination(page, size, total))
}

func (h *ChatHandler) CustomerSend(c *gin.Context) {
	var req chatSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	msg, err := h.svc.SendAsCustomer(c.Request.Context(), middleware.MustUserID(c), req.Body, toAttachments(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, msg)
}

// CustomerSocket upgrades to a websocket subscribed to the customer's room.
func (h *ChatHandler) CustomerSocket(c *gin.Context) {
	room, err := h.svc.CustomerRoom(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	h.hub.Serve(c.Writer, c.Request, "chat:"+room.ID.String())
}

// ---- Admin ----

func (h *ChatHandler) AdminRooms(c *gin.Context) {
	page, size := paginate(c)
	rooms, total, err := h.svc.ListRooms(c.Request.Context(), page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, rooms, httpx.NewPagination(page, size, total))
}

func (h *ChatHandler) AdminMessages(c *gin.Context) {
	roomID, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	page, size := paginate(c)
	msgs, total, err := h.svc.Messages(c.Request.Context(), roomID, page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	_ = h.svc.MarkRead(c.Request.Context(), roomID, models.RoleAdmin)
	httpx.List(c, msgs, httpx.NewPagination(page, size, total))
}

func (h *ChatHandler) AdminSend(c *gin.Context) {
	roomID, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req chatSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	msg, err := h.svc.SendAsAdmin(c.Request.Context(), middleware.MustUserID(c), roomID, req.Body, toAttachments(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, msg)
}

func (h *ChatHandler) AdminSocket(c *gin.Context) {
	roomID, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	h.hub.Serve(c.Writer, c.Request, "chat:"+roomID.String())
}

func toAttachments(req chatSendRequest) []services.AttachmentInput {
	out := make([]services.AttachmentInput, 0, len(req.Attachments))
	for _, a := range req.Attachments {
		out = append(out, services.AttachmentInput{URL: a.URL, MimeType: a.MimeType, SizeBytes: a.SizeBytes})
	}
	return out
}

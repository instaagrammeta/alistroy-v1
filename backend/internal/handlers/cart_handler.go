package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type CartHandler struct {
	cart      *services.CartService
	customers *services.CustomerService
}

func NewCartHandler(cart *services.CartService, cust *services.CustomerService) *CartHandler {
	return &CartHandler{cart: cart, customers: cust}
}

func (h *CartHandler) customerID(c *gin.Context) (uuid.UUID, bool) {
	cust, err := h.customers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return uuid.Nil, false
	}
	return cust.ID, true
}

func (h *CartHandler) List(c *gin.Context) {
	cid, ok := h.customerID(c)
	if !ok {
		return
	}
	items, err := h.cart.List(c.Request.Context(), cid)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, items)
}

func (h *CartHandler) Set(c *gin.Context) {
	cid, ok := h.customerID(c)
	if !ok {
		return
	}
	var req cartSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	pid, _ := uuid.Parse(req.ProductID)
	if err := h.cart.Set(c.Request.Context(), cid, pid, req.Quantity); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *CartHandler) Remove(c *gin.Context) {
	cid, ok := h.customerID(c)
	if !ok {
		return
	}
	pid, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.cart.Remove(c.Request.Context(), cid, pid); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *CartHandler) Clear(c *gin.Context) {
	cid, ok := h.customerID(c)
	if !ok {
		return
	}
	if err := h.cart.Clear(c.Request.Context(), cid); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

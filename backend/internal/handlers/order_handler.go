package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/receipt"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type OrderHandler struct {
	orders    *services.OrderService
	cart      *services.CartService
	customers *services.CustomerService
	settings  *services.SettingService
}

func NewOrderHandler(o *services.OrderService, cart *services.CartService, cust *services.CustomerService, set *services.SettingService) *OrderHandler {
	return &OrderHandler{orders: o, cart: cart, customers: cust, settings: set}
}

// ---- Customer ----

// Checkout converts the customer's cart into an order.
func (h *OrderHandler) Checkout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cust, err := h.customers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	var dd *time.Time
	if req.DeliveryDate != "" {
		if t, err := time.Parse("2006-01-02", req.DeliveryDate); err == nil {
			dd = &t
		}
	}
	order, err := h.orders.CreateFromCart(c.Request.Context(), cust.ID, req.DeliveryAddress, req.Notes, dd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, order)
}

// MyOrders lists the current customer's orders.
func (h *OrderHandler) MyOrders(c *gin.Context) {
	page, size := paginate(c)
	cust, err := h.customers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	items, total, err := h.orders.List(c.Request.Context(), repositories.ListOrdersParams{
		CustomerID: &cust.ID, Status: c.Query("status"), Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

// ---- Admin ----

func (h *OrderHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.orders.List(c.Request.Context(), repositories.ListOrdersParams{
		Search: c.Query("q"), Status: c.Query("status"),
		DriverID: optionalUUIDQuery(c, "driver_id"),
		From:     dateQuery(c, "from"), To: dateQuery(c, "to"),
		Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *OrderHandler) AdminGet(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	o, err := h.orders.Get(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, o)
}

func (h *OrderHandler) AdminCreate(c *gin.Context) {
	var req orderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	in := services.OrderInput{
		CustomerName: req.CustomerName, CustomerPhone: req.CustomerPhone,
		DeliveryAddress: req.DeliveryAddress, DiscountPercent: req.DiscountPercent,
		Notes: req.Notes, Status: req.Status,
	}
	if req.CustomerID != "" {
		if id, err := uuid.Parse(req.CustomerID); err == nil {
			in.CustomerID = &id
		}
	}
	if req.DeliveryDate != "" {
		if t, err := time.Parse("2006-01-02", req.DeliveryDate); err == nil {
			in.DeliveryDate = &t
		}
	}
	for _, it := range req.Items {
		pid, _ := uuid.Parse(it.ProductID)
		in.Items = append(in.Items, services.OrderItemInput{ProductID: pid, Quantity: it.Quantity})
	}
	order, err := h.orders.Create(c.Request.Context(), in)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, order)
}

func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req orderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	upd := services.StatusUpdate{Status: req.Status}
	if req.DriverID != "" {
		if did, err := uuid.Parse(req.DriverID); err == nil {
			upd.DriverID = &did
		}
	}
	o, err := h.orders.UpdateStatus(c.Request.Context(), id, upd)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, o)
}

func (h *OrderHandler) AdminDelete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.orders.Delete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// Receipt renders a printable HTML receipt for the order.
func (h *OrderHandler) Receipt(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	o, err := h.orders.Get(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	brand := "AliStroy CRM"
	subtitle := "Склади маводи сохтмонӣ"
	if h.settings != nil {
		if v, _ := h.settings.Get(c.Request.Context(), "site_name_ru"); v != "" {
			brand = v + " CRM"
		}
	}
	html, err := receipt.RenderHTML(brand, subtitle, o)
	if err != nil {
		httpx.Internal(c, "receipt render failed")
		return
	}
	c.Data(200, "text/html; charset=utf-8", html)
}

// ---- Driver ----

// DriverUpdateStatus lets the assigned driver advance the order status.
func (h *OrderHandler) DriverUpdateStatus(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req orderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	o, err := h.orders.UpdateStatus(c.Request.Context(), id, services.StatusUpdate{Status: req.Status})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, o)
}

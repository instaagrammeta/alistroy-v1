package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type DriverHandler struct {
	drivers *services.DriverService
	orders  *services.OrderService
}

func NewDriverHandler(d *services.DriverService, o *services.OrderService) *DriverHandler {
	return &DriverHandler{drivers: d, orders: o}
}

// ---- Admin ----

func (h *DriverHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.drivers.List(c.Request.Context(), c.Query("q"), boolQuery(c, "active"), page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *DriverHandler) AdminGet(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	d, err := h.drivers.GetByID(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, d)
}

func (h *DriverHandler) AdminCreate(c *gin.Context) {
	var req driverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	d, err := h.drivers.CreateByAdmin(c.Request.Context(), toDriverInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, d)
}

func (h *DriverHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req driverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	d, err := h.drivers.UpdateByAdmin(c.Request.Context(), id, toDriverInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, d)
}

func (h *DriverHandler) AdminDelete(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.drivers.Delete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

// ---- Driver self ----

func (h *DriverHandler) Me(c *gin.Context) {
	d, err := h.drivers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, d)
}

// MyOrders lists orders assigned to the current driver, optionally by status.
func (h *DriverHandler) MyOrders(c *gin.Context) {
	page, size := paginate(c)
	d, err := h.drivers.GetByUserID(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	items, total, err := h.orders.List(c.Request.Context(), repositories.ListOrdersParams{
		DriverID: &d.ID, Status: c.Query("status"), Page: page, Size: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func toDriverInput(req driverRequest) services.DriverInput {
	return services.DriverInput{
		FullName: req.FullName, Age: req.Age, Phone: req.Phone, PhoneAlt: req.PhoneAlt,
		WhatsApp: req.WhatsApp, Telegram: req.Telegram, Vehicle: req.Vehicle,
		PhotoURL: req.PhotoURL, Notes: req.Notes, Login: req.Login, Password: req.Password,
		Active: req.Active, OnDuty: req.OnDuty,
	}
}

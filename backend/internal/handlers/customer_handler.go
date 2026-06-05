package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type CustomerHandler struct {
	svc *services.CustomerService
}

func NewCustomerHandler(s *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: s}
}

func (h *CustomerHandler) AdminList(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.svc.List(c.Request.Context(), c.Query("q"), page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

func (h *CustomerHandler) AdminGet(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	cust, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cust)
}

func (h *CustomerHandler) AdminCreate(c *gin.Context) {
	var req customerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cust, err := h.svc.CreateByAdmin(c.Request.Context(), toCustomerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, cust)
}

func (h *CustomerHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req customerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	cust, err := h.svc.UpdateByAdmin(c.Request.Context(), id, toCustomerInput(req))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, cust)
}

func (h *CustomerHandler) AdminDelete(c *gin.Context) {
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

func toCustomerInput(req customerRequest) services.CustomerInput {
	return services.CustomerInput{
		Name: req.Name, Company: req.Company, Status: req.Status,
		Phone: req.Phone, PhoneAlt: req.PhoneAlt, Address: req.Address,
		City: req.City, Notes: req.Notes, Password: req.Password,
	}
}

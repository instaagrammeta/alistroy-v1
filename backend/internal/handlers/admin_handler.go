package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type AdminHandler struct {
	users    *services.UserService
	tracking *services.TrackingService
}

func NewAdminHandler(u *services.UserService, t *services.TrackingService) *AdminHandler {
	return &AdminHandler{users: u, tracking: t}
}

// Dashboard — combined totals
func (h *AdminHandler) Dashboard(c *gin.Context) {
	totals, err := h.tracking.GlobalTotals(c.Request.Context())
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, totals)
}

// Users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, size := paginate(c)
	items, total, err := h.users.List(c.Request.Context(), services.ListUsersInput{
		Role:     c.Query("role"),
		Search:   c.Query("q"),
		Page:     page,
		PageSize: size,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, newPagination(page, size, total))
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	u, err := h.users.Get(c.Request.Context(), id)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, u)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	var req adminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	u, err := h.users.AdminUpdate(c.Request.Context(), id, services.AdminUpdateUserInput{
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
		IsActive: req.IsActive,
		Password: req.Password,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, u)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, ok := parseUUID(c, "id")
	if !ok {
		return
	}
	if err := h.users.Delete(c.Request.Context(), id); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/oauth"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type AuthHandler struct {
	auth   *services.AuthService
	google *oauth.Google
}

func NewAuthHandler(a *services.AuthService, g *oauth.Google) *AuthHandler {
	return &AuthHandler{auth: a, google: g}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	pair, err := h.auth.RegisterCustomer(c.Request.Context(), services.RegisterCustomerInput{
		Name: req.Name, Phone: req.Phone, Email: req.Email, Password: req.Password,
		Address: req.Address, City: req.City, Locale: req.Locale,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.Created(c, pair)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	pair, err := h.auth.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, pair)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	pair, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, pair)
}

func (h *AuthHandler) Me(c *gin.Context) {
	u, err := h.auth.Me(c.Request.Context(), middleware.MustUserID(c))
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, u)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	u, err := h.auth.UpdateProfile(c.Request.Context(), middleware.MustUserID(c), services.UpdateProfileInput{
		Name: req.Name, Phone: req.Phone, Locale: req.Locale,
		Address: req.Address, City: req.City, Company: req.Company,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, u)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	if err := h.auth.ChangePassword(c.Request.Context(), middleware.MustUserID(c), req.OldPassword, req.NewPassword); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	httpx.OK(c, gin.H{"ok": true})
}

// ---- Google OAuth ----

func (h *AuthHandler) GoogleURL(c *gin.Context) {
	if !h.google.Enabled() {
		httpx.BadRequest(c, "google oauth not configured")
		return
	}
	state := c.DefaultQuery("state", "alistroy")
	httpx.OK(c, gin.H{"url": h.google.AuthURL(state)})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	if !h.google.Enabled() {
		httpx.BadRequest(c, "google oauth not configured")
		return
	}
	var req googleCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	profile, err := h.google.Exchange(c.Request.Context(), req.Code)
	if err != nil {
		httpx.Unauthorized(c, "google exchange failed")
		return
	}
	pair, needsProfile, err := h.auth.GoogleUpsert(c.Request.Context(), profile.ID, profile.Email, profile.Name, profile.Picture)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{
		"tokens":        pair,
		"needs_profile": needsProfile, // true → frontend must collect phone + address
	})
}

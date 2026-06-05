package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	pair, err := h.auth.Register(c.Request.Context(), services.RegisterInput{
		Email:      req.Email,
		Password:   req.Password,
		Name:       req.Name,
		Phone:      req.Phone,
		Role:       req.Role,
		Locale:     req.Locale,
		SellerName: req.SellerName,
		City:       req.City,
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
	pair, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
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
		Name:   req.Name,
		Phone:  req.Phone,
		Locale: req.Locale,
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

func (h *AuthHandler) Forgot(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	tok, err := h.auth.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	// Token is returned in dev/staging so admin can deliver it manually.
	httpx.OK(c, gin.H{"reset_token": tok})
}

func (h *AuthHandler) Reset(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	if err := h.auth.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, gin.H{"ok": true})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Stateless JWT — nothing to invalidate server-side.
	// Front-end should drop tokens. (Refresh token rotation/blacklisting can be
	// added later via Redis if needed.)
	httpx.OK(c, gin.H{"ok": true})
}

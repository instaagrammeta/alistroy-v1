package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	appjwt "github.com/instaagrammeta/alistroy-v1/backend/internal/jwt"
)

const (
	CtxUserID = "auth.user_id"
	CtxRole   = "auth.role"
)

// RequireAuth verifies the access token. The request is rejected if the token
// is missing or invalid.
func RequireAuth(jm *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseFromHeader(c, jm)
		if !ok {
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

// RequireRoles allows access only if the authenticated user has one of the given roles.
func RequireRoles(jm *appjwt.Manager, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseFromHeader(c, jm)
		if !ok {
			return
		}
		match := false
		for _, r := range roles {
			if claims.Role == r {
				match = true
				break
			}
		}
		if !match {
			httpx.Forbidden(c, "insufficient permissions")
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

// OptionalAuth populates user id/role if a valid token is present, but never blocks.
func OptionalAuth(jm *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}
		claims, err := jm.Parse(token)
		if err == nil && claims.Type == appjwt.AccessToken {
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxRole, claims.Role)
		}
		c.Next()
	}
}

func parseFromHeader(c *gin.Context, jm *appjwt.Manager) (*appjwt.Claims, bool) {
	token := extractToken(c)
	if token == "" {
		httpx.Unauthorized(c, "missing token")
		return nil, false
	}
	claims, err := jm.Parse(token)
	if err != nil {
		httpx.Unauthorized(c, "invalid token")
		return nil, false
	}
	if claims.Type != appjwt.AccessToken {
		httpx.Unauthorized(c, "wrong token type")
		return nil, false
	}
	return claims, true
}

func extractToken(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if h == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

// MustUserID returns the authenticated user's UUID. Caller must ensure RequireAuth ran.
func MustUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(CtxUserID)
	id, _ := v.(uuid.UUID)
	return id
}

// Role returns the authenticated user's role (empty if anonymous).
func Role(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	r, _ := v.(string)
	return r
}

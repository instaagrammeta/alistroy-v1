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

func RequireAuth(jm *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parse(c, jm)
		if !ok {
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func RequireRoles(jm *appjwt.Manager, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parse(c, jm)
		if !ok {
			return
		}
		matched := false
		for _, r := range roles {
			if claims.Role == r {
				matched = true
				break
			}
		}
		if !matched {
			httpx.Forbidden(c, "insufficient permissions")
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func OptionalAuth(jm *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tok := extract(c); tok != "" {
			if claims, err := jm.Parse(tok); err == nil && claims.Type == appjwt.AccessToken {
				c.Set(CtxUserID, claims.UserID)
				c.Set(CtxRole, claims.Role)
			}
		}
		c.Next()
	}
}

func parse(c *gin.Context, jm *appjwt.Manager) (*appjwt.Claims, bool) {
	tok := extract(c)
	if tok == "" {
		httpx.Unauthorized(c, "missing token")
		return nil, false
	}
	claims, err := jm.Parse(tok)
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

func extract(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	const p = "Bearer "
	if strings.HasPrefix(h, p) {
		return strings.TrimSpace(h[len(p):])
	}
	// allow token via query for websocket upgrades
	return c.Query("token")
}

func MustUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(CtxUserID)
	id, _ := v.(uuid.UUID)
	return id
}

func Role(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	r, _ := v.(string)
	return r
}

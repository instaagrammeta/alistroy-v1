package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

func parseUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		httpx.BadRequest(c, "invalid "+name)
		return uuid.Nil, false
	}
	return id, true
}

func optionalUUIDQuery(c *gin.Context, name string) *uuid.UUID {
	v := c.Query(name)
	if v == "" {
		return nil
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return nil
	}
	return &id
}

func intQuery(c *gin.Context, name string, def int) int {
	if v := c.Query(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func float64Query(c *gin.Context, name string) *float64 {
	if v := c.Query(name); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return &f
		}
	}
	return nil
}

func boolQuery(c *gin.Context, name string) *bool {
	if v := c.Query(name); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return &b
		}
	}
	return nil
}

func dateQuery(c *gin.Context, name string) *time.Time {
	v := c.Query(name)
	if v == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return &t
		}
	}
	return nil
}

func paginate(c *gin.Context) (page, size int) {
	page = intQuery(c, "page", 1)
	size = intQuery(c, "page_size", 20)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 20
	}
	return
}

func mapServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		httpx.NotFound(c, "resource not found")
	case errors.Is(err, services.ErrConflict):
		httpx.Conflict(c, "resource conflict")
	case errors.Is(err, services.ErrValidation):
		httpx.BadRequest(c, "validation failed")
	case errors.Is(err, services.ErrUnauthorized):
		httpx.Unauthorized(c, "unauthorized")
	case errors.Is(err, services.ErrForbidden):
		httpx.Forbidden(c, "forbidden")
	case errors.Is(err, services.ErrInvalidCredentials):
		httpx.Unauthorized(c, "invalid credentials")
	default:
		httpx.Internal(c, "internal error")
	}
}

package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

// parseUUID extracts and validates a UUID from a path param. Writes a 400 on failure.
func parseUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
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
	v := c.Query(name)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func float64Query(c *gin.Context, name string) *float64 {
	v := c.Query(name)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func boolQuery(c *gin.Context, name string) *bool {
	v := c.Query(name)
	if v == "" {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return nil
	}
	return &b
}

func paginate(c *gin.Context) (page, size int) {
	page = intQuery(c, "page", 1)
	size = intQuery(c, "page_size", 20)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return
}

func newPagination(page, size int, total int64) *httpx.Pagination {
	totalPages := 1
	if size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
		if totalPages < 1 {
			totalPages = 1
		}
	}
	return &httpx.Pagination{
		Page:       page,
		PageSize:   size,
		Total:      total,
		TotalPages: totalPages,
	}
}

// mapServiceError translates well-known service errors to HTTP responses.
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

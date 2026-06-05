package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

type ListResponse struct {
	Data       any         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPagination(page, size int, total int64) *Pagination {
	if size <= 0 {
		size = 20
	}
	totalPages := int((total + int64(size) - 1) / int64(size))
	if totalPages < 1 {
		totalPages = 1
	}
	return &Pagination{Page: page, PageSize: size, Total: total, TotalPages: totalPages}
}

func OK(c *gin.Context, data any)      { c.JSON(http.StatusOK, SuccessResponse{Data: data}) }
func Created(c *gin.Context, data any) { c.JSON(http.StatusCreated, SuccessResponse{Data: data}) }
func NoContent(c *gin.Context)         { c.Status(http.StatusNoContent) }
func List(c *gin.Context, data any, p *Pagination) {
	c.JSON(http.StatusOK, ListResponse{Data: data, Pagination: p})
}

func Error(c *gin.Context, status int, code, message string, details ...any) {
	body := ErrorBody{Code: code, Message: message}
	if len(details) > 0 {
		body.Details = details[0]
	}
	c.AbortWithStatusJSON(status, ErrorResponse{Error: body})
}
func BadRequest(c *gin.Context, msg string, d ...any) {
	Error(c, http.StatusBadRequest, "bad_request", msg, d...)
}
func Unauthorized(c *gin.Context, msg string) { Error(c, http.StatusUnauthorized, "unauthorized", msg) }
func Forbidden(c *gin.Context, msg string)    { Error(c, http.StatusForbidden, "forbidden", msg) }
func NotFound(c *gin.Context, msg string)     { Error(c, http.StatusNotFound, "not_found", msg) }
func Conflict(c *gin.Context, msg string)     { Error(c, http.StatusConflict, "conflict", msg) }
func Unprocessable(c *gin.Context, msg string, d ...any) {
	Error(c, http.StatusUnprocessableEntity, "unprocessable", msg, d...)
}
func Internal(c *gin.Context, msg string) { Error(c, http.StatusInternalServerError, "internal", msg) }

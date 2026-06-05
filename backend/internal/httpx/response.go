package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard JSON response envelopes used across the API.

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

func List(c *gin.Context, data any, p *Pagination) {
	c.JSON(http.StatusOK, ListResponse{Data: data, Pagination: p})
}

func Error(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "bad_request", message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "forbidden", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "not_found", message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "conflict", message)
}

func Internal(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "internal_error", message)
}

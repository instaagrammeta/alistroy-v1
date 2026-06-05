package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type UploadHandler struct {
	svc *services.UploadService
}

func NewUploadHandler(s *services.UploadService) *UploadHandler { return &UploadHandler{svc: s} }

// Single-file upload. Form-field "file" + optional "subdir".
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		httpx.BadRequest(c, "file is required")
		return
	}
	subdir := c.DefaultPostForm("subdir", "products")
	url, err := h.svc.SaveFile(file, subdir)
	if err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	httpx.OK(c, gin.H{"url": url})
}

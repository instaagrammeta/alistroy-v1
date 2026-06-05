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

// Upload accepts a single multipart "file" plus optional "subdir".
// The frontend tracks upload progress via the XHR/fetch upload events; the
// server simply stores the file and returns its public URL + size.
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		httpx.BadRequest(c, "file is required")
		return
	}
	subdir := c.DefaultPostForm("subdir", "misc")
	url, size, mime, err := h.svc.SaveFile(file, subdir)
	if err != nil {
		httpx.BadRequest(c, err.Error())
		return
	}
	httpx.OK(c, gin.H{"url": url, "size_bytes": size, "mime_type": mime})
}

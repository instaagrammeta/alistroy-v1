package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(s *services.ReportService) *ReportHandler { return &ReportHandler{svc: s} }

// Summary returns income/expense/purchase/profit for the resolved date range.
func (h *ReportHandler) Summary(c *gin.Context) {
	r := services.ResolveRange(c.Query("preset"), dateQuery(c, "from"), dateQuery(c, "to"))
	s, err := h.svc.FinancialSummary(c.Request.Context(), r)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.OK(c, s)
}

// Transactions lists financial transactions for the range.
func (h *ReportHandler) Transactions(c *gin.Context) {
	page, size := paginate(c)
	r := services.ResolveRange(c.Query("preset"), dateQuery(c, "from"), dateQuery(c, "to"))
	items, total, err := h.svc.Transactions(c.Request.Context(), c.Query("kind"), r, page, size)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	httpx.List(c, items, httpx.NewPagination(page, size, total))
}

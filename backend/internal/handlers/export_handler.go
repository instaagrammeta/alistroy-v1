package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/exporter"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

// ExportHandler builds Excel files honouring the same filters as the list views.
type ExportHandler struct {
	products  *services.ProductService
	catalog   *services.CatalogService
	orders    *services.OrderService
	customers *services.CustomerService
	sellers   *services.SellerService
	drivers   *services.DriverService
	report    *services.ReportService
	txRepo    *repositories.TransactionRepository
}

func NewExportHandler(
	p *services.ProductService,
	cat *services.CatalogService,
	o *services.OrderService,
	cust *services.CustomerService,
	s *services.SellerService,
	d *services.DriverService,
	r *services.ReportService,
	txRepo *repositories.TransactionRepository,
) *ExportHandler {
	return &ExportHandler{products: p, catalog: cat, orders: o, customers: cust, sellers: s, drivers: d, report: r, txRepo: txRepo}
}

const exportPageSize = 5000

func (h *ExportHandler) send(c *gin.Context, prefix string, data []byte, err error) {
	if err != nil {
		mapServiceError(c, err)
		return
	}
	name := exporter.FileName(prefix)
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *ExportHandler) Products(c *gin.Context) {
	items, _, err := h.products.ListAdmin(c.Request.Context(), services.ListInput{
		Search: c.Query("q"), Status: c.Query("status"),
		CategorySlug: c.Query("category"), SellerSlug: c.Query("seller"),
		Page: 1, Size: exportPageSize,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Products(items)
	h.send(c, "products", b, err)
}

func (h *ExportHandler) Categories(c *gin.Context) {
	items, err := h.catalog.ListCategories(c.Request.Context(), false)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Categories(items)
	h.send(c, "categories", b, err)
}

func (h *ExportHandler) Orders(c *gin.Context) {
	items, _, err := h.orders.List(c.Request.Context(), repositories.ListOrdersParams{
		Search: c.Query("q"), Status: c.Query("status"),
		From: dateQuery(c, "from"), To: dateQuery(c, "to"),
		Page: 1, Size: exportPageSize,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Orders(items)
	h.send(c, "orders", b, err)
}

func (h *ExportHandler) Customers(c *gin.Context) {
	items, _, err := h.customers.List(c.Request.Context(), c.Query("q"), 1, exportPageSize)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Customers(items)
	h.send(c, "customers", b, err)
}

func (h *ExportHandler) Sellers(c *gin.Context) {
	items, _, err := h.sellers.List(c.Request.Context(), services.ListSellersInput{
		Search: c.Query("q"), Active: boolQuery(c, "active"), Page: 1, Size: exportPageSize,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Sellers(items)
	h.send(c, "sellers", b, err)
}

func (h *ExportHandler) Drivers(c *gin.Context) {
	items, _, err := h.drivers.List(c.Request.Context(), c.Query("q"), boolQuery(c, "active"), 1, exportPageSize)
	if err != nil {
		mapServiceError(c, err)
		return
	}
	b, err := exporter.Drivers(items)
	h.send(c, "drivers", b, err)
}

func (h *ExportHandler) Transactions(c *gin.Context) {
	r := services.ResolveRange(c.Query("preset"), dateQuery(c, "from"), dateQuery(c, "to"))
	items, _, err := h.txRepo.List(c.Request.Context(), repositories.ListTxParams{
		Kind: c.Query("kind"), From: r.From, To: r.To, Page: 1, Size: exportPageSize,
	})
	if err != nil {
		mapServiceError(c, err)
		return
	}
	summary, _ := h.txRepo.Summarize(c.Request.Context(), r.From, r.To)
	b, err := exporter.Transactions(items, summary)
	h.send(c, "report", b, err)
}

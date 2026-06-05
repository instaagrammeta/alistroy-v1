package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

// SEOHandler emits sitemap.xml and robots.txt built from live data.
type SEOHandler struct {
	products  *services.ProductService
	catalog   *services.CatalogService
	sellers   *services.SellerService
	publicURL string
}

func NewSEOHandler(p *services.ProductService, cat *services.CatalogService, s *services.SellerService, publicURL string) *SEOHandler {
	return &SEOHandler{products: p, catalog: cat, sellers: s, publicURL: strings.TrimRight(publicURL, "/")}
}

func (h *SEOHandler) Robots(c *gin.Context) {
	body := "User-agent: *\nAllow: /\nDisallow: /admin\nDisallow: /seller\nDisallow: /driver\nDisallow: /me\n\nSitemap: " + h.publicURL + "/sitemap.xml\n"
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}

func (h *SEOHandler) Sitemap(c *gin.Context) {
	ctx := c.Request.Context()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	add := func(loc string) {
		b.WriteString("  <url><loc>" + h.publicURL + loc + "</loc></url>\n")
	}
	add("/")
	add("/products")
	add("/categories")
	add("/sellers")

	if cats, err := h.catalog.ListCategories(ctx, true); err == nil {
		for _, cat := range cats {
			add("/categories/" + cat.Slug)
		}
	}
	if sellers, _, err := h.sellers.List(ctx, services.ListSellersInput{Page: 1, Size: 1000}); err == nil {
		for _, s := range sellers {
			if s.Slug != "" {
				add("/sellers/" + s.Slug)
			}
		}
	}
	if prods, _, err := h.products.ListPublic(ctx, services.ListInput{Page: 1, Size: 5000}); err == nil {
		for _, p := range prods {
			add("/products/" + p.Slug)
		}
	}

	b.WriteString("</urlset>\n")
	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(b.String()))
}

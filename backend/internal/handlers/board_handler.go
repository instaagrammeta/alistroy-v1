package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/httpx"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
)

// BoardHandler builds the XMind-style tree:
//
//	Category → Subcategory → Products
//
// consumed by the admin visual board (searchable, zoomable on the frontend).
type BoardHandler struct {
	catalog  *services.CatalogService
	products *services.ProductService
}

func NewBoardHandler(c *services.CatalogService, p *services.ProductService) *BoardHandler {
	return &BoardHandler{catalog: c, products: p}
}

type boardProduct struct {
	ID     string  `json:"id"`
	NameTJ string  `json:"name_tj"`
	NameRU string  `json:"name_ru"`
	Slug   string  `json:"slug"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type boardSubcategory struct {
	ID       string         `json:"id"`
	NameTJ   string         `json:"name_tj"`
	NameRU   string         `json:"name_ru"`
	Products []boardProduct `json:"products"`
}

type boardCategory struct {
	ID            string             `json:"id"`
	NameTJ        string             `json:"name_tj"`
	NameRU        string             `json:"name_ru"`
	Slug          string             `json:"slug"`
	Subcategories []boardSubcategory `json:"subcategories"`
	Products      []boardProduct     `json:"products"` // products without a subcategory
}

func (h *BoardHandler) Tree(c *gin.Context) {
	ctx := c.Request.Context()
	cats, err := h.catalog.ListCategories(ctx, false)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	out := make([]boardCategory, 0, len(cats))
	for _, cat := range cats {
		bc := boardCategory{ID: cat.ID.String(), NameTJ: cat.NameTJ, NameRU: cat.NameRU, Slug: cat.Slug}

		// All products in this category (admin view: any status).
		prods, _, err := h.products.ListAdmin(ctx, services.ListInput{
			CategorySlug: cat.Slug, Page: 1, Size: 1000, Sort: "newest",
		})
		if err != nil {
			mapServiceError(c, err)
			return
		}

		// Bucket products by subcategory.
		bySub := map[string][]boardProduct{}
		var noSub []boardProduct
		for _, p := range prods {
			bp := boardProduct{
				ID: p.ID.String(), NameTJ: p.NameTJ, NameRU: p.NameRU,
				Slug: p.Slug, Price: p.SalePrice, Status: p.Status,
			}
			if p.SubcategoryID != nil {
				key := p.SubcategoryID.String()
				bySub[key] = append(bySub[key], bp)
			} else {
				noSub = append(noSub, bp)
			}
		}

		for _, sub := range cat.Subcategories {
			bc.Subcategories = append(bc.Subcategories, boardSubcategory{
				ID: sub.ID.String(), NameTJ: sub.NameTJ, NameRU: sub.NameRU,
				Products: bySub[sub.ID.String()],
			})
		}
		bc.Products = noSub
		out = append(out, bc)
	}
	httpx.OK(c, out)
}

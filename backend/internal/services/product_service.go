package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type ProductService struct {
	products *repositories.ProductRepository
	sellers  *repositories.SellerRepository
	cats     *repositories.CategoryRepository
	notifier *NotificationService
	defaults ContactDefaults
}

type ContactDefaults struct {
	Phone    string
	WhatsApp string
	Telegram string
}

func NewProductService(
	p *repositories.ProductRepository,
	s *repositories.SellerRepository,
	c *repositories.CategoryRepository,
	n *NotificationService,
	d ContactDefaults,
) *ProductService {
	return &ProductService{products: p, sellers: s, cats: c, notifier: n, defaults: d}
}

type ProductImageInput struct {
	URL string
	Alt string
}

type ProductInput struct {
	SellerID        uuid.UUID
	CategoryID      uuid.UUID
	SubcategoryID   *uuid.UUID
	BrandID         *uuid.UUID
	SKU             string
	NameTJ          string
	NameRU          string
	DescriptionTJ   string
	DescriptionRU   string
	Unit            string
	Currency        string
	CostPrice       float64
	SalePrice       float64
	DiscountPercent float64
	StockQuantity   int
	MinimumStock    int
	IsAvailable     *bool
	Images          []ProductImageInput

	// Admin-only:
	ContactOwner    string
	ContactPhone    string
	ContactWhatsApp string
	ContactTelegram string
	IsFeatured      *bool
	Status          string
	RejectionNote   string
}

// CreateByAdmin: admin chooses seller, sets sale price, contacts immediately.
func (s *ProductService) CreateByAdmin(ctx context.Context, in ProductInput) (*models.Product, error) {
	if in.SellerID == uuid.Nil {
		return nil, ErrValidation
	}
	p, err := s.build(ctx, in, models.ProductStatusApproved)
	if err != nil {
		return nil, err
	}
	s.applyContacts(ctx, p, in, true)
	if err := s.products.Create(ctx, p); err != nil {
		return nil, err
	}
	if len(in.Images) > 0 {
		if err := s.products.ReplaceImages(ctx, p.ID, buildImages(p.ID, in.Images)); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

// CreateBySeller: seller's price becomes cost price; status forced to pending.
func (s *ProductService) CreateBySeller(ctx context.Context, sellerID uuid.UUID, in ProductInput) (*models.Product, error) {
	in.SellerID = sellerID
	// Seller-entered price is the COST price.
	in.CostPrice = in.SalePrice
	in.SalePrice = 0
	p, err := s.build(ctx, in, models.ProductStatusPending)
	if err != nil {
		return nil, err
	}
	// Default contacts to marketplace until admin decides.
	p.ContactOwner = models.ContactOwnerAdmin
	p.ContactPhone = s.defaults.Phone
	p.ContactWhatsApp = s.defaults.WhatsApp
	p.ContactTelegram = s.defaults.Telegram
	if err := s.products.Create(ctx, p); err != nil {
		return nil, err
	}
	if len(in.Images) > 0 {
		if err := s.products.ReplaceImages(ctx, p.ID, buildImages(p.ID, in.Images)); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

func (s *ProductService) build(ctx context.Context, in ProductInput, status string) (*models.Product, error) {
	if in.NameTJ == "" || in.NameRU == "" || in.CategoryID == uuid.Nil {
		return nil, ErrValidation
	}
	if _, err := s.cats.FindByID(ctx, in.CategoryID); err != nil {
		return nil, ErrValidation
	}
	slug, err := uniqueSlug(ctx, firstNonEmpty(in.NameRU, in.NameTJ), s.products.ExistsBySlug)
	if err != nil {
		return nil, err
	}
	p := &models.Product{
		SellerID:        in.SellerID,
		CategoryID:      in.CategoryID,
		SubcategoryID:   in.SubcategoryID,
		BrandID:         in.BrandID,
		Slug:            slug,
		SKU:             in.SKU,
		NameTJ:          in.NameTJ,
		NameRU:          in.NameRU,
		DescriptionTJ:   in.DescriptionTJ,
		DescriptionRU:   in.DescriptionRU,
		Unit:            defaultStr(in.Unit, "pcs"),
		Currency:        defaultStr(in.Currency, "TJS"),
		CostPrice:       in.CostPrice,
		SalePrice:       in.SalePrice,
		DiscountPercent: in.DiscountPercent,
		StockQuantity:   in.StockQuantity,
		MinimumStock:    in.MinimumStock,
		IsAvailable:     boolOr(in.IsAvailable, true),
		Status:          status,
	}
	return p, nil
}

func (s *ProductService) UpdateBySeller(ctx context.Context, sellerID, productID uuid.UUID, in ProductInput) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	if p.SellerID != sellerID {
		return nil, ErrForbidden
	}
	// Seller-entered sale price is treated as cost price.
	if in.SalePrice > 0 {
		in.CostPrice = in.SalePrice
		in.SalePrice = 0
	}
	s.applyEditable(p, in)
	// Editing resets to pending for re-approval.
	p.Status = models.ProductStatusPending
	p.RejectionNote = ""
	if err := s.products.Save(ctx, p); err != nil {
		return nil, err
	}
	if in.Images != nil {
		if err := s.products.ReplaceImages(ctx, p.ID, buildImages(p.ID, in.Images)); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

func (s *ProductService) UpdateByAdmin(ctx context.Context, productID uuid.UUID, in ProductInput) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	s.applyEditable(p, in)
	if in.CostPrice > 0 {
		p.CostPrice = in.CostPrice
	}
	if in.SalePrice > 0 {
		p.SalePrice = in.SalePrice
	}
	if in.SellerID != uuid.Nil {
		p.SellerID = in.SellerID
	}
	if in.IsFeatured != nil {
		p.IsFeatured = *in.IsFeatured
	}
	if in.Status != "" {
		p.Status = in.Status
	}
	s.applyContacts(ctx, p, in, false)
	if err := s.products.Save(ctx, p); err != nil {
		return nil, err
	}
	if in.Images != nil {
		if err := s.products.ReplaceImages(ctx, p.ID, buildImages(p.ID, in.Images)); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

// ModerationDecision is what admin submits on approve/reject.
type ModerationDecision struct {
	Status          string
	SalePrice       float64
	ContactOwner    string
	ContactPhone    string
	ContactWhatsApp string
	ContactTelegram string
	RejectionNote   string
	IsFeatured      *bool
}

func (s *ProductService) Moderate(ctx context.Context, productID uuid.UUID, d ModerationDecision) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	switch d.Status {
	case models.ProductStatusApproved:
		p.Status = models.ProductStatusApproved
		p.RejectionNote = ""
		if d.SalePrice > 0 {
			p.SalePrice = d.SalePrice
		}
		s.applyContacts(ctx, p, ProductInput{
			ContactOwner:    d.ContactOwner,
			ContactPhone:    d.ContactPhone,
			ContactWhatsApp: d.ContactWhatsApp,
			ContactTelegram: d.ContactTelegram,
		}, true)
		if d.IsFeatured != nil {
			p.IsFeatured = *d.IsFeatured
		}
	case models.ProductStatusRejected:
		p.Status = models.ProductStatusRejected
		p.RejectionNote = d.RejectionNote
	case models.ProductStatusPending:
		p.Status = models.ProductStatusPending
		p.RejectionNote = d.RejectionNote
	default:
		return nil, ErrValidation
	}
	if err := s.products.Save(ctx, p); err != nil {
		return nil, err
	}
	// Notify seller.
	if seller, err := s.sellers.FindByID(ctx, p.SellerID); err == nil && s.notifier != nil {
		_ = s.notifier.NotifyUser(ctx, seller.UserID, models.NotifKindProduct,
			"Маҳсулот аз модератсия гузашт", "Товар прошёл модерацию",
			p.NameTJ+" → "+p.Status, p.NameRU+" → "+p.Status, "/seller/products")
	}
	return s.products.FindByID(ctx, p.ID)
}

func (s *ProductService) applyEditable(p *models.Product, in ProductInput) {
	if in.CategoryID != uuid.Nil {
		p.CategoryID = in.CategoryID
	}
	if in.SubcategoryID != nil {
		p.SubcategoryID = in.SubcategoryID
	}
	if in.BrandID != nil {
		p.BrandID = in.BrandID
	}
	if in.NameTJ != "" {
		p.NameTJ = in.NameTJ
	}
	if in.NameRU != "" {
		p.NameRU = in.NameRU
	}
	p.DescriptionTJ = firstNonEmpty(in.DescriptionTJ, p.DescriptionTJ)
	p.DescriptionRU = firstNonEmpty(in.DescriptionRU, p.DescriptionRU)
	if in.SKU != "" {
		p.SKU = in.SKU
	}
	if in.Unit != "" {
		p.Unit = in.Unit
	}
	if in.Currency != "" {
		p.Currency = in.Currency
	}
	if in.DiscountPercent > 0 {
		p.DiscountPercent = in.DiscountPercent
	}
	if in.StockQuantity != 0 {
		p.StockQuantity = in.StockQuantity
	}
	if in.MinimumStock != 0 {
		p.MinimumStock = in.MinimumStock
	}
	if in.IsAvailable != nil {
		p.IsAvailable = *in.IsAvailable
	}
}

// applyContacts resolves contact owner → concrete phone/whatsapp/telegram.
func (s *ProductService) applyContacts(ctx context.Context, p *models.Product, in ProductInput, force bool) {
	owner := in.ContactOwner
	if owner == "" && !force {
		return
	}
	if owner == "" {
		owner = p.ContactOwner
	}
	if owner != models.ContactOwnerAdmin && owner != models.ContactOwnerSeller {
		owner = models.ContactOwnerAdmin
	}
	p.ContactOwner = owner
	if owner == models.ContactOwnerAdmin {
		p.ContactPhone = firstNonEmpty(in.ContactPhone, s.defaults.Phone)
		p.ContactWhatsApp = firstNonEmpty(in.ContactWhatsApp, s.defaults.WhatsApp)
		p.ContactTelegram = firstNonEmpty(in.ContactTelegram, s.defaults.Telegram)
		return
	}
	// seller-owned: explicit values first, else seller profile.
	phone, wa, tg := in.ContactPhone, in.ContactWhatsApp, in.ContactTelegram
	if phone == "" || wa == "" || tg == "" {
		if seller, err := s.sellers.FindByID(ctx, p.SellerID); err == nil {
			phone = firstNonEmpty(phone, seller.Phone)
			wa = firstNonEmpty(wa, seller.WhatsApp, seller.Phone)
			tg = firstNonEmpty(tg, seller.Telegram)
		}
	}
	p.ContactPhone = phone
	p.ContactWhatsApp = wa
	p.ContactTelegram = tg
}

func (s *ProductService) DeleteBySeller(ctx context.Context, sellerID, productID uuid.UUID) error {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return ErrNotFound
	}
	if p.SellerID != sellerID {
		return ErrForbidden
	}
	return s.products.Delete(ctx, productID)
}
func (s *ProductService) DeleteByAdmin(ctx context.Context, id uuid.UUID) error {
	return s.products.Delete(ctx, id)
}

// ---- reads ----

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *ProductService) GetPublicBySlug(ctx context.Context, slug string) (*models.Product, error) {
	p, err := s.products.FindBySlug(ctx, slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if p.Status != models.ProductStatusApproved {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *ProductService) Related(ctx context.Context, productID uuid.UUID, limit int) ([]models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.products.Related(ctx, p, limit)
}

func (s *ProductService) PriceBounds(ctx context.Context, categorySlug string) (float64, float64, error) {
	return s.products.PriceBounds(ctx, categorySlug)
}

type ListInput struct {
	Search        string
	CategorySlug  string
	SubcatSlug    string
	SellerSlug    string
	SellerID      *uuid.UUID
	BrandIDs      []uuid.UUID
	Status        string
	OnlyAvailable bool
	IsFeatured    *bool
	LowStockOnly  bool
	MinPrice      *float64
	MaxPrice      *float64
	Sort          string
	Page          int
	Size          int
}

func (s *ProductService) ListPublic(ctx context.Context, in ListInput) ([]models.Product, int64, error) {
	in.Status = models.ProductStatusApproved
	in.OnlyAvailable = true
	return s.list(ctx, in)
}
func (s *ProductService) ListForSeller(ctx context.Context, sellerID uuid.UUID, in ListInput) ([]models.Product, int64, error) {
	in.SellerID = &sellerID
	return s.list(ctx, in)
}
func (s *ProductService) ListAdmin(ctx context.Context, in ListInput) ([]models.Product, int64, error) {
	return s.list(ctx, in)
}

func (s *ProductService) list(ctx context.Context, in ListInput) ([]models.Product, int64, error) {
	return s.products.List(ctx, repositories.ListProductsParams{
		Search:        in.Search,
		CategorySlug:  in.CategorySlug,
		SubcatSlug:    in.SubcatSlug,
		SellerSlug:    in.SellerSlug,
		SellerID:      in.SellerID,
		BrandIDs:      in.BrandIDs,
		Status:        in.Status,
		OnlyAvailable: in.OnlyAvailable,
		IsFeatured:    in.IsFeatured,
		LowStockOnly:  in.LowStockOnly,
		MinPrice:      in.MinPrice,
		MaxPrice:      in.MaxPrice,
		Sort:          in.Sort,
		Page:          in.Page,
		Size:          in.Size,
	})
}

func buildImages(productID uuid.UUID, in []ProductImageInput) []models.ProductImage {
	out := make([]models.ProductImage, 0, len(in))
	for i, im := range in {
		if strings.TrimSpace(im.URL) == "" {
			continue
		}
		out = append(out, models.ProductImage{
			ProductID: productID, URL: im.URL, Alt: im.Alt,
			SortOrder: i, IsCover: i == 0,
		})
	}
	return out
}

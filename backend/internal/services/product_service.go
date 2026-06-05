package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type ProductService struct {
	products   *repositories.ProductRepository
	sellers    *repositories.SellerRepository
	categories *repositories.CategoryRepository
	settings   *repositories.SettingRepository
	defaults   ContactDefaults
}

type ContactDefaults struct {
	Phone    string
	WhatsApp string
}

func NewProductService(
	products *repositories.ProductRepository,
	sellers *repositories.SellerRepository,
	categories *repositories.CategoryRepository,
	settings *repositories.SettingRepository,
	defaults ContactDefaults,
) *ProductService {
	return &ProductService{
		products:   products,
		sellers:    sellers,
		categories: categories,
		settings:   settings,
		defaults:   defaults,
	}
}

type ProductImageInput struct {
	URL string
	Alt string
}

type ProductInput struct {
	CategoryID    uuid.UUID
	TitleTJ       string
	TitleRU       string
	DescriptionTJ string
	DescriptionRU string
	Price         float64
	Currency      string
	Unit          string
	SKU           string
	StockQuantity int
	IsAvailable   *bool
	Images        []ProductImageInput
	// Admin-only on create/update:
	ContactType    string
	PhoneNumber    string
	WhatsAppNumber string
	IsFeatured     *bool
	Status         string
	RejectionNote  string
}

// CreateBySeller creates a product for the given seller. Status is forced to "pending".
func (s *ProductService) CreateBySeller(ctx context.Context, sellerID uuid.UUID, in ProductInput) (*models.Product, error) {
	if err := s.validateInput(in, true); err != nil {
		return nil, err
	}

	cat, err := s.categories.FindByID(ctx, in.CategoryID)
	if err != nil {
		return nil, ErrValidation
	}

	slug, err := s.uniqueSlug(ctx, in.TitleRU, in.TitleTJ)
	if err != nil {
		return nil, err
	}

	available := true
	if in.IsAvailable != nil {
		available = *in.IsAvailable
	}

	p := &models.Product{
		SellerID:      sellerID,
		CategoryID:    cat.ID,
		Slug:          slug,
		SKU:           in.SKU,
		TitleTJ:       in.TitleTJ,
		TitleRU:       in.TitleRU,
		DescriptionTJ: in.DescriptionTJ,
		DescriptionRU: in.DescriptionRU,
		Price:         in.Price,
		Currency:      defaultStr(in.Currency, "TJS"),
		Unit:          defaultStr(in.Unit, "pcs"),
		StockQuantity: in.StockQuantity,
		IsAvailable:   available,
		// Default contacts to admin/marketplace; admin can change during approval.
		ContactType:    models.ContactTypeAdmin,
		PhoneNumber:    s.defaults.Phone,
		WhatsAppNumber: s.defaults.WhatsApp,
		Status:         models.ProductStatusPending,
	}
	if err := s.products.Create(ctx, p); err != nil {
		return nil, err
	}

	if len(in.Images) > 0 {
		imgs := buildImages(p.ID, in.Images)
		if err := s.products.ReplaceImages(ctx, p.ID, imgs); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

func (s *ProductService) UpdateBySeller(ctx context.Context, sellerID, productID uuid.UUID, in ProductInput) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	if p.SellerID != sellerID {
		return nil, ErrForbidden
	}
	if err := s.validateInput(in, false); err != nil {
		return nil, err
	}
	applyProductFields(p, in)
	// Sellers cannot change moderation/contact fields.
	// Editing a product resets it to pending.
	p.Status = models.ProductStatusPending
	p.RejectionNote = ""

	if err := s.products.Update(ctx, p); err != nil {
		return nil, err
	}
	if in.Images != nil {
		imgs := buildImages(p.ID, in.Images)
		if err := s.products.ReplaceImages(ctx, p.ID, imgs); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
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

func (s *ProductService) AdminUpdate(ctx context.Context, productID uuid.UUID, in ProductInput) (*models.Product, error) {
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	if err := s.validateInput(in, false); err != nil {
		return nil, err
	}
	applyProductFields(p, in)

	// Admin-only fields:
	if in.ContactType != "" {
		if in.ContactType != models.ContactTypeAdmin && in.ContactType != models.ContactTypeSeller {
			return nil, ErrValidation
		}
		p.ContactType = in.ContactType
	}
	if in.PhoneNumber != "" {
		p.PhoneNumber = strings.TrimSpace(in.PhoneNumber)
	}
	if in.WhatsAppNumber != "" {
		p.WhatsAppNumber = strings.TrimSpace(in.WhatsAppNumber)
	}
	if in.IsFeatured != nil {
		p.IsFeatured = *in.IsFeatured
	}
	if in.Status != "" {
		p.Status = in.Status
	}
	if in.RejectionNote != "" {
		p.RejectionNote = in.RejectionNote
	}
	if err := s.products.Update(ctx, p); err != nil {
		return nil, err
	}
	if in.Images != nil {
		imgs := buildImages(p.ID, in.Images)
		if err := s.products.ReplaceImages(ctx, p.ID, imgs); err != nil {
			return nil, err
		}
	}
	return s.products.FindByID(ctx, p.ID)
}

// ModerationDecision is what admin sets on approve/reject.
type ModerationDecision struct {
	Status         string // approved|rejected
	ContactType    string // admin|seller
	PhoneNumber    string
	WhatsAppNumber string
	RejectionNote  string
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
		if d.ContactType != "" {
			if d.ContactType != models.ContactTypeAdmin && d.ContactType != models.ContactTypeSeller {
				return nil, ErrValidation
			}
			p.ContactType = d.ContactType
		}
		if d.ContactType == models.ContactTypeAdmin {
			p.PhoneNumber = firstNonEmpty(d.PhoneNumber, s.defaults.Phone)
			p.WhatsAppNumber = firstNonEmpty(d.WhatsAppNumber, s.defaults.WhatsApp)
		} else {
			// seller-routed: prefer explicit values, fall back to seller's profile.
			p.PhoneNumber = strings.TrimSpace(d.PhoneNumber)
			p.WhatsAppNumber = strings.TrimSpace(d.WhatsAppNumber)
			if p.PhoneNumber == "" || p.WhatsAppNumber == "" {
				if seller, err := s.sellers.FindByID(ctx, p.SellerID); err == nil {
					if p.PhoneNumber == "" {
						p.PhoneNumber = seller.Phone
					}
					if p.WhatsAppNumber == "" {
						p.WhatsAppNumber = seller.WhatsApp
					}
				}
			}
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
	if err := s.products.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.products.FindByID(ctx, p.ID)
}

func (s *ProductService) AdminDelete(ctx context.Context, id uuid.UUID) error {
	return s.products.Delete(ctx, id)
}

// ----- read paths -----

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

func (s *ProductService) Related(ctx context.Context, productID uuid.UUID, limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 24 {
		limit = 8
	}
	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.products.Related(ctx, p, limit)
}

type ListProductsInput struct {
	Search        string
	CategorySlug  string
	SellerSlug    string
	SellerID      *uuid.UUID
	OnlyAvailable bool
	IsFeatured    *bool
	MinPrice      *float64
	MaxPrice      *float64
	Sort          string
	Status        string // filtering for admin/seller views
	Page          int
	PageSize      int
}

func (s *ProductService) ListPublic(ctx context.Context, in ListProductsInput) ([]models.Product, int64, error) {
	in.Status = models.ProductStatusApproved
	in.OnlyAvailable = true
	return s.list(ctx, in)
}

func (s *ProductService) ListForSeller(ctx context.Context, sellerID uuid.UUID, in ListProductsInput) ([]models.Product, int64, error) {
	in.SellerID = &sellerID
	return s.list(ctx, in)
}

func (s *ProductService) ListAdmin(ctx context.Context, in ListProductsInput) ([]models.Product, int64, error) {
	return s.list(ctx, in)
}

func (s *ProductService) list(ctx context.Context, in ListProductsInput) ([]models.Product, int64, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 20
	}
	return s.products.List(ctx, repositories.ListParams{
		Search:        in.Search,
		CategorySlug:  in.CategorySlug,
		SellerSlug:    in.SellerSlug,
		SellerID:      in.SellerID,
		OnlyAvailable: in.OnlyAvailable,
		IsFeatured:    in.IsFeatured,
		MinPrice:      in.MinPrice,
		MaxPrice:      in.MaxPrice,
		Sort:          in.Sort,
		Status:        in.Status,
		Page:          in.Page,
		PageSize:      in.PageSize,
	})
}

// ----- internals -----

func (s *ProductService) validateInput(in ProductInput, requireCategory bool) error {
	if in.TitleTJ == "" || in.TitleRU == "" {
		return ErrValidation
	}
	if in.Price < 0 {
		return ErrValidation
	}
	if requireCategory && in.CategoryID == uuid.Nil {
		return ErrValidation
	}
	return nil
}

func (s *ProductService) uniqueSlug(ctx context.Context, primary, fallback string) (string, error) {
	base := Slugify(primary)
	if base == "item" {
		base = Slugify(fallback)
	}
	candidate := base
	for i := 1; i < 1000; i++ {
		exists, err := s.products.ExistsBySlug(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = base + "-" + intToStr(i)
	}
	return base + "-" + intToStr(9999), nil
}

func applyProductFields(p *models.Product, in ProductInput) {
	if in.CategoryID != uuid.Nil {
		p.CategoryID = in.CategoryID
	}
	if in.TitleTJ != "" {
		p.TitleTJ = in.TitleTJ
	}
	if in.TitleRU != "" {
		p.TitleRU = in.TitleRU
	}
	if in.DescriptionTJ != "" {
		p.DescriptionTJ = in.DescriptionTJ
	}
	if in.DescriptionRU != "" {
		p.DescriptionRU = in.DescriptionRU
	}
	if in.Price > 0 {
		p.Price = in.Price
	}
	if in.Currency != "" {
		p.Currency = in.Currency
	}
	if in.Unit != "" {
		p.Unit = in.Unit
	}
	if in.SKU != "" {
		p.SKU = in.SKU
	}
	if in.StockQuantity != 0 {
		p.StockQuantity = in.StockQuantity
	}
	if in.IsAvailable != nil {
		p.IsAvailable = *in.IsAvailable
	}
}

func buildImages(productID uuid.UUID, in []ProductImageInput) []models.ProductImage {
	out := make([]models.ProductImage, 0, len(in))
	for i, im := range in {
		out = append(out, models.ProductImage{
			ProductID: productID,
			URL:       im.URL,
			Alt:       im.Alt,
			SortOrder: i,
			IsCover:   i == 0,
		})
	}
	return out
}

func defaultStr(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

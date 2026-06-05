package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

// CatalogService groups category / subcategory / brand admin + read logic.
type CatalogService struct {
	cats   *repositories.CategoryRepository
	subs   *repositories.SubcategoryRepository
	brands *repositories.BrandRepository
}

func NewCatalogService(c *repositories.CategoryRepository, s *repositories.SubcategoryRepository, b *repositories.BrandRepository) *CatalogService {
	return &CatalogService{cats: c, subs: s, brands: b}
}

// ---------- Categories ----------

type CategoryInput struct {
	NameTJ, NameRU               string
	DescriptionTJ, DescriptionRU string
	Slug                         string
	IconURL, BannerURL           string
	SEOTitleTJ, SEOTitleRU       string
	SEODescTJ, SEODescRU         string
	SortOrder                    int
	Active                       *bool
}

func (s *CatalogService) ListCategories(ctx context.Context, onlyActive bool) ([]models.Category, error) {
	return s.cats.ListAll(ctx, onlyActive)
}

func (s *CatalogService) PopularCategories(ctx context.Context, limit int) ([]models.Category, error) {
	return s.cats.Popular(ctx, limit)
}

func (s *CatalogService) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	c, err := s.cats.FindByID(ctx, id)
	return wrapNotFound(c, err)
}

func (s *CatalogService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	c, err := s.cats.FindBySlug(ctx, slug)
	return wrapNotFound(c, err)
}

func (s *CatalogService) CreateCategory(ctx context.Context, in CategoryInput) (*models.Category, error) {
	if in.NameTJ == "" || in.NameRU == "" {
		return nil, ErrValidation
	}
	base := firstNonEmpty(in.Slug, in.NameRU, in.NameTJ)
	slug, err := uniqueSlug(ctx, base, s.cats.ExistsBySlug)
	if err != nil {
		return nil, err
	}
	c := &models.Category{
		Slug: slug, NameTJ: in.NameTJ, NameRU: in.NameRU,
		DescriptionTJ: in.DescriptionTJ, DescriptionRU: in.DescriptionRU,
		IconURL: in.IconURL, BannerURL: in.BannerURL,
		SEOTitleTJ: in.SEOTitleTJ, SEOTitleRU: in.SEOTitleRU,
		SEODescTJ: in.SEODescTJ, SEODescRU: in.SEODescRU,
		SortOrder: in.SortOrder, Active: boolOr(in.Active, true),
	}
	if err := s.cats.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CatalogService) UpdateCategory(ctx context.Context, id uuid.UUID, in CategoryInput) (*models.Category, error) {
	c, err := s.cats.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.NameTJ != "" {
		c.NameTJ = in.NameTJ
	}
	if in.NameRU != "" {
		c.NameRU = in.NameRU
	}
	c.DescriptionTJ = in.DescriptionTJ
	c.DescriptionRU = in.DescriptionRU
	if in.IconURL != "" {
		c.IconURL = in.IconURL
	}
	if in.BannerURL != "" {
		c.BannerURL = in.BannerURL
	}
	c.SEOTitleTJ = in.SEOTitleTJ
	c.SEOTitleRU = in.SEOTitleRU
	c.SEODescTJ = in.SEODescTJ
	c.SEODescRU = in.SEODescRU
	c.SortOrder = in.SortOrder
	if in.Active != nil {
		c.Active = *in.Active
	}
	if in.Slug != "" && Slugify(in.Slug) != c.Slug {
		ns := Slugify(in.Slug)
		if ex, _ := s.cats.ExistsBySlug(ctx, ns); ex {
			return nil, ErrConflict
		}
		c.Slug = ns
	}
	if err := s.cats.Save(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CatalogService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if has, _ := s.cats.HasProducts(ctx, id); has {
		return ErrConflict
	}
	if has, _ := s.cats.HasChildren(ctx, id); has {
		return ErrConflict
	}
	return s.cats.Delete(ctx, id)
}

// ---------- Subcategories ----------

type SubcategoryInput struct {
	CategoryID     uuid.UUID
	NameTJ, NameRU string
	Slug           string
	IconURL        string
	SortOrder      int
	Active         *bool
}

func (s *CatalogService) ListSubcategories(ctx context.Context, categoryID uuid.UUID, onlyActive bool) ([]models.Subcategory, error) {
	return s.subs.ListByCategory(ctx, categoryID, onlyActive)
}

func (s *CatalogService) CreateSubcategory(ctx context.Context, in SubcategoryInput) (*models.Subcategory, error) {
	if in.NameTJ == "" || in.NameRU == "" || in.CategoryID == uuid.Nil {
		return nil, ErrValidation
	}
	if _, err := s.cats.FindByID(ctx, in.CategoryID); err != nil {
		return nil, ErrValidation
	}
	base := firstNonEmpty(in.Slug, in.NameRU, in.NameTJ)
	slug, err := uniqueSlug(ctx, base, s.subs.ExistsBySlug)
	if err != nil {
		return nil, err
	}
	sc := &models.Subcategory{
		CategoryID: in.CategoryID, Slug: slug,
		NameTJ: in.NameTJ, NameRU: in.NameRU, IconURL: in.IconURL,
		SortOrder: in.SortOrder, Active: boolOr(in.Active, true),
	}
	if err := s.subs.Create(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *CatalogService) UpdateSubcategory(ctx context.Context, id uuid.UUID, in SubcategoryInput) (*models.Subcategory, error) {
	sc, err := s.subs.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.NameTJ != "" {
		sc.NameTJ = in.NameTJ
	}
	if in.NameRU != "" {
		sc.NameRU = in.NameRU
	}
	if in.IconURL != "" {
		sc.IconURL = in.IconURL
	}
	if in.CategoryID != uuid.Nil {
		sc.CategoryID = in.CategoryID
	}
	sc.SortOrder = in.SortOrder
	if in.Active != nil {
		sc.Active = *in.Active
	}
	if in.Slug != "" && Slugify(in.Slug) != sc.Slug {
		ns := Slugify(in.Slug)
		if ex, _ := s.subs.ExistsBySlug(ctx, ns); ex {
			return nil, ErrConflict
		}
		sc.Slug = ns
	}
	if err := s.subs.Save(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *CatalogService) DeleteSubcategory(ctx context.Context, id uuid.UUID) error {
	if has, _ := s.subs.HasProducts(ctx, id); has {
		return ErrConflict
	}
	return s.subs.Delete(ctx, id)
}

// ---------- Brands ----------

type BrandInput struct {
	Name      string
	Slug      string
	LogoURL   string
	SortOrder int
	Active    *bool
}

func (s *CatalogService) ListBrands(ctx context.Context, onlyActive bool) ([]models.Brand, error) {
	return s.brands.ListAll(ctx, onlyActive)
}

func (s *CatalogService) GetBrand(ctx context.Context, id uuid.UUID) (*models.Brand, error) {
	b, err := s.brands.FindByID(ctx, id)
	return wrapNotFoundBrand(b, err)
}

func (s *CatalogService) CreateBrand(ctx context.Context, in BrandInput) (*models.Brand, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, ErrValidation
	}
	slug, err := uniqueSlug(ctx, firstNonEmpty(in.Slug, in.Name), s.brands.ExistsBySlug)
	if err != nil {
		return nil, err
	}
	b := &models.Brand{Name: in.Name, Slug: slug, LogoURL: in.LogoURL, SortOrder: in.SortOrder, Active: boolOr(in.Active, true)}
	if err := s.brands.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *CatalogService) UpdateBrand(ctx context.Context, id uuid.UUID, in BrandInput) (*models.Brand, error) {
	b, err := s.brands.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.Name != "" {
		b.Name = in.Name
	}
	if in.LogoURL != "" {
		b.LogoURL = in.LogoURL
	}
	b.SortOrder = in.SortOrder
	if in.Active != nil {
		b.Active = *in.Active
	}
	if in.Slug != "" && Slugify(in.Slug) != b.Slug {
		ns := Slugify(in.Slug)
		if ex, _ := s.brands.ExistsBySlug(ctx, ns); ex {
			return nil, ErrConflict
		}
		b.Slug = ns
	}
	if err := s.brands.Save(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *CatalogService) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	if has, _ := s.brands.HasProducts(ctx, id); has {
		return ErrConflict
	}
	return s.brands.Delete(ctx, id)
}

// helpers
func boolOr(p *bool, def bool) bool {
	if p != nil {
		return *p
	}
	return def
}
func wrapNotFound(c *models.Category, err error) (*models.Category, error) {
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}
func wrapNotFoundBrand(b *models.Brand, err error) (*models.Brand, error) {
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

type CategoryInput struct {
	TitleTJ   string
	TitleRU   string
	Slug      string
	IconURL   string
	SortOrder int
	IsActive  *bool
	ParentID  *string
}

func (s *CategoryService) ListAll(ctx context.Context, onlyActive bool) ([]models.Category, error) {
	return s.repo.ListAll(ctx, onlyActive)
}

func (s *CategoryService) Popular(ctx context.Context, limit int) ([]models.Category, error) {
	if limit <= 0 || limit > 50 {
		limit = 8
	}
	return s.repo.Popular(ctx, limit)
}

func (s *CategoryService) Get(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	c, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Create(ctx context.Context, in CategoryInput) (*models.Category, error) {
	if in.TitleTJ == "" || in.TitleRU == "" {
		return nil, ErrValidation
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		slug = Slugify(in.TitleRU)
		if slug == "item" {
			slug = Slugify(in.TitleTJ)
		}
	} else {
		slug = Slugify(slug)
	}
	finalSlug := slug
	for i := 1; i < 1000; i++ {
		exists, err := s.repo.ExistsBySlug(ctx, finalSlug)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		finalSlug = slug + "-" + intToStr(i)
	}
	active := true
	if in.IsActive != nil {
		active = *in.IsActive
	}
	c := &models.Category{
		Slug:      finalSlug,
		TitleTJ:   in.TitleTJ,
		TitleRU:   in.TitleRU,
		IconURL:   in.IconURL,
		SortOrder: in.SortOrder,
		IsActive:  active,
		ParentID:  in.ParentID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, in CategoryInput) (*models.Category, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if in.TitleTJ != "" {
		c.TitleTJ = in.TitleTJ
	}
	if in.TitleRU != "" {
		c.TitleRU = in.TitleRU
	}
	if in.IconURL != "" {
		c.IconURL = in.IconURL
	}
	if in.SortOrder != 0 {
		c.SortOrder = in.SortOrder
	}
	if in.IsActive != nil {
		c.IsActive = *in.IsActive
	}
	if in.ParentID != nil {
		c.ParentID = in.ParentID
	}
	if in.Slug != "" {
		newSlug := Slugify(in.Slug)
		if newSlug != c.Slug {
			exists, err := s.repo.ExistsBySlug(ctx, newSlug)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrConflict
			}
			c.Slug = newSlug
		}
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	has, err := s.repo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if has {
		return ErrConflict
	}
	return s.repo.Delete(ctx, id)
}

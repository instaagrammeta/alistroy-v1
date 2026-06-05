package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type BannerService struct {
	repo *repositories.BannerRepository
}

func NewBannerService(r *repositories.BannerRepository) *BannerService {
	return &BannerService{repo: r}
}

type BannerInput struct {
	Position      string
	TitleTJ       string
	TitleRU       string
	DescriptionTJ string
	DescriptionRU string
	DesktopURL    string
	TabletURL     string
	MobileURL     string
	LinkURL       string
	SortOrder     int
	Active        *bool
}

func (s *BannerService) PublicGrouped(ctx context.Context) (map[string][]models.Banner, error) {
	items, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string][]models.Banner)
	for _, b := range items {
		out[b.Position] = append(out[b.Position], b)
	}
	return out, nil
}

func (s *BannerService) AdminList(ctx context.Context, position string) ([]models.Banner, error) {
	return s.repo.ListAll(ctx, position)
}

func (s *BannerService) Get(ctx context.Context, id uuid.UUID) (*models.Banner, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (s *BannerService) Create(ctx context.Context, in BannerInput) (*models.Banner, error) {
	if in.Position == "" {
		return nil, ErrValidation
	}
	b := &models.Banner{
		Position: in.Position, TitleTJ: in.TitleTJ, TitleRU: in.TitleRU,
		DescriptionTJ: in.DescriptionTJ, DescriptionRU: in.DescriptionRU,
		DesktopURL: in.DesktopURL, TabletURL: in.TabletURL, MobileURL: in.MobileURL,
		LinkURL: in.LinkURL, SortOrder: in.SortOrder, Active: boolOr(in.Active, true),
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BannerService) Update(ctx context.Context, id uuid.UUID, in BannerInput) (*models.Banner, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	if in.Position != "" {
		b.Position = in.Position
	}
	b.TitleTJ = in.TitleTJ
	b.TitleRU = in.TitleRU
	b.DescriptionTJ = in.DescriptionTJ
	b.DescriptionRU = in.DescriptionRU
	if in.DesktopURL != "" {
		b.DesktopURL = in.DesktopURL
	}
	if in.TabletURL != "" {
		b.TabletURL = in.TabletURL
	}
	if in.MobileURL != "" {
		b.MobileURL = in.MobileURL
	}
	b.LinkURL = in.LinkURL
	b.SortOrder = in.SortOrder
	if in.Active != nil {
		b.Active = *in.Active
	}
	if err := s.repo.Save(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BannerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

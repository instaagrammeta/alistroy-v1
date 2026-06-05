package services

import (
	"context"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type SettingService struct {
	repo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetAll(ctx context.Context) (map[string]string, error) {
	return s.repo.GetAll(ctx)
}

func (s *SettingService) Get(ctx context.Context, key string) (string, error) {
	v, err := s.repo.Get(ctx, key)
	if err != nil {
		if repositories.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return v, nil
}

func (s *SettingService) Update(ctx context.Context, kv map[string]string) error {
	return s.repo.SetMany(ctx, kv)
}

func (s *SettingService) Set(ctx context.Context, key, value string) error {
	return s.repo.Set(ctx, key, value)
}

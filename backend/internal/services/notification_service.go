package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/cache"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type NotificationService struct {
	repo  *repositories.NotificationRepository
	redis *cache.Client
}

func NewNotificationService(r *repositories.NotificationRepository, rc *cache.Client) *NotificationService {
	return &NotificationService{repo: r, redis: rc}
}

// NotifyUser persists a notification and publishes a realtime ping over Redis.
func (s *NotificationService) NotifyUser(ctx context.Context, userID uuid.UUID, kind, titleTJ, titleRU, bodyTJ, bodyRU, link string) error {
	n := &models.Notification{
		UserID:  userID,
		Kind:    kind,
		TitleTJ: titleTJ, TitleRU: titleRU,
		BodyTJ: bodyTJ, BodyRU: bodyRU,
		LinkURL: link,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}
	if s.redis != nil {
		_ = s.redis.Publish(ctx, "notify:"+userID.String(), n)
	}
	return nil
}

func (s *NotificationService) List(ctx context.Context, userID uuid.UUID, page, size int) ([]models.Notification, int64, error) {
	return s.repo.ListByUser(ctx, userID, page, size)
}
func (s *NotificationService) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.UnreadCount(ctx, userID)
}
func (s *NotificationService) MarkRead(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.MarkRead(ctx, userID, id)
}
func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, userID)
}

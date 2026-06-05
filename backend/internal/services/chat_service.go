package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/cache"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
)

type ChatService struct {
	chat      *repositories.ChatRepository
	customers *repositories.CustomerRepository
	redis     *cache.Client
}

func NewChatService(c *repositories.ChatRepository, cust *repositories.CustomerRepository, rc *cache.Client) *ChatService {
	return &ChatService{chat: c, customers: cust, redis: rc}
}

type AttachmentInput struct {
	URL       string
	MimeType  string
	SizeBytes int64
}

// CustomerRoom returns (creating if needed) the room for a customer user.
func (s *ChatService) CustomerRoom(ctx context.Context, userID uuid.UUID) (*models.ChatRoom, error) {
	cust, err := s.customers.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.chat.RoomByCustomer(ctx, cust.ID)
}

func (s *ChatService) RoomByID(ctx context.Context, id uuid.UUID) (*models.ChatRoom, error) {
	r, err := s.chat.RoomByID(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return r, nil
}

func (s *ChatService) ListRooms(ctx context.Context, page, size int) ([]models.ChatRoom, int64, error) {
	return s.chat.ListRooms(ctx, page, size)
}

func (s *ChatService) Messages(ctx context.Context, roomID uuid.UUID, page, size int) ([]models.ChatMessage, int64, error) {
	return s.chat.Messages(ctx, roomID, page, size)
}

// SendAsCustomer posts a message authored by the customer.
func (s *ChatService) SendAsCustomer(ctx context.Context, userID uuid.UUID, body string, atts []AttachmentInput) (*models.ChatMessage, error) {
	cust, err := s.customers.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	room, err := s.chat.RoomByCustomer(ctx, cust.ID)
	if err != nil {
		return nil, err
	}
	return s.send(ctx, room, userID, models.RoleCustomer, body, atts)
}

// SendAsAdmin posts a message into a room authored by admin.
func (s *ChatService) SendAsAdmin(ctx context.Context, adminUserID, roomID uuid.UUID, body string, atts []AttachmentInput) (*models.ChatMessage, error) {
	room, err := s.chat.RoomByID(ctx, roomID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.send(ctx, room, adminUserID, models.RoleAdmin, body, atts)
}

func (s *ChatService) send(ctx context.Context, room *models.ChatRoom, senderID uuid.UUID, role, body string, atts []AttachmentInput) (*models.ChatMessage, error) {
	if body == "" && len(atts) == 0 {
		return nil, ErrValidation
	}
	msg := &models.ChatMessage{RoomID: room.ID, SenderID: senderID, SenderRole: role, Body: body}
	var attachments []models.ChatAttachment
	for _, a := range atts {
		if a.URL == "" {
			continue
		}
		attachments = append(attachments, models.ChatAttachment{URL: a.URL, MimeType: a.MimeType, SizeBytes: a.SizeBytes})
	}
	if err := s.chat.CreateMessage(ctx, msg, attachments); err != nil {
		return nil, err
	}
	msg.Attachments = attachments
	if s.redis != nil {
		_ = s.redis.Publish(ctx, "chat:"+room.ID.String(), msg)
	}
	return msg, nil
}

func (s *ChatService) MarkRead(ctx context.Context, roomID uuid.UUID, role string) error {
	return s.chat.MarkRead(ctx, roomID, role)
}

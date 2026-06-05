package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type ChatRepository struct{ db *gorm.DB }

func NewChatRepository(db *gorm.DB) *ChatRepository { return &ChatRepository{db: db} }

// RoomByCustomer returns the customer's room, creating it lazily.
func (r *ChatRepository) RoomByCustomer(ctx context.Context, customerID uuid.UUID) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).First(&room).Error
	if err == nil {
		return &room, nil
	}
	if !IsNotFound(err) {
		return nil, err
	}
	room = models.ChatRoom{CustomerID: customerID}
	if err := r.db.WithContext(ctx).Create(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *ChatRepository) RoomByID(ctx context.Context, id uuid.UUID) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.WithContext(ctx).Preload("Customer").First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *ChatRepository) ListRooms(ctx context.Context, page, size int) ([]models.ChatRoom, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.ChatRoom{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&page, &size)
	var rooms []models.ChatRoom
	err := q.Preload("Customer").
		Order("last_message_at DESC NULLS LAST, created_at DESC").
		Limit(size).Offset(off).Find(&rooms).Error
	return rooms, total, err
}

func (r *ChatRepository) CreateMessage(ctx context.Context, m *models.ChatMessage, atts []models.ChatAttachment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		if len(atts) > 0 {
			for i := range atts {
				atts[i].MessageID = m.ID
			}
			if err := tx.Create(&atts).Error; err != nil {
				return err
			}
		}
		now := time.Now().UTC()
		updates := map[string]any{"last_message_at": now, "updated_at": now}
		if m.SenderRole == models.RoleCustomer {
			updates["unread_admin"] = gorm.Expr("unread_admin + 1")
		} else {
			updates["unread_customer"] = gorm.Expr("unread_customer + 1")
		}
		return tx.Model(&models.ChatRoom{}).Where("id = ?", m.RoomID).Updates(updates).Error
	})
}

func (r *ChatRepository) Messages(ctx context.Context, roomID uuid.UUID, page, size int) ([]models.ChatMessage, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.ChatMessage{}).Where("room_id = ?", roomID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := applyPaging(&page, &size)
	var msgs []models.ChatMessage
	err := q.Preload("Attachments").
		Order("created_at DESC").Limit(size).Offset(off).Find(&msgs).Error
	return msgs, total, err
}

// MarkRead clears the unread counter for the given side and stamps messages.
func (r *ChatRepository) MarkRead(ctx context.Context, roomID uuid.UUID, readerRole string) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// messages authored by the *other* party become read
		other := models.RoleAdmin
		col := "unread_customer"
		if readerRole == models.RoleAdmin {
			other = models.RoleCustomer
			col = "unread_admin"
		}
		if err := tx.Model(&models.ChatMessage{}).
			Where("room_id = ? AND sender_role = ? AND read_at IS NULL", roomID, other).
			Update("read_at", now).Error; err != nil {
			return err
		}
		return tx.Model(&models.ChatRoom{}).Where("id = ?", roomID).Update(col, 0).Error
	})
}

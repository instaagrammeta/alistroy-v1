package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

type SettingRepository struct{ db *gorm.DB }

func NewSettingRepository(db *gorm.DB) *SettingRepository { return &SettingRepository{db: db} }

func (r *SettingRepository) GetAll(ctx context.Context) (map[string]string, error) {
	var rows []models.Setting
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, s := range rows {
		out[s.Key] = s.Value
	}
	return out, nil
}

func (r *SettingRepository) Get(ctx context.Context, key string) (string, error) {
	var s models.Setting
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error; err != nil {
		return "", err
	}
	return s.Value, nil
}

func (r *SettingRepository) Set(ctx context.Context, key, value string) error {
	s := models.Setting{Key: key, Value: value}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&s).Error
}

func (r *SettingRepository) SetMany(ctx context.Context, kv map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for k, v := range kv {
			s := models.Setting{Key: k, Value: v}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
			}).Create(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

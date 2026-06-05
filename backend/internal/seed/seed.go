package seed

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/logger"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

// Run seeds the bootstrap admin, default settings and root categories.
// Every step is idempotent and safe on each boot.
func Run(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	if err := seedAdmin(ctx, db, cfg); err != nil {
		return err
	}
	if err := seedSettings(ctx, db, cfg); err != nil {
		return err
	}
	return seedCategories(ctx, db)
}

func seedAdmin(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	if cfg.Admin.Email == "" || cfg.Admin.Password == "" {
		logger.Warn("seed: admin email/password not set; skipping")
		return nil
	}
	var existing models.User
	err := db.WithContext(ctx).Where("email = ?", cfg.Admin.Email).First(&existing).Error
	if err == nil {
		if existing.Role != models.RoleAdmin || existing.Status != models.UserStatusActive {
			existing.Role = models.RoleAdmin
			existing.Status = models.UserStatusActive
			db.WithContext(ctx).Save(&existing)
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := &models.User{
		Email:        cfg.Admin.Email,
		Phone:        cfg.Admin.Phone,
		Login:        "admin",
		PasswordHash: string(hash),
		Name:         cfg.Admin.Name,
		Role:         models.RoleAdmin,
		Status:       models.UserStatusActive,
		Locale:       "tg",
	}
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	logger.Info("seed: admin account created", "email", cfg.Admin.Email)
	return nil
}

func seedSettings(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	defaults := map[string]string{
		models.SettingSiteNameTJ:       "AliStroy",
		models.SettingSiteNameRU:       "AliStroy",
		models.SettingTaglineTJ:        "Бозори маводи сохтмонии Тоҷикистон",
		models.SettingTaglineRU:        "Маркетплейс строительных материалов Таджикистана",
		models.SettingSEODescriptionTJ: "AliStroy — бозори бузурги маводи сохтмонии Тоҷикистон.",
		models.SettingSEODescriptionRU: "AliStroy — крупнейший маркетплейс строительных материалов Таджикистана.",
		models.SettingHeroTitleTJ:      "Маводи сохтмонӣ барои тамоми Тоҷикистон",
		models.SettingHeroTitleRU:      "Строительные материалы для всего Таджикистана",
		models.SettingHeroSubtitleTJ:   "Ҳазорон молҳо аз фурӯшандагони боэътимод.",
		models.SettingHeroSubtitleRU:   "Тысячи товаров от проверенных продавцов.",
		models.SettingFooterEmail:      "info@alistroy.tj",
		models.SettingFooterAddress:    "Душанбе, Тоҷикистон",
	}
	if cfg.Market.Phone != "" {
		defaults[models.SettingMarketplacePhone] = cfg.Market.Phone
	}
	if cfg.Market.WhatsApp != "" {
		defaults[models.SettingMarketplaceWA] = cfg.Market.WhatsApp
	}
	if cfg.Market.Telegram != "" {
		defaults[models.SettingMarketplaceTG] = cfg.Market.Telegram
	}
	if cfg.Market.TelegramUsername != "" {
		defaults[models.SettingMarketplaceTGUN] = cfg.Market.TelegramUsername
	}

	for k, v := range defaults {
		var existing models.Setting
		err := db.WithContext(ctx).Where("key = ?", k).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&models.Setting{Key: k, Value: v}).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

type seedCat struct {
	Slug, TJ, RU string
	Order        int
}

func seedCategories(ctx context.Context, db *gorm.DB) error {
	cats := []seedCat{
		{"cement", "Семент", "Цемент", 10},
		{"bricks", "Хишт", "Кирпич", 20},
		{"blocks", "Блокҳо", "Блоки", 30},
		{"sand", "Рег", "Песок", 40},
		{"gravel", "Шағал", "Щебень", 50},
		{"paint", "Ранг", "Краска", 60},
		{"plumbing", "Сантехника", "Сантехника", 70},
		{"electrical", "Электрика", "Электрика", 80},
		{"roofing", "Бомпӯшӣ", "Кровля", 90},
		{"doors", "Дарҳо", "Двери", 100},
		{"windows", "Тирезаҳо", "Окна", 110},
		{"tools", "Асбобҳо", "Инструменты", 120},
		{"hardware", "Маҳкамкунакҳо", "Метизы", 130},
		{"insulation", "Изолятсия", "Изоляция", 140},
		{"finishing", "Маводи ороишӣ", "Отделочные материалы", 150},
	}
	for _, c := range cats {
		var existing models.Category
		err := db.WithContext(ctx).Where("slug = ?", c.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&models.Category{
				Slug: c.Slug, NameTJ: c.TJ, NameRU: c.RU, SortOrder: c.Order, Active: true,
			}).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

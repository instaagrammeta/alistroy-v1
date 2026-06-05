package seed

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

// Run seeds the bootstrap admin account, default settings, and root categories.
// Safe to run on every boot — all operations are idempotent.
func Run(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	if err := seedAdmin(ctx, db, cfg); err != nil {
		return err
	}
	if err := seedSettings(ctx, db, cfg); err != nil {
		return err
	}
	if err := seedCategories(ctx, db); err != nil {
		return err
	}
	return nil
}

func seedAdmin(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	if cfg.Admin.Email == "" || cfg.Admin.Password == "" {
		log.Println("[seed] admin email/password not set, skipping admin bootstrap")
		return nil
	}
	var existing models.User
	err := db.WithContext(ctx).Where("email = ?", cfg.Admin.Email).First(&existing).Error
	if err == nil {
		// already present; ensure role is admin and active
		if existing.Role != models.RoleAdmin || !existing.IsActive {
			existing.Role = models.RoleAdmin
			existing.IsActive = true
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
		PasswordHash: string(hash),
		Name:         cfg.Admin.Name,
		Role:         models.RoleAdmin,
		Locale:       "tg",
		IsActive:     true,
	}
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		return err
	}
	log.Printf("[seed] admin account created: %s", cfg.Admin.Email)
	return nil
}

func seedSettings(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	defaults := map[string]string{
		models.SettingSiteNameTJ:       "AliStroy",
		models.SettingSiteNameRU:       "AliStroy",
		models.SettingTaglineTJ:        "Бозори маводи сохтмонии Тоҷикистон",
		models.SettingTaglineRU:        "Маркетплейс строительных материалов Таджикистана",
		models.SettingSEODescriptionTJ: "AliStroy — бозори бузурги маводи сохтмонии Тоҷикистон. Семент, хишт, ранг, асбобҳо ва дигар маводҳо аз фурӯшандагони боэътимод.",
		models.SettingSEODescriptionRU: "AliStroy — крупнейший маркетплейс строительных материалов Таджикистана. Цемент, кирпич, краска, инструменты и многое другое от проверенных продавцов.",
		models.SettingHeroTitleTJ:      "Маводи сохтмонӣ барои тамоми Тоҷикистон",
		models.SettingHeroTitleRU:      "Строительные материалы для всего Таджикистана",
		models.SettingHeroSubtitleTJ:   "Зиёда аз ҳазор молҳо аз фурӯшандагони боэътимод. Ёбед, муқоиса кунед ва зуд тамос гиред.",
		models.SettingHeroSubtitleRU:   "Тысячи товаров от проверенных продавцов. Найдите, сравните и свяжитесь напрямую.",
		models.SettingFooterEmail:      "info@alistroy.tj",
		models.SettingFooterAddress:    "Душанбе, Тоҷикистон",
	}
	if cfg.Market.Phone != "" {
		defaults[models.SettingMarketplacePhone] = cfg.Market.Phone
	}
	if cfg.Market.WhatsApp != "" {
		defaults[models.SettingMarketplaceWA] = cfg.Market.WhatsApp
	}

	for k, v := range defaults {
		var existing models.Setting
		err := db.WithContext(ctx).Where("key = ?", k).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&models.Setting{Key: k, Value: v}).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		// Existing settings are NOT overwritten — admin may have customized them.
	}
	return nil
}

type seedCategory struct {
	Slug    string
	TitleTJ string
	TitleRU string
	Order   int
}

func seedCategories(ctx context.Context, db *gorm.DB) error {
	cats := []seedCategory{
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
		{"finishing-materials", "Маводи ороишӣ", "Отделочные материалы", 150},
	}
	for _, c := range cats {
		var existing models.Category
		err := db.WithContext(ctx).Where("slug = ?", c.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&models.Category{
				Slug:      c.Slug,
				TitleTJ:   c.TitleTJ,
				TitleRU:   c.TitleRU,
				SortOrder: c.Order,
				IsActive:  true,
			}).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

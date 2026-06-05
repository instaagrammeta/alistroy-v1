package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/cache"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/database"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/handlers"
	appjwt "github.com/instaagrammeta/alistroy-v1/backend/internal/jwt"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/logger"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/oauth"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/seed"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/server"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/ws"
	"github.com/instaagrammeta/alistroy-v1/backend/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config", "err", err)
		os.Exit(1)
	}
	logger.SetDefault(logger.New(os.Stdout, cfg.LogLevel))

	// --- Postgres ---
	db, err := database.Connect(cfg.Postgres, cfg.Env)
	if err != nil {
		logger.Error("postgres", "err", err)
		os.Exit(1)
	}
	if err := migrations.Run(db); err != nil {
		logger.Error("migrations", "err", err)
		os.Exit(1)
	}

	// --- Redis ---
	redis, err := cache.New(cfg.Redis)
	if err != nil {
		logger.Error("redis", "err", err)
		os.Exit(1)
	}

	// --- Seed ---
	if err := seed.Run(context.Background(), db, cfg); err != nil {
		logger.Error("seed", "err", err)
		os.Exit(1)
	}

	// --- Repositories ---
	userRepo := repositories.NewUserRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	sellerRepo := repositories.NewSellerRepository(db)
	driverRepo := repositories.NewDriverRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	subcategoryRepo := repositories.NewSubcategoryRepository(db)
	brandRepo := repositories.NewBrandRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	bannerRepo := repositories.NewBannerRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	trackingRepo := repositories.NewTrackingRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	txRepo := repositories.NewTransactionRepository(db)

	// --- Infra services ---
	jm := appjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	googleOAuth := oauth.NewGoogle(cfg.Google)
	hub := ws.NewHub(redis)
	uploadSvc, err := services.NewUploadService(cfg.Upload.Dir, cfg.Upload.PublicBase, cfg.Upload.MaxSizeMB)
	if err != nil {
		logger.Error("upload", "err", err)
		os.Exit(1)
	}

	// --- Domain services ---
	notifySvc := services.NewNotificationService(notificationRepo, redis)
	authSvc := services.NewAuthService(userRepo, customerRepo, jm)
	userSvc := services.NewUserService(userRepo)
	catalogSvc := services.NewCatalogService(categoryRepo, subcategoryRepo, brandRepo)
	sellerSvc := services.NewSellerService(sellerRepo, userRepo)
	driverSvc := services.NewDriverService(driverRepo, userRepo)
	customerSvc := services.NewCustomerService(customerRepo, userRepo)
	productSvc := services.NewProductService(productRepo, sellerRepo, categoryRepo, notifySvc, services.ContactDefaults{
		Phone: cfg.Market.Phone, WhatsApp: cfg.Market.WhatsApp, Telegram: cfg.Market.Telegram,
	})
	orderSvc := services.NewOrderService(orderRepo, productRepo, customerRepo, driverRepo, txRepo, cartRepo, notifySvc)
	cartSvc := services.NewCartService(cartRepo, productRepo)
	favoriteSvc := services.NewFavoriteService(favoriteRepo, productRepo)
	reviewSvc := services.NewReviewService(reviewRepo, productRepo)
	bannerSvc := services.NewBannerService(bannerRepo)
	settingSvc := services.NewSettingService(settingRepo)
	trackingSvc := services.NewTrackingService(productRepo, trackingRepo)
	chatSvc := services.NewChatService(chatRepo, customerRepo, redis)
	reportSvc := services.NewReportService(txRepo, orderRepo)

	// --- Handlers / router ---
	deps := &server.Deps{
		Auth:         handlers.NewAuthHandler(authSvc, googleOAuth),
		Catalog:      handlers.NewCatalogHandler(catalogSvc),
		Product:      handlers.NewProductHandler(productSvc, sellerSvc, reviewSvc),
		Seller:       handlers.NewSellerHandler(sellerSvc, trackingSvc),
		Customer:     handlers.NewCustomerHandler(customerSvc),
		Driver:       handlers.NewDriverHandler(driverSvc, orderSvc),
		Order:        handlers.NewOrderHandler(orderSvc, cartSvc, customerSvc, settingSvc),
		Cart:         handlers.NewCartHandler(cartSvc, customerSvc),
		Favorite:     handlers.NewFavoriteHandler(favoriteSvc),
		Review:       handlers.NewReviewHandler(reviewSvc),
		Banner:       handlers.NewBannerHandler(bannerSvc),
		Setting:      handlers.NewSettingHandler(settingSvc),
		Tracking:     handlers.NewTrackingHandler(trackingSvc),
		Notification: handlers.NewNotificationHandler(notifySvc),
		Chat:         handlers.NewChatHandler(chatSvc, hub),
		NotifySocket: handlers.NewNotifySocketHandler(hub),
		Report:       handlers.NewReportHandler(reportSvc),
		Export:       handlers.NewExportHandler(productSvc, catalogSvc, orderSvc, customerSvc, sellerSvc, driverSvc, reportSvc, txRepo),
		Board:        handlers.NewBoardHandler(catalogSvc, productSvc),
		SEO:          handlers.NewSEOHandler(productSvc, catalogSvc, sellerSvc, cfg.PublicURL),
		Upload:       handlers.NewUploadHandler(uploadSvc),
		Admin:        handlers.NewAdminHandler(userSvc, trackingSvc),
		JWT:          jm,
		Cfg:          cfg,
	}
	r := server.New(deps)

	addr := cfg.Host + ":" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("AliStroy API listening", "addr", addr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown", "err", err)
	}
}

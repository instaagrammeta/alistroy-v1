package main

import (
	"context"
	"log"
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
	"github.com/instaagrammeta/alistroy-v1/backend/internal/repositories"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/seed"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/server"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/services"
	"github.com/instaagrammeta/alistroy-v1/backend/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// --- DB ---
	db, err := database.Connect(cfg.Postgres, cfg.Env)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}

	// --- Migrations ---
	if err := migrations.Run(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// --- Redis ---
	redis, err := cache.New(cfg.Redis)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	_ = redis // reserved for future caching layers

	// --- Seed (admin, settings, categories) ---
	if err := seed.Run(context.Background(), db, cfg); err != nil {
		log.Fatalf("seed: %v", err)
	}

	// --- Repositories ---
	userRepo := repositories.NewUserRepository(db)
	sellerRepo := repositories.NewSellerRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	trackingRepo := repositories.NewTrackingRepository(db)

	// --- JWT ---
	jm := appjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	// --- Services ---
	authSvc := services.NewAuthService(userRepo, sellerRepo, jm)
	userSvc := services.NewUserService(userRepo)
	categorySvc := services.NewCategoryService(categoryRepo)
	sellerSvc := services.NewSellerService(sellerRepo)
	productSvc := services.NewProductService(productRepo, sellerRepo, categoryRepo, settingRepo, services.ContactDefaults{
		Phone:    cfg.Market.Phone,
		WhatsApp: cfg.Market.WhatsApp,
	})
	reviewSvc := services.NewReviewService(reviewRepo, productRepo)
	favoriteSvc := services.NewFavoriteService(favoriteRepo, productRepo)
	settingSvc := services.NewSettingService(settingRepo)
	trackingSvc := services.NewTrackingService(productRepo, trackingRepo)
	uploadSvc, err := services.NewUploadService(cfg.Upload.Dir, cfg.Upload.PublicBase, cfg.Upload.MaxSizeMB)
	if err != nil {
		log.Fatalf("upload: %v", err)
	}

	// --- Handlers ---
	deps := &server.Deps{
		Auth:     handlers.NewAuthHandler(authSvc),
		Category: handlers.NewCategoryHandler(categorySvc),
		Seller:   handlers.NewSellerHandler(sellerSvc, productSvc, trackingSvc),
		Product:  handlers.NewProductHandler(productSvc, sellerSvc, favoriteSvc, reviewSvc),
		Review:   handlers.NewReviewHandler(reviewSvc),
		Favorite: handlers.NewFavoriteHandler(favoriteSvc),
		Tracking: handlers.NewTrackingHandler(trackingSvc),
		Setting:  handlers.NewSettingHandler(settingSvc),
		Upload:   handlers.NewUploadHandler(uploadSvc),
		Admin:    handlers.NewAdminHandler(userSvc, trackingSvc),
		JWT:      jm,
		Cfg:      cfg,
	}
	r := server.New(deps)

	addr := cfg.Host + ":" + cfg.Port
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("AliStroy API listening on %s (env=%s)", addr, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/handlers"
	appjwt "github.com/instaagrammeta/alistroy-v1/backend/internal/jwt"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

// Deps groups all the wired-up handlers. The router knows nothing about
// services or repositories — only handlers + middleware.
type Deps struct {
	Auth     *handlers.AuthHandler
	Category *handlers.CategoryHandler
	Seller   *handlers.SellerHandler
	Product  *handlers.ProductHandler
	Review   *handlers.ReviewHandler
	Favorite *handlers.FavoriteHandler
	Tracking *handlers.TrackingHandler
	Setting  *handlers.SettingHandler
	Upload   *handlers.UploadHandler
	Admin    *handlers.AdminHandler

	JWT *appjwt.Manager
	Cfg *config.Config
}

func New(d *Deps) *gin.Engine {
	if d.Cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	if d.Cfg.LogLevel == "debug" {
		r.Use(gin.Logger())
	}
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(d.Cfg.CORS.AllowedOrigins))
	r.Use(middleware.RateLimit(d.Cfg.Rate.RPS, d.Cfg.Rate.Burst))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	// ---- Public auth ----
	auth := v1.Group("/auth")
	{
		auth.POST("/register", d.Auth.Register)
		auth.POST("/login", d.Auth.Login)
		auth.POST("/refresh", d.Auth.Refresh)
		auth.POST("/forgot-password", d.Auth.Forgot)
		auth.POST("/reset-password", d.Auth.Reset)
	}

	// ---- Authenticated /me ----
	me := v1.Group("/me", middleware.RequireAuth(d.JWT))
	{
		me.GET("", d.Auth.Me)
		me.PATCH("", d.Auth.UpdateProfile)
		me.POST("/change-password", d.Auth.ChangePassword)
		me.POST("/logout", d.Auth.Logout)
	}

	// ---- Public catalog ----
	v1.GET("/categories", d.Category.List)
	v1.GET("/categories/popular", d.Category.Popular)
	v1.GET("/categories/:slug", d.Category.GetBySlug)

	v1.GET("/sellers", d.Seller.List)
	v1.GET("/sellers/top", d.Seller.Top)
	v1.GET("/sellers/:slug", d.Seller.GetBySlug)

	v1.GET("/products", d.Product.List)
	v1.GET("/products/:slug", d.Product.GetBySlug)
	v1.GET("/products/:slug/reviews", func(c *gin.Context) {
		// For simplicity, by-id endpoint exists below; alias by slug not needed.
		c.AbortWithStatusJSON(404, gin.H{"error": gin.H{"code": "not_found", "message": "use /products/id/:id/reviews"}})
	})
	// By-id read paths used after fetching product detail
	byID := v1.Group("/products/id/:id")
	{
		byID.GET("/related", d.Product.Related)
		byID.GET("/reviews", d.Review.ListForProduct)
		byID.POST("/track", middleware.OptionalAuth(d.JWT), d.Tracking.Track)
	}

	// ---- Settings (public read) ----
	v1.GET("/settings/public", d.Setting.Public)

	// ---- Authenticated user actions ----
	user := v1.Group("/", middleware.RequireAuth(d.JWT))
	{
		user.GET("/favorites", d.Favorite.List)
		user.POST("/favorites/:id", d.Favorite.Add)
		user.DELETE("/favorites/:id", d.Favorite.Remove)
		user.GET("/favorites/:id/has", d.Favorite.Has)

		user.POST("/products/id/:id/reviews", d.Review.Create)
	}

	// ---- Seller-only ----
	seller := v1.Group("/seller", middleware.RequireRoles(d.JWT, models.RoleSeller))
	{
		seller.GET("/me", d.Seller.Me)
		seller.PATCH("/me", d.Seller.UpdateMe)
		seller.GET("/stats", d.Seller.MyStats)

		seller.GET("/products", d.Product.MyList)
		seller.POST("/products", d.Product.MyCreate)
		seller.PATCH("/products/:id", d.Product.MyUpdate)
		seller.DELETE("/products/:id", d.Product.MyDelete)

		seller.POST("/upload", d.Upload.Upload)
	}

	// ---- Admin-only ----
	admin := v1.Group("/admin", middleware.RequireRoles(d.JWT, models.RoleAdmin))
	{
		admin.GET("/dashboard", d.Admin.Dashboard)
		admin.GET("/totals", d.Tracking.AdminTotals)

		admin.GET("/users", d.Admin.ListUsers)
		admin.GET("/users/:id", d.Admin.GetUser)
		admin.PATCH("/users/:id", d.Admin.UpdateUser)
		admin.DELETE("/users/:id", d.Admin.DeleteUser)

		admin.GET("/categories", d.Category.AdminList)
		admin.POST("/categories", d.Category.Create)
		admin.PATCH("/categories/:id", d.Category.Update)
		admin.DELETE("/categories/:id", d.Category.Delete)

		admin.GET("/sellers", d.Seller.AdminList)
		admin.GET("/sellers/:id", d.Seller.AdminGet)
		admin.PATCH("/sellers/:id", d.Seller.AdminUpdate)
		admin.DELETE("/sellers/:id", d.Seller.AdminDelete)

		admin.GET("/products", d.Product.AdminList)
		admin.GET("/products/:id", d.Product.AdminGet)
		admin.PATCH("/products/:id", d.Product.AdminUpdate)
		admin.POST("/products/:id/moderate", d.Product.AdminModerate)
		admin.DELETE("/products/:id", d.Product.AdminDelete)

		admin.GET("/reviews", d.Review.AdminList)
		admin.POST("/reviews/:id/moderate", d.Review.AdminModerate)
		admin.DELETE("/reviews/:id", d.Review.AdminDelete)

		admin.GET("/settings", d.Setting.AdminGet)
		admin.PATCH("/settings", d.Setting.AdminUpdate)

		admin.POST("/upload", d.Upload.Upload)
	}

	return r
}

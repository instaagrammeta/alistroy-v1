package server

import (
	"github.com/gin-gonic/gin"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/handlers"
	appjwt "github.com/instaagrammeta/alistroy-v1/backend/internal/jwt"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/middleware"
	"github.com/instaagrammeta/alistroy-v1/backend/internal/models"
)

// Deps holds every wired handler plus shared infra used by the router.
type Deps struct {
	Auth         *handlers.AuthHandler
	Catalog      *handlers.CatalogHandler
	Product      *handlers.ProductHandler
	Seller       *handlers.SellerHandler
	Customer     *handlers.CustomerHandler
	Driver       *handlers.DriverHandler
	Order        *handlers.OrderHandler
	Cart         *handlers.CartHandler
	Favorite     *handlers.FavoriteHandler
	Review       *handlers.ReviewHandler
	Banner       *handlers.BannerHandler
	Setting      *handlers.SettingHandler
	Tracking     *handlers.TrackingHandler
	Notification *handlers.NotificationHandler
	Chat         *handlers.ChatHandler
	NotifySocket *handlers.NotifySocketHandler
	Report       *handlers.ReportHandler
	Export       *handlers.ExportHandler
	Board        *handlers.BoardHandler
	SEO          *handlers.SEOHandler
	Upload       *handlers.UploadHandler
	Admin        *handlers.AdminHandler

	JWT *appjwt.Manager
	Cfg *config.Config
}

func New(d *Deps) *gin.Engine {
	if d.Cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// Allow up to 64 MB multipart uploads (product/banner images can be 50 MB).
	r.MaxMultipartMemory = 64 << 20
	r.Use(gin.Recovery())
	if d.Cfg.LogLevel == "debug" {
		r.Use(gin.Logger())
	}
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(d.Cfg.CORS.AllowedOrigins))
	r.Use(middleware.RateLimit(d.Cfg.Rate.RPS, d.Cfg.Rate.Burst))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// SEO (root-level)
	r.GET("/robots.txt", d.SEO.Robots)
	r.GET("/sitemap.xml", d.SEO.Sitemap)

	v1 := r.Group("/api/v1")

	// ---------- Public auth ----------
	auth := v1.Group("/auth")
	{
		auth.POST("/register", d.Auth.Register)
		auth.POST("/login", d.Auth.Login)
		auth.POST("/refresh", d.Auth.Refresh)
		auth.GET("/google/url", d.Auth.GoogleURL)
		auth.POST("/google/callback", d.Auth.GoogleCallback)
	}

	// ---------- Authenticated self ----------
	me := v1.Group("/me", middleware.RequireAuth(d.JWT))
	{
		me.GET("", d.Auth.Me)
		me.PATCH("", d.Auth.UpdateProfile)
		me.POST("/change-password", d.Auth.ChangePassword)
		me.POST("/logout", d.Auth.Logout)

		me.GET("/notifications", d.Notification.List)
		me.GET("/notifications/unread", d.Notification.UnreadCount)
		me.POST("/notifications/read/:id", d.Notification.MarkRead)
		me.POST("/notifications/read-all", d.Notification.MarkAllRead)
		me.GET("/notifications/socket", d.NotifySocket.Socket)
	}

	// ---------- Public catalog ----------
	// NOTE: gin does not allow a static segment and a :param to be siblings of
	// the same path node, so the by-slug routes get their own clean node and
	// all helper endpoints live under /catalog/*.
	v1.GET("/categories", d.Catalog.ListCategories)
	v1.GET("/categories/:slug", d.Catalog.GetCategoryBySlug)
	v1.GET("/brands", d.Catalog.ListBrands)
	v1.GET("/banners", d.Banner.Public)
	v1.GET("/settings/public", d.Setting.Public)

	v1.GET("/sellers", d.Seller.List)
	v1.GET("/sellers/:slug", d.Seller.GetBySlug)

	v1.GET("/products", d.Product.List)
	v1.GET("/products/:slug", d.Product.GetBySlug)

	// Catalog helper endpoints (kept off the :slug nodes).
	catalog := v1.Group("/catalog")
	{
		catalog.GET("/popular-categories", d.Catalog.PopularCategories)
		catalog.GET("/top-sellers", d.Seller.Top)
		catalog.GET("/price-bounds", d.Product.PriceBounds)
		catalog.GET("/subcategories/:id", d.Catalog.ListSubcategories)
		catalog.GET("/product/:id/related", d.Product.Related)
		catalog.GET("/product/:id/reviews", d.Review.ListForProduct)
		catalog.POST("/product/:id/track", middleware.OptionalAuth(d.JWT), d.Tracking.Track)
	}

	// ---------- Authenticated customer actions ----------
	user := v1.Group("", middleware.RequireAuth(d.JWT))
	{
		user.GET("/favorites", d.Favorite.List)
		user.POST("/favorites/:id", d.Favorite.Add)
		user.DELETE("/favorites/:id", d.Favorite.Remove)
		user.GET("/favorites/:id/has", d.Favorite.Has)

		user.POST("/catalog/product/:id/reviews", d.Review.Create)
	}

	// customer-only (must have customer profile)
	cust := v1.Group("/customer", middleware.RequireRoles(d.JWT, models.RoleCustomer))
	{
		cust.GET("/cart", d.Cart.List)
		cust.POST("/cart", d.Cart.Set)
		cust.POST("/cart/clear", d.Cart.Clear)
		cust.DELETE("/cart/:id", d.Cart.Remove)

		cust.POST("/checkout", d.Order.Checkout)
		cust.GET("/orders", d.Order.MyOrders)

		cust.GET("/chat/room", d.Chat.MyRoom)
		cust.GET("/chat/messages", d.Chat.MyMessages)
		cust.POST("/chat/messages", d.Chat.CustomerSend)
		cust.GET("/chat/socket", d.Chat.CustomerSocket)
		cust.POST("/upload", d.Upload.Upload)
	}

	// ---------- Seller-only ----------
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

	// ---------- Driver-only ----------
	driver := v1.Group("/driver", middleware.RequireRoles(d.JWT, models.RoleDriver))
	{
		driver.GET("/me", d.Driver.Me)
		driver.GET("/orders", d.Driver.MyOrders)
		driver.POST("/orders/:id/status", d.Order.DriverUpdateStatus)
	}

	// ---------- Admin-only ----------
	admin := v1.Group("/admin", middleware.RequireRoles(d.JWT, models.RoleAdmin))
	{
		admin.GET("/dashboard", d.Admin.Dashboard)
		admin.POST("/upload", d.Upload.Upload)

		// users
		admin.GET("/users", d.Admin.ListUsers)
		admin.GET("/users/:id", d.Admin.GetUser)
		admin.PATCH("/users/:id", d.Admin.UpdateUser)
		admin.DELETE("/users/:id", d.Admin.DeleteUser)

		// categories
		admin.GET("/categories", d.Catalog.AdminListCategories)
		admin.POST("/categories", d.Catalog.CreateCategory)
		admin.PATCH("/categories/:id", d.Catalog.UpdateCategory)
		admin.DELETE("/categories/:id", d.Catalog.DeleteCategory)

		// subcategories
		admin.POST("/subcategories", d.Catalog.CreateSubcategory)
		admin.PATCH("/subcategories/:id", d.Catalog.UpdateSubcategory)
		admin.DELETE("/subcategories/:id", d.Catalog.DeleteSubcategory)

		// brands
		admin.GET("/brands", d.Catalog.ListBrands)
		admin.POST("/brands", d.Catalog.CreateBrand)
		admin.PATCH("/brands/:id", d.Catalog.UpdateBrand)
		admin.DELETE("/brands/:id", d.Catalog.DeleteBrand)

		// banners
		admin.GET("/banners", d.Banner.AdminList)
		admin.POST("/banners", d.Banner.Create)
		admin.PATCH("/banners/:id", d.Banner.Update)
		admin.DELETE("/banners/:id", d.Banner.Delete)

		// sellers
		admin.GET("/sellers", d.Seller.AdminList)
		admin.GET("/sellers/:id", d.Seller.AdminGet)
		admin.POST("/sellers", d.Seller.AdminCreate)
		admin.PATCH("/sellers/:id", d.Seller.AdminUpdate)
		admin.DELETE("/sellers/:id", d.Seller.AdminDelete)

		// customers
		admin.GET("/customers", d.Customer.AdminList)
		admin.GET("/customers/:id", d.Customer.AdminGet)
		admin.POST("/customers", d.Customer.AdminCreate)
		admin.PATCH("/customers/:id", d.Customer.AdminUpdate)
		admin.DELETE("/customers/:id", d.Customer.AdminDelete)

		// drivers
		admin.GET("/drivers", d.Driver.AdminList)
		admin.GET("/drivers/:id", d.Driver.AdminGet)
		admin.POST("/drivers", d.Driver.AdminCreate)
		admin.PATCH("/drivers/:id", d.Driver.AdminUpdate)
		admin.DELETE("/drivers/:id", d.Driver.AdminDelete)

		// products
		admin.GET("/products", d.Product.AdminList)
		admin.GET("/products/:id", d.Product.AdminGet)
		admin.POST("/products", d.Product.AdminCreate)
		admin.PATCH("/products/:id", d.Product.AdminUpdate)
		admin.POST("/products/:id/moderate", d.Product.AdminModerate)
		admin.DELETE("/products/:id", d.Product.AdminDelete)

		// orders
		admin.GET("/orders", d.Order.AdminList)
		admin.GET("/orders/:id", d.Order.AdminGet)
		admin.POST("/orders", d.Order.AdminCreate)
		admin.POST("/orders/:id/status", d.Order.AdminUpdateStatus)
		admin.GET("/orders/:id/receipt", d.Order.Receipt)
		admin.DELETE("/orders/:id", d.Order.AdminDelete)

		// reviews
		admin.GET("/reviews", d.Review.AdminList)
		admin.POST("/reviews/:id/moderate", d.Review.AdminModerate)
		admin.DELETE("/reviews/:id", d.Review.AdminDelete)

		// chat
		admin.GET("/chat/rooms", d.Chat.AdminRooms)
		admin.GET("/chat/rooms/:id/messages", d.Chat.AdminMessages)
		admin.POST("/chat/rooms/:id/messages", d.Chat.AdminSend)
		admin.GET("/chat/rooms/:id/socket", d.Chat.AdminSocket)

		// reports + accounting
		admin.GET("/reports/summary", d.Report.Summary)
		admin.GET("/reports/transactions", d.Report.Transactions)

		// exports
		admin.GET("/export/products", d.Export.Products)
		admin.GET("/export/categories", d.Export.Categories)
		admin.GET("/export/orders", d.Export.Orders)
		admin.GET("/export/customers", d.Export.Customers)
		admin.GET("/export/sellers", d.Export.Sellers)
		admin.GET("/export/drivers", d.Export.Drivers)
		admin.GET("/export/report", d.Export.Transactions)

		// visual board
		admin.GET("/board", d.Board.Tree)

		// settings
		admin.GET("/settings", d.Setting.AdminGet)
		admin.PATCH("/settings", d.Setting.AdminUpdate)
	}

	return r
}

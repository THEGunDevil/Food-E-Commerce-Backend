package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/config"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/handlers"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/middleware"
	service "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, skipping...")
	}
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// Init Cloudinary
	cldURL := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	service.InitCloudinary(cldURL)

	// Load config & connect to DB
	cfg := config.LoadConfig()
	db.LocalConnect(cfg)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.Run(ctx)

	// Gin router
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://himel-s-library.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// -----------------------
	// Auth routes (public)
	// -----------------------
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/refresh", handlers.RefreshHandler)
		auth.POST("/logout", handlers.LogoutHandler)
	}

	// -----------------------
	// Users
	// -----------------------
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", middleware.AdminOnly(), handlers.GetUsersHandler)
		users.GET("/search", middleware.AdminOnly(), handlers.SearchUsersPaginatedHandler)
		users.GET("/:id", handlers.GetUserByIDHandler)
		users.PATCH("/:id", handlers.UpdateUserByIDHandler)
		users.DELETE("/:id/profile-image", handlers.DeleteProfileImage)
		users.PATCH("/:id/ban", middleware.AdminOnly(), handlers.BanUserByIDHandler)
	}

	// -----------------------
	// Addresses
	// -----------------------
	addresses := r.Group("/addresses")
	addresses.Use(middleware.AuthMiddleware())
	{
		addresses.POST("", handlers.CreateUserAddressHandler)
		addresses.GET("/user/:user_id", handlers.GetUserAddressByIDHandler)
		addresses.PATCH("/user/:user_id", handlers.UpdateUserAddressByIDHandler)
		addresses.DELETE("/user/:user_id", handlers.DeleteUserAddressByIDHandler)
	}

	// -----------------------
	// Reviews
	// -----------------------
	reviews := r.Group("/reviews")
	reviews.Use(middleware.AuthMiddleware())
	{
		reviews.POST("/:menu_id", handlers.CreateReviewHandler)
		reviews.GET("/menu/:menu_id", handlers.GetMenuItemReviewsByMenuItemIDHandler)
		reviews.GET("/review/:review_id", handlers.GetReviewByIDHandler)		
		reviews.GET("/users/:user_id", handlers.ListReviewsByUserHandler)
		// reviews.PATCH("/:review_id", handlers.UpdateReviewByIDHandler)
		reviews.DELETE("/:review_id", handlers.DeleteReviewByIDHandler)
	}

	// -----------------------
	// Orders
	// -----------------------
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("", handlers.CreateOrderHandler)
		orders.GET("/:order_id", handlers.GetOrderByIDHandler)
		orders.GET("/users/:user_id", handlers.ListOrdersByUserHandler)
		orders.PATCH("/:order_id", handlers.UpdateOrderByIDHandler)
		orders.DELETE("/:order_id", handlers.DeleteOrderByIDHandler)
	}

	// -----------------------
	// Order Items
	// -----------------------
	orderItems := r.Group("/order-items")
	orderItems.Use(middleware.AuthMiddleware())
	{
		orderItems.POST("/orders/:order_id", handlers.CreateOrderItemHandler)
		orderItems.GET("/orders/:order_id", handlers.ListOrderItemsByOrderHandler)
		orderItems.PATCH("/:order_item_id", handlers.UpdateOrderItemByIDHandler)
		orderItems.DELETE("/:order_item_id", handlers.DeleteOrderItemHandler)
	}

	// -----------------------
	// Menus
	// -----------------------
	menus := r.Group("/menus")
	{
		menus.POST("", handlers.CreateMenuItemHandler)
		menus.GET("", handlers.ListMenusHandler)
		menus.GET("/:category_id", handlers.ListMenusByCategoryHandler)
		menus.GET("/menu/:menu_id", handlers.GetMenuByMenuIDHandler)
		menus.PATCH("/menu/:menu_id", handlers.UpdateMenuByMenuIDHandler)
		menus.DELETE("/menu/:menu_id", handlers.DeleteMenuByMenuIDHandler)
	}

	// -----------------------
	// Favorites
	// -----------------------
	favorites := r.Group("/favorites")
	favorites.Use(middleware.AuthMiddleware())
	{
		favorites.POST("", handlers.CreateFavoritesHandler)
		favorites.GET("/user", handlers.ListFavoritesHandler)
		favorites.DELETE("/menu/:menu_id", handlers.DeleteFavoriteHandler)
	}

	// -----------------------
	// Delivery Zones (Admin)
	// -----------------------
	zones := r.Group("/delivery-zones")
	zones.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		zones.POST("", handlers.CreateDeliveryZoneHandler)
		zones.GET("", handlers.ListDeliveryZoneHandler)
		zones.GET("/active", handlers.ListActiveDeliveryZoneHandler)
		zones.PATCH("/:delivery_zone_id", handlers.UpdateDeliveryZoneHandler)
		zones.PATCH("/:delivery_zone_id/toggle", handlers.ToggleDeliveryZoneStatusHandler)
		zones.DELETE("/:delivery_zone_id", handlers.DeleteDeliveryZoneHandler)
	}

	// -----------------------
	// Categories (Admin)
	// -----------------------
	categories := r.Group("/categories")
	{
		categories.POST("",middleware.AdminOnly(),handlers.CreateCategoryHandler)
		categories.GET("", handlers.ListCategoriesHandler)
		categories.GET("/active", handlers.ListActiveCategoriesHandler)
		categories.GET("/:category_id", handlers.GetCategoryByIDHandler)
		categories.PATCH("/:category_id", handlers.UpdateCategoryHandler)
		categories.PATCH("/:category_id/display-order", handlers.UpdateCategoryDisplayOrderHandler)
		categories.PATCH("/:category_id/deactivate", handlers.DeactivateCategoryHandler)
		categories.DELETE("/:category_id", handlers.DeleteCategoryHandler)
	}

	// -----------------------
	// Cart
	// -----------------------
	cart := r.Group("/cart")
	cart.POST("/add-items/:user_id", handlers.CreateCartItemsHandler)
	cart.GET("/items/:user_id", handlers.ListCartItemsHandler)
	cart.Use(middleware.AuthMiddleware())
	{
		cart.PATCH("/items/:cart_item_id", handlers.UpdateCartItemHandler)
		cart.DELETE("/items/:cart_item_id", handlers.RemoveCartItemHandler)
		cart.DELETE("/users/:user_id", handlers.ClearCartHandler)
	}

	// -----------------------
	// Promotions (Admin)
	// -----------------------
	promotions := r.Group("/promotions")
	promotions.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		promotions.POST("", handlers.CreatePromotionHandler)
		promotions.GET("", handlers.ListPromotionsHandler)
		promotions.GET("/active", handlers.ListActivePromotionsHandler)
		promotions.PATCH("/:promotion_id", handlers.UpdatePromotionHandler)
		promotions.PATCH("/:promotion_id/increment", handlers.IncrementPromotionUsageHandler)
		promotions.DELETE("/:promotion_id", handlers.DeletePromotionHandler)
	}

	// -----------------------
	// Notifications
	// -----------------------
	notifications := r.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		notifications.GET("", handlers.ListNotificationsHandler)
		notifications.GET("/:notification_id", handlers.GetNotificationByIDHandler)
		notifications.PATCH("/:notification_id/read", handlers.MarkNotificationAsReadHandler)
		notifications.DELETE("/:notification_id", handlers.DeleteNotificationByIDHandler)
	}

	// -----------------------
	// Start server
	// -----------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}

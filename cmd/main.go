package main
import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"
	"eventhub/internal/events"
	"log"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()
	userRepo   := auth.NewUserRepository(db)
	eventRepo  := events.NewEventRepository(db)
	ticketRepo := events.NewTicketRepository(db)
	authService   := auth.NewAuthService(userRepo, cfg.JWTSecret)
	eventService  := events.NewEventService(eventRepo)
	ticketService := events.NewTicketService(ticketRepo)
	authHandler   := auth.NewAuthHandler(authService)
	eventHandler  := events.NewEventHandler(eventService)
	ticketHandler := events.NewTicketHandler(ticketService)
	adminHandler  := auth.NewAdminHandler(userRepo)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	limiter := rate.NewLimiter(100, 200)
	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests, slow down"})
			return
		}
		c.Next()
	})
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login",    authHandler.Login)
	}
	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)
		protected.PUT("/me", authHandler.UpdateMe)
		protected.GET("/events",     eventHandler.GetEvents)
		protected.GET("/events/:id", eventHandler.GetEventByID)
		organizer := protected.Group("/")
		organizer.Use(auth.RoleMiddleware("organizer", "admin"))
		{
			organizer.POST("/events",       eventHandler.CreateEvent)
			organizer.PUT("/events/:id",    eventHandler.UpdateEvent)
			organizer.DELETE("/events/:id", eventHandler.DeleteEvent)
		}
		protected.POST("/tickets/book", ticketHandler.BookTicket)
		protected.GET("/tickets/my",    ticketHandler.GetMyTickets)
		protected.GET("/tickets/validate/:id", ticketHandler.ValidateQR)
		adminGroup := protected.Group("/admin")
		adminGroup.Use(auth.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users",          adminHandler.GetUsers)
			adminGroup.PUT("/users/:id/role", authHandler.UpdateUserRole)
		}
	}
	log.Println("🚀 Server running on :8080")
	log.Println("   POST /auth/register  POST /auth/login")
	log.Println("   GET  /events         GET  /events/:id")
	log.Println("   POST /events         PUT  /events/:id  (org/admin)")
	log.Println("   POST /tickets/book   GET  /tickets/my")
	log.Println("   GET  /admin/users    PUT  /admin/users/:id/role")
	r.Run(":8080")
}
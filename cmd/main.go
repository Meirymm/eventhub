package main
 
import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"
	"eventhub/internal/events"
	"eventhub/internal/tickets"
	"eventhub/internal/admin"
	"github.com/gin-gonic/gin"
	"log"
)
 
func main() {
	cfg := config.Load()
 
	db := database.Connect(cfg)
	defer db.Close()
 
	// ── Repositories ──────────────────────────────────────────
	userRepo   := auth.NewUserRepository(db)
	eventRepo  := events.NewEventRepository(db)
	ticketRepo := tickets.NewTicketRepository(db)
 
	// ── Services ──────────────────────────────────────────────
	authService   := auth.NewAuthService(userRepo, cfg.JWTSecret)
	eventService  := events.NewEventService(eventRepo)
	ticketService := tickets.NewTicketService(ticketRepo, eventRepo)
	adminService  := admin.NewAdminService(userRepo)
 
	// ── Handlers ──────────────────────────────────────────────
	authHandler   := auth.NewAuthHandler(authService)
	eventHandler  := events.NewEventHandler(eventService)
	ticketHandler := tickets.NewTicketHandler(ticketService)
	adminHandler  := admin.NewAdminHandler(adminService)
 
	// ── Router ────────────────────────────────────────────────
	r := gin.Default()
 
	// Public routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login",    authHandler.Login)
	}
 
	// Protected routes (JWT required)
	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		// Auth
		protected.GET("/me", authHandler.GetMe)
 
		// Events
		protected.GET("/events",  eventHandler.GetEvents)
		protected.POST("/events", eventHandler.CreateEvent)
 
		// Tickets
		protected.POST("/tickets/book", ticketHandler.BookTicket)
		protected.GET("/tickets/my",    ticketHandler.GetMyTickets)
 
		// Admin (только admin role)
		adminGroup := protected.Group("/admin")
		adminGroup.Use(auth.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users", adminHandler.GetUsers)
		}
	}
 
	log.Println("🚀 Server running on :8080")
	log.Println("📋 Endpoints:")
	log.Println("   POST   /auth/register")
	log.Println("   POST   /auth/login")
	log.Println("   GET    /me")
	log.Println("   GET    /events")
	log.Println("   POST   /events")
	log.Println("   POST   /tickets/book")
	log.Println("   GET    /tickets/my")
	log.Println("   GET    /admin/users")
 
	r.Run(":8080")
}
 
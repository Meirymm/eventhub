package main

import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"
	"eventhub/internal/events"
	"log"

	"github.com/gin-gonic/gin"
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

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)

		protected.GET("/events", eventHandler.GetEvents)
		protected.POST("/events", eventHandler.CreateEvent)

		protected.POST("/tickets/book", ticketHandler.BookTicket)
		protected.GET("/tickets/my", ticketHandler.GetMyTickets)

		adminGroup := protected.Group("/admin")
		adminGroup.Use(auth.RoleMiddleware("admin"))
		{
			adminGroup.GET("/users", adminHandler.GetUsers)
		}
	}

	log.Println("🚀 Server running on :8080")
	log.Println("📋 Endpoints:")
	log.Println("   POST  /auth/register")
	log.Println("   POST  /auth/login")
	log.Println("   GET   /me")
	log.Println("   GET   /events")
	log.Println("   POST  /events")
	log.Println("   POST  /tickets/book")
	log.Println("   GET   /tickets/my")
	log.Println("   GET   /admin/users")

	r.Run(":8080")
}

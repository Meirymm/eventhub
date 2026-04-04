package main

import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"
	"eventhub/internal/events" 

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	
	userRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	
	eventRepo := events.NewEventRepository(db)
	eventService := events.NewEventService(eventRepo)
	eventHandler := events.NewEventHandler(eventService)

	r := gin.Default()

	
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}


	r.GET("/events", eventHandler.GetEvents)

	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)
		
	}

	r.Run(":8080")
}

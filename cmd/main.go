package main

import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"
	"eventhub/internal/events" // Добавляем твой пакет

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()

	// --- Логика Dev 1 (Auth) ---
	userRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// --- Твоя логика Dev 2 (Events) ---
	eventRepo := events.NewEventRepository(db)
	eventService := events.NewEventService(eventRepo)
	eventHandler := events.NewEventHandler(eventService)

	r := gin.Default()

	// Группа регистрации/логина
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Твои открытые маршруты (посмотреть ивенты может любой)
	r.GET("/events", eventHandler.GetEvents)

	// Защищенные маршруты (нужен токен)
	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)
		// Если захочешь сделать создание ивента только для залогиненных:
		// protected.POST("/events", eventHandler.CreateEvent)
	}

	r.Run(":8080")
}

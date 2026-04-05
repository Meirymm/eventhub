package main

import (
	"eventhub/config"
	"eventhub/internal/auth"
	"eventhub/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	defer db.Close()
	userRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)
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
	}
	r.Run(":8080")
}

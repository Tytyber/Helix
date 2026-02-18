package main

import (
	"api/internal/auth"
	"api/internal/chat"
	"api/internal/config"
	"api/internal/database"
	"api/internal/group"
	"api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := database.New(cfg.DatabaseURL)

	authService := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	chatService := chat.NewService(db)
	chatHandler := chat.NewHandler(chatService)

	groupService := group.NewService(db)
	groupHandler := group.NewHandler(groupService)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	api := r.Group("/api/v1")
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
	}

	protected := api.Group("/protected")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.Me)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.POST("/reputation", authHandler.AddReputation)
		protected.POST("/chats", chatHandler.CreateChat)
		protected.GET("/chats/:id/messages", chatHandler.GetMessages)
		protected.POST("/chats/:id/messages", chatHandler.SendMessage)

		// Группы
		protected.POST("/groups", groupHandler.CreateGroup)
		protected.POST("/groups/:id/join", groupHandler.JoinGroup)
		protected.POST("/groups/:id/leave", groupHandler.LeaveGroup)
		protected.GET("/groups/:id/messages", groupHandler.GetMessages)
		protected.POST("/groups/:id/messages", groupHandler.SendMessage)
	}

	// Личные чаты

	r.Run(cfg.Port)
}

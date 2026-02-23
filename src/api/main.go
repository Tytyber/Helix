package main

import (
	"api/internal/auth"
	"api/internal/calls"
	"api/internal/chat"
	"api/internal/config"
	"api/internal/database"
	"api/internal/feed"
	friends "api/internal/friend"
	"api/internal/group"
	"api/internal/hub"
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

	hubService := hub.NewService(db)
	hubHandler := hub.NewHandler(hubService)

	feedService := feed.NewService(db)
	feedHandler := feed.NewHandler(feedService)

	friendsService := friends.NewService(db)
	friendsHandler := friends.NewHandler(friendsService)

	callsService := calls.NewService(db)
	callsHandler := calls.NewHandler(callsService)

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
		//Профиль
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

		//Хабы
		protected.GET("/hubs", hubHandler.GetAllHubs)
		protected.GET("/hubs/:id", hubHandler.GetHub)
		protected.POST("/hubs/:id/join", hubHandler.JoinHub)
		protected.POST("/hubs/:id/leave", hubHandler.LeaveHub)
		protected.GET("/hubs/:id/posts", hubHandler.GetHubPosts)
		protected.POST("/hubs/:id/posts", hubHandler.CreatePost)

		//лента
		protected.GET("/feed", feedHandler.UserFeed)
		protected.GET("/feed/top-hubs", feedHandler.TopHubs)

		//друзья
		protected.POST("/friends/request", friendsHandler.SendRequest)
		protected.POST("/friends/accept", friendsHandler.AcceptRequest)
		protected.GET("/friends", friendsHandler.ListFriends)
		protected.GET("/friends/requests", friendsHandler.IncomingRequests)
		protected.POST("/friends/decline", friendsHandler.DeclineRequest)

		//звонки
		protected.POST("/calls", callsHandler.CreateCall)
		protected.PUT("/calls/:id/status", callsHandler.UpdateCallStatus)
	}

	r.Run(cfg.Port)
}

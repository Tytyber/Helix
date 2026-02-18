package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

type updateProfileRequest struct {
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Status   string `json:"status,omitempty"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type repRequest struct {
	UserID string `json:"user_id" binding:"required"` // кому ставим
	Score  int    `json:"score" binding:"required"`   // +1 или -1
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	accessToken, refreshToken, err := h.Service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	accessToken, refreshToken, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var username, email, avatar, bio, status string
	var reputation int
	err := h.Service.DB.QueryRow(c, `
		SELECT username, email, avatar_url, bio, status, reputation
		FROM users
		WHERE id=$1`, userID).Scan(&username, &email, &avatar, &bio, &status, &reputation)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"id":         userID,
			"username":   username,
			"email":      email,
			"avatar":     avatar,
			"bio":        bio,
			"status":     status,
			"reputation": reputation,
		},
		"error": nil,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	userID, err := h.Service.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	h.Service.DeleteRefreshToken(req.RefreshToken)

	accessToken, _ := GenerateAccessToken(userID, h.Service.JWTSecret)
	refreshToken, _ := GenerateRefreshToken()
	h.Service.SaveRefreshToken(userID, refreshToken, 7*24*time.Hour)

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	h.Service.DeleteRefreshToken(req.RefreshToken)
	c.JSON(200, gin.H{"message": "logged out"})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	_, err := h.Service.DB.Exec(c, `
		UPDATE users
		SET username = COALESCE($1, username),
			avatar_url = COALESCE($2, avatar_url),
			bio = COALESCE($3, bio),
			status = COALESCE($4, status)
		WHERE id = $5
	`, nullIfEmpty(req.Username), nullIfEmpty(req.Avatar), nullIfEmpty(req.Bio), nullIfEmpty(req.Status), userID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "profile updated"})
}

// вспомогательная функция для COALESCE
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (h *Handler) AddReputation(c *gin.Context) {
	fromUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var req repRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if req.UserID == fromUserID {
		c.JSON(400, gin.H{"error": "cannot rate yourself"})
		return
	}

	if err := h.Service.AddReputation(req.UserID, fromUserID.(string), req.Score); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "reputation updated"})
}

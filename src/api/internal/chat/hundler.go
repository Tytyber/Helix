package chat

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

type createChatRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

func (h *Handler) CreateChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req createChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	chatID, err := h.Service.CreateChat(userID.(string), req.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"chat_id": chatID})
}

type sendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Service.SendMessage(chatID, userID.(string), req.Content); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "sent"})
}

func (h *Handler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	msgs, err := h.Service.GetMessages(chatID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"messages": msgs})
}

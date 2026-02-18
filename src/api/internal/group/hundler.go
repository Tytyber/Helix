package group

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

type createGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) CreateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	groupID, err := h.Service.CreateGroup(req.Name, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// добавляем создателя в группу
	h.Service.JoinGroup(groupID, userID.(string))

	c.JSON(200, gin.H{"group_id": groupID})
}

type joinLeaveRequest struct{}

func (h *Handler) JoinGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID := c.Param("id")

	if err := h.Service.JoinGroup(groupID, userID.(string)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "joined"})
}

func (h *Handler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID := c.Param("id")

	if err := h.Service.LeaveGroup(groupID, userID.(string)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "left"})
}

type sendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *Handler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID := c.Param("id")

	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Service.SendMessage(groupID, userID.(string), req.Content); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "sent"})
}

func (h *Handler) GetMessages(c *gin.Context) {
	groupID := c.Param("id")
	msgs, err := h.Service.GetMessages(groupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"messages": msgs})
}

package friends

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) SendRequest(c *gin.Context) {
	userID := c.GetString("userID")
	var req struct {
		ToID string `json:"to_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := h.Service.SendRequest(c, userID, req.ToID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request sent"})
}

func (h *Handler) AcceptRequest(c *gin.Context) {
	userID := c.GetString("userID")
	var req struct {
		FromID string `json:"from_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := h.Service.AcceptRequest(c, req.FromID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend added"})
}

func (h *Handler) ListFriends(c *gin.Context) {
	userID := c.GetString("userID")
	friends, err := h.Service.ListFriends(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

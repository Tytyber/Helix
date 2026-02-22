package calls

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateCall(c *gin.Context) {
	callerID := c.GetString("userID")
	var req struct {
		CalleeID string `json:"callee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	call, err := h.Service.CreateCall(c, callerID, req.CalleeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, call)
}

func (h *Handler) UpdateCallStatus(c *gin.Context) {
	callIDStr := c.Param("id")
	callID, _ := strconv.Atoi(callIDStr)

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Service.UpdateCallStatus(c, callID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

package hub

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

func (h *Handler) GetAllHubs(c *gin.Context) {
	hubs, err := h.Service.GetAllHubs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hubs": hubs})
}

func (h *Handler) GetHub(c *gin.Context) {
	hubID, _ := strconv.Atoi(c.Param("id"))
	hub, err := h.Service.GetHub(c, hubID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hub not found"})
		return
	}
	c.JSON(http.StatusOK, hub)
}

func (h *Handler) JoinHub(c *gin.Context) {
	hubID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetString("userID")
	if err := h.Service.JoinHub(c, hubID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "joined"})
}

func (h *Handler) LeaveHub(c *gin.Context) {
	hubID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetString("userID")
	if err := h.Service.LeaveHub(c, hubID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "left"})
}

func (h *Handler) GetHubPosts(c *gin.Context) {
	hubID, _ := strconv.Atoi(c.Param("id"))
	posts, err := h.Service.GetHubPosts(c, hubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) CreatePost(c *gin.Context) {
	hubID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetString("userID")
	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
		return
	}
	post, err := h.Service.CreatePost(c, hubID, userID, body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

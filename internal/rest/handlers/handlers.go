package handlers

import (
	repository "lab5/internal/storage/repo"
	"lab5/internal/storage/repo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repository.Neo4jRepository
}

func NewHandler(repo *repository.Neo4jRepository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) GetAllNodes(c *gin.Context) {
	nodes, err := h.Repo.GetAllNodesWithAttributes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения узлов"})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func (h *Handler) GetAllRelationships(c *gin.Context) {
	relationships, err := h.Repo.GetAllRelationships(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения связей"})
		return
	}
	c.JSON(http.StatusOK, relationships)
}

func (h *Handler) GetNodeWithRelationships(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID узла"})
		return
	}

	relationships, err := h.Repo.GetNodeWithRelationships(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения узла и связей"})
		return
	}
	c.JSON(http.StatusOK, relationships)
}

func (h *Handler) AddNodeAndRelationships(c *gin.Context) {
	var data models.NodeWithRelationships
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	if err := h.Repo.AddNodeAndRelationships(c.Request.Context(), data.Node, data.Relationships); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления узла и связей"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Узел и связи добавлены"})
}

func (h *Handler) DeleteNodeAndRelationships(c *gin.Context) {
	nodeID, err := strconv.Atoi(c.Param("node_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID узла"})
		return
	}

	if err := h.Repo.DeleteNodeAndRelationships(c.Request.Context(), nodeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления узла и связей"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Узел и связи удалены"})
}

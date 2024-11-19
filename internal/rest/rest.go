package rest

import (
	"lab5/config"
	"lab5/internal/rest/auth"
	"lab5/internal/rest/handlers"
	repository "lab5/internal/storage/repo"
	"log"

	"github.com/gin-gonic/gin"
)

func Init(cfg config.Config, repo *repository.Neo4jRepository) {
	// Create Gin router
	router := gin.Default()

	// Middleware for token verification
	authMiddleware := auth.VerifyToken()

	// Handlers
	handler := handlers.NewHandler(repo)

	// Define routes
	router.GET("/nodes", handler.GetAllNodes)
	router.GET("/relationships", handler.GetAllRelationships)
	router.GET("/node/:node_id", handler.GetNodeWithRelationships)
	router.POST("/node", authMiddleware, handler.AddNodeAndRelationships)
	router.DELETE("/node/:node_id", authMiddleware, handler.DeleteNodeAndRelationships)

	// Start server
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

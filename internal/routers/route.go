package routes

import (
	"lab5/internal/handlers"
	"lab5/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, neo4jService *services.Neo4jService, authToken string) {
	graphController := &handlers.GraphController{Neo4j: neo4jService}

	router.HandleFunc("/nodes", graphController.GetAllNodes).Methods(http.MethodGet)
	// Аналогично можно добавить остальные обработчики
}

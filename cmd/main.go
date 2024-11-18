package main

import (
	"log"
	"net/http"
	"project/config"
	"project/routes"
	"project/services"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	neo4jService, err := services.NewNeo4jService(cfg.Neo4jURI, cfg.Neo4jUsername, cfg.Neo4jPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer neo4jService.Close()

	router := mux.NewRouter()
	routes.RegisterRoutes(router, neo4jService, cfg.AuthToken)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

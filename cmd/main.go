package main

import (
	"lab5/config"
	"lab5/internal/rest"
	repository "lab5/internal/storage/repo"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to Neo4j
	neo4jRepo, err := repository.NewNeo4jRepository(cfg.Neo4j.URI, cfg.Neo4j.Username, cfg.Neo4j.Password)
	if err != nil {
		log.Fatalf("Ошибка подключения к Neo4j: %v", err)
	}
	defer neo4jRepo.Close()

	rest.Init(cfg, neo4jRepo)
}

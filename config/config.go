package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	Neo4jURI      string
	Neo4jUsername string
	Neo4jPassword string
	AuthToken     string
}

func LoadConfig() Config {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	return Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		Neo4jURI:      "bolt://" + getEnv("NEO4J_HOST", "localhost") + ":" + getEnv("NEO4J_BOLT_PORT", "7687"),
		Neo4jUsername: getEnv("NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", "securepassword"),
		AuthToken:     getEnv("AUTH_TOKEN", "your-secret-token"),
	}
}

// getEnv возвращает значение переменной среды или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

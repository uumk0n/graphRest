package config

type Config struct {
	Neo4jURI      string
	Neo4jUsername string
	Neo4jPassword string
	AuthToken     string
}

func LoadConfig() Config {
	return Config{
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "securepassword",
		AuthToken:     "your-secret-token",
	}
}

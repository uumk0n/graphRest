package services

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jService struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jService(uri, username, password string) (*Neo4jService, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jService{Driver: driver}, nil
}

func (s *Neo4jService) Close() {
	s.Driver.Close(context.Background())
}

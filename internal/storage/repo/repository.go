package repository

import (
	"context"
	"lab5/internal/storage/repo/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRepository struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jRepository(uri, username, password string) (*Neo4jRepository, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jRepository{Driver: driver}, nil
}

func (r *Neo4jRepository) Close() error {
	return r.Driver.Close(context.Background())
}

// GET всех узлов с атрибутами id, label
func (r *Neo4jRepository) GetAllNodesWithAttributes(ctx context.Context) ([]models.Node, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (n) RETURN n.id AS id, labels(n) AS label"
	nodes, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var nodes []models.Node
		for result.Next(ctx) {
			record := result.Record()
			nodes = append(nodes, models.Node{
				ID:   int(record.Values[0].(int64)),
				Name: record.Values[1].([]string)[0],
			})
		}
		return nodes, nil
	})
	if err != nil {
		return nil, err
	}
	return nodes.([]models.Node), nil
}

// GET узла и всех его связей с узлами на конце связей (узлы и связи - со всеми доступными атрибутами)
func (r *Neo4jRepository) GetNodeWithAllRelationships(ctx context.Context, nodeID int) ([]models.Relationship, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (n)-[r]->(m) WHERE n.id=$id
		RETURN n {.*}, type(r) AS relationship_type, m {.*} AS end_node
	`
	params := map[string]any{"id": nodeID}
	relationships, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var rels []models.Relationship
		for result.Next(ctx) {
			record := result.Record()
			rels = append(rels, models.Relationship{
				Node:             record.Values[0].(map[string]any),
				RelationshipType: record.Values[1].(string),
				EndNode:          record.Values[2].(map[string]any),
			})
		}
		return rels, nil
	})
	if err != nil {
		return nil, err
	}
	return relationships.([]models.Relationship), nil
}

// POST добавление узла и связей и/или сегмента графа
func (r *Neo4jRepository) AddGraphSegment(ctx context.Context, node models.Node, relationships []models.Relationship) error {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MERGE (n {id: $id})
		SET n += $node
		WITH n
		UNWIND $relationships AS rel
		MERGE (m {id: rel.end_node_id})
		SET m += rel.end_node
		MERGE (n)-[r:rel.relationship_type]->(m)
	`
	params := map[string]any{
		"id":            node.ID,
		"node":          node.Attributes, // Map of attributes
		"relationships": relationships,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// DELETE узла и связей и/или сегмента графа
func (r *Neo4jRepository) DeleteGraphSegment(ctx context.Context, nodeID int) error {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := "MATCH (n) WHERE n.id=$id DETACH DELETE n"
	params := map[string]any{"id": nodeID}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// GET сегмента графа (определите ограничения сегмента самостоятельно)
func (r *Neo4jRepository) GetGraphSegment(ctx context.Context, limit int) ([]models.Node, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	query := `
		MATCH (n)-[r]->(m)
		RETURN n.id AS start_node_id, labels(n) AS label, type(r) AS relationship_type, m.id AS end_node_id
		LIMIT $limit
	`
	params := map[string]any{"limit": limit}
	nodes, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var segment []models.Node
		for result.Next(ctx) {
			record := result.Record()
			segment = append(segment, models.Node{
				ID:   int(record.Values[0].(int64)),
				Name: record.Values[1].([]string)[0],
				Relationships: []models.Relationship{
					{
						RelationshipType: record.Values[2].(string),
						EndNodeID:        int(record.Values[3].(int64)),
					},
				},
			})
		}
		return segment, nil
	})
	if err != nil {
		return nil, err
	}
	return nodes.([]models.Node), nil
}
func (r *Neo4jRepository) GetAllRelationships(ctx context.Context) ([]models.Relationship, error) {
	query := `
		MATCH (n)-[r]->(m)
		RETURN n.id AS start_node_id, type(r) AS relationship_type, m.id AS end_node_id, m {.*} AS end_node
	`
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, err
		}

		var relationships []models.Relationship
		for records.Next(ctx) {
			record := records.Record()
			relationships = append(relationships, models.Relationship{
				StartNodeID:      int(record.Values[0].(int64)),
				RelationshipType: record.Values[1].(string),
				EndNodeID:        int(record.Values[2].(int64)),
				EndNode:          record.Values[3].(map[string]interface{}),
			})
		}
		return relationships, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]models.Relationship), nil
}

// Метод: GetNodeWithRelationships
func (r *Neo4jRepository) GetNodeWithRelationships(ctx context.Context, nodeID int) ([]models.Relationship, error) {
	query := `
		MATCH (n)-[r]->(m) WHERE n.id=$id
		RETURN n {.*} AS node, type(r) AS relationship_type, m {.*} AS end_node
	`
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{"id": nodeID}
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var relationships []models.Relationship
		for records.Next(ctx) {
			record := records.Record()
			relationships = append(relationships, models.Relationship{
				StartNodeID:      nodeID,
				RelationshipType: record.Values[1].(string),
				EndNode:          record.Values[2].(map[string]interface{}),
			})
		}
		return relationships, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]models.Relationship), nil
}

// Метод: AddNodeAndRelationships
func (r *Neo4jRepository) AddNodeAndRelationships(ctx context.Context, node models.Node, relationships []models.Relationship) error {
	query := `
		MERGE (n:User {id: $id, name: $name})
		WITH n
		UNWIND $relationships as rel
		MERGE (m:User {id: rel.end_node_id})
		MERGE (n)-[r:rel.relationship_type]->(m)
	`
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"id":            node.ID,
		"name":          node.Name,
		"relationships": relationships,
	}
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

// Метод: DeleteNodeAndRelationships
func (r *Neo4jRepository) DeleteNodeAndRelationships(ctx context.Context, nodeID int) error {
	query := "MATCH (n) WHERE n.id=$id DETACH DELETE n"
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{"id": nodeID}
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})
	return err
}

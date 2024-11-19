package handlers

import (
	"encoding/json"
	"lab5/internal/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphController struct {
	Neo4j *services.Neo4jService
}

// GET всех узлов с атрибутами id, label
func (gc *GraphController) GetAllNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := gc.Neo4j.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	nodes, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "MATCH (n) RETURN n.id AS id, labels(n)[0] AS label", nil)
		if err != nil {
			return nil, err
		}
		var nodes []map[string]interface{}
		for result.Next(ctx) {
			record := result.Record()
			nodes = append(nodes, map[string]interface{}{
				"id":    record.Values[0],
				"label": record.Values[1],
			})
		}
		return nodes, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

// GET узла и всех его связей с узлами на конце связей
func (gc *GraphController) GetNodeWithRelationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := gc.Neo4j.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	id := mux.Vars(r)["id"]

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (n)-[r]->(m)
			WHERE n.id = $id
			RETURN n, r, m`
		params := map[string]interface{}{"id": id}
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var nodes []map[string]interface{}
		for result.Next(ctx) {
			record := result.Record()
			nodes = append(nodes, map[string]interface{}{
				"node":           record.Values[0],
				"relationship":   record.Values[1],
				"connected_node": record.Values[2],
			})
		}
		return nodes, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// POST добавление узла и связей и/или сегмента графа
func (gc *GraphController) CreateNodeAndRelations(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Node          map[string]interface{}   `json:"node"`
		Relationships []map[string]interface{} `json:"relationships"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	session := gc.Neo4j.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) error {
		// Create the node
		query := "CREATE (n:Node {id: $id, label: $label})"
		_, err := tx.Run(ctx, query, requestData.Node)
		if err != nil {
			return err
		}

		// Create relationships
		for _, rel := range requestData.Relationships {
			_, err := tx.Run(ctx, "MATCH (a), (b) WHERE a.id = $startId AND b.id = $endId CREATE (a)-[r:CONNECTED]->(b)", rel)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Node and relationships created"))
}

// DELETE удаление узла и связей и/или сегмента графа
func (gc *GraphController) DeleteNodeAndRelations(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	ctx := r.Context()
	session := gc.Neo4j.Driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) error {
		// Delete relationships first
		_, err := tx.Run(ctx, "MATCH (n)-[r]->() WHERE n.id = $id DELETE r", map[string]interface{}{"id": id})
		if err != nil {
			return err
		}

		// Then delete the node
		_, err = tx.Run(ctx, "MATCH (n) WHERE n.id = $id DELETE n", map[string]interface{}{"id": id})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Node and its relationships deleted"))
}

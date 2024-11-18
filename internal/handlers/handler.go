package handlers

import (
	"encoding/json"
	"lab5/internal/services"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphController struct {
	Neo4j *services.Neo4jService
}

// GET всех узлов
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

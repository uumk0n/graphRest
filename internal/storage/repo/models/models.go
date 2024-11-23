package models

type Node struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Attributes    map[string]string `json:"attributes"`
	Relationships []Relationship    `json:"relationships"` // Add this field
}

type Relationship struct {
	StartNodeID      int            `json:"start_node_id"`
	EndNodeID        int            `json:"end_node_id"`
	RelationshipType string         `json:"relationship_type"`
	EndNode          map[string]any `json:"end_node"`
	Node             map[string]any `json:"node"` // Add this field
}

type NodeWithRelationships struct {
	Node          Node           `json:"node"`
	Relationships []Relationship `json:"relationships"`
}

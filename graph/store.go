package graph

import "fmt"

// GraphStore is the main entry point for graph operations.
// It wraps any Storage backend (ObjectBox, in-memory, etc.).
//
// Usage:
//
//	// For testing:
//	   store := graph.New(graph.NewMemoryStorage())
//
//	// For production (future ObjectBox backend):
//	   store := graph.New(objectboxstorage.New(obxStore))
type GraphStore struct {
	storage Storage
}

// New creates a new GraphStore with the given storage backend.
func New(storage Storage) *GraphStore {
	return &GraphStore{storage: storage}
}

// AddNode inserts a new node and returns it with the assigned ID.
func (g *GraphStore) AddNode(label string, properties map[string]string) (*Node, error) {
	node := &Node{Label: label, Properties: properties}
	id, err := g.storage.PutNode(node)
	if err != nil {
		return nil, fmt.Errorf("AddNode: %w", err)
	}
	node.Id = id
	return node, nil
}

// AddEdge inserts a directed, weighted edge between two nodes.
func (g *GraphStore) AddEdge(fromId, toId uint64, label string, weight float64) (*Edge, error) {
	edge := &Edge{FromId: fromId, ToId: toId, Label: label, Weight: weight}
	id, err := g.storage.PutEdge(edge)
	if err != nil {
		return nil, fmt.Errorf("AddEdge: %w", err)
	}
	edge.Id = id
	return edge, nil
}

// GetNode retrieves a node by its ID.
func (g *GraphStore) GetNode(id uint64) (*Node, error) {
	return g.storage.GetNode(id)
}

// Neighbors returns all nodes directly reachable from nodeId (1 hop).
func (g *GraphStore) Neighbors(nodeId uint64) ([]*Node, error) {
	edges, err := g.storage.EdgesFrom(nodeId)
	if err != nil {
		return nil, err
	}
	var neighbors []*Node
	for _, e := range edges {
		n, err := g.storage.GetNode(e.ToId)
		if err != nil {
			continue
		}
		neighbors = append(neighbors, n)
	}
	return neighbors, nil
}

// Close closes the underlying storage.
func (g *GraphStore) Close() error {
	return g.storage.Close()
}

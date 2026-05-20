package graph

import (
	"fmt"
	"sync"
)

// MemoryStorage is a thread-safe in-memory implementation of Storage.
// It mimics what ObjectBox would provide as the actual backend.
// Use this for unit tests, demos, and as a reference for ObjectBoxStorage.
type MemoryStorage struct {
	mu         sync.RWMutex
	nodes      map[uint64]*Node
	edges      map[uint64]*Edge
	nextNodeId uint64
	nextEdgeId uint64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		nodes: make(map[uint64]*Node),
		edges: make(map[uint64]*Edge),
	}
}

func (m *MemoryStorage) PutNode(node *Node) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if node.Id == 0 {
		m.nextNodeId++
		node.Id = m.nextNodeId
	}
	clone := *node
	m.nodes[node.Id] = &clone
	return node.Id, nil
}

func (m *MemoryStorage) GetNode(id uint64) (*Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.nodes[id]
	if !ok {
		return nil, fmt.Errorf("node %d not found", id)
	}
	return n, nil
}

func (m *MemoryStorage) AllNodes() ([]*Node, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Node, 0, len(m.nodes))
	for _, n := range m.nodes {
		result = append(result, n)
	}
	return result, nil
}

func (m *MemoryStorage) RemoveNode(id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.nodes, id)
	return nil
}

func (m *MemoryStorage) PutEdge(edge *Edge) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if edge.Id == 0 {
		m.nextEdgeId++
		edge.Id = m.nextEdgeId
	}
	clone := *edge
	m.edges[edge.Id] = &clone
	return edge.Id, nil
}

func (m *MemoryStorage) GetEdge(id uint64) (*Edge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.edges[id]
	if !ok {
		return nil, fmt.Errorf("edge %d not found", id)
	}
	return e, nil
}

func (m *MemoryStorage) EdgesFrom(nodeId uint64) ([]*Edge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*Edge
	for _, e := range m.edges {
		if e.FromId == nodeId {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *MemoryStorage) EdgesTo(nodeId uint64) ([]*Edge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*Edge
	for _, e := range m.edges {
		if e.ToId == nodeId {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *MemoryStorage) AllEdges() ([]*Edge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Edge, 0, len(m.edges))
	for _, e := range m.edges {
		result = append(result, e)
	}
	return result, nil
}

func (m *MemoryStorage) RemoveEdge(id uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.edges, id)
	return nil
}

func (m *MemoryStorage) Close() error {
	return nil
}

package graph

// Storage defines the contract that any backend must fulfill to work
// with the graph layer.
//
// Current implementations:
//   - MemoryStorage  — in-memory, for testing and demos (this repo)
//   - ObjectBoxStorage — objectbox-go backend (planned, see ROADMAP)
//
// This interface is the bridge between the graph algorithms and the
// storage engine, following the same pattern as pgvector (PostgreSQL
// extension) and MongoDB Atlas Vector Search.
type Storage interface {
	// Node operations
	PutNode(node *Node) (uint64, error)
	GetNode(id uint64) (*Node, error)
	AllNodes() ([]*Node, error)
	RemoveNode(id uint64) error

	// Edge operations
	PutEdge(edge *Edge) (uint64, error)
	GetEdge(id uint64) (*Edge, error)
	EdgesFrom(nodeId uint64) ([]*Edge, error)
	EdgesTo(nodeId uint64) ([]*Edge, error)
	AllEdges() ([]*Edge, error)
	RemoveEdge(id uint64) error

	Close() error
}

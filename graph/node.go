package graph

// Node represents a vertex in the graph.
//
// ObjectBox mapping (when using objectbox-go backend):
//
//	//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen
//	type Node struct {
//		Id         uint64
//		Label      string
//		Properties string                                          // JSON-encoded map
//		Embedding  []float32 `objectbox:"hnsw:dimensions=384"`   // HNSW index
//	}
type Node struct {
	Id         uint64
	Label      string
	Properties map[string]string
	// Embedding is the vector representation of this node.
	// In ObjectBoxStorage it is indexed with ObjectBox's HNSW algorithm
	// for O(log n) approximate nearest-neighbor search.
	// MemoryStorage falls back to an exact O(n) cosine similarity scan.
	Embedding []float32
}

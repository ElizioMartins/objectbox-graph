package graph

// Node represents a vertex in the graph.
//
// ObjectBox mapping (when using objectbox-go backend):
//
//	//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen
//	type Node struct {
//		Id         uint64
//		Label      string
//		Properties string // JSON-encoded map
//	}
type Node struct {
	Id         uint64
	Label      string
	Properties map[string]string
}

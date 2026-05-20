//go:build objectbox

// Package objectboxstorage provides an ObjectBox-backed Storage implementation
// for objectbox-graph.
//
// Prerequisites:
//
//	1. CGO enabled (requires GCC/MinGW-w64 on Windows)
//	2. ObjectBox C library installed:
//	     bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh)
//	   or on Windows:
//	     pwsh -File install.ps1
//	3. Add the Go module:
//	     go get github.com/objectbox/objectbox-go/objectbox
//
// After setup, generate the Box/Model code:
//
//	go generate ./objectboxstorage/...
//
// Then build and test:
//
//	go test -tags objectbox ./...
package objectboxstorage

// NodeEntity is the ObjectBox-persisted representation of a graph Node.
//
// Properties is stored as a JSON string because ObjectBox does not natively
// support map types — a deliberate trade-off for on-device compactness.
// For hot-path properties, promote them to dedicated fields.
//
//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen
type NodeEntity struct {
	Id         uint64
	Label      string
	Properties string // JSON-encoded map[string]string
}

// EdgeEntity is the ObjectBox-persisted representation of a directed,
// weighted graph edge.
//
// FromId and ToId are stored as plain uint64 foreign keys.
// They *could* be ToOne<NodeEntity> ObjectBox relations, but plain IDs
// let us query via EdgeEntity_.FromId.Equals() with zero relation overhead,
// which is crucial for deep traversals (BFS depth > 3).
//
// This design decision is the core of the graph layer's performance story:
// ObjectBox stores the edges; the graph layer owns the traversal.
type EdgeEntity struct {
	Id     uint64
	FromId uint64  // would be ToOne<NodeEntity> in a pure ObjectBox model
	ToId   uint64  // would be ToOne<NodeEntity> in a pure ObjectBox model
	Label  string
	Weight float64
}

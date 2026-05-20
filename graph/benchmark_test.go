package graph_test

import (
	"fmt"
	"testing"

	"github.com/ElizioMartins/objectbox-graph/graph"
)

// buildChainGraph creates: node0 → node1 → ... → node(n-1)
// Useful for measuring worst-case BFS/DFS/Dijkstra (linear topology).
func buildChainGraph(store *graph.GraphStore, n int) []uint64 {
	ids := make([]uint64, n)
	for i := 0; i < n; i++ {
		node, _ := store.AddNode("Node", map[string]string{"name": fmt.Sprintf("N%d", i)})
		ids[i] = node.Id
	}
	for i := 0; i < n-1; i++ {
		store.AddEdge(ids[i], ids[i+1], "link", 1.0)
	}
	return ids
}

// buildFanGraph creates a hub connected to fanOut leaves (star topology).
// Useful for measuring wide BFS (many neighbors at depth 1).
func buildFanGraph(store *graph.GraphStore, fanOut int) (hub uint64, leaves []uint64) {
	h, _ := store.AddNode("Hub", map[string]string{"name": "hub"})
	hub = h.Id
	leaves = make([]uint64, fanOut)
	for i := 0; i < fanOut; i++ {
		leaf, _ := store.AddNode("Leaf", map[string]string{"name": fmt.Sprintf("L%d", i)})
		store.AddEdge(hub, leaf.Id, "connects", 1.0)
		leaves[i] = leaf.Id
	}
	return
}

// --- BFS benchmarks ---

func BenchmarkBFS_Chain_100(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.BFS(ids[0], 100) //nolint:errcheck
	}
}

func BenchmarkBFS_Chain_1000(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.BFS(ids[0], 1000) //nolint:errcheck
	}
}

func BenchmarkBFS_Fan_1000(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	hub, _ := buildFanGraph(store, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.BFS(hub, 1) //nolint:errcheck
	}
}

// --- DFS benchmarks ---

func BenchmarkDFS_Chain_100(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.DFS(ids[0], 100) //nolint:errcheck
	}
}

func BenchmarkDFS_Chain_1000(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.DFS(ids[0], 1000) //nolint:errcheck
	}
}

// --- Dijkstra benchmarks ---

func BenchmarkDijkstra_Chain_100(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.ShortestPath(ids[0], ids[len(ids)-1]) //nolint:errcheck
	}
}

func BenchmarkDijkstra_Chain_1000(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	ids := buildChainGraph(store, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.ShortestPath(ids[0], ids[len(ids)-1]) //nolint:errcheck
	}
}

// --- AddNode / AddEdge throughput ---

func BenchmarkAddNode(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.AddNode("Node", map[string]string{"name": fmt.Sprintf("N%d", i)}) //nolint:errcheck
	}
}

func BenchmarkAddEdge(b *testing.B) {
	store := graph.New(graph.NewMemoryStorage())
	a, _ := store.AddNode("A", nil)
	z, _ := store.AddNode("Z", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.AddEdge(a.Id, z.Id, "link", 1.0) //nolint:errcheck
	}
}

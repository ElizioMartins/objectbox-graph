package main

import (
	"fmt"
	"log"

	"github.com/ElizioMartins/objectbox-graph/graph"
)

// Social recommendation graph example:
//
//	Alice ──[follows, w=1.0]──► Bob
//	Alice ──[follows, w=1.0]──► Carol
//	Bob   ──[follows, w=1.0]──► Dave
//	Carol ──[follows, w=2.0]──► Dave
//	Dave  ──[follows, w=1.0]──► Eve
//
// This demo uses MemoryStorage (objectbox-go integration is the next step).
func main() {
	store := graph.New(graph.NewMemoryStorage())
	defer store.Close()

	// --- Build the graph ---
	alice, _ := store.AddNode("Person", map[string]string{"name": "Alice"})
	bob, _ := store.AddNode("Person", map[string]string{"name": "Bob"})
	carol, _ := store.AddNode("Person", map[string]string{"name": "Carol"})
	dave, _ := store.AddNode("Person", map[string]string{"name": "Dave"})
	eve, _ := store.AddNode("Person", map[string]string{"name": "Eve"})

	store.AddEdge(alice.Id, bob.Id, "follows", 1.0)
	store.AddEdge(alice.Id, carol.Id, "follows", 1.0)
	store.AddEdge(bob.Id, dave.Id, "follows", 1.0)
	store.AddEdge(carol.Id, dave.Id, "follows", 2.0)
	store.AddEdge(dave.Id, eve.Id, "follows", 1.0)

	// --- BFS: Who does Alice reach within 2 hops? ---
	fmt.Println("=== BFS from Alice (depth=2) ===")
	bfsNodes, err := store.BFS(alice.Id, 2)
	if err != nil {
		log.Fatal(err)
	}
	for i, n := range bfsNodes {
		fmt.Printf("  [%d] %s (id=%d)\n", i, n.Properties["name"], n.Id)
	}

	// --- Dijkstra: Cheapest path from Alice to Eve ---
	fmt.Println("\n=== Shortest path: Alice \u2192 Eve ===")
	result, err := store.ShortestPath(alice.Id, eve.Id)
	if err != nil {
		log.Fatal(err)
	}
	if result == nil {
		fmt.Println("  No path found")
	} else {
		for i, n := range result.Nodes {
			if i > 0 {
				fmt.Print(" \u2192 ")
			}
			fmt.Print(n.Properties["name"])
		}
		fmt.Printf("\n  Total cost: %.1f\n", result.TotalCost)
	}

	// --- DFS from Bob ---
	fmt.Println("\n=== DFS from Bob (depth=3) ===")
	dfsNodes, _ := store.DFS(bob.Id, 3)
	for i, n := range dfsNodes {
		fmt.Printf("  [%d] %s\n", i, n.Properties["name"])
	}

	// Suppress unused variable warning
	_ = eve
}

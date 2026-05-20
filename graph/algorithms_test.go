package graph

import (
	"testing"
)

func setupSocialGraph(t *testing.T) (*GraphStore, map[string]uint64) {
	t.Helper()
	store := New(NewMemoryStorage())

	alice, _ := store.AddNode("Person", map[string]string{"name": "Alice"})
	bob, _ := store.AddNode("Person", map[string]string{"name": "Bob"})
	carol, _ := store.AddNode("Person", map[string]string{"name": "Carol"})
	dave, _ := store.AddNode("Person", map[string]string{"name": "Dave"})
	eve, _ := store.AddNode("Person", map[string]string{"name": "Eve"})

	//  Alice → Bob   (w=1)
	//  Alice → Carol (w=1)
	//  Bob   → Dave  (w=1)
	//  Carol → Dave  (w=2)
	//  Dave  → Eve   (w=1)
	store.AddEdge(alice.Id, bob.Id, "follows", 1.0)
	store.AddEdge(alice.Id, carol.Id, "follows", 1.0)
	store.AddEdge(bob.Id, dave.Id, "follows", 1.0)
	store.AddEdge(carol.Id, dave.Id, "follows", 2.0)
	store.AddEdge(dave.Id, eve.Id, "follows", 1.0)

	return store, map[string]uint64{
		"alice": alice.Id, "bob": bob.Id, "carol": carol.Id,
		"dave": dave.Id, "eve": eve.Id,
	}
}

func TestBFS_depth2(t *testing.T) {
	store, ids := setupSocialGraph(t)
	nodes, err := store.BFS(ids["alice"], 2)
	if err != nil {
		t.Fatal(err)
	}
	// Depth 0: Alice; Depth 1: Bob, Carol; Depth 2: Dave
	if len(nodes) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(nodes))
	}
}

func TestDFS_depth3(t *testing.T) {
	store, ids := setupSocialGraph(t)
	nodes, err := store.DFS(ids["bob"], 3)
	if err != nil {
		t.Fatal(err)
	}
	// Bob → Dave → Eve
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes (Bob, Dave, Eve), got %d", len(nodes))
	}
}

func TestShortestPath_AliceToEve(t *testing.T) {
	store, ids := setupSocialGraph(t)
	result, err := store.ShortestPath(ids["alice"], ids["eve"])
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("expected a path, got nil")
	}
	// Alice → Bob → Dave → Eve = cost 3.0
	// Alice → Carol → Dave → Eve = cost 4.0
	if result.TotalCost != 3.0 {
		t.Errorf("expected cost 3.0, got %.1f", result.TotalCost)
	}
	if len(result.Nodes) != 4 {
		t.Errorf("expected 4 nodes in path, got %d", len(result.Nodes))
	}
}

func TestShortestPath_NoPath(t *testing.T) {
	store, ids := setupSocialGraph(t)
	// Eve has no outgoing edges
	result, err := store.ShortestPath(ids["eve"], ids["alice"])
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Error("expected nil result for unreachable path")
	}
}

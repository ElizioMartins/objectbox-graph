package graph_test

import (
	"math"
	"testing"

	"github.com/ElizioMartins/objectbox-graph/graph"
)

// --- CosineSimilarity unit tests ---

func TestCosineSimilarity_Identical(t *testing.T) {
	a := []float32{1, 0, 0}
	if got := graph.CosineSimilarity(a, a); math.Abs(got-1.0) > 1e-6 {
		t.Errorf("identical vectors: want 1.0, got %f", got)
	}
}

func TestCosineSimilarity_Orthogonal(t *testing.T) {
	a := []float32{1, 0, 0}
	b := []float32{0, 1, 0}
	if got := graph.CosineSimilarity(a, b); math.Abs(got) > 1e-6 {
		t.Errorf("orthogonal: want 0.0, got %f", got)
	}
}

func TestCosineSimilarity_Opposite(t *testing.T) {
	a := []float32{1, 0, 0}
	c := []float32{-1, 0, 0}
	if got := graph.CosineSimilarity(a, c); math.Abs(got+1.0) > 1e-6 {
		t.Errorf("opposite vectors: want -1.0, got %f", got)
	}
}

func TestCosineSimilarity_ZeroVector(t *testing.T) {
	if got := graph.CosineSimilarity([]float32{}, []float32{}); got != 0 {
		t.Errorf("empty vectors: want 0, got %f", got)
	}
}

// --- Vector search tests ---

func setupVectorGraph(t *testing.T) (*graph.GraphStore, map[string]*graph.Node) {
	t.Helper()
	store := graph.New(graph.NewMemoryStorage())

	// 3D embeddings: [AI/ML relevance, Math relevance, Culinary relevance]
	concepts := map[string][]float32{
		"Machine Learning": {0.95, 0.70, 0.05},
		"Deep Learning":    {0.90, 0.65, 0.05},
		"Neural Networks":  {0.85, 0.55, 0.10},
		"Cooking":          {0.05, 0.10, 0.99},
		"Recipes":          {0.02, 0.05, 0.97},
	}
	nodes := make(map[string]*graph.Node)
	for name, emb := range concepts {
		n, _ := store.AddNode("Concept", map[string]string{"name": name})
		n.Embedding = emb
		store.UpdateNode(n) //nolint:errcheck
		nodes[name] = n
	}

	// ML → Deep Learning → Neural Networks
	store.AddEdge(nodes["Machine Learning"].Id, nodes["Deep Learning"].Id, "includes", 1.0)   //nolint:errcheck
	store.AddEdge(nodes["Deep Learning"].Id, nodes["Neural Networks"].Id, "uses", 1.0)        //nolint:errcheck
	store.AddEdge(nodes["Cooking"].Id, nodes["Recipes"].Id, "contains", 1.0)                  //nolint:errcheck
	return store, nodes
}

func TestFindSimilarNodes_ReturnsTopK(t *testing.T) {
	store, _ := setupVectorGraph(t)
	query := []float32{0.92, 0.68, 0.05} // clearly in ML space
	results, err := store.FindSimilarNodes(query, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
	for _, r := range results {
		name := r.Node.Properties["name"]
		if name == "Cooking" || name == "Recipes" {
			t.Errorf("food concepts should not appear in ML top-2, got: %s", name)
		}
	}
}

func TestFindSimilarNodes_ScoresDescending(t *testing.T) {
	store, _ := setupVectorGraph(t)
	results, _ := store.FindSimilarNodes([]float32{0.92, 0.68, 0.05}, 5, 0)
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("scores not descending at index %d: %.4f > %.4f",
				i, results[i].Score, results[i-1].Score)
		}
	}
}

func TestFindSimilarNodes_MinScore(t *testing.T) {
	store, _ := setupVectorGraph(t)
	results, _ := store.FindSimilarNodes([]float32{0.92, 0.68, 0.05}, 10, 0.90)
	for _, r := range results {
		if r.Score < 0.90 {
			t.Errorf("result below minScore 0.90: %s score=%.4f",
				r.Node.Properties["name"], r.Score)
		}
	}
}

func TestVectorTraversal_ExpandsContext(t *testing.T) {
	store, nodes := setupVectorGraph(t)
	query := []float32{0.94, 0.70, 0.05} // closest to Machine Learning

	contextNodes, err := store.VectorTraversal(query, 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	// VectorSearch finds ML (top-1), BFS depth=2:
	// ML → Deep Learning → Neural Networks
	foundML, foundDL, foundNN := false, false, false
	for _, n := range contextNodes {
		switch n.Id {
		case nodes["Machine Learning"].Id:
			foundML = true
		case nodes["Deep Learning"].Id:
			foundDL = true
		case nodes["Neural Networks"].Id:
			foundNN = true
		}
	}
	if !foundML {
		t.Error("expected Machine Learning in traversal result")
	}
	if !foundDL {
		t.Error("expected Deep Learning in traversal (1 hop from ML)")
	}
	if !foundNN {
		t.Error("expected Neural Networks in traversal (2 hops from ML)")
	}
}

func TestVectorTraversal_IsolatesUnrelatedClusters(t *testing.T) {
	store, nodes := setupVectorGraph(t)
	query := []float32{0.94, 0.70, 0.05} // ML space

	contextNodes, _ := store.VectorTraversal(query, 1, 2)

	// Cooking and Recipes should NOT appear — different cluster, not reachable
	for _, n := range contextNodes {
		if n.Id == nodes["Cooking"].Id || n.Id == nodes["Recipes"].Id {
			t.Errorf("food cluster leaked into ML traversal: %s", n.Properties["name"])
		}
	}
}

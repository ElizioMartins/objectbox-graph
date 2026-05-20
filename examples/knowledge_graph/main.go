package main

import (
	"fmt"
	"math"

	"github.com/ElizioMartins/objectbox-graph/graph"
)

// Knowledge graph modeling an AI/ML concept hierarchy.
//
// Each node has a 4-dimensional mock embedding representing its position in:
//   [AI/ML relevance, Math relevance, Programming relevance, Biology relevance]
//
// In production, replace mockEmbed() with a real local model:
//   - ONNX runtime + all-MiniLM-L6-v2  (384 dims, runs on Raspberry Pi)
//   - llama.cpp embeddings              (offline, mobile-capable)
//   - ObjectBox's planned on-device embedding support
//
// Storage: MemoryStorage (no persistence).
// In Phase 2+, swap for objectboxstorage.New(obxStore) to persist the
// knowledge graph across device restarts — the key advantage over Neo4j
// for edge/IoT deployments.
func main() {
	store := graph.New(graph.NewMemoryStorage())
	defer store.Close()

	// --- Build knowledge graph ---

	concepts := []string{
		"Machine Learning", "Deep Learning", "Neural Networks", "Backpropagation",
		"Linear Algebra", "Gradient Descent", "Python", "TensorFlow",
		"Neuroscience", "Computer Vision",
	}
	nodes := make(map[string]*graph.Node)
	for _, c := range concepts {
		n, _ := store.AddNode("Concept", map[string]string{"name": c})
		n.Embedding = mockEmbed(c)
		store.UpdateNode(n) //nolint:errcheck
		nodes[c] = n
	}

	relations := [][3]string{
		{"Deep Learning", "Machine Learning", "is_subset_of"},
		{"Neural Networks", "Deep Learning", "powers"},
		{"Backpropagation", "Neural Networks", "trains"},
		{"Gradient Descent", "Backpropagation", "used_by"},
		{"Linear Algebra", "Neural Networks", "required_by"},
		{"Linear Algebra", "Gradient Descent", "required_by"},
		{"Python", "TensorFlow", "language_of"},
		{"TensorFlow", "Deep Learning", "implements"},
		{"Neuroscience", "Neural Networks", "inspired"},
		{"Computer Vision", "Deep Learning", "uses"},
	}
	for _, r := range relations {
		store.AddEdge(nodes[r[0]].Id, nodes[r[1]].Id, r[2], 1.0) //nolint:errcheck
	}

	// --- Query 1: Vector Search ---
	fmt.Println("=== Vector Search: most similar to \"Backpropagation\" ===")
	query := mockEmbed("Backpropagation")
	similar, _ := store.FindSimilarNodes(query, 4, 0)
	for i, s := range similar {
		fmt.Printf("  [%d] %-22s score=%.4f\n", i+1, s.Node.Properties["name"], s.Score)
	}

	// --- Query 2: Vector + Graph Traversal (on-device RAG context) ---
	fmt.Println("\n=== Vector + Graph: context for \"I want to understand neural networks\" ===")
	// Query embedding similar to Neural Networks
	learnQuery := normalize([]float32{0.87, 0.58, 0.50, 0.28})
	contextNodes, _ := store.VectorTraversal(learnQuery, 2, 2)
	fmt.Println("  RAG context nodes retrieved on-device:")
	seen := make(map[string]bool)
	for _, n := range contextNodes {
		name := n.Properties["name"]
		if !seen[name] {
			seen[name] = true
			fmt.Printf("    • %s\n", name)
		}
	}

	// --- Query 3: Shortest Path ---
	fmt.Println("\n=== Graph: How does Python connect to Machine Learning? ===")
	result, _ := store.ShortestPath(nodes["Python"].Id, nodes["Machine Learning"].Id)
	if result != nil {
		for i, n := range result.Nodes {
			if i > 0 {
				fmt.Print(" \u2192 ")
			}
			fmt.Print(n.Properties["name"])
		}
		fmt.Printf("\n  Hops: %d\n", len(result.Nodes)-1)
	}

	// --- Query 4: BFS from a concept ---
	fmt.Println("\n=== BFS from \"Linear Algebra\" (depth=2): what does it touch? ===")
	bfsNodes, _ := store.BFS(nodes["Linear Algebra"].Id, 2)
	for _, n := range bfsNodes {
		fmt.Printf("  • %s\n", n.Properties["name"])
	}
}

// mockEmbed returns a 4-dimensional normalized embedding for a concept.
// Dimensions: [AI/ML, Math, Programming, Biology]
func mockEmbed(concept string) []float32 {
	raw := map[string][]float32{
		"Machine Learning":  {0.95, 0.70, 0.60, 0.10},
		"Deep Learning":     {0.92, 0.65, 0.55, 0.15},
		"Neural Networks":   {0.90, 0.60, 0.50, 0.30},
		"Backpropagation":   {0.80, 0.85, 0.70, 0.05},
		"Linear Algebra":    {0.40, 0.95, 0.60, 0.10},
		"Gradient Descent":  {0.75, 0.90, 0.65, 0.05},
		"Python":            {0.30, 0.30, 0.98, 0.10},
		"TensorFlow":        {0.85, 0.50, 0.90, 0.05},
		"Neuroscience":      {0.50, 0.40, 0.10, 0.95},
		"Computer Vision":   {0.88, 0.60, 0.70, 0.20},
	}
	if v, ok := raw[concept]; ok {
		return normalize(v)
	}
	return nil
}

func normalize(v []float32) []float32 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	norm := float32(math.Sqrt(sum))
	if norm == 0 {
		return v
	}
	result := make([]float32, len(v))
	for i, x := range v {
		result[i] = x / norm
	}
	return result
}

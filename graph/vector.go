package graph

import (
	"math"
	"sort"
)

// ScoredNode pairs a Node with its vector similarity score.
type ScoredNode struct {
	Node  *Node
	Score float64 // cosine similarity in [-1, 1]; higher is more similar
}

// VectorBackend is an optional capability extension for Storage backends
// that support native approximate nearest-neighbor (ANN) search.
//
// If the storage backend implements this interface, FindSimilarNodes delegates
// to it instead of the O(n) linear cosine-scan fallback.
//
// ObjectBoxStorage implements VectorBackend via ObjectBox's HNSW index,
// reducing vector search from O(n) to O(log n) — the same advantage that
// pgvector brings to PostgreSQL.
type VectorBackend interface {
	NearestNeighbors(query []float32, maxResults int) ([]*ScoredNode, error)
}

// CosineSimilarity computes the cosine similarity between two float32 vectors.
// Returns values in [-1, 1]. Returns 0 for zero-length or mismatched vectors.
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// FindSimilarNodes finds nodes whose Embedding is most similar to the query vector.
//
// If the storage backend implements VectorBackend (e.g., ObjectBoxStorage with
// ObjectBox HNSW index), the query is O(log n).
// Otherwise, falls back to an exact O(n) cosine similarity scan over AllNodes.
//
// minScore filters out results below the threshold (use 0 to return all).
func (g *GraphStore) FindSimilarNodes(query []float32, maxResults int, minScore float64) ([]*ScoredNode, error) {
	// Capability check: use HNSW if the backend supports it.
	// ObjectBoxStorage satisfies VectorBackend; MemoryStorage does not.
	if vb, ok := g.storage.(VectorBackend); ok {
		return vb.NearestNeighbors(query, maxResults)
	}

	// Fallback: O(n) exact scan.
	nodes, err := g.storage.AllNodes()
	if err != nil {
		return nil, err
	}
	var scored []*ScoredNode
	for _, n := range nodes {
		if len(n.Embedding) == 0 {
			continue
		}
		score := CosineSimilarity(query, n.Embedding)
		if score >= minScore {
			scored = append(scored, &ScoredNode{Node: n, Score: score})
		}
	}
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	if maxResults > 0 && len(scored) > maxResults {
		scored = scored[:maxResults]
	}
	return scored, nil
}

// VectorTraversal finds the top-k most similar nodes to the query embedding,
// then expands each via BFS up to depth levels.
//
// This is the core primitive for on-device RAG with graph context:
//
//  1. Vector search finds semantically similar entry-point nodes
//  2. BFS from each entry-point retrieves their neighborhood
//  3. The union is a rich context window — no cloud API required
//
// Example: knowledge graph where nodes have embeddings from a local
// sentence-transformer. The query "how does backpropagation work?" finds
// "Backpropagation" via vector similarity, then BFS retrieves
// "Gradient Descent", "Neural Networks", "Loss Function" — all on-device.
func (g *GraphStore) VectorTraversal(query []float32, topK int, depth int) ([]*Node, error) {
	similar, err := g.FindSimilarNodes(query, topK, 0)
	if err != nil {
		return nil, err
	}
	seen := make(map[uint64]bool)
	var result []*Node
	for _, s := range similar {
		bfsNodes, err := g.BFS(s.Node.Id, depth)
		if err != nil {
			continue
		}
		for _, n := range bfsNodes {
			if !seen[n.Id] {
				seen[n.Id] = true
				result = append(result, n)
			}
		}
	}
	return result, nil
}

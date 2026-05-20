//go:build objectbox

package objectboxstorage

import (
	"github.com/ElizioMartins/objectbox-graph/graph"
)

// NearestNeighbors implements graph.VectorBackend using ObjectBox's HNSW index.
//
// This makes FindSimilarNodes O(log n) instead of the O(n) linear scan used
// by MemoryStorage — the same performance advantage that pgvector brings
// to PostgreSQL for vector similarity search.
//
// The graph layer detects this capability automatically via interface check:
//
//	if vb, ok := g.storage.(VectorBackend); ok {
//	    return vb.NearestNeighbors(query, maxResults)  // O(log n)
//	}
//	// else: O(n) scan fallback
func (s *ObjectBoxStorage) NearestNeighbors(query []float32, maxResults int) ([]*graph.ScoredNode, error) {
	// ObjectBox HNSW query:
	// - NearestNeighbors(vector, count) finds approximate nearest neighbors
	// - WithScores() attaches the Euclidean distance to each result
	q, err := s.nodeBox.Query(
		NodeEntity_.Embedding.NearestNeighbors(query, uint32(maxResults)).WithScores(),
	)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// FindWithScores returns []objectbox.EntityWithScore[NodeEntity]
	results, err := q.FindWithScores()
	if err != nil {
		return nil, err
	}

	scored := make([]*graph.ScoredNode, len(results))
	for i, r := range results {
		// ObjectBox returns squared Euclidean distance.
		// For unit-normalized embeddings: similarity ≈ 1 - (distance² / 2)
		similarity := 1.0 - float64(r.Score)/2.0
		scored[i] = &graph.ScoredNode{
			Node:  nodeEntityToGraph(r.Object),
			Score: similarity,
		}
	}
	return scored, nil
}

package graph

import (
	"container/heap"
	"math"
)

// BFS performs a breadth-first traversal starting from startId,
// visiting nodes up to maxDepth levels deep.
// Returns nodes in visit order (level by level).
func (g *GraphStore) BFS(startId uint64, maxDepth int) ([]*Node, error) {
	type item struct {
		id    uint64
		depth int
	}
	visited := make(map[uint64]bool)
	queue := []item{{startId, 0}}
	var result []*Node

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if visited[cur.id] || cur.depth > maxDepth {
			continue
		}
		visited[cur.id] = true

		node, err := g.storage.GetNode(cur.id)
		if err != nil {
			continue
		}
		result = append(result, node)

		if cur.depth < maxDepth {
			edges, err := g.storage.EdgesFrom(cur.id)
			if err != nil {
				continue
			}
			for _, e := range edges {
				if !visited[e.ToId] {
					queue = append(queue, item{e.ToId, cur.depth + 1})
				}
			}
		}
	}
	return result, nil
}

// DFS performs a depth-first traversal starting from startId,
// visiting nodes up to maxDepth levels deep.
func (g *GraphStore) DFS(startId uint64, maxDepth int) ([]*Node, error) {
	visited := make(map[uint64]bool)
	var result []*Node
	g.dfsHelper(startId, 0, maxDepth, visited, &result)
	return result, nil
}

func (g *GraphStore) dfsHelper(id uint64, depth, maxDepth int, visited map[uint64]bool, result *[]*Node) {
	if visited[id] || depth > maxDepth {
		return
	}
	visited[id] = true

	node, err := g.storage.GetNode(id)
	if err != nil {
		return
	}
	*result = append(*result, node)

	edges, err := g.storage.EdgesFrom(id)
	if err != nil {
		return
	}
	for _, e := range edges {
		g.dfsHelper(e.ToId, depth+1, maxDepth, visited, result)
	}
}

// PathResult holds the result of a shortest path query.
type PathResult struct {
	Nodes     []*Node
	TotalCost float64
}

// ShortestPath finds the minimum-cost path between two nodes using Dijkstra's algorithm.
// Edge weights are used as costs. Returns nil if no path exists.
func (g *GraphStore) ShortestPath(fromId, toId uint64) (*PathResult, error) {
	dist := map[uint64]float64{fromId: 0}
	prev := make(map[uint64]uint64)

	pq := &priorityQueue{}
	heap.Push(pq, &pqItem{id: fromId, cost: 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem)

		if cur.id == toId {
			break
		}
		if cur.cost > dist[cur.id] {
			continue // stale entry
		}

		edges, err := g.storage.EdgesFrom(cur.id)
		if err != nil {
			continue
		}
		for _, e := range edges {
			newCost := dist[cur.id] + e.Weight
			if d, ok := dist[e.ToId]; !ok || newCost < d {
				dist[e.ToId] = newCost
				prev[e.ToId] = cur.id
				heap.Push(pq, &pqItem{id: e.ToId, cost: newCost})
			}
		}
	}

	if _, ok := dist[toId]; !ok {
		return nil, nil // no path found
	}

	// Reconstruct path by walking backwards through prev
	var ids []uint64
	for cur := toId; ; cur = prev[cur] {
		ids = append([]uint64{cur}, ids...)
		if cur == fromId {
			break
		}
		if _, ok := prev[cur]; !ok {
			break
		}
	}

	var nodes []*Node
	for _, id := range ids {
		n, err := g.storage.GetNode(id)
		if err != nil {
			continue
		}
		nodes = append(nodes, n)
	}

	totalCost := dist[toId]
	if math.IsInf(totalCost, 1) {
		totalCost = 0
	}
	return &PathResult{Nodes: nodes, TotalCost: totalCost}, nil
}

// --- Priority Queue for Dijkstra ---

type pqItem struct {
	id   uint64
	cost float64
	idx  int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx = i
	pq[j].idx = j
}
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.idx = n
	*pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

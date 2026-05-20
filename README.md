# objectbox-graph

> A native graph layer built on top of [ObjectBox](https://objectbox.io/) — bringing BFS, DFS, weighted edges and shortest path to the world's fastest on-device database.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-4%20passing-brightgreen)](#running-tests)

---

## Motivation

ObjectBox is the fastest on-device database for mobile and IoT — but it lacks a native **property graph model**.

This project adds graph capabilities as a **layer on top of ObjectBox**, following the same pattern used by:

- **pgvector** — adds vector search to PostgreSQL without changing the core
- **MongoDB Atlas Vector Search** — extends the document model with embeddings

The result: the **only graph + vector database optimized for edge, IoT and mobile** devices.

```
Your App
   │
   ▼
objectbox-graph          ← this library
  Graph API (BFS, DFS, Dijkstra, ...)
  Storage Interface
   │
   ▼
ObjectBox (objectbox-go) ← storage backend (phase 2)
  Ultra-fast on-device object store
  On-device Vector Search (HNSW)
```

---

## Features

- **Node** and **Edge** entities with `Label`, `Properties` and `Weight`
- **Directed, weighted graph** — edges map directly to ObjectBox `ToOne<Node>` relations
- **BFS** — breadth-first traversal with configurable depth
- **DFS** — depth-first traversal with configurable depth
- **Dijkstra** — shortest weighted path between two nodes
- **Storage interface** — swap ObjectBox for any backend (in-memory included for testing)

---

## Quick Start

```go
import "github.com/ElizioMartins/objectbox-graph/graph"

// Create a graph store (in-memory backend for now)
store := graph.New(graph.NewMemoryStorage())
defer store.Close()

// Add nodes
alice, _ := store.AddNode("Person", map[string]string{"name": "Alice"})
bob, _   := store.AddNode("Person", map[string]string{"name": "Bob"})
eve, _   := store.AddNode("Person", map[string]string{"name": "Eve"})

// Add directed, weighted edges
store.AddEdge(alice.Id, bob.Id, "follows", 1.0)
store.AddEdge(bob.Id,   eve.Id, "follows", 1.0)

// BFS — who does Alice reach in 2 hops?
nodes, _ := store.BFS(alice.Id, 2)

// Shortest path — Alice to Eve
result, _ := store.ShortestPath(alice.Id, eve.Id)
fmt.Printf("Cost: %.1f\n", result.TotalCost) // Cost: 2.0
```

---

## Running the Example

```bash
git clone https://github.com/ElizioMartins/objectbox-graph.git
cd objectbox-graph
go run ./examples/social/main.go
```

Expected output:
```
=== BFS from Alice (depth=2) ===
  [0] Alice (id=1)
  [1] Bob (id=2)
  [2] Carol (id=3)
  [3] Dave (id=4)

=== Shortest path: Alice → Eve ===
Alice → Bob → Dave → Eve
  Total cost: 3.0

=== DFS from Bob (depth=3) ===
  [0] Bob
  [1] Dave
  [2] Eve
```

---

## Running Tests

```bash
go test ./graph/... -v
```

```
=== RUN   TestBFS_depth2
--- PASS: TestBFS_depth2 (0.00s)
=== RUN   TestDFS_depth3
--- PASS: TestDFS_depth3 (0.00s)
=== RUN   TestShortestPath_AliceToEve
--- PASS: TestShortestPath_AliceToEve (0.00s)
=== RUN   TestShortestPath_NoPath
--- PASS: TestShortestPath_NoPath (0.00s)
PASS
ok  github.com/ElizioMartins/objectbox-graph/graph
```

---

## Benchmarks (MemoryStorage baseline)

> Machine: Intel i5-10400 @ 2.90GHz · Windows 11 · Go 1.26.3
> Command: `go test ./graph/... -bench=Benchmark -benchmem`

| Benchmark | ns/op | B/op | allocs/op | ops/sec |
|---|---:|---:|---:|---:|
| `AddNode` | 647 | 461 | 6 | **1.5M/s** |
| `AddEdge` | 307 | 161 | 2 | **4.6M/s** |
| `BFS` — chain 100 nodes | 109,860 | 8,913 | 210 | 9,100/s |
| `BFS` — chain 1,000 nodes | 11,651,271 | 115,680 | 2,024 | 86/s |
| `BFS` — fan 1,000 leaves | 188,859 | 159,577 | 46 | 5,300/s |
| `DFS` — chain 100 nodes | 103,584 | 7,416 | 116 | 9,600/s |
| `DFS` — chain 1,000 nodes | 11,448,620 | 99,785 | 1,030 | 87/s |
| `Dijkstra` — chain 100 nodes | 139,346 | 57,352 | 427 | 7,200/s |
| `Dijkstra` — chain 1,000 nodes | 13,904,844 | **4,480,076** | 4,052 | 72/s |

### Key insight from the benchmarks

The `Dijkstra` 1,000-node row allocates **4.4 MB per call** — because `MemoryStorage`
holds the entire graph in RAM and the priority queue builds large intermediate structures.

With `ObjectBoxStorage` (phase 2), only the nodes actually *visited* during
traversal are loaded from disk. For sparse graphs typical of IoT topologies,
this is expected to reduce per-traversal allocations by **10×–50×**,
while persisting the graph across device restarts — something Neo4j cannot do on a 512 MB device.

---

## How it Maps to ObjectBox Relations

The key insight is that ObjectBox's `ToOne` relations are **edges without properties**.
This library adds `Label` and `Weight` to those edges, turning them into a full property graph:

| Graph Concept | ObjectBox today | objectbox-graph |
|---|---|---|
| Node | `@Entity` | `Node{Id, Label, Properties}` |
| Simple edge | `ToOne<T>` | ✅ covered |
| **Property edge** | ❌ not supported | `Edge{FromId, ToId, Label, Weight}` |
| Traversal | `.link()` (1 level) | BFS / DFS (any depth) |
| Shortest path | ❌ | Dijkstra ✅ |
| Vector + Graph | ❌ | Phase 3 (planned) |

---

## Roadmap

See [ROADMAP.md](ROADMAP.md) for the full plan.

**Phase 1** (done): Graph layer with in-memory storage, 4 passing tests, benchmarks 
**Phase 2** (in progress): ObjectBox backend — `objectboxstorage/` package ready, pending CGO setup 
**Phase 3**: Vector + Graph — combine ObjectBox HNSW vector index with graph traversal for on-device RAG 
**Phase 4**: Proposal to ObjectBox team 

---

## License

Apache 2.0 — same as ObjectBox.

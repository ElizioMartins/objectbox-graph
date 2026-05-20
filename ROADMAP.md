# objectbox-graph Roadmap

## Phase 1 — Graph Layer (current)

- [x] `Node` and `Edge` entities with label and weight
- [x] `Storage` interface (decoupled from ObjectBox)
- [x] `MemoryStorage` — in-memory backend for dev/test
- [x] `GraphStore` API: `AddNode`, `AddEdge`, `Neighbors`
- [x] Algorithms: BFS, DFS, Dijkstra (shortest path)
- [x] Unit tests
- [x] Social graph example

## Phase 2 — ObjectBox Backend

- [ ] `ObjectBoxStorage` implementing `Storage` interface
- [ ] `Node` and `Edge` as proper objectbox-go entities
- [ ] `ToOne<Node>` for `FromId` and `ToId` on Edge
- [ ] On-device persistence (Android, IoT, desktop)
- [ ] Benchmark: objectbox-graph vs Neo4j (embedded)

## Phase 3 — Vector + Graph (the unique differentiator)

- [ ] Add `Embedding []float32` field to `Node`
- [ ] Vector similarity search using ObjectBox HNSW index
- [ ] `VectorGraph` queries: find similar nodes + traverse neighbors
- [ ] On-device RAG with graph context (no cloud required)
- [ ] Example: knowledge graph + semantic search on mobile

## Phase 4 — Proposal to ObjectBox Team

- [ ] Open GitHub Issue on objectbox-go with POC link
- [ ] Write RFC document
- [ ] Contact: contact@objectbox.io

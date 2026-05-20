package graph

// Edge represents a directed, weighted connection between two Nodes.
//
// ObjectBox mapping (when using objectbox-go backend):
//
//	FromId and ToId become ToOne<Node> relations:
//
//		type Edge struct {
//			Id     uint64
//			From   objectbox.ToOne `objectbox:"link"`
//			To     objectbox.ToOne `objectbox:"link"`
//			Label  string
//			Weight float64
//		}
//
// This is the key structural gap vs. ObjectBox's built-in relations:
// native edges here carry Label and Weight, enabling a true property graph.
type Edge struct {
	Id     uint64
	FromId uint64  // maps to ToOne<Node> in ObjectBox
	ToId   uint64  // maps to ToOne<Node> in ObjectBox
	Label  string
	Weight float64
}

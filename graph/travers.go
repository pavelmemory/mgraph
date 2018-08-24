package graph

// Describes rules how vertexes in graph must be searched
type TraversingStrategy interface {
	// Starts searching from provided vertex
	Search(vtx *Vertex) []*Vertex
}

// Search throw the graph with provided `strategy` and apply `action` to found data
func (vtx *Vertex) TraverseWith(strategy TraversingStrategy, action Action) (err error) {
	fvtxs := strategy.Search(vtx)
	for _, fvtx := range fvtxs {
		if err = action(fvtx); err != nil {
			return
		}
	}
	return
}

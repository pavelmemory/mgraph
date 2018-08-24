package graph

// Search type
type SearchAlgorithm int

const (
	// Depth-First Search
	DFS = SearchAlgorithm(iota)
	// Breadth-First Search
	BFS
)

func (sd SearchAlgorithm) StartAt(vtx *Vertex) GraphIterator {
	switch sd {
	case DFS:
		dfs := &dfSearcherWrapper{}
		dfs.searcher = &dfSearcher{vtx: vtx, vtxs: vtx.Outcoming().Vertexes(), wrapper: dfs}
		return dfs
	case BFS:
		return &bfSearcher{check: []*Vertex{vtx}}
	default:
		return badSearcher{}
	}
}

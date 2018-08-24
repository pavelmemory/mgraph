package graph

// Iterates by all vertexes of the graph
type GraphIterator interface {
	// reports if call to `Next` method will return next vertex of the graph
	HasNext() bool
	// returns next vertex of the graph if method `HasNext` returned `true` or panics otherwise
	Next() *Vertex
}

// Depth-First Search implementation structs
type dfSearcherWrapper struct {
	searcher *dfSearcher
}

func (dfsw *dfSearcherWrapper) Next() *Vertex {
	return dfsw.searcher.Next()
}

type dfSearcher struct {
	vtx      *Vertex
	vtxs     []*Vertex
	index    int
	previous *dfSearcher
	wrapper  *dfSearcherWrapper
}

func (dfs *dfSearcher) Next() *Vertex {
	if dfs.index < len(dfs.vtxs) {
		vtx := dfs.vtxs[dfs.index]
		dfs.index++
		dfs.wrapper.searcher = &dfSearcher{vtx: vtx, previous: dfs, vtxs: vtx.Outcoming().Vertexes(), wrapper: dfs.wrapper}
		return dfs.wrapper.Next()
	} else {
		vtx := dfs.vtx
		if dfs.previous != nil {
			dfs.wrapper.searcher = dfs.previous
		} else {
			dfs.wrapper.searcher = nil
		}
		return vtx
	}
}

func (dfs *dfSearcherWrapper) HasNext() bool {
	return dfs.searcher != nil
}

// Breadth-First Search implementation structs
type bfSearcher struct {
	check []*Vertex
	index int
}

func (bfs *bfSearcher) Next() *Vertex {
	vtx := bfs.check[bfs.index]
	bfs.index++
	bfs.check = append(bfs.check, vtx.Outcoming().Vertexes()...)
	return vtx
}

func (bfs *bfSearcher) HasNext() bool {
	return bfs.index < len(bfs.check)
}

// Stub implementation used when invalid value of `SearchDirection` type was used
type badSearcher struct{}

func (badSearcher) Next() *Vertex {
	return nil
}

func (badSearcher) HasNext() bool {
	return false
}

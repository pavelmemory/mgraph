package graph

// returns first N found vertexes
func FindN(algorithm SearchAlgorithm, number uint, isFound VertexPredicate) TraversingStrategy {
	return &findNFirstSearch{algorithm: algorithm, number: number, isFound: isFound}
}

// returns first found vertex
func FindFirst(algorithm SearchAlgorithm, isFound VertexPredicate) TraversingStrategy {
	return FindN(algorithm, 1, isFound)
}

type findNFirstSearch struct {
	algorithm SearchAlgorithm
	number    uint
	isFound   VertexPredicate
}

func (nfs *findNFirstSearch) Search(vtx *Vertex) (found []*Vertex) {
	for iterator := nfs.algorithm.StartAt(vtx); nfs.number != 0 && iterator.HasNext(); {
		nVtx := iterator.Next()
		if nfs.isFound(nVtx) {
			found = append(found, nVtx)
			nfs.number--
		}
	}
	return
}

// returns all vertexes that suits to predicate
func FindAll(isFound VertexPredicate) TraversingStrategy {
	return findAll(isFound)
}

type findAll VertexPredicate

func (isFound findAll) Search(vtx *Vertex) (found []*Vertex) {
	for iterator := BFS.StartAt(vtx); iterator.HasNext(); {
		nVtx := iterator.Next()
		if isFound(nVtx) {
			found = append(found, nVtx)
		}
	}
	return
}

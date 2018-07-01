package graph

type (
	DataPredicate  func(d interface{}) bool
	VertexPredicate   func(vtx *Vertex) bool
	EdgePredicate   func(edge *Edge) bool
	Action      func(vtx *Vertex) error
	GroupEdgesAction func(groupKey []byte, edges []*Edge) error
	GroupVertexesAction func(groupKey []byte, vtxs []*Vertex) error
	EdgesGrouper     func(edge *Edge) []byte
	VertexesGrouper     func(vtx *Vertex) []byte
	VertexSelector func (vtx *Vertex) bool

	GroupedEdges struct {
		GroupKey []byte
		Edges    []*Edge
	}

	GroupedVertexes struct {
		GroupKey []byte
		Vertexes    []*Vertex
	}
)

func GoOverEdge(edgeSelector EdgePredicate) PathOverEdge {
	return PathOverEdge{selectors: []EdgePredicate{edgeSelector}}
}

type PathOverEdge struct{
	selectors []EdgePredicate
}

func (poe PathOverEdge) GoOverEdge(edgeSelector EdgePredicate) PathOverEdge {
	poe.selectors = append(poe.selectors, edgeSelector)
	return poe
}

func (poe PathOverEdge) GroupVertexesWith(vtxGrouper VertexesGrouper) CompleteGrouperOverEdgesPath {
	return CompleteGrouperOverEdgesPath{pathOverEdge: poe, vtxGrouper: vtxGrouper}
}

type CompleteGrouperOverEdgesPath struct {
	pathOverEdge PathOverEdge
	vtxGrouper VertexesGrouper
}

func EdgeAttributeEqualsTo(data interface{}) EdgePredicate {
	return func(edge *Edge) bool {
		return edge.attributes == data
	}
}

type Edge struct {
	vertex     *Vertex
	attributes interface{}
}

func (edge *Edge) Vertex() *Vertex {
	return edge.vertex
}

func (edge *Edge) Attributes() interface{} {
	return edge.attributes
}

type Vertex struct {
	incoming, outcoming EdgeSet
	adjacent            VertexSet
	data                interface{}
}

func VertexWith(data interface{}) *Vertex {
	return &Vertex{
		data:      data,
		incoming:  NewEdgeSet(),
		outcoming: NewEdgeSet(),
		adjacent:  NewVertexSet()}
}

func (vtx *Vertex) Data() interface{} {
	return vtx.data
}

func (fromVtx *Vertex) EdgesTo(toVtx *Vertex) EdgeSet {
	es := NewEdgeSet()
	for iterator := fromVtx.outcoming.Iterator(); iterator.HasNext(); {
		edge := iterator.Next()
		if edge.vertex == toVtx {
			es.put(edge)
		}
	}
	return es
}

func (fromVtx *Vertex) EdgeTo(toVtx *Vertex) *Vertex {
	fromVtx.EdgeToWith(toVtx, nil)
	return fromVtx
}

func (fromVtx *Vertex) EdgeToWith(toVtx *Vertex, attributes interface{}) *Vertex {
	fromVtx.outcoming.put(&Edge{vertex: toVtx, attributes: attributes})
	fromVtx.adjacent.put(toVtx)

	toVtx.incoming.put(&Edge{vertex: fromVtx, attributes: attributes})
	toVtx.adjacent.put(fromVtx)
	return fromVtx
}

func (oneVtx *Vertex) Edge(anotherVtx *Vertex) *Vertex {
	oneVtx.EdgeWith(anotherVtx, nil)
	return oneVtx
}

func (oneVtx *Vertex) EdgeWith(anotherVtx *Vertex, attributes interface{}) *Vertex {
	oneVtx.EdgeToWith(anotherVtx, attributes)
	anotherVtx.EdgeToWith(oneVtx, attributes)
	return oneVtx
}

func (vtx *Vertex) Outcoming() EdgeSet {
	return vtx.OutcomingWhich(nil)
}

func (vtx *Vertex) OutcomingWhich(predicate EdgePredicate) EdgeSet {
	return vtx.directedWhich(vtx.outcoming, predicate)
}

func (vtx *Vertex) Incoming() EdgeSet {
	return vtx.incoming
}

func (vtx *Vertex) IncomingWhich(predicate EdgePredicate) EdgeSet {
	return vtx.directedWhich(vtx.incoming, predicate)
}

func (vtx *Vertex) directedWhich(es EdgeSet, predicate EdgePredicate) EdgeSet {
	if predicate == nil {
		return es
	}
	esw := NewEdgeSet()
	for iterator := es.Iterator(); iterator.HasNext(); {
		edge := iterator.Next()
		if predicate(edge) {
			esw.put(edge)
		}
	}
	return esw
}

func (vtx *Vertex) Adjacent() VertexSet {
	return vtx.adjacent
}

func (vtx *Vertex) GroupedBy(defineGroup EdgesGrouper) (groups []GroupedEdges) {
	_ = vtx.GroupBy(defineGroup, func(groupKey []byte, edges []*Edge) error {
		groups = append(groups, GroupedEdges{GroupKey: groupKey, Edges: edges})
		return nil
	})
	return
}

func (vtx *Vertex) GroupedOutcomingBy(defineGroup EdgesGrouper) (groups []GroupedEdges) {
	_ = vtx.GroupOutcomingBy(defineGroup, func(groupKey []byte, edges []*Edge) error {
		groups = append(groups, GroupedEdges{GroupKey: groupKey, Edges: edges})
		return nil
	})
	return
}

func (vtx *Vertex) GroupedIncomingBy(defineGroup EdgesGrouper) (groups []GroupedEdges) {
	_ = vtx.GroupIncomingBy(defineGroup, func(groupKey []byte, edges []*Edge) error {
		groups = append(groups, GroupedEdges{GroupKey: groupKey, Edges: edges})
		return nil
	})
	return
}

func (vtx *Vertex) GroupBy(defineGroup EdgesGrouper, action GroupEdgesAction) error {
	grouped := groupEdgesInto(groupEdges(vtx.outcoming, defineGroup), vtx.incoming, defineGroup)
	return applyEdgeGroupAction(grouped, action)
}

func (vtx *Vertex) GroupOutcomingBy(defineGroup EdgesGrouper, action GroupEdgesAction) error {
	return applyEdgeGroupAction(groupEdges(vtx.outcoming, defineGroup), action)
}

func (vtx *Vertex) GroupIncomingBy(defineGroup EdgesGrouper, action GroupEdgesAction) error {
	return applyEdgeGroupAction(groupEdges(vtx.incoming, defineGroup), action)
}

func groupEdges(es EdgeSet, defineGroup EdgesGrouper) map[string][]*Edge {
	grouped := map[string][]*Edge{}
	return groupEdgesInto(grouped, es, defineGroup)
}

func groupEdgesInto(grouped map[string][]*Edge, es EdgeSet, defineGroup EdgesGrouper) map[string][]*Edge {
	for iterator := es.Iterator(); iterator.HasNext(); {
		edge := iterator.Next()
		gkey := string(defineGroup(edge))
		grouped[gkey] = append(grouped[gkey], edge)
	}
	return grouped
}

func (vtx *Vertex) GroupVertexes(pathGrouper CompleteGrouperOverEdgesPath) []GroupedVertexes {
	// TODO: get rid of cycling over closed graph paths
	currentVtxs := NewVertexSet(vtx)
	for _, pathSelector := range pathGrouper.pathOverEdge.selectors {
		nextVtxs := NewVertexSet()
		for currentVtxIterator := currentVtxs.Iterator(); currentVtxIterator.HasNext(); {
			currentVtx := currentVtxIterator.Next()
			for _, iterator := range []EdgeSetIterator{currentVtx.incoming.Iterator(), currentVtx.outcoming.Iterator()} {
				for iterator.HasNext() {
					edge := iterator.Next()
					if pathSelector(edge) {
						nextVtxs.put(edge.vertex)
					}
				}
			}
		}
		currentVtxs = nextVtxs
	}
	return currentVtxs.GroupedBy(pathGrouper.vtxGrouper)
}

func applyEdgeGroupAction(grouped map[string][]*Edge, action GroupEdgesAction) error {
	for gkey, edges := range grouped {
		if err := action([]byte(gkey), edges); err != nil {
			return err
		}
	}
	return nil
}

type VertexSet struct {
	set map[*Vertex]int
	order map[int]*Vertex
}

func NewVertexSet(vtxs ...*Vertex) (vs VertexSet) {
	vs.set = map[*Vertex]int{}
	vs.order = map[int]*Vertex{}
	for _, vtx := range vtxs {
		vs.put(vtx)
	}
	return
}

func (vs VertexSet) Contains(vtx *Vertex) (found bool) {
	_, found = vs.set[vtx]
	return
}

func (vs VertexSet) ContainsData(predicate DataPredicate) *Vertex {
	for iterator := vs.Iterator(); iterator.HasNext(); {
		vtx := iterator.Next()
		if predicate(vtx.Data()) {
			return vtx
		}
	}
	return nil
}

func (vs VertexSet) Len() int {
	return len(vs.set)
}

func (vs VertexSet) OutcomingWhich(predicate EdgePredicate) (res EdgeSet) {
	for vtx := range vs.set {
		es := vtx.OutcomingWhich(predicate)
		res = res.Merge(es)
	}
	return res
}

func (vs VertexSet) GroupedBy(defineGroup VertexesGrouper) (groups []GroupedVertexes) {
	_ = vs.GroupBy(defineGroup, func(groupKey []byte, vtxs []*Vertex) error {
		groups = append(groups, GroupedVertexes{GroupKey: groupKey, Vertexes: vtxs})
		return nil
	})
	return
}

func (vs VertexSet) GroupBy(defineGroup VertexesGrouper, action GroupVertexesAction) error {
	grouped := map[string][]*Vertex{}
	for vtx := range vs.set {
		gkey := string(defineGroup(vtx))
		grouped[gkey] = append(grouped[gkey], vtx)
	}
	return applyVertexGroupAction(grouped, action)
}

func applyVertexGroupAction(grouped map[string][]*Vertex, action GroupVertexesAction) error {
	for gkey, vtxs := range grouped {
		if err := action([]byte(gkey), vtxs); err != nil {
			return err
		}
	}
	return nil
}

func (vs VertexSet) put(vtx *Vertex) {
	vs.set[vtx] = len(vs.set)
	vs.order[len(vs.order)] = vtx
}

func (vs VertexSet) Iterator() VertexSetIterator {
	return VertexSetIterator{vs: vs}
}

type VertexSetIterator struct {
	current int
	vs      VertexSet
}

func (vi VertexSetIterator) HasNext() bool {
	return vi.current < len(vi.vs.order)
}

func (vi *VertexSetIterator) Next() (vtx *Vertex) {
	vtx = vi.vs.order[vi.current]
	vi.current++
	return
}

func NewEdgeSet(edges ...*Edge) (es EdgeSet) {
	es.container = map[int]*Edge{}
	for _, edge := range edges {
		es.put(edge)
	}
	return
}

type EdgeSet struct {
	container map[int]*Edge
}

func (es EdgeSet) Len() int {
	return len(es.container)
}

func (es EdgeSet) GroupedBy(defineGroup EdgesGrouper) (groups []GroupedEdges) {
	_ = es.GroupBy(defineGroup, func(groupKey []byte, edges []*Edge) error {
		groups = append(groups, GroupedEdges{GroupKey: groupKey, Edges: edges})
		return nil
	})
	return
}

func (es EdgeSet) GroupBy(defineGroup EdgesGrouper, action GroupEdgesAction) error {
	return applyEdgeGroupAction(groupEdges(es, defineGroup), action)
}

func (es EdgeSet) Vertexes() (vtxs []*Vertex) {
	for iterator := es.Iterator(); iterator.HasNext(); {
		edge := iterator.Next()
		vtxs = append(vtxs, edge.Vertex())
	}
	return
}

func (es EdgeSet) VertexesSet() VertexSet {
	vs := NewVertexSet()
	for iterator := es.Iterator(); iterator.HasNext(); {
		edge := iterator.Next()
		vs.put(edge.Vertex())
	}
	return vs
}

func (es EdgeSet) put(edge *Edge) {
	es.container[len(es.container)] = edge
}

func (es EdgeSet) Iterator() EdgeSetIterator {
	return EdgeSetIterator{es: es}
}

func (es EdgeSet) Merge(withEs EdgeSet) EdgeSet {
	merged := NewEdgeSet()
	for order, edge := range es.container {
		merged.container[order] = edge
	}
	delta := len(merged.container)
	for order, edge := range withEs.container {
		merged.container[order+delta] = edge
	}
	return merged
}

type EdgeSetIterator struct {
	current int
	es      EdgeSet
}

func (ei EdgeSetIterator) HasNext() bool {
	return ei.current < len(ei.es.container)
}

func (ei *EdgeSetIterator) Next() (edge *Edge) {
	edge = ei.es.container[ei.current]
	ei.current++
	return
}

package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pavelmemory/mgraph/graph"
)

func TestGraph(t *testing.T) {
	type testType struct {
		Field string
	}

	v1 := graph.VertexWith(testType{Field: "data1"})
	if v1.Adjacent().Len() != 0 {
		t.Error("must not have siblings")
	}
	if v1.Incoming().Len() != 0 {
		t.Error("must not have incoming edges")
	}
	if v1.Outcoming().Len() != 0 {
		t.Error("must not have outcoming edges")
	}

	expectedAmountOfEdges := 5
	v2 := graph.VertexWith(testType{Field: "data2"})
	for i := 0; i < expectedAmountOfEdges; i++ {
		v1.Edge(v2)
	}
	adjacent := v1.Adjacent()
	if adjacent.Len() != 1 {
		t.Errorf("wrong amount of adjacent vertexes: %d", adjacent.Len())
	}
	actualAmountOfEdges := v1.EdgesTo(v2).Len()
	if actualAmountOfEdges != expectedAmountOfEdges {
		t.Errorf("wrong amount of edges to vertex: %d", expectedAmountOfEdges)
	}
}

func TestGraph_EdgeTo(t *testing.T) {
	v0 := graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(graph.VertexWith(2))

	dstVtxs := v0.Outcoming().Vertexes()
	if len(dstVtxs) != 2 {
		t.Fatalf("incorrect amount of vertexes connected with outcoming edges: %d", len(dstVtxs))
	}
	amountOfIncomingEdges := v0.Incoming().Len()
	if amountOfIncomingEdges != 0 {
		t.Fatalf("must not have incoming edges: %d", amountOfIncomingEdges)
	}

	v1 := graph.VertexWith(3).
		EdgeTo(graph.VertexWith(4)).
		EdgeTo(graph.VertexWith(5))
	v0.EdgeTo(v1)

	dstVtxs = v1.Outcoming().Vertexes()
	if len(dstVtxs) != 2 {
		t.Fatalf("incorrect amount of vertexes connected with outcoming edges: %d", len(dstVtxs))
	}
	amountOfIncomingEdges = v1.Incoming().Len()
	if amountOfIncomingEdges != 1 {
		t.Fatalf("must have only one incoming edge: %d", amountOfIncomingEdges)
	}

	dstVtxs = v0.Outcoming().Vertexes()
	if len(dstVtxs) != 3 {
		t.Fatalf("incorrect amount of vertexes connected with outcoming edges: %d", len(dstVtxs))
	}
}

func TestGraphDFS_FindN(t *testing.T) {
	v0 := graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(graph.VertexWith(4).
			EdgeTo(graph.VertexWith(2)).
			EdgeTo(graph.VertexWith(3))).
		EdgeTo(graph.VertexWith(5))

	var searchChecks int
	var found bool

	err := v0.TraverseWith(
		graph.FindN(graph.DFS, 1, func(vtx *graph.Vertex) bool {
			searchChecks++
			return vtx.Data().(int) == 3
		}),
		func(vtx *graph.Vertex) error {
			found = true
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if searchChecks != 3 {
		t.Errorf("wrong amount of search checks: %d", searchChecks)
	}
	if !found {
		t.Error("data was not found")
	}
}

func TestGraphBFS_FindN(t *testing.T) {
	v0 := graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(graph.VertexWith(4).
			EdgeTo(graph.VertexWith(2)).
			EdgeTo(graph.VertexWith(3))).
		EdgeTo(graph.VertexWith(5))

	var searchChecks int
	var found bool

	err := v0.TraverseWith(
		graph.FindN(graph.BFS, 1, func(vtx *graph.Vertex) bool {
			searchChecks++
			return vtx.Data().(int) == 2
		}),
		func(vtx *graph.Vertex) error {
			found = true
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if searchChecks != 5 {
		t.Errorf("wrong amount of search checks: %d", searchChecks)
	}
	if !found {
		t.Error("data was not found")
	}
}

func TestGraph_FindAll(t *testing.T) {
	v0 := graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(graph.VertexWith(4).
			EdgeTo(graph.VertexWith(1)).
			EdgeTo(graph.VertexWith(3))).
		EdgeTo(graph.VertexWith(1))

	var searchChecks int
	var found int
	searchStrategy := graph.FindAll(func(vtx *graph.Vertex) bool {
		searchChecks++
		return vtx.Data().(int) == 1
	})

	err := v0.TraverseWith(searchStrategy, func(vtx *graph.Vertex) error {
		found++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if searchChecks != 6 {
		t.Errorf("wrong amount of search checks: %d", searchChecks)
	}
	if found != 3 {
		t.Error("data was not found")
	}
}

func TestGraph_GroupOutcoming(t *testing.T) {
	v := graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(graph.VertexWith(2).
			EdgeTo(graph.VertexWith(4)).
			EdgeTo(graph.VertexWith(5))).
		EdgeTo(graph.VertexWith(3))

	var even int
	var odd int
	var groups int
	err := v.GroupOutcomingBy(func(edge *graph.Edge) []byte {
		if edge.Vertex().Data().(int)%2 == 0 {
			even++
			return []byte{0}
		}
		odd++
		return []byte{1}
	}, func(groupKey []byte, edges []*graph.Edge) error {
		groups++
		switch {
		case bytes.Equal(groupKey, []byte{0}):
			even -= len(edges)
		case bytes.Equal(groupKey, []byte{1}):
			odd -= len(edges)
		default:
			t.Error("unexpected group key:", groupKey)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if groups != 2 {
		t.Error("unexpected amount of groups:", groups)
	}
	if even != 0 {
		t.Error("unexpected amount of elements in 'even' group:", even)
	}
	if odd != 0 {
		t.Error("unexpected amount of elements in 'odd' group:", odd)
	}
}

func TestGraph_Group(t *testing.T) {
	v := graph.VertexWith(2).
		EdgeTo(graph.VertexWith(4)).
		EdgeTo(graph.VertexWith(5))
	graph.VertexWith(0).
		EdgeTo(graph.VertexWith(1)).
		EdgeTo(v).
		EdgeTo(graph.VertexWith(3))

	t.Run("GroupBy", func(t *testing.T) {
		even := 2
		odd := 1
		var groups int
		err := v.GroupBy(func(edge *graph.Edge) []byte {
			if edge.Vertex().Data().(int)%2 == 0 {
				return []byte{0}
			}
			return []byte{1}
		}, func(groupKey []byte, edges []*graph.Edge) error {
			groups++
			switch {
			case bytes.Equal(groupKey, []byte{0}):
				even -= len(edges)
			case bytes.Equal(groupKey, []byte{1}):
				odd -= len(edges)
			default:
				t.Error("unexpected group key:", groupKey)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if groups != 2 {
			t.Error("unexpected amount of groups:", groups)
		}
		if even != 0 {
			t.Error("unexpected amount of elements in 'even' group:", even)
		}
		if odd != 0 {
			t.Error("unexpected amount of elements in 'odd' group:", odd)
		}
	})

	t.Run("GroupedBy", func(t *testing.T) {
		groups := v.GroupedBy(func(edge *graph.Edge) []byte {
			if edge.Vertex().Data().(int)%2 == 0 {
				return []byte{0}
			}
			return []byte{1}
		})
		if len(groups) != 2 {
			t.Error("unexpected amount of groups:", groups)
		}
		if !bytes.Equal(groups[0].GroupKey, []byte{0}) {
			t.Error("group has unexpected key:", groups[0].GroupKey)
		}
		if !bytes.Equal(groups[1].GroupKey, []byte{1}) {
			t.Error("group has unexpected key:", groups[1].GroupKey)
		}
		if len(groups[0].Edges) != 2 {
			t.Errorf("group has unexpected edges: %+v", groups[0].Edges)
		}
		if len(groups[1].Edges) != 1 {
			t.Errorf("group has unexpected edges: %+v", groups[1].Edges)
		}
	})
}

func TestGraph_GroupGroup(t *testing.T) {
	type Simulation struct {
		PK       string
		Baseline bool
	}

	type Product struct {
		ID   int
		Code string
	}

	type Attribute struct {
		Label string
	}

	type Metric struct {
		Id    string
		Value float64
	}

	const (
		simulation = "simulation"
		metric     = "metric"
		product    = "product"
		brand      = "brand"
	)
	opportunity := graph.VertexWith(nil)
	baseline := graph.VertexWith(&Simulation{PK: "1", Baseline: true})
	working := graph.VertexWith(&Simulation{PK: "2"})
	opportunity.EdgeToWith(baseline, simulation)
	opportunity.EdgeToWith(working, simulation)

	brandA := graph.VertexWith(&Attribute{Label: "A"})
	brandB := graph.VertexWith(&Attribute{Label: "B"})

	prod1 := graph.VertexWith(&Product{ID: 10, Code: "Coke"}).
		EdgeToWith(brandA, brand).
		EdgeToWith(graph.VertexWith(&Metric{Id: "units", Value: 10.0}), metric)

	prod2 := graph.VertexWith(&Product{ID: 11, Code: "Pepsi"}).
		EdgeToWith(brandA, brand).
		EdgeToWith(graph.VertexWith(&Metric{Id: "units", Value: 1.0}), metric)

	prod3 := graph.VertexWith(&Product{ID: 11, Code: "Cherry Juice"}).
		EdgeToWith(brandB, brand).
		EdgeToWith(graph.VertexWith(&Metric{Id: "units", Value: 2.0}), metric)

	prod4 := graph.VertexWith(&Product{ID: 11, Code: "Lemonnello"}).
		EdgeToWith(brandB, brand).
		EdgeToWith(graph.VertexWith(&Metric{Id: "units", Value: 4.0}), metric)

	baseline.EdgeToWith(prod1, product)
	baseline.EdgeToWith(prod2, product)
	baseline.EdgeToWith(prod3, product)

	working.EdgeToWith(prod2, product)
	working.EdgeToWith(prod3, product)
	working.EdgeToWith(prod4, product)

	groupedProductVertexes := opportunity.
		OutcomingWhich(graph.EdgeAttributeEqualsTo(simulation)).
		VertexesSet().
		OutcomingWhich(graph.EdgeAttributeEqualsTo(product)).
		VertexesSet().
		GroupedBy(func(vtx *graph.Vertex) []byte {
			for _, vtx := range vtx.OutcomingWhich(graph.EdgeAttributeEqualsTo(brand)).Vertexes() {
				return []byte(vtx.Data().(*Attribute).Label)
			}
			return nil
		})

	for _, perGroupProductVertexes := range groupedProductVertexes {
		fmt.Println(string(perGroupProductVertexes.GroupKey))
		for _, vtx := range perGroupProductVertexes.Vertexes {
			p := vtx.Data().(*Product)
			fmt.Println(p.ID, p.Code)
		}
	}

	groupedProductEdges := opportunity.
		OutcomingWhich(graph.EdgeAttributeEqualsTo(simulation)).
		VertexesSet().
		OutcomingWhich(graph.EdgeAttributeEqualsTo(product)).
		GroupedBy(func(edge *graph.Edge) []byte {
			return []byte(edge.Vertex().Data().(*Product).Code)
		})
	for _, groupedProductEdge := range groupedProductEdges {
		fmt.Println(string(groupedProductEdge.GroupKey))
		for _, edge := range groupedProductEdge.Edges {
			fmt.Println(edge.Attributes().(string))
		}
	}

	groupedProductVertexesByMetrics := opportunity.GroupVertexes(
		graph.GoOverEdge(func(edge *graph.Edge) bool {
			attr, ok := edge.Attributes().(string)
			return ok && attr == simulation
		}).GoOverEdge(func(edge *graph.Edge) bool {
			attr, ok := edge.Attributes().(string)
			return ok && attr == product
		}).GroupVertexesWith(func(vtx *graph.Vertex) []byte {
			iter := vtx.OutcomingWhich(func(edge *graph.Edge) bool {
				attr, ok := edge.Attributes().(string)
				return ok && attr == metric
			}).VertexesSet().Iterator()
			return []byte(iter.Next().Data().(*Metric).Id)
		}))

	for _, groupedProductVtx := range groupedProductVertexesByMetrics {
		fmt.Println("Products group:", string(groupedProductVtx.GroupKey))
		for _, vtx := range groupedProductVtx.Vertexes {
			fmt.Println("Group item:", *(vtx.Data().(*Product)))
		}
	}

	groupedMetricVertexesByBrand := opportunity.GroupVertexes(
		graph.GoOverEdge(graph.EdgeAttributeEqualsTo(simulation)).
			GoOverEdge(graph.EdgeAttributeEqualsTo(product)).
			GoOverEdge(graph.EdgeAttributeEqualsTo(metric)).
			GroupVertexesWith(func(vtx *graph.Vertex) []byte {
				return vtx.GroupVertexes(
					graph.GoOverEdge(graph.EdgeAttributeEqualsTo(metric)).
						GoOverEdge(graph.EdgeAttributeEqualsTo(brand)).
						GroupVertexesWith(func(vtx *graph.Vertex) []byte {
							return []byte(vtx.Data().(*Attribute).Label)
						}))[0].GroupKey
			}))

	for _, groupedMetricVertexByBrand := range groupedMetricVertexesByBrand {
		fmt.Println("Metrics group:", string(groupedMetricVertexByBrand.GroupKey))
		for _, vtx := range groupedMetricVertexByBrand.Vertexes {
			fmt.Println("Group item:", *(vtx.Data().(*Metric)))
		}
	}
}

func TestExistVertexes(t *testing.T) {
	type Product struct {
		Brand, Vendor string
	}

	type Metric struct {
		Name  string
		Value float64
	}

	type Market struct {
		Id string
	}

	type Period struct {
		Id string
	}

	const (
		attributeProduct = "product"
		attributeMetric  = "metric"
		attributeMarket  = "market"
		attributePeriod  = "period"
	)

	metricDist02 := graph.VertexWith(Metric{"distribution", 0.2})
	metricDist03 := graph.VertexWith(Metric{"distribution", 0.3})

	marketChild := graph.
		VertexWith(Market{"child"}).
		EdgeToWith(metricDist02, attributeMetric).
		EdgeToWith(graph.VertexWith(Period{"2000"}), attributePeriod)

	marketParent := graph.
		VertexWith(Market{Id: "parent"}).
		EdgeToWith(metricDist03, attributeMetric).
		EdgeToWith(graph.VertexWith(Period{"2001"}), attributePeriod)

	productCoke := graph.
		VertexWith(Product{Brand: "Coca-Cola", Vendor: "Nestle"}).
		EdgeToWith(marketParent, attributeMarket)

	productPepsi := graph.
		VertexWith(Product{Brand: "PepsiCo", Vendor: "Nestle"}).
		EdgeToWith(marketChild, attributeMarket).
		EdgeToWith(marketParent, attributeMarket)

	assortment := graph.
		VertexWith("Assortment of drinks").
		EdgeToWith(productCoke, attributeProduct).
		EdgeToWith(productPepsi, attributeProduct)

	periodIdIs := func(id string) func(vtx *graph.Vertex) bool {
		return func(vtx *graph.Vertex) bool {
			period, ok := vtx.Data().(Period)
			return ok && period.Id == id
		}
	}

	gg := assortment.GroupVertexes(graph.
		GoOverEdge(graph.EdgeAttributeEqualsTo(attributeProduct)).
		GroupVertexesWith(
			func(vtx *graph.Vertex) []byte {
				if vtx.ExistVertexes(graph.
					GoOverEdge(graph.EdgeAttributeEqualsTo(attributeMarket)).
					GoOverEdge(graph.EdgeAttributeEqualsTo(attributePeriod)).
					ExistVertexesWith(periodIdIs("2000")),
				) {
					return []byte("with 2000")
				}
				return []byte("without 2000")
			}))

	if len(gg) != 2 {
		t.Error("incorrect number of groups", len(gg))
		return
	}

	check := func(v1 graph.GroupedVertexes, v2 graph.GroupedVertexes) {
		if len(v1.Vertexes) != 1 {
			t.Errorf("incorrect number of vertexes '%d'", len(v1.Vertexes))
			return
		}
		if v1.Vertexes[0].Data() != productPepsi.Data() {
			t.Errorf("unexpected data in vertex %+v", v1.Vertexes[0].Data())
		}
		if len(v2.Vertexes) != 1 {
			t.Errorf("incorrect number of vertexes '%d'", len(v2.Vertexes))
			return
		}
		if v2.Vertexes[0].Data() != productCoke.Data() {
			t.Errorf("unexpected data in vertex %+v", v2.Vertexes[0].Data())
		}
	}

	switch {
	case bytes.Equal(gg[0].GroupKey, []byte("with 2000")):
		check(gg[0], gg[1])
	case bytes.Equal(gg[0].GroupKey, []byte("without 2000")):
		check(gg[1], gg[0])
	default:
		t.Errorf("unexpected group key '%s'", string(gg[0].GroupKey))
	}
}

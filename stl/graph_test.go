package stl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TopoSortableData struct {
	Key     string
	Formers []string
}

func (d TopoSortableData) GetKey() string {
	return d.Key
}

func (d TopoSortableData) GetFormers() []string {
	return d.Formers
}

func assertTopoOrder[T comparable](t *testing.T, graph Graph[T], order []T) {
	t.Helper()

	index := make(map[T]int, len(order))
	for i, node := range order {
		index[node] = i
	}

	for from, tos := range graph {
		fromIdx, fromOK := index[from]
		if !fromOK {
			continue
		}
		for _, to := range tos {
			toIdx, toOK := index[to]
			if !toOK {
				continue
			}
			assert.Less(t, fromIdx, toIdx, "edge %v -> %v violated in order %v", from, to, order)
		}
	}
}

func TestTopoSortByFormers(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"b", "c"}},
		{Key: "b", Formers: []string{"c"}},
		{Key: "c", Formers: []string{}},
	}

	result := TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, nil)
	assert.Equal(t, []string{"c", "b", "a"}, Map(result, TopoSortableData.GetKey))
}

func TestTopoSortWithOrder(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a"},
		{Key: "b"},
		{Key: "c"},
		{Key: "d"},
		{Key: "e"},
	}

	result := TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, Map(result, TopoSortableData.GetKey))
}

func TestTopoSortDataByFormersWithMaxHeap(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a"},
		{Key: "b"},
		{Key: "c"},
	}

	result := TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMaxHeapAsContainer[string])
	assert.Equal(t, []string{"c", "b", "a"}, Map(result, TopoSortableData.GetKey))
}

func TestTopoSortDataByFormersIgnoresUnknownFormerInOutput(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"x"}},
		{Key: "b", Formers: []string{"a"}},
	}

	result := TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"a", "b"}, Map(result, TopoSortableData.GetKey))
}

func TestTopoSortDataByFormersWithCycle(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"b"}},
		{Key: "b", Formers: []string{"a"}},
		{Key: "c", Formers: []string{}},
	}

	result := TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"c"}, Map(result, TopoSortableData.GetKey))
}

func TestGraphTopoSortEmptyAndSingle(t *testing.T) {
	assert.Empty(t, Graph[string]{}.TopoSort(nil))
	assert.Equal(t, []string{"a"}, Graph[string]{"a": nil}.TopoSort(NewMinHeapAsContainer[string]))
}

func TestGraphTopoSortWithMinHeap(t *testing.T) {
	graph := Graph[string]{
		"a": {"c"},
		"b": {"c"},
		"c": nil,
	}

	order := graph.TopoSort(NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"a", "b", "c"}, order)
	assertTopoOrder(t, graph, order)
}

func TestGraphTopoSortDefaultStackIsValid(t *testing.T) {
	graph := Graph[int]{
		1: {3},
		2: {3},
		3: {4},
		4: nil,
	}

	order := graph.TopoSort(nil)
	assert.Len(t, order, 4)
	assertTopoOrder(t, graph, order)
}

func TestGraphTopoSortWithCycle(t *testing.T) {
	graph := Graph[string]{
		"a": {"b"},
		"b": {"a"},
		"c": {"a"},
	}

	order := graph.TopoSort(NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"c"}, order)
}

func TestGraphTopoSortIncludesTargetOnlyNodes(t *testing.T) {
	// "b" only appears as an edge target.
	graph := Graph[string]{
		"a": {"b"},
	}

	order := graph.TopoSort(NewMinHeapAsContainer[string])
	assert.Equal(t, []string{"a", "b"}, order)
	assertTopoOrder(t, graph, order)
}

func flattenLayers[T any](layers [][]T) []T {
	var out []T
	for _, layer := range layers {
		out = append(out, layer...)
	}
	return out
}

func assertTopoLayers[T comparable](t *testing.T, graph Graph[T], layers [][]T) {
	t.Helper()

	layerOf := make(map[T]int)
	for i, layer := range layers {
		for _, node := range layer {
			layerOf[node] = i
		}
	}

	for from, tos := range graph {
		fromLayer, fromOK := layerOf[from]
		if !fromOK {
			continue
		}
		for _, to := range tos {
			toLayer, toOK := layerOf[to]
			if !toOK {
				continue
			}
			assert.Less(t, fromLayer, toLayer, "edge %v -> %v violated across layers %v", from, to, layers)
			assert.NotEqual(t, fromLayer, toLayer, "intra-layer edge %v -> %v in layer %d", from, to, fromLayer)
		}
	}
}

func TestGraphTopoSortLayersEmptyAndSingle(t *testing.T) {
	assert.Empty(t, Graph[string]{}.TopoSortLayers(nil))
	assert.Equal(t, [][]string{{"a"}}, Graph[string]{"a": nil}.TopoSortLayers(NewMinHeapAsContainer[string]))
}

func TestGraphTopoSortLayersWithMinHeap(t *testing.T) {
	graph := Graph[string]{
		"a": {"c"},
		"b": {"c"},
		"c": nil,
	}

	layers := graph.TopoSortLayers(NewMinHeapAsContainer[string])
	assert.Equal(t, [][]string{{"a", "b"}, {"c"}}, layers)
	assertTopoLayers(t, graph, layers)
}

func TestGraphTopoSortLayersDefaultStackIsValid(t *testing.T) {
	graph := Graph[int]{
		1: {3},
		2: {3},
		3: {4},
		4: nil,
	}

	layers := graph.TopoSortLayers(nil)
	assert.Len(t, flattenLayers(layers), 4)
	assertTopoLayers(t, graph, layers)
}

func TestGraphTopoSortLayersWithCycle(t *testing.T) {
	graph := Graph[string]{
		"a": {"b"},
		"b": {"a"},
		"c": {"a"},
	}

	layers := graph.TopoSortLayers(NewMinHeapAsContainer[string])
	assert.Equal(t, [][]string{{"c"}}, layers)

	order := graph.TopoSort(NewMinHeapAsContainer[string])
	assert.ElementsMatch(t, order, flattenLayers(layers))
}

func TestGraphTopoSortLayersIncludesTargetOnlyNodes(t *testing.T) {
	graph := Graph[string]{
		"a": {"b"},
	}

	layers := graph.TopoSortLayers(NewMinHeapAsContainer[string])
	assert.Equal(t, [][]string{{"a"}, {"b"}}, layers)
	assertTopoLayers(t, graph, layers)

	order := graph.TopoSort(NewMinHeapAsContainer[string])
	assert.ElementsMatch(t, order, flattenLayers(layers))
}

func TestTopoSortDataByFormersLayers(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"b", "c"}},
		{Key: "b", Formers: []string{"c"}},
		{Key: "c", Formers: []string{}},
	}

	layers := TopoSortDataByFormersLayers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	assert.Equal(t, [][]string{{"c"}, {"b"}, {"a"}}, Map(layers, func(layer []TopoSortableData) []string {
		return Map(layer, TopoSortableData.GetKey)
	}))
}

func TestTopoSortDataByFormersLayersIgnoresUnknownFormer(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"x"}},
		{Key: "b", Formers: []string{"a"}},
	}

	layers := TopoSortDataByFormersLayers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	keys := Map(flattenLayers(layers), TopoSortableData.GetKey)
	assert.Equal(t, []string{"a", "b"}, keys)

	oneDim := Map(
		TopoSortDataByFormers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string]),
		TopoSortableData.GetKey,
	)
	assert.ElementsMatch(t, oneDim, keys)
}

func TestTopoSortDataByFormersLayersWithCycle(t *testing.T) {
	datas := []TopoSortableData{
		{Key: "a", Formers: []string{"b"}},
		{Key: "b", Formers: []string{"a"}},
		{Key: "c", Formers: []string{}},
	}

	layers := TopoSortDataByFormersLayers(datas, TopoSortableData.GetKey, TopoSortableData.GetFormers, NewMinHeapAsContainer[string])
	assert.Equal(t, [][]string{{"c"}}, Map(layers, func(layer []TopoSortableData) []string {
		return Map(layer, TopoSortableData.GetKey)
	}))
}

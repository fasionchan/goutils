package stl

type Graph[T comparable] map[T][]T

func (g Graph[T]) TopoSort(container func(capacity int) Container[T]) []T {
	inDegree := make(map[T]int)
	for node := range g {
		inDegree[node] = 0
	}

	for _, neighbors := range g {
		for _, neighbor := range neighbors {
			inDegree[neighbor]++
		}
	}

	if container == nil {
		container = NewStackAsContainer[T]
	}

	q := container(len(inDegree))
	for node, degree := range inDegree {
		if degree == 0 {
			q.Push(node)
		}
	}

	result := make([]T, 0)
	for !q.IsEmpty() {
		node := q.Pop()
		result = append(result, node)
		for _, neighbor := range g[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				q.Push(neighbor)
			}
		}
	}

	return result
}

// TopoSortLayers returns Kahn wave layers: each layer is the batch of nodes
// that are indegree-zero in the same wave (drained fully before processing
// outgoing edges). Within a layer, order follows continuous Pop from container;
// nil container defaults to NewStackAsContainer, same as TopoSort.
// Flattening layers need not equal TopoSort with the same container: the
// one-dimensional sort mixes newly unlocked nodes into the same queue immediately.
func (g Graph[T]) TopoSortLayers(container func(capacity int) Container[T]) [][]T {
	inDegree := make(map[T]int)
	for node := range g {
		inDegree[node] = 0
	}

	for _, neighbors := range g {
		for _, neighbor := range neighbors {
			inDegree[neighbor]++
		}
	}

	if container == nil {
		container = NewStackAsContainer[T]
	}

	q := container(len(inDegree))
	for node, degree := range inDegree {
		if degree == 0 {
			q.Push(node)
		}
	}

	layers := make([][]T, 0)
	for !q.IsEmpty() {
		layer := make([]T, 0)
		for !q.IsEmpty() {
			layer = append(layer, q.Pop())
		}
		for _, node := range layer {
			for _, neighbor := range g[node] {
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					q.Push(neighbor)
				}
			}
		}
		layers = append(layers, layer)
	}

	return layers
}

func graphFromFormers[
	Datas ~[]Data,
	Keys ~[]Key,
	Data any,
	Key comparable,
](datas Datas, getKey func(Data) Key, getFormerKeys func(Data) Keys) (Graph[Key], map[Key]Data) {
	mapping := MappingByKey(datas, getKey)

	graph := make(Graph[Key])
	for _, data := range datas {
		key := getKey(data)
		if _, ok := graph[key]; !ok {
			graph[key] = nil
		}

		for _, former := range getFormerKeys(data) {
			if _, ok := mapping[former]; !ok {
				continue
			}

			graph[former] = append(graph[former], key)
		}
	}

	return graph, mapping
}

func TopoSortDataByFormers[
	Datas ~[]Data,
	Keys ~[]Key,
	Data any,
	Key comparable,
](datas Datas, getKey func(Data) Key, getFormerKeys func(Data) Keys, container func(capacity int) Container[Key]) []Data {
	graph, mapping := graphFromFormers(datas, getKey, getFormerKeys)
	return MapValuesByKeys(mapping, graph.TopoSort(container)...)
}

func TopoSortDataByFormersLayers[
	Datas ~[]Data,
	Keys ~[]Key,
	Data any,
	Key comparable,
](datas Datas, getKey func(Data) Key, getFormerKeys func(Data) Keys, container func(capacity int) Container[Key]) [][]Data {
	graph, mapping := graphFromFormers(datas, getKey, getFormerKeys)
	keyLayers := graph.TopoSortLayers(container)

	result := make([][]Data, 0, len(keyLayers))
	for _, keys := range keyLayers {
		result = append(result, MapValuesByKeys(mapping, keys...))
	}
	return result
}

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

func TopoSortDataByFormers[
	Datas ~[]Data,
	Keys ~[]Key,
	Data any,
	Key comparable,
](datas Datas, getKey func(Data) Key, getFormerKeys func(Data) Keys, container func(capacity int) Container[Key]) []Data {
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

	return MapValuesByKeys(mapping, graph.TopoSort(container)...)
}

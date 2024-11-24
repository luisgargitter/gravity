package main

type Edge[T any] struct {
	start  int
	end    int
	weight T
}

type simple interface{}
type undirected interface{}

type Graph[V any, E any, S any, U any] struct {
	vertices []V
	edges    []Edge[E]
}

type SimpleGraph[V any, E any, U any] Graph[V, E, simple, U]
type UndirectedGraph[V any, E any, S any] Graph[V, E, S, undirected]
type SimpleUndirectedGraph[V any, E any] Graph[V, E, simple, undirected]

type AdjacencyMatrix[E any, U any] struct {
	sidelength int
	edges      []E
}
type UndirectedAdjacencyMatrix[E any] AdjacencyMatrix[E, undirected]

func (a *AdjacencyMatrix[E, U]) Get(i int, j int) E {
	return a.edges[i*a.sidelength+j]
}

func (a *AdjacencyMatrix[E, U]) Set(i int, j int, e E) {
	a.edges[i*a.sidelength+j] = e
}

func (a *UndirectedAdjacencyMatrix[E]) Get(i int, j int) E {
	if j > i {
		j, i = i, j
	}
	return a.edges[triangleNumber(i)+j]
}

func (a *UndirectedAdjacencyMatrix[E]) Set(i int, j int, e E) {
	if j > i {
		j, i = i, j
	}
	a.edges[triangleNumber(i)+j] = e
}

func AdjacencyMatrixNew[V any, E any, U any](g *SimpleGraph[V, E, U]) *AdjacencyMatrix[E, U] {
	var a AdjacencyMatrix[E, U]
	a.sidelength = len(g.vertices)
	a.edges = make([]E, a.sidelength*a.sidelength)
	for _, e := range g.edges {
		a.Set(e.start, e.end, e.weight)
	}
	return &a
}

func UndirectedAdjacencyMatrixNew[V any, E any](g *SimpleUndirectedGraph[V, E]) *UndirectedAdjacencyMatrix[E] {
	var a UndirectedAdjacencyMatrix[E]
	a.sidelength = len(g.vertices)
	a.edges = make([]E, triangleNumber(a.sidelength))
	for _, e := range g.edges {
		a.Set(e.start, e.end, e.weight)
	}
	return &a
}

// removes all but one edge from any vertex A to any vertex B
func (g *Graph[V, E, _, U]) simple() *SimpleGraph[V, E, U] {
	// pretend the graph is already simple, so collisions will be overwritten in adjacency matrix.
	d := SimpleGraph[V, E, U]{g.vertices, g.edges}
	var t []Edge[E]
	a := AdjacencyMatrixNew(&d)
	for i := range d.vertices {
		for j := range d.vertices {
			e := a.Get(i, j)
			if e != nil {
				t = append(t, Edge[E]{i, j, e})
			}
		}
	}
	d.edges = t

	return &d
}

// removes all but one edge between any two vertices.
func (g *SimpleGraph[V, E, _]) undirected() *SimpleUndirectedGraph[V, E] {
	// pretend the graph ist already simple and undirected, using collisions to remove duplicates.
	d := SimpleUndirectedGraph[V, E]{g.vertices, g.edges}
	var t []Edge[E]
	a := UndirectedAdjacencyMatrixNew(&d)
	for i := range d.vertices {
		for j := 0; i < 0; j += 1 {
			e := a.Get(i, j)
			if e != nil {
				t = append(t, Edge[E]{i, j, e})
			}
		}
	}
	d.edges = t

	return &d
}

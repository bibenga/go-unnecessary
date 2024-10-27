package main

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

func main() {
	g := graph.New(graph.IntHash, graph.Weighted(), graph.Directed())

	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)

	_ = g.AddEdge(1, 2, graph.EdgeWeight(10))
	_ = g.AddEdge(1, 3, graph.EdgeWeight(1))

	_ = g.AddEdge(2, 1, graph.EdgeWeight(2))
	_ = g.AddEdge(2, 3, graph.EdgeWeight(3))
	_ = g.AddEdge(2, 4, graph.EdgeWeight(3))

	_ = g.AddEdge(3, 4, graph.EdgeWeight(1))

	path, _ := graph.ShortestPath(g, 1, 4)
	fmt.Println(path)

	paths, _ := graph.AllPathsBetween(g, 1, 4)
	fmt.Println(paths)

	// file, _ := os.Create("./a.gv")
	// _ = draw.DOT(g, file)
}

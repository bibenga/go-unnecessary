package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	gs := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	gs.AddNode(simple.Node(1))
	gs.AddNode(simple.Node(2))
	gs.AddNode(simple.Node(3))
	gs.AddNode(simple.Node(4))
	gs.AddNode(simple.Node(5))

	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(1), gs.Node(2), 1))
	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(1), gs.Node(5), 1))

	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(2), gs.Node(3), 1))
	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(2), gs.Node(4), 3))
	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(2), gs.Node(5), 1))

	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(3), gs.Node(1), 1))
	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(3), gs.Node(4), 1))

	gs.SetWeightedEdge(gs.NewWeightedEdge(gs.Node(4), gs.Node(5), 1))

	pt := path.DijkstraAllFrom(simple.Node(1), gs)
	for _, id := range []int64{2, 3, 4} {
		p, w := pt.AllTo(id)
		fmt.Println("p=", p, ";", "w=", w)
		if math.IsInf(w, -1) {
			fmt.Printf("negative cycle in path to %c path:%c\n", id, p)
		}
	}

	fmt.Println("--------------------")
	pta := path.YenKShortestPaths(gs, 10, 1, simple.Node(1), simple.Node(5))
	fmt.Println(pta)

	// fmt.Println("--------------------")
	// pta2, _ := path.AStar(simple.Node(1), simple.Node(4), gs, func(x, y graph.Node) float64 {
	// 	return 1
	// })
	// fmt.Println(pta2.)

	fmt.Println("====================")

	gm := multi.NewWeightedDirectedGraph()

	gm.AddNode(multi.Node(1))
	gm.AddNode(multi.Node(2))
	gm.AddNode(multi.Node(3))
	gm.AddNode(multi.Node(4))

	gm.SetWeightedLine(gm.NewWeightedLine(gm.Node(1), gm.Node(2), 1))
	gm.SetWeightedLine(gm.NewWeightedLine(gm.Node(2), gm.Node(3), 1))
	gm.SetWeightedLine(gm.NewWeightedLine(gm.Node(3), gm.Node(1), 1))
	gm.SetWeightedLine(gm.NewWeightedLine(gm.Node(3), gm.Node(4), 1))

	pt2, _ := path.BellmanFordFrom(simple.Node(1), gm)
	for _, id := range []int64{2, 3, 4} {
		p, w := pt2.To(id)
		fmt.Println("p=", p, ";", "w=", w)
		if math.IsInf(w, -1) {
			fmt.Printf("negative cycle in path to %c path:%c\n", id, p)
		}
	}
}

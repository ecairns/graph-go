package app

import "fmt"

func PrintGraph(graph Graph) {
	fmt.Printf("\nGraph: %v (%v)\n",
		graph.Name,
		graph.Id,
	)

	fmt.Printf("Nodes Len: %v\n", len(graph.Nodes))
	for i := range graph.Nodes {
		fmt.Printf("Node: %v (%v)\n",
			graph.Nodes[i].Name,
			graph.Nodes[i].Id,
		)
	}

	fmt.Printf("Edges Len: %v\n", len(graph.Edges))
	for i := range graph.Edges {
		fmt.Printf("Edge: (%v-%v)(%v)\tCost: %v\n",
			graph.Edges[i].From,
			graph.Edges[i].To,
			graph.Edges[i].Id,
			graph.Edges[i].Cost,
		)
	}
}

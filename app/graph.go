package app

import (
	"errors"
	"fmt"
)

type Graph struct {
	Id    string `xml:"id"`
	Name  string `xml:"name"`
	Nodes []Node `xml:"nodes>node"`
	Edges []Edge `xml:"edges>node"`
}

func (this *Graph) Save() bool {
	var query string
	var err error

	errs := this.Validate()
	if errs != nil {
		return false
	}

	query = "INSERT  INTO graph (id, name) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM graph WHERE id=$3);"
	_, err = DB.Exec(query, this.Id, this.Name, this.Id)
	if err != nil {
		fmt.Printf("graph INSERT ERROR: %v\n", err)
		return false
	}

	for i := range this.Nodes {
		var node = this.Nodes[i]
		saved := node.Save()
		if saved {
			query = "INSERT INTO graph_node(graph_id, node_id) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM graph_node where graph_id=$3 AND node_id=$4)"
			_, err = DB.Exec(query, this.Id, node.Id, this.Id, node.Id)
			if err != nil {
				fmt.Printf("graph_node INSERT ERROR: %v\n", err)
				return false
			}
		}
	}

	for i := range this.Edges {
		var edge = this.Edges[i]
		saved := edge.Save()
		if saved {
			query = "INSERT INTO graph_edge(graph_id, edge_id) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM graph_edge where graph_id=$3 AND edge_id=$4)"
			_, err = DB.Exec(query, this.Id, edge.Id, this.Id, edge.Id)
			if err != nil {
				fmt.Printf("graph_edge INSERT ERROR: %v\n", err)
				return false
			}
		}
	}

	return true
}

func (this *Graph) Validate() []error {
	var errs []error
	var nodes = make(map[string]Node)
	var edges = make(map[string]Edge)

	if this.Id == "" {
		errs = append(errs, errors.New("Graph ID cannot be empty"))
	}

	if this.Name == "" {
		errs = append(errs, errors.New("Graph Name cannot be empty"))
	}

	if len(this.Nodes) == 0 {
		errs = append(errs, errors.New("Nodes cannot be empty"))
	} else {

		for i := range this.Nodes {
			var node = this.Nodes[i]

			tmpErrs := node.Validate()
			if tmpErrs != nil {
				errs = append(errs, tmpErrs...)
			} else {
				//Check if node already exists
				_, ok := nodes[node.Id]
				if ok {
					errs = append(errs, errors.New(fmt.Sprintf("Node Id: %s is already defined.", node.Id)))
				} else {
					nodes[node.Id] = node
				}
			}
		}

		for i := range this.Edges {
			var edge = this.Edges[i]

			tmpErrs := edge.Validate()
			if tmpErrs != nil {
				errs = append(errs, tmpErrs...)
			} else {
				_, ok := edges[edge.Id]
				if ok {
					errs = append(errs, errors.New(fmt.Sprintf("Edge ID: %s is already defined.", edge.Id)))
				} else {
					edges[edge.Id] = edge

					_, ok := nodes[edge.From]
					if !ok {
						errs = append(errs, errors.New(fmt.Sprintf("Edge From Id `%s` does not exist", edge.From)))
					}
					_, ok = nodes[edge.To]
					if !ok {
						errs = append(errs, errors.New(fmt.Sprintf("Edge To Id `%s` does not exist", edge.To)))
					}
				}
			}
		}
	}

	return errs
}

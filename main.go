package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Usage() {
	fmt.Println("USAGE:\n\t[URL | LocalFile]\n\n")
}
func main() {
	if len(os.Args) != 2 {
		Usage()
	} else {
		fileLoc := os.Args[1]

		errs := Parse(fileLoc)
		if errs != nil {
			log.Println("Graph file is invalid.")
			log.Println(errs)
		}
	}
}

func Parse(fileLoc string) []error {
	var graph Graph
	var errs []error

	if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
		graph, errs = ParseGraphFile(fileLoc)
		if errs != nil {
			return errs
		}

	} else if isValidUrl(fileLoc) {
		fmt.Println("Retrieve from URL")
		graph, errs = ParseUrl(fileLoc)
		if errs != nil {
			return errs
		}

	} else {
		return append(errs, errors.New("Invalid location for graph file"))
	}

	PrintGraph(graph)

	return nil
}

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

func ParseUrl(url string) (Graph, []error) {
	var graph Graph
	var errs []error

	// Create Temp File:  This will create a filename like /tmp/graph-123456.xml
	file, err := ioutil.TempFile(os.TempDir(), "graph-*.xml")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	// Clean up temp file after function ends
	defer os.Remove(file.Name())

	err = DownloadFile(file, url)
	if err != nil {
		return graph, append(errs, err)
	}

	graph, errs = ParseGraphFile(file.Name())
	if errs != nil {
		return graph, errs
	}

	return graph, nil
}

func ParseGraphFile(file string) (Graph, []error) {
	var errs []error
	var graph Graph

	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		return graph, append(errs, err)
	}

	xml.Unmarshal(byteValue, &graph)
	validationErrs := graph.Validate()
	if validationErrs != nil {
		return graph, validationErrs
	}

	return graph, nil
}

func DownloadFile(file *os.File, url string) error {
	fmt.Printf("Downloading graph xml: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Invalid status downloading xml: %v", resp.Status))
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func isValidUrl(toTest string) bool {
	u, err := url.Parse(toTest)
	fmt.Printf("Hostname: %v\n", u.Hostname())

	if err != nil {
		return false
	}

	if u.Hostname() == "" {
		return false
	}

	return true
}

type Node struct {
	Id   string `xml:"id"`
	Name string `xml:"name"`
}

func (this *Node) Validate() []error {
	var errs []error

	if this.Id == "" {
		errs = append(errs, errors.New("Node ID cannot be empty"))
	}

	if this.Name == "" {
		this.Name = this.Id
	}

	return errs
}

type Edge struct {
	Id   string  `xml:"id"`
	From string  `xml:"from"`
	To   string  `xml:"to"`
	Cost float32 `xml:"cost"` //default will be 0
}

func (this *Edge) Validate() []error {
	var errs []error

	if this.Id == "" {
		errs = append(errs, errors.New("Edge ID cannot be empty"))
	}
	if this.From == "" {
		errs = append(errs, errors.New("Edge From cannot be empty"))
	}
	if this.To == "" {
		errs = append(errs, errors.New("Edge To cannot be empty"))
	}
	if this.Cost < 0 {
		errs = append(errs, errors.New("Edge Cost cannot be negative"))

	}
	return errs
}

type Graph struct {
	Id    string `xml:"id"`
	Name  string `xml:"name"`
	Nodes []Node `xml:"nodes>node"`
	Edges []Edge `xml:"edges>node"`
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
			var edge= this.Edges[i]

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

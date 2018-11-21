package app

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func Parse(fileLoc string) (Graph, []error) {
	var graph Graph
	var errs []error

	if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
		graph, errs = ParseGraphFile(fileLoc)
		if errs != nil {
			return graph, errs
		}

	} else if isValidUrl(fileLoc) {
		graph, errs = ParseUrl(fileLoc)
		if errs != nil {
			return graph, errs
		}

	} else {
		return graph, append(errs, errors.New("Invalid location for graph file"))
	}

	return graph, nil
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

	err = download(file, url)
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

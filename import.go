package main

import (
	"fmt"
	"log"
	"os"
	"uniregistry/app"
)

func ImporterUsage() {
	fmt.Println("USAGE:\n\t[URL | LocalFile]\n\n")
}

func main() {
	if len(os.Args) != 2 {
		ImporterUsage()
	} else {
		app.DbInit()
		defer app.DbClose()

		fileLoc := os.Args[1]

		graph, errs := app.Parse(fileLoc)
		if errs != nil {
			log.Println("Graph file is invalid.")
			log.Println(errs)
		} else {
			graph.Save()

			app.PrintGraph(graph)
		}
	}
}

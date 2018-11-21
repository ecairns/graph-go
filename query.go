package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"uniregistry/app"
)

func QueryUsage() {
	fmt.Println("USAGE:\n\t[URL | LocalFile]\n\n")
}
func main() {
	if len(os.Args) != 2 {
		QueryUsage()
	} else {
		app.DbInit()
		defer app.DbClose()
		fileLoc := os.Args[1]

		if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
			queryStr, err := ioutil.ReadFile(fileLoc)
			if err != nil {
				panic(err)
			}
			qr, err := app.ParseQueryRequest(queryStr)
			if err != nil {
				panic(err)
			}

			jsonString, err := app.Query(qr)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(jsonString)
		}
	}
}

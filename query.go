package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"uniregistry/app"
)

func QueryUsage() {
	fmt.Println("USAGE:\n\t[STDIN | LocalFile]\n\n")
}
func main() {
	var fileLoc string
	var queryStr []byte

	if len(os.Args) == 2 {
		fileLoc = os.Args[1]
	}
	app.DbInit()
	defer app.DbClose()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		queryStr, _ = ioutil.ReadAll(os.Stdin)

	} else if _, err := os.Stat(fileLoc); !os.IsNotExist(err) {
		queryStr, err = ioutil.ReadFile(fileLoc)
		if err != nil {
			panic(err)
		}
	} else {
		QueryUsage()
	}

	if len(queryStr) > 0 {
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

package main

import (
	"docapi/format"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		println("Usage: docapi <path/to/project>")
		return
	}

	err := format.NewOpenAPI(args[0]).Generate()
	if err != nil {
		println(err.Error())
		return
	}
}

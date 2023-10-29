package main

import (
	"docapi/generator"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		println("Usage: api-doc <path/to/project>")
		return
	}
	doc := generator.Generate(args[0])
	println(doc)
}

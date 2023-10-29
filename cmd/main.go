package main

import (
	"docapi/generator"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		println("Usage: docapi <path/to/project>")
		return
	}

	gen := generator.New()
	out, err := gen.Run(args[0])
	if err != nil {
		return
	}

	println(out)
}

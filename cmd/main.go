package main

import (
	"docapi/exporter/openapi"
	"docapi/generator"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		println("Usage: docapi <path/to/project>")
		return
	}

	gen := generator.New()
	intermediateGen, err := gen.Run(args[0])
	if err != nil {
		return
	}

	e := openapi.NewOpenAPIExporter()
	out, err := e.Export(intermediateGen)
	if err != nil {
		return
	}

	j, err := yaml.Marshal(out)
	if err != nil {
		return
	}

	fmt.Print(string(j))
}

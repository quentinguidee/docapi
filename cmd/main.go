package main

import (
	"docapi/exporter/openapi"
	"docapi/generator"
	"log"
	"os"
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
		log.Fatal(err)
		return
	}

	e := openapi.NewOpenAPIExporter()
	err = e.Export(intermediateGen)
	if err != nil {
		log.Fatal(err)
		return
	}
}

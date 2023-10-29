package openapi

import (
	"docapi/generator"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Exporter struct{}

func NewOpenAPIExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(inter generator.IntermediateGen) error {
	out := Format{
		Openapi: "3.0.0",
		Info: FormatInfo{
			Title:   inter.Api.Title,
			Desc:    inter.Api.Desc,
			Version: inter.Api.Version,
		},
		Paths: map[string]FormatRoutes{},
	}

	for _, route := range inter.Api.Routes {
		method := strings.ToLower(route.Method)

		if out.Paths[route.Path] == nil {
			out.Paths[route.Path] = FormatRoutes{}
		}
		out.Paths[route.Path][method] = FormatRoute{
			Summary:   route.Summary,
			Responses: map[int]FormatResponse{},
		}
		for _, resp := range route.Responses {
			response := FormatResponse{}
			if resp.Type != "" {
				response.Content = map[string]FormatContent{
					"application/json": {
						Schema: FormatSchema{
							Type:       resp.Type,
							Properties: map[string]FormatSchema{},
							Ref:        fmt.Sprintf("#/components/schemas/%s", resp.Ref),
						},
					},
				}
			}
			out.Paths[route.Path][method].Responses[resp.Code] = response
		}
	}

	out.Components.Schemas = map[string]FormatSchema{}
	for name, tp := range inter.Types {
		out.Components.Schemas[name] = FormatSchema{
			Type:       tp.Type,
			Properties: map[string]FormatSchema{},
		}
		for fieldName, field := range tp.Fields {
			out.Components.Schemas[name].Properties[fieldName] = FormatSchema{
				Type: field.Type,
			}
		}
	}

	file, err := os.OpenFile("openapi.yaml", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	y, err := yaml.Marshal(out)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(file, string(y))
	if err != nil {
		return err
	}

	return nil
}

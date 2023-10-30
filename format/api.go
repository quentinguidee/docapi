package format

import (
	"fmt"

	"github.com/quentinguidee/docapi/collector"
	"github.com/quentinguidee/docapi/types"
)

type api struct {
	types.Format
	alias             string
	filename          string
	referencedSchemas []string
	routes            map[string]string
	tempHandler       types.FormatRoute
	handlers          map[string]types.FormatRoute
	handlerMethods    map[string]string
}

func newAPI(id string) *api {
	return &api{
		alias: id,
		Format: types.Format{
			Openapi: "3.0.0",
		},
		routes:         map[string]string{},
		handlers:       map[string]types.FormatRoute{},
		handlerMethods: map[string]string{},
	}
}

func (a *api) LinkResponses() error {
	for path, routes := range a.Paths {
		for method, route := range routes {
			for code, resp := range route.Responses {
				if resp.Ref != "" || resp.Description != "" {
					continue
				}
				a.Paths[path][method].Responses[code] = types.FormatResponse{
					Ref: fmt.Sprintf("#/components/responses/%s", code),
				}
			}
		}
	}
	return nil
}

func (a *api) schemaFromStruct(tp collector.Struct) types.FormatSchema {
	schema := types.FormatSchema{
		Type: tp.Type,
	}
	for fieldName, field := range tp.Fields {
		schema.SetProperty(fieldName, a.schemaFromAlias(field.Type))
	}
	return schema
}

func (a *api) schemaFromAlias(name string) types.FormatSchema {
	if isDefaultType(name) {
		if name == "bool" {
			name = "boolean"
		}
		return types.FormatSchema{
			Type: name,
		}
	} else {
		s := types.FormatSchema{
			Ref: fmt.Sprintf("#/components/schemas/%s", name),
		}
		a.referencedSchemas = append(a.referencedSchemas, name)
		return s
	}
}

func isDefaultType(name string) bool {
	switch name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128",
		"string", "bool", "byte", "rune":
		return true
	default:
		return false
	}
}

package format

import (
	"fmt"
	"strings"

	"github.com/quentinguidee/docapi/collector"
	"github.com/quentinguidee/docapi/types"
)

type api struct {
	types.Format
	alias          string
	filename       string
	routes         map[string]string
	tempHandler    types.FormatRoute
	handlers       map[string]types.FormatRoute
	handlerMethods map[string]string
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
				if resp.Ref.Name() != "" || resp.Description != "" {
					continue
				}
				a.Paths[path][method].Responses[code] = types.FormatResponse{
					Ref: types.CreateRef(types.RefResponse, code),
				}
			}
		}
	}
	return nil
}

func (a *api) CollectComponents(structs map[string]collector.Struct, aliases map[string]string, maps map[string]collector.Map) error {
	it := 0
	itComponents := a.GetReferencedComponents()
	done := 0
	count := len(itComponents)

	// The loop handles the case where a schema references another schema.
	for {
		for _, comp := range itComponents {
			if s, ok := structs[comp]; ok {
				a.Components.SetSchema(comp, a.schemaFromStruct(s))
			} else if alias, ok := aliases[comp]; ok {
				a.Components.SetSchema(comp, a.schemaFromAlias(alias))
			} else if m, ok := maps[comp]; ok {
				a.Components.SetSchema(comp, a.schemaFromMap(m))
			} else {
				a.Components.SetSchema(comp, types.FormatSchema{
					Type: "string",
				})
			}
		}

		done = count
		itComponents = a.GetReferencedComponents()
		count = len(itComponents)

		if count == done {
			break
		}

		it += 1
		if it > 100 {
			return fmt.Errorf("too many iterations")
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

func (a *api) schemaFromMap(tp collector.Map) types.FormatSchema {
	return types.FormatSchema{
		Type: "object",
		Properties: map[string]types.FormatSchema{
			tp.Key:   a.schemaFromAlias(tp.Key),
			tp.Value: a.schemaFromAlias(tp.Value),
		},
	}
}

func (a *api) schemaFromAlias(name string) types.FormatSchema {
	if strings.HasPrefix(name, "[]") {
		child := a.schemaFromAlias(name[2:])
		return types.FormatSchema{
			Type:  "array",
			Items: &child,
		}
	} else if isDefaultType(name) {
		if name == "bool" {
			name = "boolean"
		}
		return types.FormatSchema{
			Type: name,
		}
	} else {
		return types.FormatSchema{
			Ref: types.CreateRef(types.RefSchema, name),
		}
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

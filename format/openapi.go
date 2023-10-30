package format

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/quentinguidee/docapi/collector"
	"github.com/quentinguidee/docapi/types"

	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	types.Format
	path              string
	referencedSchemas []string
	servers           map[string]*types.FormatServer
}

func NewOpenAPI(path string) *OpenAPI {
	return &OpenAPI{
		Format: types.Format{
			Openapi: "3.0.0",
		},
		path:    path,
		servers: map[string]*types.FormatServer{},
	}
}

func (f *OpenAPI) Generate() error {
	err := f.CollectCommands(f.path)
	if err != nil {
		return err
	}

	err = f.CollectComponents(f.path)
	if err != nil {
		return err
	}

	err = f.LinkResponses()
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(f.Format)
	if err != nil {
		return err
	}
	return os.WriteFile("openapi.yaml", out, 0644)
}

func (f *OpenAPI) CollectCommands(path string) error {
	commands, err := collector.NewCommandsCollector().Run(path)
	if err != nil {
		return err
	}

	handlers := map[string]types.FormatRoute{}
	handlerMethods := map[string]string{}
	routes := map[string]string{}
	routeAliases := map[string]string{}

	var tempHandler types.FormatRoute
	var tempHandlerID string

	for _, cmd := range commands {
		switch cmd.Type {
		case types.CmdTitle:
			f.Info.Title = strings.Join(cmd.Args, " ")
		case types.CmdDescription:
			f.Info.Desc = strings.Join(cmd.Args, " ")
		case types.CmdVersion:
			f.Info.Version = cmd.Args[0]
		case types.CmdUrl:
			f.collectUrl(cmd.Args)
		case types.CmdUrlVar:
			f.collectUrlVar(cmd.Args)
		case types.CmdCode:
			f.collectCode(cmd)
		case types.CmdRoute:
			routes[cmd.Args[1]] = cmd.Args[0]
			if cmd.ServerAlias != "" {
				routeAliases[cmd.Args[1]] = cmd.ServerAlias
			}
		case types.CmdBegin:
			tempHandler = types.FormatRoute{}
			tempHandlerID = cmd.Args[0]
		case types.CmdMethod:
			handlerMethods[tempHandlerID] = strings.ToLower(cmd.Args[0])
		case types.CmdSummary:
			tempHandler.Summary = strings.Join(cmd.Args, " ")
		case types.CmdTags:
			tempHandler.Tags = cmd.Args
		case types.CmdBody:
			f.collectBody(&tempHandler, cmd)
		case types.CmdQuery:
			f.collectQuery(&tempHandler, cmd)
		case types.CmdResponse:
			f.collectResponse(&tempHandler, cmd)
		case types.CmdEnd:
			handlers[tempHandlerID] = tempHandler
		}
	}

	f.Paths = map[string]types.FormatRoutes{}
	for handlerID, route := range routes {
		if f.Paths[route] == nil {
			f.Paths[route] = types.FormatRoutes{}
		}
		method := handlerMethods[handlerID]
		alias := routeAliases[handlerID]
		handler := handlers[handlerID]
		if alias != "" {
			handler.AddServer(*f.servers[alias])
		}
		f.Paths[route][method] = handler
	}

	return nil
}

func (f *OpenAPI) collectUrl(args []string) {
	server := types.FormatServer{
		Url: args[1],
	}
	f.AddServer(server)
	f.servers[args[0]] = &server
}

func (f *OpenAPI) collectCode(cmd types.Command) {
	code := cmd.Args[0]
	args := cmd.Args[1:]
	ref := ""
	resp := types.FormatResponse{}
	if strings.HasPrefix(args[0], "{") {
		ref = args[0][1 : len(args[0])-1]
		args = args[1:]
		resp.Content = map[string]types.FormatContent{
			"application/json": {
				Schema: types.FormatSchema{
					Ref: fmt.Sprintf("#/components/schemas/%s", ref),
				},
			},
		}
		f.referencedSchemas = append(f.referencedSchemas, ref)
	}
	resp.Description = strings.Join(args, " ")
	f.Components.SetResponse(code, resp)
}

func (f *OpenAPI) collectBody(tempHandler *types.FormatRoute, cmd types.Command) {
	component := cmd.Args[0]
	component = component[1 : len(component)-1]
	description := cmd.Args[1:]

	tempHandler.RequestBody = types.FormatRequestBody{
		Description: strings.Join(description, " "),
		Required:    true,
		Content: map[string]types.FormatContent{
			"application/json": {
				Schema: types.FormatSchema{
					Ref: fmt.Sprintf("#/components/schemas/%s", component),
				},
			},
		},
	}

	f.referencedSchemas = append(f.referencedSchemas, component)
}

func (f *OpenAPI) collectQuery(tempHandler *types.FormatRoute, cmd types.Command) {
	component := cmd.Args[1]
	component = component[1 : len(component)-1]
	schema := f.schemaFromAlias(component)
	tempHandler.AddParameter(types.FormatParameter{
		In:          "query",
		Name:        cmd.Args[0],
		Description: strings.Join(cmd.Args[2:], " "),
		Required:    true,
		Schema:      schema,
	})
}

func (f *OpenAPI) collectResponse(tempHandler *types.FormatRoute, cmd types.Command) {
	if len(cmd.Args) <= 1 {
		tempHandler.SetResponse(cmd.Args[0], types.FormatResponse{})
		return
	}

	resp := types.FormatResponse{
		Description: strings.Join(cmd.Args[2:], " "),
		Content:     map[string]types.FormatContent{},
	}

	content := types.FormatContent{}
	component := cmd.Args[1]
	component = component[1 : len(component)-1]
	if strings.HasPrefix(component, "[]") {
		component = component[2:]
		content.Schema.Type = "array"
		content.Schema.Items = types.FormatItems{
			Ref: fmt.Sprintf("#/components/schemas/%s", component),
		}
	} else {
		content.Schema.Ref = fmt.Sprintf("#/components/schemas/%s", component)
	}
	f.referencedSchemas = append(f.referencedSchemas, component)

	resp.Content["application/json"] = content
	tempHandler.SetResponse(cmd.Args[0], resp)
}

func (f *OpenAPI) CollectComponents(path string) error {
	structs, aliases, err := collector.NewTypesCollector().Run(path)
	if err != nil {
		return err
	}

	it := 0
	// The loop handles the case where a schema references another schema.
	for len(f.referencedSchemas) > 0 {
		referencedSchemas := slices.Clone(f.referencedSchemas)
		f.referencedSchemas = []string{}

		for structName, s := range structs {
			if !slices.Contains(referencedSchemas, structName) {
				continue
			}
			f.Components.SetSchema(structName, f.schemaFromStruct(s))
		}

		for aliasName, alias := range aliases {
			if !slices.Contains(referencedSchemas, aliasName) {
				continue
			}
			f.Components.SetSchema(aliasName, types.FormatSchema{
				Type: alias,
			})
		}

		it += 1
		if it > 100 {
			return fmt.Errorf("too many iterations")
		}
	}

	return nil
}

func (f *OpenAPI) LinkResponses() error {
	for path, routes := range f.Paths {
		for method, route := range routes {
			for code, resp := range route.Responses {
				if resp.Ref != "" || resp.Description != "" {
					continue
				}
				r := f.Components.Responses[code]
				f.Paths[path][method].Responses[code] = r
			}
		}
	}
	return nil
}

func (f *OpenAPI) schemaFromStruct(tp collector.Struct) types.FormatSchema {
	schema := types.FormatSchema{
		Type: tp.Type,
	}
	for fieldName, field := range tp.Fields {
		schema.SetProperty(fieldName, f.schemaFromAlias(field.Type))
	}
	return schema
}

func (f *OpenAPI) schemaFromAlias(name string) types.FormatSchema {
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
		f.referencedSchemas = append(f.referencedSchemas, name)
		return s
	}
}

func (f *OpenAPI) collectUrlVar(args []string) {
	var (
		alias        = args[0]
		name         = args[1]
		defaultValue = args[2]
		description  = strings.Join(args[3:], " ")
	)

	v := types.FormatServerVariable{
		Default:     defaultValue,
		Description: description,
	}

	if f.servers == nil {
		f.servers = map[string]*types.FormatServer{}
	}
	server := f.servers[alias]
	server.SetVariable(name, v)
	for i, s := range f.Servers {
		if s.Url == server.Url {
			f.Servers[i] = *server
			break
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

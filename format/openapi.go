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
	path string
	apis map[string]*api
}

type api struct {
	types.Format
	filename          string
	referencedSchemas []string
	routes            map[string]string
	tempHandler       types.FormatRoute
}

func newServer() *api {
	return &api{
		Format: types.Format{
			Openapi: "3.0.0",
		},
		routes: map[string]string{},
	}
}

func NewOpenAPI(path string) *OpenAPI {
	return &OpenAPI{
		path: path,
		apis: map[string]*api{},
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

	for _, a := range f.apis {
		err = a.LinkResponses()
		if err != nil {
			return err
		}

		out, err := yaml.Marshal(a.Format)
		if err != nil {
			return err
		}

		name := fmt.Sprintf("openapi.%s.yaml", a.filename)
		err = os.WriteFile(name, out, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *OpenAPI) CollectCommands(path string) error {
	commands, err := collector.NewCommandsCollector().Run(path)
	if err != nil {
		return err
	}

	handlers := map[string]types.FormatRoute{}
	handlerMethods := map[string]string{}

	// get all aliases
	var aliases []string
	for _, cmd := range commands {
		if cmd.ServerAlias != "" {
			aliases = append(aliases, cmd.ServerAlias)
		}
	}

	// initialize servers
	for _, alias := range aliases {
		f.apis[alias] = newServer()
	}

	for _, cmd := range commands {
		// servers impacted by the command
		var apis []*api

		alias := cmd.ServerAlias
		if cmd.ServerAlias == "" {
			for _, a := range f.apis {
				apis = append(apis, a)
			}
		} else {
			apis = append(apis, f.apis[alias])
		}

		for _, a := range apis {
			switch cmd.Type {
			case types.CmdTitle:
				a.Info.Title = strings.Join(cmd.Args, " ")
			case types.CmdDescription:
				a.Info.Desc = strings.Join(cmd.Args, " ")
			case types.CmdVersion:
				a.Info.Version = cmd.Args[0]
			case types.CmdFilename:
				a.filename = cmd.Args[0]
			case types.CmdUrl:
				a.collectUrl(cmd.Args)
			case types.CmdUrlVar:
				a.collectUrlVar(cmd.Args)
			case types.CmdCode:
				a.collectCode(cmd)
			case types.CmdRoute:
				a.routes[cmd.Args[1]] = cmd.Args[0]
			case types.CmdBegin:
				a.tempHandler = types.FormatRoute{
					OperationId: cmd.Args[0],
				}
			case types.CmdMethod:
				handlerMethods[a.tempHandler.OperationId] = strings.ToLower(cmd.Args[0])
			case types.CmdSummary:
				a.tempHandler.Summary = strings.Join(cmd.Args, " ")
			case types.CmdTags:
				a.tempHandler.Tags = cmd.Args
			case types.CmdBody:
				a.collectBody(&a.tempHandler, cmd)
			case types.CmdQuery:
				a.collectQuery(&a.tempHandler, cmd)
			case types.CmdResponse:
				a.collectResponse(&a.tempHandler, cmd)
			case types.CmdEnd:
				handlers[a.tempHandler.OperationId] = a.tempHandler
			}
		}
	}

	for _, a := range f.apis {
		a.Paths = map[string]types.FormatRoutes{}
		for handlerID, route := range a.routes {
			if a.Paths[route] == nil {
				a.Paths[route] = types.FormatRoutes{}
			}
			method := handlerMethods[handlerID]
			handler := handlers[handlerID]
			a.Paths[route][method] = handler
		}
	}

	return nil
}

func (a *api) collectUrl(args []string) {
	a.AddServer(types.FormatServer{
		Url: args[0],
	})
}

func (a *api) collectCode(cmd types.Command) {
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
		a.referencedSchemas = append(a.referencedSchemas, ref)
	}
	resp.Description = strings.Join(args, " ")
	a.Components.SetResponse(code, resp)
}

func (a *api) collectBody(tempHandler *types.FormatRoute, cmd types.Command) {
	component := cmd.Args[0]
	component = component[1 : len(component)-1]
	description := cmd.Args[1:]

	tempHandler.RequestBody = types.FormatRequestBody{
		Description: strings.Join(description, " "),
		Required:    true,
		Content: map[string]types.FormatContent{
			"application/json": {
				Schema: a.schemaFromAlias(component),
			},
		},
	}

	a.referencedSchemas = append(a.referencedSchemas, component)
}

func (a *api) collectQuery(tempHandler *types.FormatRoute, cmd types.Command) {
	component := cmd.Args[1]
	component = component[1 : len(component)-1]
	schema := a.schemaFromAlias(component)
	tempHandler.AddParameter(types.FormatParameter{
		In:          "query",
		Name:        cmd.Args[0],
		Description: strings.Join(cmd.Args[2:], " "),
		Required:    true,
		Schema:      schema,
	})
}

func (a *api) collectResponse(tempHandler *types.FormatRoute, cmd types.Command) {
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
	a.referencedSchemas = append(a.referencedSchemas, component)

	resp.Content["application/json"] = content
	tempHandler.SetResponse(cmd.Args[0], resp)
}

func (f *OpenAPI) CollectComponents(path string) error {
	structs, aliases, err := collector.NewTypesCollector().Run(path)
	if err != nil {
		return err
	}

	for _, a := range f.apis {
		it := 0
		// The loop handles the case where a schema references another schema.
		for len(a.referencedSchemas) > 0 {
			referencedSchemas := slices.Clone(a.referencedSchemas)
			a.referencedSchemas = []string{}

			for structName, s := range structs {
				if !slices.Contains(referencedSchemas, structName) {
					continue
				}
				a.Components.SetSchema(structName, a.schemaFromStruct(s))
			}

			for aliasName, alias := range aliases {
				if !slices.Contains(referencedSchemas, aliasName) {
					continue
				}
				a.Components.SetSchema(aliasName, types.FormatSchema{
					Type: alias,
				})
			}

			it += 1
			if it > 100 {
				return fmt.Errorf("too many iterations")
			}
		}
	}

	return nil
}

func (a *api) LinkResponses() error {
	for path, routes := range a.Paths {
		for method, route := range routes {
			for code, resp := range route.Responses {
				if resp.Ref != "" || resp.Description != "" {
					continue
				}
				r := a.Components.Responses[code]
				a.Paths[path][method].Responses[code] = r
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

func (a *api) collectUrlVar(args []string) {
	var (
		name         = args[0]
		defaultValue = args[1]
		description  = strings.Join(args[2:], " ")
	)
	v := types.FormatServerVariable{
		Default:     defaultValue,
		Description: description,
	}
	a.Servers[len(a.Servers)-1].SetVariable(name, v)
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

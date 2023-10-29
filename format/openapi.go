package format

import (
	"docapi/collector"
	"docapi/types"
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	types.Format
	path                 string
	referencedComponents []string
}

func NewOpenAPI(path string) *OpenAPI {
	return &OpenAPI{
		Format: types.Format{
			Openapi: "3.0.0",
		},
		path: path,
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
		case types.CmdRoute:
			routes[cmd.Args[1]] = cmd.Args[0]
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
			tempHandler.RequestBody = types.FormatRequestBody{
				Description: strings.Join(cmd.Args[1:], " "),
				Required:    true,
				Content: map[string]types.FormatContent{
					"application/json": {
						Schema: types.FormatSchema{
							Ref: fmt.Sprintf("#/components/schemas/%s", cmd.Args[0]),
						},
					},
				},
			}
			f.referencedComponents = append(f.referencedComponents, cmd.Args[0])
		case types.CmdResponse:
			if tempHandler.Responses == nil {
				tempHandler.Responses = map[string]types.FormatResponse{}
			}
			if len(cmd.Args) > 1 {
				resp := types.FormatResponse{
					Description: strings.Join(cmd.Args[2:], " "),
				}
				content := types.FormatContent{}
				if strings.HasPrefix(cmd.Args[1], "[]") {
					content.Schema.Type = "array"
					component := cmd.Args[1][2:]
					content.Schema.Items = types.FormatItems{
						Ref: fmt.Sprintf("#/components/schemas/%s", component),
					}
					f.referencedComponents = append(f.referencedComponents, component)
				} else {
					content.Schema.Ref = cmd.Args[1]
				}
				resp.Content = map[string]types.FormatContent{}
				resp.Content["application/json"] = content
				tempHandler.Responses[cmd.Args[0]] = resp
			}
		case types.CmdEnd:
			handlers[tempHandlerID] = tempHandler
		}
	}

	f.Paths = map[string]types.FormatRoutes{}
	for handlerID, route := range routes {
		if f.Paths[route] == nil {
			f.Paths[route] = types.FormatRoutes{}
		}
		f.Paths[route][handlerMethods[handlerID]] = handlers[handlerID]
	}

	return nil
}

func (f *OpenAPI) CollectComponents(path string) error {
	tps, err := collector.NewTypesCollector().Run(path)
	if err != nil {
		return err
	}

	for tpName, tp := range tps {
		if !slices.Contains(f.referencedComponents, tpName) {
			continue
		}
		if f.Components.Schemas == nil {
			f.Components.Schemas = map[string]types.FormatSchema{}
		}
		schema := types.FormatSchema{
			Type: tp.Type,
		}
		for fieldName, field := range tp.Fields {
			if schema.Properties == nil {
				schema.Properties = map[string]types.FormatSchema{}
			}
			schema.Properties[fieldName] = types.FormatSchema{
				Type: field.Type,
			}
		}
		f.Components.Schemas[tpName] = schema
	}

	return nil
}

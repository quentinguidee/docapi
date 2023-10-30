package format

import (
	"fmt"
	"os"
	"slices"

	"github.com/quentinguidee/docapi/collector"
	"github.com/quentinguidee/docapi/types"

	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	path string
	apis []*api
}

func NewOpenAPI(path string) *OpenAPI {
	return &OpenAPI{
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

	// get all aliases
	var aliases []string
	for _, cmd := range commands {
		if cmd.ServerAlias != "" {
			aliases = append(aliases, cmd.ServerAlias)
		}
	}

	// initialize servers
	for _, alias := range aliases {
		f.apis = append(f.apis, newAPI(alias))
	}

	for _, a := range f.apis {
		cv := NewCommandsVisitor(a)
		for _, cmd := range commands {
			if cmd.ServerAlias != a.alias && cmd.ServerAlias != "" {
				continue
			}

			err := cv.Visit(cmd)
			if err != nil {
				return err
			}
		}
	}

	for _, a := range f.apis {
		a.Paths = map[string]types.FormatRoutes{}
		for handlerID, route := range a.routes {
			if a.Paths[route] == nil {
				a.Paths[route] = types.FormatRoutes{}
			}
			method := a.handlerMethods[handlerID]
			handler := a.handlers[handlerID]
			a.Paths[route][method] = handler
		}
	}

	return nil
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

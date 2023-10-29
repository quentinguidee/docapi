package generator

import (
	"docapi/generator/collector"
	"docapi/types"
	"encoding/json"
)

type Generator struct {
	api   *collector.ApiCollector
	types *collector.TypesCollector
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Run(path string) (string, error) {
	apiCollector := collector.NewAPICollector()
	err := apiCollector.Run(path)
	if err != nil {
		return "", err
	}
	g.api = apiCollector

	typesCollector := collector.NewTypesCollector()
	err = typesCollector.Run(path)
	if err != nil {
		return "", err
	}
	g.types = typesCollector

	return g.Output()
}

func (g *Generator) Output() (string, error) {
	out := struct {
		Api   interface{} `json:"api"`
		Types interface{} `json:"types"`
	}{}

	apis, err := g.api.Output()
	if err != nil {
		return "", err
	}
	out.Api = apis

	allTypes, err := g.types.Output()
	if err != nil {
		return "", err
	}

	keptTypes := map[string]types.Value{}
	for _, r := range g.api.Handlers {
		for _, resp := range r.Responses {
			if resp.Ref != "" {
				keptTypes[resp.Ref] = allTypes[resp.Ref]
			}
		}
	}
	out.Types = keptTypes

	s, err := json.MarshalIndent(out, "", "  ")
	return string(s), err
}
